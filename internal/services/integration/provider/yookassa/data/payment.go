package data

import "time"

type Amount struct {
	Currency string
	Value    float64
}

type PaymentMethodCard struct {
	First6      string
	Last4       string
	ExpiryYear  string
	ExpiryMonth string
	CardType    string
	IssuerName  string
}

type PaymentMethod struct {
	Type          string
	ID            string
	AccountNumber string
	Card          PaymentMethodCard
	Saved         bool
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
	Confirmation      PaymentConfirmation
	PaymentMethodData PaymentMethod
	Amount            Amount
	Description       string
	PaymentMethodID   string
	Capture           bool
	SavePaymentMethod bool
}

type CreatePaymentResponse struct {
	ID              string
	ConfirmationURL string
	Status          string
	CapturedAt      time.Time
	Cancellation    Cancellation
	IncomeAmount    Amount
	PaymentMethod   PaymentMethod
}

type Cancellation struct {
	Party  string
	Reason string
}

type GetPaymentResponse struct {
	CapturedAt    time.Time
	ID            string
	Status        string
	Cancellation  Cancellation
	IncomeAmount  Amount
	PaymentMethod PaymentMethod
}
