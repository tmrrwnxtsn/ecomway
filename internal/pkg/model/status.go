package model

import "time"

const (
	OperationFailReasonTimeout                = "Timeout"
	OperationFailReasonManual                 = "Manual"
	OperationFailReasonWrongCodeLimitExceeded = "Code confirm attempts has been exceeded"
)

type GetOperationStatusData struct {
	CreatedAt      time.Time
	ExternalID     string
	ExternalSystem string
	ExternalMethod string
	Currency       string
	OperationType  OperationType
	OperationID    int64
	UserID         int64
	Amount         int64
}

type GetOperationStatusResult struct {
	ExternalID     string
	ExternalStatus OperationExternalStatus
	ProcessedAt    time.Time
	FailReason     string
	NewAmount      int64
	Tool           *Tool
}
