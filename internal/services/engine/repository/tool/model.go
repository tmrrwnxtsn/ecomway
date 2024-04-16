package tool

import (
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type dbTool struct {
	ID             int64          `db:"id"`
	UserID         int64          `db:"user_id"`
	ExternalMethod string         `db:"external_method"`
	Type           *string        `db:"type"`
	Details        map[string]any `db:"details"`
	Displayed      string         `db:"displayed"`
	Fake           bool           `db:"fake"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

func toolToDB(t *model.Tool) dbTool {
	dbT := dbTool{
		ID:             t.ID,
		UserID:         t.UserID,
		ExternalMethod: t.ExternalMethod,
		Displayed:      t.Displayed,
		Fake:           t.Fake,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}

	if t.Type != "" {
		dbT.Type = (*string)(&t.Type)
	}

	if len(t.Details) > 0 {
		dbT.Details = t.Details
	}

	return dbT
}

func toolFromDB(dbT dbTool) *model.Tool {
	t := &model.Tool{
		ID:             dbT.ID,
		UserID:         dbT.UserID,
		ExternalMethod: dbT.ExternalMethod,
		Displayed:      dbT.Displayed,
		Fake:           dbT.Fake,
		CreatedAt:      dbT.CreatedAt,
		UpdatedAt:      dbT.UpdatedAt,
	}

	if dbT.Type != nil {
		t.Type = model.ToolType(*dbT.Type)
	}

	if len(dbT.Details) > 0 {
		t.Details = dbT.Details
	}

	return t
}
