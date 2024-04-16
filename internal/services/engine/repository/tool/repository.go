package tool

import "github.com/jackc/pgx/v4/pgxpool"

const (
	toolTable     = "tool"
	toolTableAbbr = "t"
)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{
		conn: conn,
	}
}
