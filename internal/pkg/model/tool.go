package model

import "time"

type ToolType string

const (
	ToolTypeBankCard ToolType = "BANK_CARD"
	ToolTypeWallet   ToolType = "WALLET"
)

type ToolStatus string

const (
	ToolStatusActive                 ToolStatus = "ACTIVE"
	ToolStatusRemovedByClient        ToolStatus = "REMOVED_BY_CLIENT"
	ToolStatusPendingRecovery        ToolStatus = "PENDING_RECOVERY"
	ToolStatusRemovedByAdministrator ToolStatus = "REMOVED_BY_ADMINISTRATOR"
)

type Tool struct {
	ID             string
	UserID         string
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

func (t Tool) CanBeRecovered() bool {
	return t.Status != ToolStatusRemovedByAdministrator
}

func (t Tool) Removed() bool {
	return t.Status == ToolStatusRemovedByClient ||
		t.Status == ToolStatusRemovedByAdministrator ||
		t.Status == ToolStatusPendingRecovery
}
