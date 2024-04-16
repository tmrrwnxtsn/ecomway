package channel

import (
	"fmt"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const (
	channelCodeBankCard = "card"
)

type bankCardChannel struct {
	baseChannel
}

func newBankCardChannel(cfg config.YooKassaChannelConfig) bankCardChannel {
	return bankCardChannel{
		baseChannel: newBaseChannel(cfg),
	}
}

func (c bankCardChannel) CreatePaymentRequest(d model.CreatePaymentData) data.CreatePaymentRequest {
	request := c.baseChannel.CreatePaymentRequest(d)

	if d.ToolID != 0 && d.Tool != nil {
		token, ok := d.Tool.Details["token"].(string)
		if ok {
			request.PaymentMethodID = token
		}
	}

	return request
}

func (c bankCardChannel) PaymentTool(d model.GetOperationStatusData, resp data.GetPaymentResponse) *model.Tool {
	if resp.PaymentMethod.Type != c.paymentMethodType || !resp.PaymentMethod.Saved {
		return nil
	}

	displayed := fmt.Sprintf("%v******%v", resp.PaymentMethod.Card.First6, resp.PaymentMethod.Card.Last4)

	return &model.Tool{
		UserID:         d.UserID,
		ExternalMethod: d.ExternalMethod,
		Displayed:      displayed,
		Type:           model.ToolTypeBankCard,
		Details: map[string]any{
			"token":        resp.PaymentMethod.ID,
			"first6":       resp.PaymentMethod.Card.First6,
			"last4":        resp.PaymentMethod.Card.Last4,
			"expiry_year":  resp.PaymentMethod.Card.ExpiryYear,
			"expiry_month": resp.PaymentMethod.Card.ExpiryMonth,
			"card_type":    resp.PaymentMethod.Card.CardType,
			"bank_name":    resp.PaymentMethod.Card.IssuerName,
		},
	}
}
