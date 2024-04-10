package operation

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) updateOne(ctx context.Context, dbTx pgx.Tx, op *model.Operation) error {
	_, err := dbTx.Exec(ctx, fmt.Sprintf(`
UPDATE %v
SET status          = $2,
    external_id     = $3,
    external_status = $4,
    updated_at      = NOW()
WHERE id = $1
`, operationTable),
		op.ID,
		op.Status,
		op.ExternalID,
		op.ExternalStatus,
	)
	if err != nil {
		return err
	}

	_, err = dbTx.Exec(ctx, fmt.Sprintf(`
UPDATE %v
SET tool_id           = $2,
    additional        = $3,
    fail_reason       = $4,
    confirmation_code = $5,
    processed_at      = $6
WHERE operation_id = $1
`, operationMetadataTable),
		op.ID,
		op.ToolID,
		op.Additional,
		op.FailReason,
		op.ConfirmationCode,
		op.ProcessedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
