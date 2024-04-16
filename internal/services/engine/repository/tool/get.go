package tool

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) All(ctx context.Context, userID int64) ([]*model.Tool, error) {
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

func (r *Repository) dbGetAll(ctx context.Context, userID int64) ([]dbTool, error) {
	var dbTools []dbTool
	err := pgxscan.Select(ctx, r.conn, &dbTools, fmt.Sprintf(`
SELECT id,
       user_id,
       external_method,
       type,
       details,
       displayed,
       fake,
       created_at,
       updated_at
FROM %v 
WHERE user_id = $1
`, toolTable), userID)
	if err != nil {
		return nil, err
	}
	return dbTools, nil
}
