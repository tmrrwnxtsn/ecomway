package yookassa

import (
	"context"
	"fmt"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (i *Integration) CreatePayout(ctx context.Context, payoutData model.CreatePayoutData) (model.CreatePayoutResult, error) {
	var result model.CreatePayoutResult

	channel, err := i.channelResolver.Channel(payoutData.ExternalMethod)
	if err != nil {
		return result, fmt.Errorf("resolving channel: %w", err)
	}

	request := channel.CreatePayoutRequest(payoutData)

	response, err := i.apiClient.CreatePayout(ctx, request)
	if err != nil {
		return result, fmt.Errorf("sending external system request: %w", err)
	}

	opExternalStatus, err := i.resolveExternalStatus(model.OperationTypePayout, payoutData.ExternalMethod, time.Time{}, response.Status, false)
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
