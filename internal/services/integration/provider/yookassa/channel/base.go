package channel

import (
	"fmt"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const defaultPaymentTimeoutToFailed = 90 * time.Minute

type baseChannel struct {
	code                   string
	paymentMethodType      string
	paymentTimeoutToFailed time.Duration
}

func newBaseChannel(cfg config.YooKassaChannelConfig) baseChannel {
	paymentTimeoutToFailed := time.Duration(cfg.PaymentTimeoutToFailedMin) * time.Minute
	if cfg.PaymentTimeoutToFailedMin == 0 {
		paymentTimeoutToFailed = defaultPaymentTimeoutToFailed
	}
	return baseChannel{
		code:                   cfg.Code,
		paymentMethodType:      cfg.PaymentMethodType,
		paymentTimeoutToFailed: paymentTimeoutToFailed,
	}
}

func (c baseChannel) CreatePaymentRequest(d model.CreatePaymentData) data.CreatePaymentRequest {
	return data.CreatePaymentRequest{
		Confirmation: data.PaymentConfirmation{
			Type:      data.PaymentConfirmationTypeRedirect,
			ReturnURL: d.ReturnURLs.Common,
			Locale:    getLocale(d.LangCode),
		},
		PaymentMethodData: data.PaymentMethod{
			Type: c.paymentMethodType,
		},
		Amount: data.PaymentAmount{
			Currency: d.Currency,
			Value:    convert.CentsToBase(d.Amount),
		},
		Description:       getDescription(d.LangCode, d.OperationID),
		Capture:           true,
		SavePaymentMethod: true,
	}
}

func (c baseChannel) PaymentTool(d model.GetOperationStatusData, _ data.GetPaymentResponse) *model.Tool {
	return &model.Tool{
		UserID:         d.UserID,
		ExternalMethod: d.ExternalMethod,
		Fake:           true,
	}
}

func (c baseChannel) PaymentTimeoutToFailed() time.Duration {
	return c.paymentTimeoutToFailed
}

func getDescription(langCode string, operationID int64) string {
	descriptionsFmt := map[string]string{
		"en": "Order payment №%v",
		"ru": "Оплата заказа №%v",
	}
	descriptionFmt, ok := descriptionsFmt[langCode]
	if !ok {
		descriptionFmt = descriptionsFmt["ru"]
	}
	return fmt.Sprintf(descriptionFmt, operationID)
}

func getLocale(langCode string) string {
	locales := map[string]string{
		"en": "en_US",
		"ru": "ru_RU",
	}
	locale, ok := locales[langCode]
	if !ok {
		locale = locales["ru"]
	}
	return locale
}
