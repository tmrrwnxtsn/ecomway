package model

import "time"

type ReturnURLs struct {
	Common  string
	Success string
	Fail    string
}

type CreatePaymentData struct {
	Tool           *Tool
	ReturnURLs     ReturnURLs
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

type CreatePaymentResult struct {
	ProcessedAt    time.Time
	RedirectURL    string
	ExternalID     string
	ExternalStatus OperationExternalStatus
	FailReason     string
	Status         OperationStatus
	OperationID    int64
	NewAmount      int64
	Tool           *Tool
}

type SuccessPaymentData struct {
	ProcessedAt    time.Time
	ExternalID     string
	ExternalStatus OperationExternalStatus
	OperationID    int64
	NewAmount      int64
	Tool           *Tool
}

type FailPaymentData struct {
	ExternalID     string
	ExternalStatus OperationExternalStatus
	FailReason     string
	OperationID    int64
}
