package payout

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) Fail(ctx context.Context, data model.FailPayoutData) error {
	return s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
		func(ctx context.Context, op *model.Operation) error {
			if op.Type != model.OperationTypePayout {
				return fmt.Errorf("payout fail called for unresolved type: %q", op.Type)
			}

			if op.Status == model.OperationStatusFailed {
				slog.Warn(
					"duplicate payout fail called",
					"operation_id", op.ID,
				)
				return nil
			}

			if op.Status == model.OperationStatusSuccess {
				return errors.New("payout SUCCESS to FAILED not allowed")
			}

			// TODO: осуществлять уведомление E-commerce системы
			slog.Info(
				"ecommerce system has been notified successfully",
				"operation_id", op.ID,
			)

			op.Status = model.OperationStatusFailed
			op.FailReason = data.FailReason
			op.ExternalID = data.ExternalID
			op.ExternalStatus = data.ExternalStatus

			if !op.ProcessedAt.IsZero() {
				op.ProcessedAt = time.Time{}
			}

			return nil
		},
	)
}
