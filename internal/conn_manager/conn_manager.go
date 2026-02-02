package conn_manager

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionManager struct {
	host string
	port int
	password string
	username string
}

func InitConnManager() ConnectionManager {
	return ConnectionManager{
		host: "localhost",
		port: 8080,
		password: "admin",
		username: "admin",
	}
}

func (m ConnectionManager) Init() tea.Cmd {
	return nil
}

func (m ConnectionManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConnectionManager) View() string {
	return fmt.Sprintf("Conection Manager Here\n")
}
