package operation

import (
	"context"
	"fmt"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) Create(ctx context.Context, op *model.Operation) error {
	dbOp := operationToDB(op)
	if dbOp.ID != 0 {
		return fmt.Errorf("creating operation with existing ID: %v", dbOp.ID)
	}

	operationID, err := r.dbCreate(ctx, dbOp)
	if err != nil {
		return err
	}

	// проставляем идентификатор созданной операции, чтобы использовать его в дальнейшем
	op.ID = operationID
	return nil
}

func (r *Repository) dbCreate(ctx context.Context, dbOp dbOperation) (int64, error) {
	dbTX, err := r.conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer r.dbRollback(ctx, dbTX)

	if err = dbTX.QueryRow(ctx, fmt.Sprintf(`
INSERT INTO %v (user_id,
                type,
                currency,
                amount,
                status,
                external_id,
                external_system,
                external_method)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id`,
		operationTable),
		dbOp.UserID,
		dbOp.Type,
		dbOp.Currency,
		dbOp.Amount,
		dbOp.Status,
		dbOp.ExternalID,
		dbOp.ExternalSystem,
		dbOp.ExternalMethod,
	).Scan(
		&dbOp.ID,
	); err != nil {
		return 0, err
	}

	_, err = dbTX.Exec(ctx, fmt.Sprintf(`
INSERT INTO %v (operation_id,
                tool_id,
                additional,
                fail_reason,
                confirmation_code,
                processed_at)
VALUES ($1, $2, $3, $4, $5, $6)`,
		operationMetadataTable),
		dbOp.ID,
		dbOp.ToolID,
		dbOp.Additional,
		dbOp.FailReason,
		dbOp.ConfirmationCode,
		dbOp.ProcessedAt)
	if err != nil {
		return 0, err
	}

	if err = dbTX.Commit(ctx); err != nil {
		return 0, err
	}

	return dbOp.ID, nil
}
