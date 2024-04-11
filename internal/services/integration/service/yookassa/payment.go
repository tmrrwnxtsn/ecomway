package yookassa

import (
	"context"
	"errors"
	"math/rand"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (i *Integration) CreatePayment(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error) {
	r := rand.Int()
	if r%2 == 0 {
		return model.CreatePaymentResult{
			RedirectURL: "google.com",
		}, nil
	} else {
		return model.CreatePaymentResult{}, errors.New("some ps error")
	}
}
