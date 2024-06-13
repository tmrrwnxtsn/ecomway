package operation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Repository interface {
	All(ctx context.Context, criteria model.OperationCriteria) ([]*model.Operation, error)
	AllForReport(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error)
	GetOneWithoutLock(ctx context.Context, criteria model.OperationCriteria) (*model.Operation, error)
	AcquireOneLocked(ctx context.Context, criteria model.OperationCriteria, script model.ScriptAcquiredFor) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) All(ctx context.Context, criteria model.OperationCriteria) ([]*model.Operation, error) {
	return s.repository.All(ctx, criteria)
}

func (s *Service) AllForReport(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error) {
	return s.repository.AllForReport(ctx, criteria)
}

func (s *Service) GetOne(ctx context.Context, criteria model.OperationCriteria) (*model.Operation, error) {
	operation, err := s.repository.GetOneWithoutLock(ctx, criteria)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var description string
			switch {
			case criteria.ID != nil:
				description = fmt.Sprintf("operation with id %v not found", *criteria.ID)
			case criteria.ExternalID != nil:
				description = fmt.Sprintf("operation with external id %v not found", *criteria.ExternalID)
			default:
				description = "operation not found"
			}
			return nil, perror.NewInternal().WithCode(
				perror.CodeObjectNotFound,
			).WithDescription(
				description,
			)
		}
		return nil, err
	}
	return operation, nil
}

func (s *Service) ChangeStatus(ctx context.Context, id int64, newStatus model.OperationStatus, newExternalStatus model.OperationExternalStatus) (model.OperationChangeStatusResult, error) {
	var result model.OperationChangeStatusResult
	if err := s.repository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &id},
		func(ctx context.Context, op *model.Operation) error {
			switch newStatus {
			case model.OperationStatusNew:
				switch newExternalStatus {
				case model.OperationExternalStatusSuccess:
					if op.Type == model.OperationTypePayment {
						result = model.OperationChangeStatusResultSuccessPayment
						return nil
					} else {
						return fmt.Errorf("can't process payout in %q status - %q external status", newStatus, newExternalStatus)
					}
				case model.OperationExternalStatusFailed:
					if op.Type == model.OperationTypePayment {
						result = model.OperationChangeStatusResultFailPayment
						return nil
					} else {
						result = model.OperationChangeStatusResultFailPayout
						return nil
					}
				}
			case model.OperationStatusConfirmed, model.OperationStatusPending:
				switch newExternalStatus {
				case model.OperationExternalStatusSuccess:
					if op.Type == model.OperationTypePayout {
						result = model.OperationChangeStatusResultSuccessPayout
						return nil
					} else {
						return fmt.Errorf("can't process payment in %q status - %q external status", newStatus, newExternalStatus)
					}
				case model.OperationExternalStatusFailed:
					if op.Type == model.OperationTypePayout {
						result = model.OperationChangeStatusResultFailPayout
						return nil
					} else {
						return fmt.Errorf("can't process payment in %q status - %q external status", newStatus, newExternalStatus)
					}
				}
			}
			return nil
		},
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", perror.NewInternal().WithCode(
				perror.CodeObjectNotFound,
			).WithDescription(
				fmt.Sprintf("payout with id %v not found", id),
			)
		}
		return "", err
	}
	return result, nil
}
