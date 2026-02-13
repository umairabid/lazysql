package conn_manager

import (
	"fmt"
)

func connectWithDatabase(host string, port string, username string, password string) {
	fmt.Printf("Connecting to database at %s:%s with username %s and password %s\n", host, port, username, password)
}
