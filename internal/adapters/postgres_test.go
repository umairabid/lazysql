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
