package translate

import "github.com/salihzain/nakedi18n"

const defaultLang = "ru"

const (
	KeyInvalidAPIKey            = "INVALID_API_KEY"
	KeyPaymentSuccessful        = "PAYMENT_SUCCESSFUL"
	KeyPaymentRejected          = "PAYMENT_REJECTED"
	KeyObjectNotFound           = "OBJECT_NOT_FOUND"
	KeyUnexpectedError          = "UNEXPECTED_ERROR"
	KeyToolRemoved              = "TOOL_REMOVED"
	KeyExternalStatusUnknown    = "EXTERNAL_STATUS_UNKNOWN"
	KeyExternalStatusPending    = "EXTERNAL_STATUS_PENDING"
	KeyExternalStatusSuccess    = "EXTERNAL_STATUS_SUCCESS"
	KeyExternalStatusFailed     = "EXTERNAL_STATUS_FAILED"
	KeyToolRecovered            = "TOOL_RECOVERED"
	KeyUnresolvedStatusConflict = "UNRESOLVED_STATUS_CONFLICT"
	KeyForbiddenOnRemovedTool   = "FORBIDDEN_ON_REMOVED_TOOL"
	KeyWrongConfirmationCode    = "WRONG_CONFIRMATION_CODE"
	KeyWrongCodeLimitExceeded   = "WRONG_CODE_LIMIT_EXCEEDED"
	KeyPayoutConfirmed          = "PAYOUT_CONFIRMED"
)

type Translator struct {
	translateFunc func(string, string, ...any) string
}

func NewTranslator(lang ...string) *Translator {
	translations := map[string]map[string]string{
		"en": {
			KeyInvalidAPIKey:            "Authorization header must contain valid API key.",
			KeyPaymentSuccessful:        "Payment successful!",
			KeyPaymentRejected:          "Payment rejected.",
			KeyObjectNotFound:           "Object not found.",
			KeyUnexpectedError:          "Unexpected error occurred. Try again later.",
			KeyToolRemoved:              "Payment tool has been removed.",
			KeyExternalStatusUnknown:    "Could not retrieve operation status on payment system side.",
			KeyExternalStatusPending:    "Operation is still being processed by the payment system.",
			KeyExternalStatusSuccess:    "Operation has been processed successfully by the payment system.",
			KeyExternalStatusFailed:     "Operation has been failed by the payment system.",
			KeyToolRecovered:            "Payment tool is ready for recovery.",
			KeyUnresolvedStatusConflict: "Not able to perform the action for the object status.",
			KeyForbiddenOnRemovedTool:   "Not able to perform action on restricted payment tool.",
			KeyWrongConfirmationCode:    "Wrong confirmation code.",
			KeyWrongCodeLimitExceeded:   "The maximum number of possible confirmation attempts has been exceeded. The payout was rejected.",
			KeyPayoutConfirmed:          "Payout successfully confirmed.",
		},
		"ru": {
			KeyInvalidAPIKey:            "Заголовок авторизации должен содержать корректный API ключ.",
			KeyPaymentSuccessful:        "Оплата успешна!",
			KeyPaymentRejected:          "Оплата неуспешна.",
			KeyObjectNotFound:           "Искомый объект не найден.",
			KeyUnexpectedError:          "Произошла непредвиденная ошибка. Повторите попытку позже.",
			KeyToolRemoved:              "Платежное средство успешно удалено.",
			KeyExternalStatusUnknown:    "Невозможно автоматически проверить статус операции на стороне платежной системы.",
			KeyExternalStatusPending:    "Операция на стороне платежной системы еще находится в обработке.",
			KeyExternalStatusSuccess:    "Операция на стороне платеной системы имеет успешный статус.",
			KeyExternalStatusFailed:     "Операция на стороне платеной системы отклонена.",
			KeyToolRecovered:            "Платежное средство готово для восстановления.",
			KeyUnresolvedStatusConflict: "Целевое действие невозможно для данного статуса объекта.",
			KeyForbiddenOnRemovedTool:   "Невозможно осуществить операцию с использованием удаленного платежного средства.",
			KeyWrongConfirmationCode:    "Неверный код подтверждения.",
			KeyWrongCodeLimitExceeded:   "Превышено максимальное количество возможных попыток для подтверждения. Выплата отклонена.",
			KeyPayoutConfirmed:          "Вывод средств успешно подтвержден.",
		},
	}

	i18nInstance := nakedi18n.NewNakedI18n(defaultLang, lang, true, translations)
	translateFunc := i18nInstance.UseNakedI18n(nil, true)

	return &Translator{
		translateFunc: translateFunc,
	}
}

func (t *Translator) Translate(lang, key string, _ ...any) string {
	return t.translateFunc(lang, key)
}
