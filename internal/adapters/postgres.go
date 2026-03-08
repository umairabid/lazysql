package adapters

import (
	"fmt"
	"database/sql"
)

type Postgres struct {
	db *sql.DB
}

func InitPostgres(db *sql.DB) Postgres {
	return Postgres{db: db}
}

func (p Postgres) GetDatabases() ([]string, error) {
	rows, err := p.db.Query("SELECT datname FROM pg_database;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println(rows)
	return []string{"a"}, nil
}

func (p Postgres) GetTables() ([]string, error) {
	return []string{"a"}, nil
}

func (p Postgres) Query(query string) (*sql.Rows, error) {
	return p.db.Query(query)
}

func (p Postgres) InpsectRows(rows *sql.Rows) ([][]string, error) {
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	result := [][]string{columns}

	vals := make([]any, len(columns))
	valPtrs := make([]any, len(columns))
	for i := range vals {
		valPtrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(valPtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		row := make([]string, len(columns))
		for i, val := range vals {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
