package channel

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const defaultPaymentTimeoutToFailed = 20 * time.Minute

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
	var paymentMethodID string
	if d.Tool != nil {
		token, ok := d.Tool.Details["token"].(string)
		if ok {
			paymentMethodID = token
		} else {
			slog.Warn(
				"token not found for given tool",
				"tool_id", d.Tool.ID,
				"user_id", d.UserID,
				"external_method", d.ExternalMethod,
				"operation_id", d.OperationID,
			)
		}
	}

	return data.CreatePaymentRequest{
		Confirmation: data.PaymentConfirmation{
			Type:      data.PaymentConfirmationTypeRedirect,
			ReturnURL: d.ReturnURLs.Common,
			Locale:    locale(d.LangCode),
		},
		PaymentMethodData: data.PaymentMethod{
			Type: c.paymentMethodType,
		},
		Amount: data.Amount{
			Currency: d.Currency,
			Value:    convert.CentsToBase(d.Amount),
		},
		Description:       description(model.OperationTypePayment, d.LangCode, d.OperationID),
		PaymentMethodID:   paymentMethodID,
		Capture:           true,
		SavePaymentMethod: true,
	}
}

func (c baseChannel) CreatePayoutRequest(d model.CreatePayoutData) data.CreatePayoutRequest {
	var paymentMethodID string
	if d.Tool != nil {
		token, ok := d.Tool.Details["token"].(string)
		if ok {
			paymentMethodID = token
		} else {
			slog.Warn(
				"token not found for given tool",
				"tool_id", d.Tool.ID,
				"user_id", d.UserID,
				"external_method", d.ExternalMethod,
				"operation_id", d.OperationID,
			)
		}
	}

	return data.CreatePayoutRequest{
		Amount: data.Amount{
			Currency: d.Currency,
			Value:    convert.CentsToBase(d.Amount),
		},
		Description:     description(model.OperationTypePayout, d.LangCode, d.OperationID),
		PaymentMethodID: paymentMethodID,
	}
}

func (c baseChannel) PaymentTool(_ string, _ string, _ data.PaymentMethod) *model.Tool {
	return &model.Tool{}
}

func (c baseChannel) PaymentTimeoutToFailed() time.Duration {
	return c.paymentTimeoutToFailed
}

func description(opType model.OperationType, langCode string, operationID int64) string {
	var descriptionsFmt map[string]string
	switch opType {
	case model.OperationTypePayment:
		descriptionsFmt = map[string]string{
			"en": "Replenishment of balance №%v",
			"ru": "Пополнение баланса №%v",
		}
	case model.OperationTypePayout:
		descriptionsFmt = map[string]string{
			"en": "Withdrawal of funds №%v",
			"ru": "Вывод средств №%v",
		}
	}
	descriptionFmt, ok := descriptionsFmt[langCode]
	if !ok {
		descriptionFmt = descriptionsFmt["ru"]
	}
	return fmt.Sprintf(descriptionFmt, operationID)
}

func locale(langCode string) string {
	locales := map[string]string{
		"en": "en_US",
		"ru": "ru_RU",
	}
	loc, ok := locales[langCode]
	if !ok {
		loc = locales["ru"]
	}
	return loc
}
