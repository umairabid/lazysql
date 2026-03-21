package adapters

import (
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
