package yookassa

import (
	"context"
	"errors"
	"fmt"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

func (i *Integration) getPaymentStatus(ctx context.Context, statusData model.GetOperationStatusData) (model.GetOperationStatusResult, error) {
	var result model.GetOperationStatusResult

	response, err := i.apiClient.GetPayment(ctx, statusData.ExternalID)
	if err != nil {
		if errors.Is(err, data.ErrPaymentNotFound) {
			var opExternalStatus model.OperationExternalStatus
			opExternalStatus, err = i.resolveExternalStatus(statusData.OperationType, statusData.ExternalMethod, statusData.CreatedAt, "", true)
			if err != nil {
				return result, fmt.Errorf("resolving external system operation status: %w", err)
			}

			result.ExternalStatus = opExternalStatus
			if opExternalStatus == model.OperationExternalStatusFailed {
				result.FailReason = model.OperationFailReasonTimeout
			}

			return result, nil
		} else {
			return result, err
		}
	}

	opExternalStatus, err := i.resolveExternalStatus(statusData.OperationType, statusData.ExternalMethod, statusData.CreatedAt, response.Status, false)
	if err != nil {
		return result, fmt.Errorf("resolving external system operation status: %w", err)
	}

	result.ExternalID = response.ID
	result.ExternalStatus = opExternalStatus

	switch opExternalStatus {
	case model.OperationExternalStatusSuccess:
		if response.IncomeAmount.Currency != statusData.Currency {
			return result, fmt.Errorf(
				"status response currency (%v) differs from operation currency (%v)",
				response.IncomeAmount.Currency, statusData.Currency,
			)
		}

		ch, err := i.channelResolver.Channel(statusData.ExternalMethod)
		if err != nil {
			return result, fmt.Errorf("resolving channel: %w", err)
		}

		result.ProcessedAt = response.CapturedAt.UTC()
		result.NewAmount = convert.BaseToCents(response.IncomeAmount.Value)
		result.Tool = ch.PaymentTool(statusData.UserID, statusData.ExternalMethod, response.PaymentMethod)
	case model.OperationExternalStatusFailed:
		result.FailReason = fmt.Sprintf("%v: %v", response.Cancellation.Party, response.Cancellation.Reason)
	}

	return result, nil
}
