package channel

import "github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"

const (
	channelCodeCard = "card"
)

type cardChannel struct {
	baseChannel
}

func newCardChannel(cfg config.YooKassaChannelConfig) cardChannel {
	return cardChannel{
		baseChannel: newBaseChannel(cfg),
	}
}
