package convert

import (
	"time"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func OperationTypeFromProto(opType pb.OperationType) model.OperationType {
	switch opType {
	case pb.OperationType_PAYMENT:
		return model.OperationTypePayment
	case pb.OperationType_PAYOUT:
		return model.OperationTypePayout
	default:
		return ""
	}
}

func OperationStatusFromProto(opStatus pb.OperationStatus) model.OperationStatus {
	switch opStatus {
	case pb.OperationStatus_OPERATION_STATUS_NEW:
		return model.OperationStatusNew
	case pb.OperationStatus_OPERATION_STATUS_SUCCESS:
		return model.OperationStatusSuccess
	case pb.OperationStatus_OPERATION_STATUS_FAILED:
		return model.OperationStatusFailed
	case pb.OperationStatus_OPERATION_STATUS_CONFIRMED:
		return model.OperationStatusConfirmed
	default:
		return ""
	}
}

func OperationExternalStatusFromProto(opExternalStatus pb.OperationExternalStatus) model.OperationExternalStatus {
	switch opExternalStatus {
	case pb.OperationExternalStatus_OPERATION_EXTERNAL_STATUS_PENDING:
		return model.OperationExternalStatusPending
	case pb.OperationExternalStatus_OPERATION_EXTERNAL_STATUS_SUCCESS:
		return model.OperationExternalStatusSuccess
	case pb.OperationExternalStatus_OPERATION_EXTERNAL_STATUS_FAILED:
		return model.OperationExternalStatusFailed
	case pb.OperationExternalStatus_OPERATION_EXTERNAL_STATUS_UNKNOWN:
		return model.OperationExternalStatusUnknown
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
		Currency: commission.GetCurrency(),
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

func ReturnURLsFromProto(returnURLs *pb.ReturnURLs) model.ReturnURLs {
	return model.ReturnURLs{
		Common:  returnURLs.GetCommon(),
		Success: returnURLs.GetSuccess(),
		Fail:    returnURLs.GetFail(),
	}
}

func ToolTypeFromProto(toolType pb.ToolType) model.ToolType {
	switch toolType {
	case pb.ToolType_BANK_CARD:
		return model.ToolTypeBankCard
	case pb.ToolType_WALLET:
		return model.ToolTypeWallet
	default:
		return ""
	}
}

func ToolStatusFromProto(toolStatus pb.ToolStatus) model.ToolStatus {
	switch toolStatus {
	case pb.ToolStatus_ACTIVE:
		return model.ToolStatusActive
	case pb.ToolStatus_REMOVED_BY_USER:
		return model.ToolStatusRemovedByClient
	case pb.ToolStatus_PENDING_RECOVERY:
		return model.ToolStatusPendingRecovery
	case pb.ToolStatus_REMOVED_BY_ADMINISTRATOR:
		return model.ToolStatusRemovedByAdministrator
	default:
		return ""
	}
}

func ToolFromProto(tool *pb.Tool) *model.Tool {
	result := &model.Tool{
		ID:             tool.GetId(),
		UserID:         tool.GetUserId(),
		ExternalMethod: tool.GetExternalMethod(),
		Displayed:      tool.GetDisplayed(),
		Name:           tool.GetName(),
		Status:         ToolStatusFromProto(tool.GetStatus()),
		Fake:           tool.GetFake(),
		CreatedAt:      time.Unix(tool.GetCreatedAt(), 0).UTC(),
		UpdatedAt:      time.Unix(tool.GetUpdatedAt(), 0).UTC(),
	}

	if tool.Type != nil {
		result.Type = ToolTypeFromProto(tool.GetType())
	}

	if tool.Details != nil {
		result.Details = tool.Details.AsMap()
	}

	return result
}
