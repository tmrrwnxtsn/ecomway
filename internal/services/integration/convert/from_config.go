package convert

import (
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
)

func MethodsFromConfig(methods []config.MethodConfig) []model.Method {
	result := make([]model.Method, 0, len(methods))

	for _, pbMethod := range methods {
		method := MethodFromProto(pbMethod)

		result = append(result, method)
	}

	return result
}

func MethodFromProto(method config.MethodConfig) model.Method {
	return model.Method{
		ID:             method.ID,
		DisplayedName:  method.DisplayedName,
		ExternalSystem: method.ExternalSystem,
		ExternalMethod: method.ExternalMethod,
		Limits:         LimitsFromProto(method.Limits),
		Commission:     CommissionFromProto(method.Commission),
	}
}

func LimitsFromProto(limits config.MethodLimitsConfig) model.Limits {
	return model.Limits{
		Currency:  limits.Currency,
		MinAmount: convert.BaseToCents(limits.MinAmount),
		MaxAmount: convert.BaseToCents(limits.MaxAmount),
	}
}

func CommissionFromProto(commission config.MethodCommissionConfig) model.Commission {
	return model.Commission{
		Type:     model.CommissionType(commission.Type),
		Currency: commission.Currency,
		Percent:  commission.Percent,
		Absolute: commission.Absolute,
		Message:  commission.Message,
	}
}
