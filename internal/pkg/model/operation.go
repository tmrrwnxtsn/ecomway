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
	OperationStatusNew       OperationStatus = "NEW"
	OperationStatusConfirmed OperationStatus = "CONFIRMED"
	OperationStatusFailed    OperationStatus = "FAILED"
	OperationStatusSuccess   OperationStatus = "SUCCESS"
	OperationStatusPending   OperationStatus = "PENDING"
)

type OperationExternalStatus string

const (
	OperationExternalStatusUnknown OperationExternalStatus = "UNKNOWN"
	OperationExternalStatusPending OperationExternalStatus = "PENDING"
	OperationExternalStatusSuccess OperationExternalStatus = "SUCCESS"
	OperationExternalStatusFailed  OperationExternalStatus = "FAILED"
)

type OperationChangeStatusResult string

const (
	OperationChangeStatusResultSuccessPayment = "SUCCESS_PAYMENT"
	OperationChangeStatusResultFailPayment    = "FAIL_PAYMENT"
	OperationChangeStatusResultSuccessPayout  = "SUCCESS_PAYOUT"
	OperationChangeStatusResultFailPayout     = "FAIL_PAYOUT"
)

type Operation struct {
	ID             int64
	UserID         string
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

	ToolID               string
	Additional           map[string]any
	FailReason           string
	ConfirmationCode     string
	ProcessedAt          time.Time
	ConfirmationAttempts int
}

type ReportOperation struct {
	ID             int64
	UserID         string
	Type           OperationType
	Currency       string
	Amount         int64
	Status         OperationStatus
	ExternalID     string
	ExternalSystem string
	ExternalMethod string
	ExternalStatus OperationExternalStatus
	ToolDisplayed  string
	FailReason     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ProcessedAt    time.Time
}

type ScriptAcquiredFor func(ctx context.Context, op *Operation) error

type OperationCriteria struct {
	ID         *int64
	UserID     *string
	ExternalID *string

	Types           *[]OperationType
	Statuses        *[]OperationStatus
	StatusesByType  map[OperationType][]OperationStatus
	ExternalSystems *[]string
	CreatedAtFrom   time.Time
	CreatedAtTo     time.Time
	MaxCount        int64
}
