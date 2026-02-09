package conn_manager

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionList struct {
}

func InitConnectionList() ConnectionList {
	return ConnectionList{}
}

func (m ConnectionList) Init() tea.Cmd	 {
	return nil
}

func (m ConnectionList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConnectionList) View() string {
	return ""
}
