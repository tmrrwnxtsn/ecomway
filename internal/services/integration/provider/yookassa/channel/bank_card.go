package channel

import (
	"fmt"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const channelCodeBankCard = "card"

type bankCardChannel struct {
	baseChannel
}

func newBankCardChannel(cfg config.YooKassaChannelConfig) bankCardChannel {
	return bankCardChannel{
		baseChannel: newBaseChannel(cfg),
	}
}

func (c bankCardChannel) PaymentTool(userID int64, externalMethod string, method data.PaymentMethod) *model.Tool {
	if method.Type != c.paymentMethodType || !method.Saved {
		return nil
	}

	displayed := fmt.Sprintf("%v******%v", method.Card.First6, method.Card.Last4)

	return &model.Tool{
		ID:             method.ID,
		UserID:         userID,
		ExternalMethod: externalMethod,
		Displayed:      displayed,
		Name:           "Bank card", // TODO: можно сохранять локаль в additional и сюда писать название по локали
		Type:           model.ToolTypeBankCard,
		Details: map[string]any{
			"token":        method.ID,
			"first6":       method.Card.First6,
			"last4":        method.Card.Last4,
			"expiry_year":  method.Card.ExpiryYear,
			"expiry_month": method.Card.ExpiryMonth,
			"card_type":    method.Card.CardType,
			"bank_name":    method.Card.IssuerName,
		},
	}
}
