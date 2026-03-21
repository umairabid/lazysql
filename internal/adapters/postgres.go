package adapters

import (
	"database/sql"
	"fmt"
)

type Postgres struct {
	dbConnection    *DbConnection
	db              *sql.DB
	currentDatabase string
}

func InitPostgres(dbConnection *DbConnection) Postgres {
	return Postgres{dbConnection: dbConnection, db: nil, currentDatabase: "postgres"}
}

func (p Postgres) execute(database string, query string, params ...any) (*sql.Rows, error) {
	var err error

	if p.db != nil && p.currentDatabase != database {
		p.db.Close()
		p.db = nil
	}

	if p.db == nil {
		p.db, err = sql.Open(p.dbConnection.Driver, p.dbConnection.String(database))

		if err != nil {
			return nil, err
		}
	}

	return p.db.Query(query, params...)
}

func (p Postgres) GetDatabases() ([]string, error) {
	rows, err := p.execute("postgres", "SELECT datname FROM pg_database WHERE NOT datistemplate;")
	if err != nil {
		return nil, err
	}
	var databases []string
	for rows.Next() {
		var datname string
		if err := rows.Scan(&datname); err != nil {
			return nil, err
		}
		databases = append(databases, datname)
	}

	return databases, rows.Err()
}

func (p Postgres) GetTables(database string) ([]string, error) {
	rows, err := p.execute(database, "SELECT table_name FROM information_schema.tables WHERE table_catalog = $1 AND table_schema NOT IN ('pg_catalog', 'information_schema');", database)
	if err != nil {
		return nil, err
	}
	var tables []string
	for rows.Next() {
		var tablename string
		if err := rows.Scan(&tablename); err != nil {
			return nil, err
		}
		tables = append(tables, tablename)
	}

	return tables, rows.Err()
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
				row[i] = ""
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
