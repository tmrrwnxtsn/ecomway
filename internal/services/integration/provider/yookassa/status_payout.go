package yookassa

import (
	"context"
	"errors"
	"fmt"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

func (i *Integration) getPayoutStatus(ctx context.Context, statusData model.GetOperationStatusData) (model.GetOperationStatusResult, error) {
	var result model.GetOperationStatusResult

	response, err := i.apiClient.GetPayout(ctx, statusData.ExternalID)
	if err != nil {
		if errors.Is(err, data.ErrPayoutNotFound) {
			var opExternalStatus model.OperationExternalStatus
			opExternalStatus, err = i.resolveExternalStatus(statusData.OperationType, statusData.ExternalMethod, statusData.CreatedAt, "", true)
			if err != nil {
				return result, fmt.Errorf("resolving external system operation status: %w", err)
			}

			result.ExternalStatus = opExternalStatus

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
		result.ProcessedAt = response.CapturedAt.UTC()
	case model.OperationExternalStatusFailed:
		result.FailReason = fmt.Sprintf("%v: %v", response.Cancellation.Party, response.Cancellation.Reason)
	}

	return result, nil
}
