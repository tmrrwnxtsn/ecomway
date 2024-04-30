package model

import (
	"context"
	"time"
)

type OperationType string

const (
	OperationTypePayment OperationType = "payment"
	OperationTypePayout  OperationType = "payout"
)

type OperationStatus string

const (
	OperationStatusNew     OperationStatus = "NEW"
	OperationStatusFailed  OperationStatus = "FAILED"
	OperationStatusSuccess OperationStatus = "SUCCESS"
)

type OperationExternalStatus string

const (
	OperationExternalStatusUnknown OperationExternalStatus = "UNKNOWN"
	OperationExternalStatusPending OperationExternalStatus = "PENDING"
	OperationExternalStatusSuccess OperationExternalStatus = "SUCCESS"
	OperationExternalStatusFailed  OperationExternalStatus = "FAILED"
)

type Operation struct {
	ID             int64
	UserID         int64
	Type           OperationType
	Currency       string
	Amount         int64
	Status         OperationStatus
	ExternalID     string
	ExternalSystem string
	ExternalMethod string
	ExternalStatus OperationExternalStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time

	ToolID           string
	Additional       map[string]any
	FailReason       string
	ConfirmationCode string
	ProcessedAt      time.Time
}

type ScriptAcquiredFor func(ctx context.Context, op *Operation) error

type OperationCriteria struct {
	ID         *int64
	UserID     *int64
	ExternalID *string

	Types           *[]OperationType
	Statuses        *[]OperationStatus
	StatusesByType  map[OperationType][]OperationStatus
	ExternalSystems *[]string
	CreatedAtFrom   time.Time
	CreatedAtTo     time.Time
	MaxCount        int64
}
