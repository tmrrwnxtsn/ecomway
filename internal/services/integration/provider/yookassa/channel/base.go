package channel

import (
	"fmt"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

type baseChannel struct {
	code              string
	paymentMethodType string
}

func newBaseChannel(cfg config.YooKassaChannelConfig) baseChannel {
	return baseChannel{
		code:              cfg.Code,
		paymentMethodType: cfg.PaymentMethodType,
	}
}

func (c baseChannel) CreatePaymentRequest(d model.CreatePaymentData) data.CreatePaymentRequest {
	return data.CreatePaymentRequest{
		Confirmation: data.PaymentConfirmation{
			Type:      data.PaymentConfirmationTypeRedirect,
			ReturnURL: d.ReturnURLs.Common,
			Locale:    getLocale(d.LangCode),
		},
		PaymentMethod: data.PaymentMethod{
			Type: c.paymentMethodType,
		},
		Amount: data.PaymentAmount{
			Currency: d.Currency,
			Amount:   convert.CentsToBase(d.Amount),
		},
		Description: getDescription(d.LangCode, d.OperationID),
		Capture:     true,
	}
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
