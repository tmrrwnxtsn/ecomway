package channel

import (
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const channelCodeWallet = "wallet"

type walletChannel struct {
	baseChannel
}

func newWalletChannel(cfg config.YooKassaChannelConfig) walletChannel {
	return walletChannel{
		baseChannel: newBaseChannel(cfg),
	}
}

func (c walletChannel) CreatePaymentRequest(d model.CreatePaymentData) data.CreatePaymentRequest {
	request := c.baseChannel.CreatePaymentRequest(d)

	if request.PaymentMethodID != "" {
		request.Confirmation = data.PaymentConfirmation{}
	}

	return request
}

func (c walletChannel) PaymentTool(userID int64, externalMethod string, method data.PaymentMethod) *model.Tool {
	if method.Type != c.paymentMethodType || !method.Saved {
		return nil
	}

	return &model.Tool{
		ID:             method.AccountNumber,
		UserID:         userID,
		ExternalMethod: externalMethod,
		Displayed:      method.AccountNumber,
		Name:           "Wallet number", // TODO: можно сохранять локаль в additional и сюда писать название по локали
		Type:           model.ToolTypeWallet,
		Details: map[string]any{
			"token":  method.ID,
			"number": method.AccountNumber,
		},
	}
}
