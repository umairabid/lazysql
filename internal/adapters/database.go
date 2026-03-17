package adapters

import (
	"fmt"

	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbConnection struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
	Driver   string
}

type Database interface {
	GetDatabases() ([]string, error)
	GetTables(string) ([]string, error)
}

func (c DbConnection) String() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s", c.Username, c.Password, c.Host, c.Port)
}

func (c DbConnection) InitConnection() (Database, error) {
	err := c.validateConnection()
	var db *sql.DB
	if err != nil {
		return nil, err
	}

	db, err = sql.Open(c.Driver, c.String())
	if err != nil {
		return nil, err
	}

	var greeting string
	err = db.QueryRow("select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return nil, err
	}
	if c.Driver == "pgx" {
		return InitPostgres(&c), nil
	}
	defer db.Close()
	return nil, fmt.Errorf("unsupported driver: %s", c.Driver)
}

func (c DbConnection) validateConnection() error {
	var errorMessage string
	if c.Host == "" {
		errorMessage += "Host is required. "
	}
	if c.Port == "" {
		errorMessage += "Port is required. "
	}
	if c.Username == "" {
		errorMessage += "Username is required. "
	}
	if c.Password == "" {
		errorMessage += "Password is required. "
	}
	if c.Driver == "" {
		errorMessage += "Driver is required. "
	}
	if errorMessage != "" {
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}
