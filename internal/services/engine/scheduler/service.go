package scheduler

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type OperationService interface {
	All(ctx context.Context, criteria model.OperationCriteria) ([]*model.Operation, error)
}

type IntegrationClient interface {
	GetOperationStatus(ctx context.Context, data model.GetOperationStatusData) (model.GetOperationStatusResult, error)
}

type PaymentService interface {
	Success(ctx context.Context, data model.SuccessPaymentData) error
	Fail(ctx context.Context, data model.FailPaymentData) error
}

type PayoutService interface {
	Fail(ctx context.Context, data model.FailPayoutData) error
	RequestPayout(ctx context.Context, data model.CreatePayoutData) error
}
