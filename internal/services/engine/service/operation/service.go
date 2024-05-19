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
