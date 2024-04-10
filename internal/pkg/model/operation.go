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
	OperationStatusNew    OperationStatus = "NEW"
	OperationStatusFailed OperationStatus = "FAILED"
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
	ExternalStatus string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	ToolID           int64
	Additional       map[string]any
	FailReason       string
	ConfirmationCode string
	ProcessedAt      time.Time
}

type ScriptAcquiredFor func(ctx context.Context, op *Operation) error

type OperationCriteria struct {
	ID            *int64
	UserID        *int64
	Types         *[]string
	Statuses      *[]string
	ExternalID    *string
	CreatedAtFrom time.Time
	CreatedAtTo   time.Time
}