package operation

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func (r *Repository) dbUpdateOne(ctx context.Context, dbTx pgx.Tx, dbOp dbOperation) error {
	_, err := dbTx.Exec(ctx, fmt.Sprintf(`
UPDATE %v
SET status          = $2,
    external_id     = $3,
    external_status = $4,
    updated_at      = NOW()
WHERE id = $1
`, operationTable),
		dbOp.ID,
		dbOp.Status,
		dbOp.ExternalID,
		dbOp.ExternalStatus,
	)
	if err != nil {
		return err
	}

	_, err = dbTx.Exec(ctx, fmt.Sprintf(`
UPDATE %v
SET tool_id               = $2,
    additional            = $3,
    fail_reason           = $4,
    confirmation_code     = $5,
    processed_at          = $6,
    confirmation_attempts = $7
WHERE operation_id = $1
`, operationMetadataTable),
		dbOp.ID,
		dbOp.ToolID,
		dbOp.Additional,
		dbOp.FailReason,
		dbOp.ConfirmationCode,
		dbOp.ProcessedAt,
		dbOp.ConfirmationAttempts,
	)
	if err != nil {
		return err
	}

	return nil
}
