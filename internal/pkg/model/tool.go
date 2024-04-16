package model

import "time"

type ToolType string

const (
	ToolTypeBankCard = "BANK_CARD"
	ToolTypeWallet   = "WALLET"
)

type Tool struct {
	ID             int64
	UserID         int64
	ExternalMethod string
	Displayed      string
	Type           ToolType
	Details        map[string]any
	Fake           bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
