package payout

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) ResendCode(ctx context.Context, opID int64, userID string, langCode string) error {
	confirmationCode := s.codeManager.GenerateCode()

	var email string
	if err := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &opID},
		func(ctx context.Context, op *model.Operation) error {
			if op.Type != model.OperationTypePayout {
				return fmt.Errorf("resend code called for unresolved type: %q", op.Type)
			}

			if op.UserID != userID {
				return fmt.Errorf("invalid user id: %v", userID)
			}

			if op.Status != model.OperationStatusNew {
				return perror.NewInternal().WithCode(
					perror.CodeUnresolvedStatusConflict,
				).WithDescription(
					fmt.Sprintf("cannot resend confirmation code with status %q", op.Status),
				)
			}

			op.ConfirmationCode = confirmationCode
			email, _ = op.Additional["email"].(string)

			return nil
		},
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return perror.NewInternal().WithCode(
				perror.CodeObjectNotFound,
			).WithDescription(
				fmt.Sprintf("payout with id %v not found", opID),
			)
		}
		return err
	}

	if err := s.codeManager.SendCode(ctx, opID, email, confirmationCode, langCode); err != nil {
		if saveErr := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &opID},
			func(ctx context.Context, op *model.Operation) error {
				op.Status = model.OperationStatusFailed
				op.FailReason = fmt.Sprintf("sending confirmation code: %v", err)
				return nil
			},
		); saveErr != nil {
			err = multierror.Append(err, saveErr)
		}
		return err
	}

	return nil
}
