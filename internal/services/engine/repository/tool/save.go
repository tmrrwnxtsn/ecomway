package tool

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) Save(ctx context.Context, tool *model.Tool) error {
	if tool == nil || tool.ID == "" {
		return errors.New("saving nil tool or tool with empty ID")
	}

	dbT := toolToDB(tool)

	if err := r.dbSave(ctx, dbT); err != nil {
		return err
	}

	tool.UpdatedAt = time.Now().UTC()

	return nil
}

func (r *Repository) dbSave(ctx context.Context, dbT dbTool) error {
	updateStmtBuilder := new(strings.Builder)
	if dbT.Name != "" {
		updateStmtBuilder.WriteString("name = $7, ")
	}
	if dbT.Status != "" {
		updateStmtBuilder.WriteString("status = $8, ")
	}
	updateStmtBuilder.WriteString("updated_at = NOW()")

	_, err := r.conn.Exec(ctx, fmt.Sprintf(`
INSERT INTO %[1]v (id,
                user_id,
                external_method,
                type,
                details,
                displayed,
                name,
                fake)
VALUES ($1, $2, $3, $4, $5, $6, $7, $9)
ON CONFLICT (id, user_id, external_method) 
	DO UPDATE SET 
		%[2]v 
	WHERE %[1]v.id = $1 AND %[1]v.user_id = $2 AND %[1]v.external_method = $3`,
		toolTable, updateStmtBuilder.String()),
		dbT.ID,
		dbT.UserID,
		dbT.ExternalMethod,
		dbT.Type,
		dbT.Details,
		dbT.Displayed,
		dbT.Name,
		dbT.Status,
		dbT.Fake,
	)
	return err
}
