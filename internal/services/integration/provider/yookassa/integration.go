package yookassa

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/api"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/channel"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const ExternalSystem = "yookassa"

type APIClient interface {
	CreatePayment(ctx context.Context, request data.CreatePaymentRequest) (data.CreatePaymentResponse, error)
	GetPayment(ctx context.Context, paymentID string) (data.GetPaymentResponse, error)
}

type ChannelResolver interface {
	Channel(externalMethod string) (channel.Channel, error)
}

type Integration struct {
	apiClient       APIClient
	channelResolver ChannelResolver
	paymentMethods  []model.Method
	payoutMethods   []model.Method
}

func NewIntegration(cfg *config.YooKassaConfig) *Integration {
	if cfg == nil {
		return nil
	}

	apiClient := api.NewClient(api.ClientOptions{
		BaseURL:   cfg.API.BaseURL,
		ShopID:    cfg.API.ShopID,
		SecretKey: cfg.API.SecretKey,
	})

	channelResolver := channel.NewResolver(cfg.Channels)

	return &Integration{
		apiClient:       apiClient,
		channelResolver: channelResolver,
		paymentMethods:  convert.MethodsFromConfig(cfg.Methods.Payment),
		payoutMethods:   convert.MethodsFromConfig(cfg.Methods.Payout),
	}
}
