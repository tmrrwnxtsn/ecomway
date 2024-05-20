package convert

import (
	"google.golang.org/protobuf/types/known/structpb"

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

func OperationStatusToProto(opStatus model.OperationStatus) pb.OperationStatus {
	switch opStatus {
	case model.OperationStatusNew:
		return pb.OperationStatus_OPERATION_STATUS_NEW
	case model.OperationStatusSuccess:
		return pb.OperationStatus_OPERATION_STATUS_SUCCESS
	case model.OperationStatusFailed:
		return pb.OperationStatus_OPERATION_STATUS_FAILED
	case model.OperationStatusConfirmed:
		return pb.OperationStatus_OPERATION_STATUS_CONFIRMED
	default:
		return -1
	}
}

func OperationExternalStatusToProto(opExternalStatus model.OperationExternalStatus) pb.OperationExternalStatus {
	switch opExternalStatus {
	case model.OperationExternalStatusPending:
		return pb.OperationExternalStatus_OPERATION_EXTERNAL_STATUS_PENDING
	case model.OperationExternalStatusSuccess:
		return pb.OperationExternalStatus_OPERATION_EXTERNAL_STATUS_SUCCESS
	case model.OperationExternalStatusFailed:
		return pb.OperationExternalStatus_OPERATION_EXTERNAL_STATUS_FAILED
	case model.OperationExternalStatusUnknown:
		return pb.OperationExternalStatus_OPERATION_EXTERNAL_STATUS_UNKNOWN
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
	result := &pb.Commission{
		Type:     CommissionTypeToProto(commission.Type),
		Percent:  commission.Percent,
		Absolute: commission.Absolute,
		Message:  commission.Message,
	}

	if commission.Currency != "" {
		result.Currency = &commission.Currency
	}

	return result
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
	pbReturnURLs := &pb.ReturnURLs{
		Common: returnURLs.Common,
	}

	if returnURLs.Success != "" {
		pbReturnURLs.Success = &returnURLs.Success
	}

	if returnURLs.Fail != "" {
		pbReturnURLs.Fail = &returnURLs.Fail
	}

	return pbReturnURLs
}

func ToolTypeToProto(toolType model.ToolType) pb.ToolType {
	switch toolType {
	case model.ToolTypeBankCard:
		return pb.ToolType_BANK_CARD
	case model.ToolTypeWallet:
		return pb.ToolType_WALLET
	default:
		return -1
	}
}

func ToolStatusToProto(toolStatus model.ToolStatus) pb.ToolStatus {
	switch toolStatus {
	case model.ToolStatusActive:
		return pb.ToolStatus_ACTIVE
	case model.ToolStatusRemovedByClient:
		return pb.ToolStatus_REMOVED_BY_USER
	case model.ToolStatusPendingRecovery:
		return pb.ToolStatus_PENDING_RECOVERY
	case model.ToolStatusRemovedByAdministrator:
		return pb.ToolStatus_REMOVED_BY_ADMINISTRATOR
	default:
		return -1
	}
}

func ToolToProto(tool *model.Tool) *pb.Tool {
	result := &pb.Tool{
		Id:             tool.ID,
		UserId:         tool.UserID,
		ExternalMethod: tool.ExternalMethod,
		Displayed:      tool.Displayed,
		Name:           tool.Name,
		Status:         ToolStatusToProto(tool.Status),
		Fake:           tool.Fake,
		CreatedAt:      tool.CreatedAt.UTC().Unix(),
		UpdatedAt:      tool.UpdatedAt.UTC().Unix(),
	}

	if tool.Type != "" {
		pbToolType := ToolTypeToProto(tool.Type)
		result.Type = &pbToolType
	}

	if tool.Details != nil {
		pbDetails, err := structpb.NewStruct(tool.Details)
		if err == nil {
			result.Details = pbDetails
		}
	}

	return result
}
