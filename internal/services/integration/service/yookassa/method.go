package yookassa

import (
	"context"
	"fmt"
	"slices"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (i *Integration) AvailableMethods(_ context.Context, _ int64, txType model.TransactionType) ([]model.Method, error) {
	switch txType {
	case model.TransactionTypePayment:
		return slices.Clone(i.paymentMethods), nil
	case model.TransactionTypePayout:
		return slices.Clone(i.payoutMethods), nil
	default:
		return nil, fmt.Errorf("unresolved transaction type: %q", txType)
	}
}
