package repository

import (
	"database/sql"

	"github.com/lassulfi/imersao11-consolidacao/internal/infra/db"
)

type Repository struct {
	dbConn *sql.DB
	*db.Queries
}
