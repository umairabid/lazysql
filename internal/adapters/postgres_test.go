package adapters

import (
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

  databse, err := dbConnection.ConnectWithDatabase()
  if err != nil {
	t.Fatalf("Failed to connect to database: %v", err)
  }
  postgres := databse.(Postgres)
  rows, err := postgres.Query("SELECT * FROM pg_database;")
  if err != nil {
	  t.Fatalf("Failed to execute query: %v", err)
  }
  result, err := postgres.InpsectRows(rows)
  fmt.Println(err)
  fmt.Println(result)

}
