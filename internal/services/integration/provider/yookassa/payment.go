package yookassa

import (
	"context"
	"fmt"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (i *Integration) CreatePayment(ctx context.Context, paymentData model.CreatePaymentData) (model.CreatePaymentResult, error) {
	var result model.CreatePaymentResult

	channel, err := i.channelResolver.Channel(paymentData.ExternalMethod)
	if err != nil {
		return result, fmt.Errorf("resolving channel: %w", err)
	}

	request := channel.CreatePaymentRequest(paymentData)

	response, err := i.apiClient.CreatePayment(ctx, request)
	if err != nil {
		// TODO: создавать кастомную ошибку для проброса в engine (причина отклонения) и gateway (отображение ошибки)
		//var errorResponse *data.ErrorResponse
		//if errors.As(err, &errorResponse) {
		//	pmnterror.NewExternal()
		//	failReason := errorResponse.Error()
		//	log.Warnf("creating invoice: %v", failReason)
		//	return common.ErrorResult(failReason), nil
		//} else {
		//	err = log.ErrorfErr("creating invoice: %v", err)
		//	return common.TechnicalErrorResult(err), nil
		//}
		return result, fmt.Errorf("sending external system request: %w", err)
	}

	opExternalStatus, err := i.resolveExternalStatus(model.OperationTypePayment, paymentData.ExternalMethod, time.Time{}, response.Status, false)
	if err != nil {
		return result, fmt.Errorf("resolving external system operation status: %w", err)
	}

	result.ExternalID = response.ID
	result.ExternalStatus = opExternalStatus

	switch opExternalStatus {
	case model.OperationExternalStatusSuccess:
		if response.IncomeAmount.Currency != paymentData.Currency {
			return result, fmt.Errorf(
				"create payment response currency (%v) differs from operation currency (%v)",
				response.IncomeAmount.Currency, paymentData.Currency,
			)
		}

		ch, err := i.channelResolver.Channel(paymentData.ExternalMethod)
		if err != nil {
			return result, fmt.Errorf("resolving channel: %w", err)
		}

		result.ProcessedAt = response.CapturedAt.UTC()
		result.NewAmount = convert.BaseToCents(response.IncomeAmount.Value)
		result.Tool = ch.PaymentTool(paymentData.UserID, paymentData.ExternalMethod, response.PaymentMethod)
	case model.OperationExternalStatusFailed:
		result.FailReason = fmt.Sprintf("%v: %v", response.Cancellation.Party, response.Cancellation.Reason)
	default:
		result.RedirectURL = response.ConfirmationURL
	}

	return result, nil
}
