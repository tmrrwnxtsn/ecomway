package convert

import (
	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func OperationTypeToProto(opType model.OperationType) pb.OperationType {
	switch opType {
	case model.OperationTypePayment:
		return pb.OperationType_PAYMENT
	case model.OperationTypePayout:
		return pb.OperationType_PAYOUT
	default:
		return -1
	}
}

func OperationExternalStatusToProto(opExternalStatus model.OperationExternalStatus) pb.OperationExternalStatus {
	switch opExternalStatus {
	case model.OperationExternalStatusPending:
		return pb.OperationExternalStatus_PENDING
	case model.OperationExternalStatusSuccess:
		return pb.OperationExternalStatus_SUCCESS
	case model.OperationExternalStatusFailed:
		return pb.OperationExternalStatus_FAILED
	case model.OperationExternalStatusUnknown:
		return pb.OperationExternalStatus_UNKNOWN
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

func ReturnURLsToProto(returnURLs model.ReturnURLs) *pb.ReturnURLs {
	return &pb.ReturnURLs{
		Common:  returnURLs.Common,
		Success: returnURLs.Success,
		Fail:    returnURLs.Fail,
	}
}
