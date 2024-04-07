package convert

import (
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
)

func MethodsFromConfig(methods []config.MethodConfig) []model.Method {
	result := make([]model.Method, 0, len(methods))

	for _, pbMethod := range methods {
		method := MethodFromConfig(pbMethod)

		result = append(result, method)
	}

	return result
}

func MethodFromConfig(method config.MethodConfig) model.Method {
	return model.Method{
		ID:             method.ID,
		DisplayedName:  method.DisplayedName,
		ExternalSystem: method.ExternalSystem,
		ExternalMethod: method.ExternalMethod,
		Limits:         LimitsFromConfig(method.Limits),
		Commission:     CommissionFromConfig(method.Commission),
	}
}

func LimitsFromConfig(limits map[string]config.MethodLimitsConfig) map[string]model.Limits {
	if len(limits) == 0 {
		return nil
	}

	result := make(map[string]model.Limits, len(limits))
	for currency, l := range limits {
		result[currency] = model.Limits{
			MinAmount: convert.BaseToCents(l.MinAmount),
			MaxAmount: convert.BaseToCents(l.MaxAmount),
		}
	}
	return result
}

func CommissionFromConfig(commission config.MethodCommissionConfig) model.Commission {
	return model.Commission{
		Type:     model.CommissionType(commission.Type),
		Currency: commission.Currency,
		Percent:  commission.Percent,
		Absolute: commission.Absolute,
		Message:  commission.Message,
	}
}
