package yookassa

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

	switch statusData.OperationType {
	case model.OperationTypePayment:
		return i.getPaymentStatus(ctx, statusData)
	case model.OperationTypePayout:
		return i.getPayoutStatus(ctx, statusData)
	default:
		return result, fmt.Errorf("unresolved operation type: %q", statusData.OperationType)
	}
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

	switch yooKassaStatus {
	case data.StatusPending:
		return model.OperationExternalStatusPending, nil
	case data.StatusSucceeded:
		return model.OperationExternalStatusSuccess, nil
	case data.StatusCanceled:
		return model.OperationExternalStatusFailed, nil
	default:
		return "", fmt.Errorf("unresolved status: %q", yooKassaStatus)
	}
}
