package yookassa

import (
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/convert"
)

const ExternalSystem = "yookassa"

type Integration struct {
	paymentMethods []model.Method
	payoutMethods  []model.Method
}

func NewIntegration(cfg *config.YooKassaConfig) *Integration {
	if cfg == nil {
		return nil
	}
	return &Integration{
		paymentMethods: convert.MethodsFromConfig(cfg.Methods.Payment),
		payoutMethods:  convert.MethodsFromConfig(cfg.Methods.Payout),
	}
}
