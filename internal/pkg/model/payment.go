package model

import "time"

type ReturnURLs struct {
	Common  string
	Success string
	Fail    string
}

type CreatePaymentData struct {
	ReturnURLs     ReturnURLs
	AdditionalData map[string]any
	ExternalSystem string
	ExternalMethod string
	Currency       string
	LangCode       string
	UserID         int64
	Amount         int64
	OperationID    int64
}

type CreatePaymentResult struct {
	RedirectURL    string
	ExternalID     string
	ExternalStatus OperationExternalStatus
	OperationID    int64
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
