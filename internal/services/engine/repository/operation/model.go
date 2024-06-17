package operation

import (
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type dbOperation struct {
	ID             int64     `db:"id"`
	UserID         string    `db:"user_id"`
	Type           string    `db:"type"`
	Currency       string    `db:"currency"`
	Amount         float64   `db:"amount"`
	Status         string    `db:"status"`
	ExternalID     *string   `db:"external_id"`
	ExternalSystem string    `db:"external_system"`
	ExternalMethod string    `db:"external_method"`
	ExternalStatus *string   `db:"external_status"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`

	ToolID               *string        `db:"tool_id"`
	Additional           map[string]any `db:"additional"`
	FailReason           *string        `db:"fail_reason"`
	ConfirmationCode     *string        `db:"confirmation_code"`
	ProcessedAt          *time.Time     `db:"processed_at"`
	ConfirmationAttempts *int           `db:"confirmation_attempts"`
}

func operationToDB(op *model.Operation) dbOperation {
	dbOp := dbOperation{
		ID:             op.ID,
		UserID:         op.UserID,
		Type:           string(op.Type),
		Currency:       op.Currency,
		Amount:         convert.CentsToBase(op.Amount),
		Status:         string(op.Status),
		ExternalSystem: op.ExternalSystem,
		ExternalMethod: op.ExternalMethod,
		CreatedAt:      op.CreatedAt,
		UpdatedAt:      op.UpdatedAt,
	}

	if op.ExternalID != "" {
		dbOp.ExternalID = &op.ExternalID
	}

	if op.ExternalStatus != "" {
		dbOp.ExternalStatus = (*string)(&op.ExternalStatus)
	}

	if len(op.Additional) > 0 {
		dbOp.Additional = op.Additional
	}

	if op.ToolID != "" {
		dbOp.ToolID = &op.ToolID
	}

	if op.FailReason != "" {
		dbOp.FailReason = &op.FailReason
	}

	if op.ConfirmationCode != "" {
		dbOp.ConfirmationCode = &op.ConfirmationCode
	}

	if !op.ProcessedAt.IsZero() {
		dbOp.ProcessedAt = &op.ProcessedAt
	}

	if op.ConfirmationAttempts != 0 {
		dbOp.ConfirmationAttempts = &op.ConfirmationAttempts
	}

	return dbOp
}

func operationFromDB(dbOp dbOperation) *model.Operation {
	op := &model.Operation{
		ID:             dbOp.ID,
		UserID:         dbOp.UserID,
		Type:           model.OperationType(dbOp.Type),
		Currency:       dbOp.Currency,
		Amount:         convert.BaseToCents(dbOp.Amount),
		Status:         model.OperationStatus(dbOp.Status),
		ExternalSystem: dbOp.ExternalSystem,
		ExternalMethod: dbOp.ExternalMethod,
		CreatedAt:      dbOp.CreatedAt,
		UpdatedAt:      dbOp.UpdatedAt,
	}

	if dbOp.ExternalID != nil {
		op.ExternalID = *dbOp.ExternalID
	}

	if dbOp.ExternalStatus != nil {
		op.ExternalStatus = model.OperationExternalStatus(*dbOp.ExternalStatus)
	}

	if len(dbOp.Additional) > 0 {
		op.Additional = dbOp.Additional
	}

	if dbOp.ToolID != nil {
		op.ToolID = *dbOp.ToolID
	}

	if dbOp.FailReason != nil {
		op.FailReason = *dbOp.FailReason
	}

	if dbOp.ConfirmationCode != nil {
		op.ConfirmationCode = *dbOp.ConfirmationCode
	}

	if dbOp.ProcessedAt != nil {
		op.ProcessedAt = *dbOp.ProcessedAt
	}

	if dbOp.ConfirmationAttempts != nil {
		op.ConfirmationAttempts = *dbOp.ConfirmationAttempts
	}

	return op
}

type dbReportOperation struct {
	ID             int64      `db:"id"`
	UserID         string     `db:"user_id"`
	Type           string     `db:"type"`
	Currency       string     `db:"currency"`
	Amount         float64    `db:"amount"`
	Status         string     `db:"status"`
	ExternalID     *string    `db:"external_id"`
	ExternalSystem string     `db:"external_system"`
	ExternalMethod string     `db:"external_method"`
	ExternalStatus *string    `db:"external_status"`
	ToolDisplayed  *string    `db:"tool_displayed"`
	FailReason     *string    `db:"fail_reason"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	ProcessedAt    *time.Time `db:"processed_at"`
}

func reportOperationFromDB(dbOp dbReportOperation) model.ReportOperation {
	op := model.ReportOperation{
		ID:             dbOp.ID,
		UserID:         dbOp.UserID,
		Type:           model.OperationType(dbOp.Type),
		Currency:       dbOp.Currency,
		Amount:         convert.BaseToCents(dbOp.Amount),
		Status:         model.OperationStatus(dbOp.Status),
		ExternalSystem: dbOp.ExternalSystem,
		ExternalMethod: dbOp.ExternalMethod,
		CreatedAt:      dbOp.CreatedAt,
		UpdatedAt:      dbOp.UpdatedAt,
	}

	if dbOp.ExternalID != nil {
		op.ExternalID = *dbOp.ExternalID
	}

	if dbOp.ExternalStatus != nil {
		op.ExternalStatus = model.OperationExternalStatus(*dbOp.ExternalStatus)
	}

	if dbOp.ToolDisplayed != nil {
		op.ToolDisplayed = *dbOp.ToolDisplayed
	}

	if dbOp.FailReason != nil {
		op.FailReason = *dbOp.FailReason
	}

	if dbOp.ProcessedAt != nil {
		op.ProcessedAt = *dbOp.ProcessedAt
	}

	return op
}
