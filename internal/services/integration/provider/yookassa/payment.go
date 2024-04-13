package yookassa

import (
	"context"
	"fmt"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
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

	if response.Status != data.PaymentStatusPending {
		return result, fmt.Errorf("unresolved status on payment creation: %q", response.Status)
	}

	result.RedirectURL = response.ConfirmationURL
	result.ExternalID = response.ID
	result.ExternalStatus = model.OperationExternalStatusPending

	return result, nil
}
