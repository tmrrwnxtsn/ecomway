package payment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) Fail(ctx context.Context, data model.FailPaymentData) error {
	return s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
		func(ctx context.Context, op *model.Operation) error {
			if op.Type != model.OperationTypePayment {
				return fmt.Errorf("payment fail called for unresolved type: %q", op.Type)
			}

			if op.Status == model.OperationStatusFailed {
				slog.Warn(
					"duplicate payment fail called",
					"operation_id", op.ID,
				)
				return nil
			}

			if op.Status == model.OperationStatusSuccess {
				return errors.New("payment SUCCESS to FAILED not allowed")
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
