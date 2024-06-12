package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/gammazero/workerpool"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/config"
)

type RequestPayoutsTask struct {
	interval           time.Duration
	operationBatchSize int64
	maxWorkers         int
	operationService   OperationService
	payoutService      PayoutService
}

func NewRequestPayoutsTask(
	cfg config.SchedulerTaskConfig,
	operationService OperationService,
	payoutService PayoutService,
) *RequestPayoutsTask {
	return &RequestPayoutsTask{
		interval:           time.Duration(cfg.Interval) * time.Second,
		operationBatchSize: cfg.OperationBatchSize,
		maxWorkers:         cfg.MaxWorkers,
		operationService:   operationService,
		payoutService:      payoutService,
	}
}

func (t *RequestPayoutsTask) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				t.execute(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (t *RequestPayoutsTask) execute(ctx context.Context) {
	criteria := model.OperationCriteria{
		Statuses: &[]model.OperationStatus{model.OperationStatusConfirmed},
		Types:    &[]model.OperationType{model.OperationTypePayout},
		MaxCount: t.operationBatchSize,
	}

	log := slog.Default().With("task", "request_payouts")

	operations, err := t.operationService.All(ctx, criteria)
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

			log := log.With("operation_id", operation.ID)

			if operation.Status != model.OperationStatusConfirmed {
				return
			}

			requestPayoutData := model.CreatePayoutData{
				AdditionalData: operation.Additional,
				ExternalSystem: operation.ExternalSystem,
				ExternalMethod: operation.ExternalMethod,
				Currency:       operation.Currency,
				ToolID:         operation.ToolID,
				UserID:         operation.UserID,
				Amount:         operation.Amount,
				OperationID:    operation.ID,
			}

			if err := t.payoutService.RequestPayout(ctx, requestPayoutData); err != nil {
				log.Error(
					"failed to request payout",
					"error", err,
				)
				return
			}
		})
	}

	wp.StopWait()
}
