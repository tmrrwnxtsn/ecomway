package model

import "time"

type ToolType string

const (
	ToolTypeBankCard ToolType = "BANK_CARD"
	ToolTypeWallet   ToolType = "WALLET"
)

type ToolStatus string

const (
	ToolStatusActive ToolStatus = "ACTIVE"
)

type Tool struct {
	ID             string
	UserID         int64
	ExternalMethod string
	Displayed      string
	Name           string
	Status         ToolStatus
	Type           ToolType
	Details        map[string]any
	Fake           bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
