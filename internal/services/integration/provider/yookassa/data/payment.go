package data

import "time"

const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusCanceled  = "canceled"
)

type PaymentAmount struct {
	Currency string
	Value    float64
}

type PaymentMethod struct {
	Type string
}

const (
	PaymentConfirmationTypeRedirect = "redirect"
)

type PaymentConfirmation struct {
	Type      string
	ReturnURL string
	Locale    string
}

type CreatePaymentRequest struct {
	Confirmation  PaymentConfirmation
	PaymentMethod PaymentMethod
	Amount        PaymentAmount
	Description   string
	Capture       bool
}

type CreatePaymentResponse struct {
	ID              string
	ConfirmationURL string
	Status          string
}

type PaymentCancellation struct {
	Party  string
	Reason string
}

type GetPaymentResponse struct {
	CapturedAt   time.Time
	ID           string
	Status       string
	Cancellation PaymentCancellation
	IncomeAmount PaymentAmount
}
