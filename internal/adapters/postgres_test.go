package adapters

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestPostgresGetDatabases(t *testing.T) {
	dbConnection := DbConnection{
		Name:     "testdb",
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		Driver:   "pgx",
	}

	databse, err := dbConnection.InitConnection()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	postgres := databse.(Postgres)
	result, err := postgres.GetDatabases()
	if err != nil {
		t.Fatalf("Failed to get databases: %v", err)
	}
	if len(result) == 0 {
		t.Fatalf("Expected at least one database, got 0")
	}
}

func TestPostgresGetTables(t *testing.T) {
	dbConnection := DbConnection{
		Name:     "testdb",
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		Driver:   "pgx",
	}

	database, err := dbConnection.InitConnection()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	postgres := database.(Postgres)
	result, err := postgres.GetTables("postgres")
	if err != nil {
		t.Fatalf("Failed to get tables: %v", err)
	}
	if len(result) == 0 {
		t.Fatalf("Expected at least one table, got 0")
	}
}

func TestPostgresGetTableItem(t *testing.T) {
	const testDBName = "lazysql_test_db"
	const testTable = "users"

	dbConnection := DbConnection{
		Name:     "testdb",
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		Driver:   "pgx",
	}

	database, err := dbConnection.InitConnection()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	postgres := database.(Postgres)

	// Setup: create test database using a direct connection to the default postgres db
	adminDB, err := sql.Open(dbConnection.Driver, dbConnection.String("postgres"))
	if err != nil {
		t.Fatalf("Failed to open admin connection: %v", err)
	}
	defer adminDB.Close()

	adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if _, err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName)); err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", testDBName))
	})

	// Connect to the new test database and create a table with data
	testDB, err := sql.Open(dbConnection.Driver, dbConnection.String(testDBName))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	if _, err := testDB.Exec(fmt.Sprintf(`
		CREATE TABLE %s (
			id    SERIAL PRIMARY KEY,
			name  TEXT NOT NULL,
			email TEXT
		)`, testTable)); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	if _, err := testDB.Exec(fmt.Sprintf(
		`INSERT INTO %s (name, email) VALUES ('Alice', 'alice@example.com'), ('Bob', NULL)`,
		testTable,
	)); err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	t.Run("data", func(t *testing.T) {
		result, err := postgres.GetTableItem(testDBName, testTable, "data")
		if err != nil {
			t.Fatalf("GetTableItem(data) failed: %v", err)
		}
		// result[0] is the header row; result[1..] are data rows
		if len(result) != 3 {
			t.Fatalf("Expected 3 rows (1 header + 2 data), got %d", len(result))
		}
		if len(result[0]) != 3 {
			t.Errorf("Expected 3 columns (id, name, email), got %d", len(result[0]))
		}
	})

	t.Run("schema", func(t *testing.T) {
		result, err := postgres.GetTableItem(testDBName, testTable, "schema")
		if err != nil {
			t.Fatalf("GetTableItem(schema) failed: %v", err)
		}
		// result[0] is the header row (column_name, data_type); result[1..] are the table's columns
		if len(result) != 4 {
			t.Fatalf("Expected 4 rows (1 header + 3 columns), got %d", len(result))
		}
		headers := result[0]
		if headers[0] != "column_name" || headers[1] != "data_type" {
			t.Errorf("Unexpected schema headers: %v", headers)
		}
	})

	t.Run("indexes", func(t *testing.T) {
		result, err := postgres.GetTableItem(testDBName, testTable, "indexes")
		if err != nil {
			t.Fatalf("GetTableItem(indexes) failed: %v", err)
		}
		// SERIAL PRIMARY KEY auto-creates an index
		if len(result) < 2 {
			t.Fatalf("Expected at least 2 rows (1 header + 1 index), got %d", len(result))
		}
		headers := result[0]
		if headers[0] != "indexname" || headers[1] != "indexdef" {
			t.Errorf("Unexpected indexes headers: %v", headers)
		}
	})

	t.Run("unknown item returns error", func(t *testing.T) {
		_, err := postgres.GetTableItem(testDBName, testTable, "unknown")
		if err == nil {
			t.Error("Expected an error for unknown item type, got nil")
		}
	})
}
