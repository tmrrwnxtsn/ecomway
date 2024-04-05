package method

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type IntegrationClient interface {
	AvailableMethods(ctx context.Context, userID int64, txType model.TransactionType) ([]model.Method, error)
}

type Service struct {
	integrationClient IntegrationClient
}

func NewService(integrationClient IntegrationClient) *Service {
	return &Service{
		integrationClient: integrationClient,
	}
}

func (s *Service) AvailableMethods(ctx context.Context, userID int64, txType model.TransactionType) ([]model.Method, error) {
	return s.integrationClient.AvailableMethods(ctx, userID, txType)
}
