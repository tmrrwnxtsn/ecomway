package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/gammazero/workerpool"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/config"
)

type OperationRepository interface {
	All(ctx context.Context, criteria model.OperationCriteria) ([]*model.Operation, error)
}

type IntegrationClient interface {
	GetOperationStatus(ctx context.Context, data model.GetOperationStatusData) (model.GetOperationStatusResult, error)
}

type PaymentService interface {
	Success(ctx context.Context, data model.SuccessPaymentData) error
	Fail(ctx context.Context, data model.FailPaymentData) error
}

type FinalizeOperationsTask struct {
	interval                 time.Duration
	operationBatchSize       int64
	maxWorkers               int
	actualizeStatusIntervals map[time.Duration]time.Duration
	externalSystemLifetime   map[string]time.Duration
	operationRepository      OperationRepository
	integrationClient        IntegrationClient
	paymentService           PaymentService
}

func NewFinalizeOperationsTask(
	cfg config.SchedulerTaskConfig,
	operationRepository OperationRepository,
	integrationClient IntegrationClient,
	paymentService PaymentService,
) *FinalizeOperationsTask {
	task := &FinalizeOperationsTask{
		interval:                 time.Duration(cfg.Interval) * time.Second,
		operationBatchSize:       cfg.OperationBatchSize,
		maxWorkers:               cfg.MaxWorkers,
		actualizeStatusIntervals: make(map[time.Duration]time.Duration, len(cfg.ActualizeStatusIntervals)),
		externalSystemLifetime:   make(map[string]time.Duration, len(cfg.ExternalSystemLifetime)),
		operationRepository:      operationRepository,
		integrationClient:        integrationClient,
		paymentService:           paymentService,
	}

	for slot, interval := range cfg.ActualizeStatusIntervals {
		task.actualizeStatusIntervals[time.Duration(slot)*time.Second] = time.Duration(interval) * time.Second
	}

	for externalSystem, lifetime := range cfg.ExternalSystemLifetime {
		task.externalSystemLifetime[externalSystem] = time.Duration(lifetime) * time.Second
	}

	return task
}

func (t *FinalizeOperationsTask) Start(ctx context.Context) {
	for externalSystem := range t.externalSystemLifetime {
		go func(externalSystem string) {
			ticker := time.NewTicker(t.interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					t.execute(ctx, externalSystem)
				case <-ctx.Done():
					return
				}
			}
		}(externalSystem)
	}
}

func (t *FinalizeOperationsTask) execute(ctx context.Context, externalSystem string) {
	criteria := model.OperationCriteria{
		StatusesByType: map[model.OperationType][]model.OperationStatus{
			model.OperationTypePayment: {model.OperationStatusNew},
		},
		ExternalSystems: &[]string{externalSystem},
		MaxCount:        t.operationBatchSize,
	}

	log := slog.Default().With("task", "finalize_operations")

	operations, err := t.operationRepository.All(ctx, criteria)
	if err != nil {
		log.Error(
			"failed to receive operations by criteria",
			"error", err,
		)
		return
	}

	wp := workerpool.New(t.maxWorkers)

	for _, operation := range operations {
		operation := operation

		wp.Submit(func() {
			ctx := context.Background()

			if !t.needsActualizeExternalStatus(operation) {
				return
			}

			data := model.GetOperationStatusData{
				CreatedAt:      operation.CreatedAt,
				ExternalID:     operation.ExternalID,
				ExternalSystem: operation.ExternalSystem,
				ExternalMethod: operation.ExternalMethod,
				Currency:       operation.Currency,
				OperationType:  operation.Type,
				OperationID:    operation.ID,
				UserID:         operation.UserID,
				Amount:         operation.Amount,
			}

			result, err := t.integrationClient.GetOperationStatus(ctx, data)
			if err != nil {
				log.Error(
					"failed to get operation external status",
					"error", err,
					"operation_id", operation.ID,
				)
				return
			}

			switch result.ExternalStatus {
			case model.OperationExternalStatusSuccess:
				data := model.SuccessPaymentData{
					ProcessedAt:    result.ProcessedAt,
					ExternalID:     result.ExternalID,
					ExternalStatus: result.ExternalStatus,
					OperationID:    operation.ID,
					NewAmount:      result.NewAmount,
					Tool:           result.Tool,
				}

				if err = t.paymentService.Success(ctx, data); err != nil {
					log.Error(
						"failed to success payment",
						"error", err,
						"operation_id", operation.ID,
					)
					return
				}
			case model.OperationExternalStatusFailed:
				data := model.FailPaymentData{
					ExternalID:     result.ExternalID,
					ExternalStatus: result.ExternalStatus,
					FailReason:     result.FailReason,
					OperationID:    operation.ID,
				}

				if err = t.paymentService.Fail(ctx, data); err != nil {
					log.Error(
						"failed to fail payment",
						"error", err,
						"operation_id", operation.ID,
					)
					return
				}
			}

		})
	}

	wp.StopWait()
}

func (t *FinalizeOperationsTask) needsActualizeExternalStatus(op *model.Operation) bool {
	if time.Since(op.UpdatedAt) < t.externalSystemLifetime[op.ExternalSystem] {
		return false
	}

	var (
		currSlot       time.Duration
		resultInterval time.Duration
	)
	for slot, interval := range t.actualizeStatusIntervals {
		if time.Since(op.CreatedAt) < slot {
			continue
		}

		if slot > currSlot {
			currSlot = slot
			resultInterval = interval
		}
	}
	if currSlot == 0 {
		return true
	}

	return time.Since(op.UpdatedAt) > resultInterval
}
