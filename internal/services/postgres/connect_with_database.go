package postgres

import (
	"fmt"

	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Connection struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
	Driver   string
}

func (c Connection) String() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s", c.Username, c.Password, c.Host, c.Port)
}

func ConnectWithDatabase(connection Connection) (*sql.DB, error) {
	err := validateConnection(connection)
	var db *sql.DB
	if err != nil {
		return nil, err
	}

	db, err = sql.Open(connection.Driver, connection.String())
	if err != nil {
		return nil, err
	}
	var greeting string
	err = db.QueryRow("select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db, nil
}

func validateConnection(connection Connection) error {
	var errorMessage string
	if connection.Host == "" {
		errorMessage += "Host is required. "
	}
	if connection.Port == "" {
		errorMessage += "Port is required. "
	}
	if connection.Username == "" {
		errorMessage += "Username is required. "
	}
	if connection.Password == "" {
		errorMessage += "Password is required. "
	}
	if connection.Driver == "" {
		errorMessage += "Driver is required. "
	}
	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}
	return nil
}
