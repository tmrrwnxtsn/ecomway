package tool

import (
	"context"
	"errors"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) Update(ctx context.Context, tool *model.Tool) error {
	if tool == nil {
		return errors.New("updating nil tool")
	}

	dbT := toolToDB(tool)

	updatedDBT, err := r.dbUpdate(ctx, dbT)
	if err != nil {
		return err
	}

	tool = toolFromDB(updatedDBT)

	return nil
}

func (r *Repository) dbUpdate(ctx context.Context, dbT dbTool) (dbTool, error) {
	var updated dbTool

	if err := r.conn.QueryRow(ctx, `
UPDATE tool
SET name = $4, status = $5, updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND external_method = $3
RETURNING id, user_id, external_method, type, details, displayed, name, status, fake, created_at, updated_at`,
		dbT.ID,
		dbT.UserID,
		dbT.ExternalMethod,
		dbT.Name,
		dbT.Status,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.ExternalMethod,
		&updated.Type,
		&updated.Details,
		&updated.Displayed,
		&updated.Name,
		&updated.Status,
		&updated.Fake,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	); err != nil {
		return dbTool{}, err
	}

	return updated, nil
}
