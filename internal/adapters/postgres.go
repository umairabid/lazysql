package adapters

import (
	"database/sql"
)

type Postgres struct {
	db *sql.DB
}

func InitPostgres(db *sql.DB) Postgres {
	return Postgres{db: db}
}

func (p Postgres) getTables() ([]string, error) {
	return []string{"a"}, nil
}
