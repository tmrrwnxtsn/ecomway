package v1

import (
	"strconv"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

const (
	toolTypeBankCard = "card"
	toolTypeWallet   = "wallet"
)

type toolDetails struct {
	// Тип банковской карты
	CardType string `json:"card_type,omitempty" example:"Visa"`
	// Владелец банковской карты
	CardHolder string `json:"card_holder,omitempty" example:"Ivanov Ivan"`
	// Срок действия банковской карты (месяц, MM)
	ExpiryMonth int64 `json:"expiry_month,omitempty" example:"10"`
	// Срок действия банковской карты (год, YYYY)
	ExpiryYear int64 `json:"expiry_year,omitempty" example:"2023"`
	// Название банка, выпустившего банковскую карту
	BankName string `json:"bank_name,omitempty" example:"Sberbank"`
	// Номер электронного кошелька
	WalletNumber string `json:"wallet_number,omitempty" example:"410011758831136"`
}

type tool struct {
	// Идентификатор платежного инструмента
	ID int64 `json:"id" example:"14124" validate:"required"`
	// Тип платежного инструмента:
	// * Банковская карта - "card"
	// * Электронный кошелек - "wallet"
	Type string `json:"type" example:"card" validate:"required"`
	// Значение платежного инструмента, например:
	// * Маскированная банковская карта
	// * Номер электронного кошелька
	// * Адрес электронной почты
	// * и т.д.
	Caption string `json:"caption" example:"444444******4444" validate:"required"`
	// Дополнительная информация о платежном инструменте
	Details *toolDetails `json:"details,omitempty"`
}

func (h *Handler) tool(item *model.Tool) tool {
	t := tool{
		ID:      item.ID,
		Caption: item.Displayed,
	}

	switch item.Type {
	case model.ToolTypeBankCard:
		t.Type = toolTypeBankCard

		cardType, _ := item.Details["card_type"].(string)
		cardHolder, _ := item.Details["card_holder"].(string)
		bankName, _ := item.Details["bank_name"].(string)
		expiryMonthStr, _ := item.Details["expiry_month"].(string)
		expiryYearStr, _ := item.Details["expiry_year"].(string)

		if cardType+cardHolder+bankName+expiryMonthStr+expiryYearStr == "" {
			break
		}

		expiryMonth, _ := strconv.ParseInt(expiryMonthStr, 10, 64)
		expiryYear, _ := strconv.ParseInt(expiryYearStr, 10, 64)

		t.Details = &toolDetails{
			CardType:    cardType,
			CardHolder:  cardHolder,
			ExpiryMonth: expiryMonth,
			ExpiryYear:  expiryYear,
			BankName:    bankName,
		}
	case model.ToolTypeWallet:
		t.Type = toolTypeWallet

		walletNumber, _ := item.Details["number"].(string)
		if walletNumber == "" {
			break
		}

		t.Details = &toolDetails{
			WalletNumber: walletNumber,
		}
	}

	return t
}

func (h *Handler) tools(items []*model.Tool) []tool {
	var tools []tool
	if itemNum := len(items); itemNum > 0 {
		tools = make([]tool, 0, itemNum)
		for _, item := range items {
			tools = append(tools, h.tool(item))
		}
	}
	return tools
}
