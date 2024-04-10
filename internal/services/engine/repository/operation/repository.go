package operation

import "github.com/jackc/pgx/v4/pgxpool"

const (
	operationTable     = "operation"
	operationTableAbbr = "op"

	operationMetadataTable     = "operation_metadata"
	operationMetadataTableAbbr = "op_meta"
)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{
		conn: conn,
	}
}
