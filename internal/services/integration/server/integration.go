package server

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Integration interface {
	AvailableMethods(ctx context.Context, userID int64, txType model.TransactionType) ([]model.Method, error)
}
