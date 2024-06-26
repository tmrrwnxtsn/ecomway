package tool

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/pgxscan"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) All(ctx context.Context, userID string) ([]*model.Tool, error) {
	dbTools, err := r.dbGetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	tools := make([]*model.Tool, 0, len(dbTools))
	for _, dbT := range dbTools {
		tools = append(tools, toolFromDB(dbT))
	}
	return tools, nil
}

func (r *Repository) GetOne(ctx context.Context, id string, userID string, externalMethod string) (*model.Tool, error) {
	dbT, err := r.dbGetOne(ctx, id, userID, externalMethod)
	if err != nil {
		return nil, err
	}
	return toolFromDB(dbT), nil
}

func (r *Repository) dbGetOne(ctx context.Context, id string, userID string, externalMethod string) (dbTool, error) {
	var dbT dbTool

	err := pgxscan.Get(ctx, r.conn, &dbT, `
SELECT id,
       user_id,
       external_method,
       type,
       details,
       displayed,
       name,
       status,
       fake,
       created_at,
       updated_at
FROM tool 
WHERE id = $1 AND user_id = $2 AND external_method = $3
`, id, userID, externalMethod)
	if err != nil {
		if pgxscan.NotFound(err) {
			return dbT, sql.ErrNoRows
		}
		return dbT, err
	}

	return dbT, nil
}

func (r *Repository) dbGetAll(ctx context.Context, userID string) ([]dbTool, error) {
	var dbTools []dbTool
	err := pgxscan.Select(ctx, r.conn, &dbTools, `
SELECT id,
       user_id,
       external_method,
       type,
       details,
       displayed,
       name,
       status,
       fake,
       created_at,
       updated_at
FROM tool
WHERE user_id = $1
`, userID)
	if err != nil {
		return nil, err
	}
	return dbTools, nil
}
