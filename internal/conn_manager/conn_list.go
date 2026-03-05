package conn_manager

import (
	tea "github.com/charmbracelet/bubbletea"
	postgres "app.lazygit/internal/services/postgres"
	"github.com/charmbracelet/lipgloss"
)

type ConnectionList struct {
	connections        []postgres.Connection
	selectedConnection int
}

func InitConnectionList(connections []postgres.Connection) ConnectionList {
	return ConnectionList{connections: connections, selectedConnection: 0}
}

func (m ConnectionList) Init() tea.Cmd {
	return nil
}

func (m ConnectionList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedConnection > 0 {
				m.selectedConnection--
			}
			cmd = m.changeSelectedConnection()
		case "down", "j":
			if m.selectedConnection < len(m.connections)-1 {
				m.selectedConnection++
			}
			cmd = m.changeSelectedConnection()
		}
	}
	return m, cmd
}

func (m ConnectionList) View() string {
	var result string

	normalStyle := lipgloss.NewStyle().
		Padding(0, 2)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("229")).
		Padding(0, 2)

	for i, conn := range m.connections {
		if i == m.selectedConnection {
			result += selectedStyle.Render(conn.Name) + "\n"
		} else {
			result += normalStyle.Render(conn.Name) + "\n"
		}
	}

	return result
}

func (m ConnectionList) changeSelectedConnection() tea.Cmd {
	return func() tea.Msg {
		return SelectedConnectionMsg(m.connections[m.selectedConnection])
	}
}
