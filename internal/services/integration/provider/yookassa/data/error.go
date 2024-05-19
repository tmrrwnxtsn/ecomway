package data

import (
	"errors"
	"fmt"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrPayoutNotFound  = errors.New("payout not found")
)

type ErrorResponse struct {
	ID          string
	Code        string
	Description string
}

func (r ErrorResponse) Error() string {
	errorMessage := r.Code
	if r.Description != "" {
		errorMessage = fmt.Sprintf("%s: %s", errorMessage, r.Description)
	}
	return fmt.Sprintf("%s (%s)", errorMessage, r.ID)
}
