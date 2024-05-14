package tool

import (
	"context"
	"errors"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) Create(ctx context.Context, tool *model.Tool) error {
	if tool == nil {
		return errors.New("creating nil tool")
	}

	dbT := toolToDB(tool)

	createdDBT, err := r.dbCreate(ctx, dbT)
	if err != nil {
		return err
	}

	tool = toolFromDB(createdDBT)

	return nil
}

func (r *Repository) dbCreate(ctx context.Context, dbT dbTool) (dbTool, error) {
	var created dbTool

	if err := r.conn.QueryRow(ctx, `
INSERT INTO tool (id, user_id, external_method, type, details, displayed, name)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, external_method, type, details, displayed, name, status, fake, created_at, updated_at`,
		dbT.ID,
		dbT.UserID,
		dbT.ExternalMethod,
		dbT.Type,
		dbT.Details,
		dbT.Displayed,
		dbT.Name,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.ExternalMethod,
		&created.Type,
		&created.Details,
		&created.Displayed,
		&created.Name,
		&created.Status,
		&created.Fake,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return dbTool{}, err
	}

	return created, nil
}
