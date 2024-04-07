package method

import (
	"context"
	"fmt"
	"slices"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type IntegrationClient interface {
	AllAvailableMethods(ctx context.Context, txType model.TransactionType, currency string) ([]model.Method, error)
	AvailableMethodsByExternalSystem(ctx context.Context, txType model.TransactionType, currency, externalSystem string) ([]model.Method, error)
}

type Service struct {
	integrationClient IntegrationClient
}

func NewService(integrationClient IntegrationClient) *Service {
	return &Service{
		integrationClient: integrationClient,
	}
}

func (s *Service) All(ctx context.Context, txType model.TransactionType, currency string) ([]model.Method, error) {
	methods, err := s.integrationClient.AllAvailableMethods(ctx, txType, currency)
	if err != nil {
		return nil, err
	}

	var result []model.Method
	for _, method := range methods {
		if _, found := method.Limits[currency]; !found {
			continue
		}

		result = append(result, method)
	}
	return result, nil
}

func (s *Service) GetOne(ctx context.Context, txType model.TransactionType, currency, externalSystem, externalMethod string) (*model.Method, error) {
	methods, err := s.integrationClient.AvailableMethodsByExternalSystem(ctx, txType, currency, externalSystem)
	if err != nil {
		return nil, err
	}

	methodIdx := slices.IndexFunc(methods, func(m model.Method) bool {
		return m.ExternalSystem == externalSystem && m.ExternalMethod == externalMethod
	})
	if methodIdx == -1 {
		return nil, fmt.Errorf(
			"%v method not found for external system %q, external method %q",
			txType, externalSystem, externalMethod,
		)
	}

	method := methods[methodIdx]

	if _, found := method.Limits[currency]; !found {
		return nil, fmt.Errorf(
			"%v external method %q does not support currency %q",
			txType, externalMethod, currency)
	}

	return &method, nil
}
