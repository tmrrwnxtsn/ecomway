package model

type CreatePaymentData struct {
	AdditionalData map[string]any
	ExternalSystem string
	ExternalMethod string
	Currency       string
	LangCode       string
	UserID         int64
	Amount         int64
}

type CreatePaymentResult struct {
	RedirectURL string
}
