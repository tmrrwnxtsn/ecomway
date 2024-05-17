package method

import (
	"context"
	"fmt"
	"slices"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type IntegrationClient interface {
	AllAvailableMethods(ctx context.Context, opType model.OperationType, currency string) ([]model.Method, error)
	AvailableMethodsByExternalSystem(ctx context.Context, opType model.OperationType, currency, externalSystem string) ([]model.Method, error)
}

type Service struct {
	integrationClient IntegrationClient
}

func NewService(integrationClient IntegrationClient) *Service {
	return &Service{
		integrationClient: integrationClient,
	}
}

func (s *Service) All(ctx context.Context, opType model.OperationType, currency string) ([]model.Method, error) {
	methods, err := s.integrationClient.AllAvailableMethods(ctx, opType, currency)
	if err != nil {
		return nil, err
	}

	result := slices.DeleteFunc(methods, func(method model.Method) bool {
		_, found := method.Limits[currency]
		return !found
	})

	// TODO: сортировка по часто используемым методам пользователя и избранным

	return result, nil
}

func (s *Service) GetOne(ctx context.Context, opType model.OperationType, currency, externalSystem, externalMethod string) (*model.Method, error) {
	methods, err := s.integrationClient.AvailableMethodsByExternalSystem(ctx, opType, currency, externalSystem)
	if err != nil {
		return nil, err
	}

	methodIdx := slices.IndexFunc(methods, func(m model.Method) bool {
		return m.ExternalSystem == externalSystem && m.ExternalMethod == externalMethod
	})
	if methodIdx == -1 {
		return nil, perror.NewInternal().WithCode(
			perror.CodeObjectNotFound,
		).WithDescription(
			fmt.Sprintf(
				"%v method not found for external system %q, external method %q",
				opType, externalSystem, externalMethod,
			),
		)
	}

	method := methods[methodIdx]

	if _, found := method.Limits[currency]; !found {
		return nil, fmt.Errorf(
			"%v external method %q does not support currency %q",
			opType, externalMethod, currency)
	}

	return &method, nil
}
