package conn_manager

import (
	"fmt"

	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Connection struct {
	name	 string
	host     string
	port     string
	username string
	password string
	driver   string
}

func (c Connection) String() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s", c.username, c.password, c.host, c.port)
}

func connectWithDatabase(connection Connection) (*sql.DB, error) {
	err := validateConnection(connection)
	if err != nil {
		return nil, err
	}

	var db *sql.DB
	db, err = sql.Open(connection.driver, connection.String())
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
	if connection.host == "" {
		errorMessage += "Host is required. "
	}
	if connection.port == "" {
		errorMessage += "Port is required. "
	}
	if connection.username == "" {
		errorMessage += "Username is required. "
	}
	if connection.password == "" {
		errorMessage += "Password is required. "
	}
	if connection.driver == "" {
		errorMessage += "Driver is required. "
	}
	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}
	return nil
}
