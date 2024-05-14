package payment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) Success(ctx context.Context, data model.SuccessPaymentData) error {
	return s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
		func(ctx context.Context, op *model.Operation) error {
			if op.Type != model.OperationTypePayment {
				return fmt.Errorf("payment success called for unresolved type: %q", op.Type)
			}

			if op.Status == model.OperationStatusSuccess {
				slog.Warn(
					"duplicate payment success called",
					"operation_id", op.ID,
				)
				return nil
			}

			if op.Status == model.OperationStatusFailed {
				return errors.New("payment FAILED to SUCCESS not allowed")
			}

			if data.Tool != nil {
				if op.ToolID != "" && data.Tool.ID != op.ToolID {
					return fmt.Errorf("operation tool %q differs from payment tool %q", op.ToolID, data.Tool.ID)
				}

				tool, err := s.toolRepository.GetOne(ctx, data.Tool.ID, data.Tool.UserID, data.Tool.ExternalMethod)
				switch {
				case errors.Is(err, sql.ErrNoRows):
					if err = s.toolRepository.Create(ctx, data.Tool); err != nil {
						return fmt.Errorf("create tool into db: %w", err)
					}
					op.ToolID = data.Tool.ID
				case err != nil:
					return fmt.Errorf("get tool from db: %w", err)
				case tool.CanBeRecovered():
					tool.Status = model.ToolStatusActive

					if err = s.toolRepository.Update(ctx, tool); err != nil {
						return fmt.Errorf("update tool into db: %w", err)
					}
					op.ToolID = data.Tool.ID

					slog.Info(
						"payment tool has been recovered",
						"operation_id", op.ID,
						"tool_id", op.ToolID,
					)
				}
			}

			op.Status = model.OperationStatusSuccess
			op.FailReason = ""
			op.ExternalID = data.ExternalID
			op.ExternalStatus = data.ExternalStatus

			if data.NewAmount > 0 {
				slog.Info(
					"payment amount changed",
					"operation_id", op.ID,
					"old_amount", op.Amount,
					"new_amount", data.NewAmount,
				)
				op.Amount = data.NewAmount
			}

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
