package operation

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
)

func (r *Repository) dbRollback(ctx context.Context, dbTX pgx.Tx) {
	err := dbTX.Rollback(ctx)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		log.Printf("rollback db transaction: %v", err)
	}
}
