package limit

import (
	"fmt"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ValidateAmount(amount int64, currency string, method *model.Method) error {
	limits, found := method.Limits[currency]
	if !found {
		return fmt.Errorf("currency %q not supported", currency)
	}

	if amount < limits.MinAmount {
		return fmt.Errorf("amount is less than minimal: %v", limits.MinAmount)
	}

	if amount > limits.MaxAmount {
		return fmt.Errorf("amount is more than maximum: %v", limits.MaxAmount)
	}

	return nil
}
