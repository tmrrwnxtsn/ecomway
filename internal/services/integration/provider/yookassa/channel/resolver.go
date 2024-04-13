package channel

import (
	"fmt"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

type Channel interface {
	CreatePaymentRequest(data model.CreatePaymentData) data.CreatePaymentRequest
	PaymentTimeoutToFailed() time.Duration
}

type Resolver struct {
	channels map[string]Channel
}

func NewResolver(channelsConfigs map[string]config.YooKassaChannelConfig) Resolver {
	return Resolver{
		channels: channels(channelsConfigs),
	}
}

func channels(channelsConfigs map[string]config.YooKassaChannelConfig) map[string]Channel {
	var result map[string]Channel
	if channelsNum := len(channelsConfigs); channelsNum > 0 {
		result = make(map[string]Channel, channelsNum)
		for externalMethod, channelCfg := range channelsConfigs {
			result[externalMethod] = newBaseChannel(channelCfg)
		}
	}
	return result
}

func (r Resolver) Channel(externalMethod string) (Channel, error) {
	channel, found := r.channels[externalMethod]
	if !found {
		return nil, fmt.Errorf("channel not found for external method %q", externalMethod)
	}
	return channel, nil
}
