package model

import "time"

type CreatePayoutData struct {
	Tool           *Tool
	AdditionalData map[string]any
	ExternalSystem string
	ExternalMethod string
	Currency       string
	LangCode       string
	ToolID         string
	UserID         string
	Amount         int64
	OperationID    int64
}

type CreatePayoutResult struct {
	ProcessedAt    time.Time
	ExternalID     string
	ExternalStatus OperationExternalStatus
	FailReason     string
	Status         OperationStatus
	OperationID    int64
}

type ConfirmPayoutData struct {
	ConfirmationCode string
	LangCode         string
	UserID           string
	OperationID      int64
}

type SuccessPayoutData struct {
	ProcessedAt    time.Time
	ExternalID     string
	ExternalStatus OperationExternalStatus
	OperationID    int64
	Tool           *Tool
}

type FailPayoutData struct {
	ExternalID     string
	ExternalStatus OperationExternalStatus
	FailReason     string
	OperationID    int64
}
