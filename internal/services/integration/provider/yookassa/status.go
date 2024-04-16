package yookassa

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

func (i *Integration) GetOperationStatus(ctx context.Context, statusData model.GetOperationStatusData) (model.GetOperationStatusResult, error) {
	var result model.GetOperationStatusResult

	if statusData.ExternalID == "" {
		slog.Warn("external id is empty", "operation_id", statusData.OperationID)

		opExternalStatus, err := i.resolveExternalStatus(statusData.OperationType, statusData.ExternalMethod, statusData.CreatedAt, "", true)
		if err != nil {
			return result, fmt.Errorf("resolving external system operation status: %w", err)
		}

		result.ExternalStatus = opExternalStatus
		if opExternalStatus == model.OperationExternalStatusFailed {
			result.FailReason = model.OperationFailReasonTimeout
		}

		return result, nil
	}

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

func (i *Integration) resolveExternalStatus(
	opType model.OperationType,
	externalMethod string,
	createdAt time.Time,
	yooKassaStatus string,
	opNotFound bool,
) (model.OperationExternalStatus, error) {
	if opNotFound {
		if opType == model.OperationTypePayment {
			ch, err := i.channelResolver.Channel(externalMethod)
			if err != nil {
				return "", fmt.Errorf("resolving channel: %w", err)
			}

			opTimeoutToFailed := ch.PaymentTimeoutToFailed()

			nowTime := time.Now().UTC()
			createdAtTime := createdAt.UTC()

			if createdAtTime.Add(opTimeoutToFailed).Before(nowTime) {
				return model.OperationExternalStatusFailed, nil
			}
		}
		return model.OperationExternalStatusUnknown, nil
	}

	switch opType {
	case model.OperationTypePayment:
		switch yooKassaStatus {
		case data.PaymentStatusPending:
			return model.OperationExternalStatusPending, nil
		case data.PaymentStatusSucceeded:
			return model.OperationExternalStatusSuccess, nil
		case data.PaymentStatusCanceled:
			return model.OperationExternalStatusFailed, nil
		default:
			return "", fmt.Errorf("unresolved status: %q", yooKassaStatus)
		}
	default:
		return "", fmt.Errorf("unresolved operation type: %q", opType)
	}
}
