package tool

import (
	"context"
	"fmt"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) Save(ctx context.Context, tool *model.Tool) error {
	dbT := toolToDB(tool)

	toolID, err := r.dbSave(ctx, dbT)
	if err != nil {
		return err
	}

	// проставляем идентификатор созданного инструмента, чтобы использовать его в дальнейшем
	tool.ID = toolID
	tool.UpdatedAt = time.Now().UTC()

	return nil
}

func (r *Repository) dbSave(ctx context.Context, dbT dbTool) (int64, error) {
	toolIDPlaceholder := "DEFAULT"
	if dbT.ID != 0 {
		toolIDPlaceholder = "$1"
	}

	var toolID int64
	if err := r.conn.QueryRow(ctx, fmt.Sprintf(`
INSERT INTO %[1]v (id,
                user_id,
                external_method,
                type,
                details,
                displayed,
                fake)
VALUES (%[2]v, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) 
	DO UPDATE SET 
		updated_at = NOW() 
	WHERE %[1]v.id = $1 AND %[1]v.user_id = $2 AND %[1]v.external_method = $3
RETURNING id`,
		toolTable, toolIDPlaceholder),
		dbT.ID,
		dbT.UserID,
		dbT.ExternalMethod,
		dbT.Type,
		dbT.Details,
		dbT.Displayed,
		dbT.Fake,
	).Scan(
		&toolID,
	); err != nil {
		return 0, err
	}

	return toolID, nil
}
