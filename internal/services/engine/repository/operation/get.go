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

const defaultOperationMaxCount = 10000

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
			slog.Info(
				"operation status changed",
				"operation_id", op.ID,
				"old_status", oldStatus,
				"new_status", op.Status,
			)
		}
	}()

	return script(ctx, op)
}

func (r *Repository) GetOneWithoutLock(ctx context.Context, criteria model.OperationCriteria) (*model.Operation, error) {
	dbOp, err := r.dbGetOne(ctx, r.conn, criteria, false)
	if err != nil {
		return nil, err
	}

	op := operationFromDB(dbOp)
	return op, nil
}

func (r *Repository) All(ctx context.Context, criteria model.OperationCriteria) ([]*model.Operation, error) {
	dbOps, err := r.dbGetAll(ctx, criteria)
	if err != nil {
		return nil, err
	}

	ops := make([]*model.Operation, 0, len(dbOps))
	for _, dbOp := range dbOps {
		ops = append(ops, operationFromDB(dbOp))
	}
	return ops, nil
}

func (r *Repository) AllForReport(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error) {
	dbOps, err := r.dbGetAllForReport(ctx, criteria)
	if err != nil {
		return nil, err
	}

	ops := make([]model.ReportOperation, 0, len(dbOps))
	for _, dbOp := range dbOps {
		ops = append(ops, reportOperationFromDB(dbOp))
	}
	return ops, nil
}

func (r *Repository) dbGetOne(ctx context.Context, db pgxscan.Querier, criteria model.OperationCriteria, withLock bool) (dbOperation, error) {
	var dbOp dbOperation

	whereStmt, args, err := r.whereStmt(criteria)
	if err != nil {
		return dbOp, err
	}

	forUpdateStmt := "FOR UPDATE"
	if !withLock {
		forUpdateStmt = ""
	}

	err = pgxscan.Get(ctx, db, &dbOp, fmt.Sprintf(`
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
       %[4]v.processed_at,
       %[4]v.confirmation_attempts
FROM %[1]v %[3]v
         JOIN %[2]v %[4]v on %[3]v.id = %[4]v.operation_id
%v %v
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

func (r *Repository) dbGetAll(ctx context.Context, criteria model.OperationCriteria) ([]dbOperation, error) {
	if criteria.MaxCount == 0 {
		criteria.MaxCount = defaultOperationMaxCount
	}

	whereStmt, args, err := r.whereStmt(criteria)
	if err != nil {
		return nil, err
	}

	var dbOps []dbOperation
	err = pgxscan.Select(ctx, r.conn, &dbOps, fmt.Sprintf(`
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
       %[4]v.processed_at,
       %[4]v.confirmation_attempts
FROM %[1]v %[3]v
         JOIN %[2]v %[4]v on %[3]v.id = %[4]v.operation_id
%v ORDER BY random() LIMIT %v
`, operationTable, operationMetadataTable, operationTableAbbr, operationMetadataTableAbbr, whereStmt, criteria.MaxCount),
		args...)
	if err != nil {
		return nil, err
	}

	return dbOps, nil
}

func (r *Repository) dbGetAllForReport(ctx context.Context, criteria model.OperationCriteria) ([]dbReportOperation, error) {
	if criteria.MaxCount == 0 {
		criteria.MaxCount = defaultOperationMaxCount
	}

	whereStmt, args, _ := r.whereStmt(criteria)

	var dbOps []dbReportOperation
	err := pgxscan.Select(ctx, r.conn, &dbOps, fmt.Sprintf(`
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
       %[6]v.displayed tool_displayed,
       %[4]v.fail_reason,
       %[3]v.created_at,
       %[3]v.updated_at,
       %[4]v.processed_at
FROM %[1]v %[3]v
         JOIN %[2]v %[4]v on %[3]v.id = %[4]v.operation_id
         LEFT JOIN %[5]v %[6]v on %[4]v.tool_id = %[6]v.id AND %[3]v.user_id = %[6]v.user_id AND %[3]v.external_method = %[6]v.external_method
%v LIMIT %v
`, operationTable, operationMetadataTable, operationTableAbbr, operationMetadataTableAbbr, toolTable, toolTableAbbr,
		whereStmt, criteria.MaxCount),
		args...)
	if err != nil {
		return nil, err
	}

	return dbOps, nil
}
