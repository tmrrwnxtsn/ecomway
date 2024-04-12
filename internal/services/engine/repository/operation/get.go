package operation

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) AcquireOneLocked(ctx context.Context, criteria model.OperationCriteria, script model.ScriptAcquiredFor) (err error) {
	dbTX, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(dbTX pgx.Tx, ctx context.Context) {
		commitErr := dbTX.Commit(ctx)
		if commitErr != nil {
			err = multierror.Append(err, commitErr)
			return
		}
	}(dbTX, ctx)

	dbOp, err := r.dbGetOne(ctx, dbTX, criteria, true)
	if err != nil {
		return err
	}

	// на время выполнения действий над операцией преобразуем её в общую структуру (модель)
	op := operationFromDB(dbOp)

	oldStatus := op.Status
	defer func() {
		dbOp = operationToDB(op)

		updateErr := r.dbUpdateOne(ctx, dbTX, dbOp)
		if updateErr != nil {
			err = multierror.Append(err, updateErr)
			return
		}

		if oldStatus != op.Status {
			slog.Info("operation status changed", "old_status", oldStatus, "new_status", op.Status)
		}
	}()

	return script(ctx, op)
}

func (r *Repository) dbGetOne(ctx context.Context, dbTX pgx.Tx, criteria model.OperationCriteria, withLock bool) (dbOperation, error) {
	var dbOp dbOperation

	whereStmt, args, err := r.whereStmt(criteria)
	if err != nil {
		return dbOp, err
	}

	forUpdateStmt := "FOR UPDATE"
	if !withLock {
		forUpdateStmt = ""
	}

	err = pgxscan.Get(ctx, dbTX, &dbOp, fmt.Sprintf(`
SELECT %[3]v.id,
       %[3]v.user_id,
       %[3]v.type,
       %[3]v.currency,
       %[3]v.amount,
       %[3]v.status,
       %[3]v.external_id,
       %[3]v.external_system,
       %[3]v.external_method,
       %[3]v.external_status,
       %[3]v.created_at,
       %[3]v.updated_at,
       %[4]v.tool_id,
       %[4]v.additional,
       %[4]v.fail_reason,
       %[4]v.confirmation_code,
       %[4]v.processed_at
FROM %[1]v %[3]v
         JOIN %[2]v %[4]v on %[3]v.id = %[4]v.operation_id
WHERE %v %v
`, operationTable, operationMetadataTable, operationTableAbbr, operationMetadataTableAbbr, whereStmt, forUpdateStmt),
		args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return dbOp, sql.ErrNoRows
		}
		return dbOp, err
	}

	return dbOp, nil
}
