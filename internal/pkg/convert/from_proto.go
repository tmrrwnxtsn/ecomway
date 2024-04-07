package convert

import (
	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func TransactionTypeFromProto(txType pb.TransactionType) model.TransactionType {
	switch txType {
	case pb.TransactionType_PAYMENT:
		return model.TransactionTypePayment
	case pb.TransactionType_PAYOUT:
		return model.TransactionTypePayout
	default:
		return ""
	}
}

func MethodsFromProto(methods []*pb.Method) []model.Method {
	result := make([]model.Method, 0, len(methods))

	for _, pbMethod := range methods {
		method := MethodFromProto(pbMethod)

		result = append(result, method)
	}

	return result
}

func MethodFromProto(method *pb.Method) model.Method {
	return model.Method{
		ID:             method.GetId(),
		DisplayedName:  method.GetDisplayedName(),
		ExternalSystem: method.GetExternalSystem(),
		ExternalMethod: method.GetExternalMethod(),
		Limits:         LimitsFromProto(method.GetLimits()),
		Commission:     CommissionFromProto(method.GetCommission()),
	}
}

func LimitsFromProto(limits map[string]*pb.Limits) map[string]model.Limits {
	if len(limits) == 0 {
		return nil
	}

	result := make(map[string]model.Limits, len(limits))
	for currency, l := range limits {
		result[currency] = model.Limits{
			MinAmount: l.GetMinAmount(),
			MaxAmount: l.GetMaxAmount(),
		}
	}
	return result
}

func CommissionFromProto(commission *pb.Commission) model.Commission {
	if commission == nil {
		return model.Commission{}
	}
	return model.Commission{
		Type:     CommissionTypeFromProto(commission.Type),
		Currency: commission.Currency,
		Percent:  commission.Percent,
		Absolute: commission.Absolute,
		Message:  commission.Message,
	}
}

func CommissionTypeFromProto(commissionType pb.CommissionType) model.CommissionType {
	switch commissionType {
	case pb.CommissionType_PERCENT:
		return model.CommissionTypePercent
	case pb.CommissionType_FIXED:
		return model.CommissionTypeFixed
	case pb.CommissionType_COMBINED:
		return model.CommissionTypeCombined
	case pb.CommissionType_TEXT:
		return model.CommissionTypeText
	default:
		return ""
	}
}
