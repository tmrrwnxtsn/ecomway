package payout

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) Success(ctx context.Context, data model.SuccessPayoutData) error {
	return s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
		func(ctx context.Context, op *model.Operation) error {
			if op.Type != model.OperationTypePayout {
				return fmt.Errorf("payout success called for unresolved type: %q", op.Type)
			}

			if op.Status == model.OperationStatusSuccess {
				slog.Warn(
					"duplicate payout success called",
					"operation_id", op.ID,
				)
				return nil
			}

			if op.Status == model.OperationStatusFailed {
				return errors.New("payout FAILED to SUCCESS not allowed")
			}

			if data.Tool != nil {
				if err := s.toolRepository.Save(ctx, data.Tool); err != nil {
					return fmt.Errorf("save tool: %w", err)
				}
			}

			op.Status = model.OperationStatusSuccess
			op.FailReason = ""
			op.ExternalID = data.ExternalID
			op.ExternalStatus = data.ExternalStatus

			if !data.ProcessedAt.IsZero() {
				op.ProcessedAt = data.ProcessedAt.UTC()
			} else {
				op.ProcessedAt = time.Now().UTC()
			}

			// TODO: осуществлять уведомление E-commerce системы
			slog.Info(
				"ecommerce system has been notified successfully",
				"operation_id", op.ID,
			)

			return nil
		},
	)
}
