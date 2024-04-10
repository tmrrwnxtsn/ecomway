package yookassa

import (
	"context"
	"fmt"
	"slices"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (i *Integration) AvailableMethods(_ context.Context, opType model.OperationType, _ string) ([]model.Method, error) {
	switch opType {
	case model.OperationTypePayment:
		return slices.Clone(i.paymentMethods), nil
	case model.OperationTypePayout:
		return slices.Clone(i.payoutMethods), nil
	default:
		return nil, fmt.Errorf("unresolved operation type: %q", opType)
	}
}
