package operation

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) Create(ctx context.Context, op *model.Operation) error {
	dbTX, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(dbTX pgx.Tx, ctx context.Context) {
		if err = dbTX.Rollback(ctx); err != nil {
			log.Printf("rollback db transaction: %v", err)
		}
	}(dbTX, ctx)

	dbOp := operationToDB(op)
	if dbOp.ID != 0 {
		return fmt.Errorf("creating operation with existing ID: %v", dbOp.ID)
	}

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
		return err
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
		return err
	}

	return dbTX.Commit(ctx)
}
