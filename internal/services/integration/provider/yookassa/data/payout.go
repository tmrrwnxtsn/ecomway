package data

import "time"

type CreatePayoutRequest struct {
	Amount          Amount
	Description     string
	PaymentMethodID string
}

type CreatePayoutResponse struct {
	ID           string
	Status       string
	CapturedAt   time.Time
	Cancellation Cancellation
}

type GetPayoutResponse struct {
	CapturedAt   time.Time
	ID           string
	Status       string
	Cancellation Cancellation
}
