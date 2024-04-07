package convert

import (
	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func TransactionTypeToProto(txType model.TransactionType) pb.TransactionType {
	switch txType {
	case model.TransactionTypePayment:
		return pb.TransactionType_PAYMENT
	case model.TransactionTypePayout:
		return pb.TransactionType_PAYOUT
	default:
		return -1
	}
}

func MethodsToProto(methods []model.Method) []*pb.Method {
	result := make([]*pb.Method, 0, len(methods))

	for _, method := range methods {
		pbMethod := MethodToProto(method)

		result = append(result, pbMethod)
	}

	return result
}

func MethodToProto(method model.Method) *pb.Method {
	return &pb.Method{
		Id:             method.ID,
		DisplayedName:  method.DisplayedName,
		ExternalSystem: method.ExternalSystem,
		ExternalMethod: method.ExternalMethod,
		Limits:         LimitsToProto(method.Limits),
		Commission:     CommissionToProto(method.Commission),
	}
}

func LimitsToProto(limits map[string]model.Limits) map[string]*pb.Limits {
	if len(limits) == 0 {
		return nil
	}

	result := make(map[string]*pb.Limits, len(limits))
	for currency, l := range limits {
		result[currency] = &pb.Limits{
			MinAmount: l.MinAmount,
			MaxAmount: l.MaxAmount,
		}
	}
	return result
}

func CommissionToProto(commission model.Commission) *pb.Commission {
	return &pb.Commission{
		Type:     CommissionTypeToProto(commission.Type),
		Currency: commission.Currency,
		Percent:  commission.Percent,
		Absolute: commission.Absolute,
		Message:  commission.Message,
	}
}

func CommissionTypeToProto(commissionType model.CommissionType) pb.CommissionType {
	switch commissionType {
	case model.CommissionTypePercent:
		return pb.CommissionType_PERCENT
	case model.CommissionTypeFixed:
		return pb.CommissionType_FIXED
	case model.CommissionTypeCombined:
		return pb.CommissionType_COMBINED
	case model.CommissionTypeText:
		return pb.CommissionType_TEXT
	default:
		return -1
	}
}
