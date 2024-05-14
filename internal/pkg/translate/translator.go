package translate

import "github.com/salihzain/nakedi18n"

const defaultLang = "ru"

const (
	KeyInvalidAPIKey     = "INVALID_API_KEY"
	KeyPaymentSuccessful = "PAYMENT_SUCCESSFUL"
	KeyPaymentRejected   = "PAYMENT_REJECTED"
	KeyObjectNotFound    = "OBJECT_NOT_FOUND"
	KeyUnexpectedError   = "UNEXPECTED_ERROR"
	KeyToolRemoved       = "TOOL_REMOVED"
)

type Translator struct {
	translateFunc func(string, string, ...any) string
}

func NewTranslator(lang ...string) *Translator {
	translations := map[string]map[string]string{
		"en": {
			KeyInvalidAPIKey:     "Authorization header must contain valid API key.",
			KeyPaymentSuccessful: "Payment successful!",
			KeyPaymentRejected:   "Payment rejected.",
			KeyObjectNotFound:    "Object not found.",
			KeyUnexpectedError:   "Unexpected error occurred. Try again later.",
			KeyToolRemoved:       "Payment tool has been removed.",
		},
		"ru": {
			KeyInvalidAPIKey:     "Заголовок авторизации должен содержать корректный API ключ.",
			KeyPaymentSuccessful: "Оплата успешна!",
			KeyPaymentRejected:   "Оплата неуспешна.",
			KeyObjectNotFound:    "Искомый объект не найден.",
			KeyUnexpectedError:   "Произошла непредвиденная ошибка. Повторите попытку позже.",
			KeyToolRemoved:       "Платежное средство успешно удалено.",
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
