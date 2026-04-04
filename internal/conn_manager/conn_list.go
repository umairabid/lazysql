package conn_manager

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionList struct {
	connections        []adapters.DbConnection
	selectedConnection int
	viewport           viewport.Model
	layout             utils.ConnectionManagerLayout
}

func InitConnectionList(connections []adapters.DbConnection, layout utils.ConnectionManagerLayout) ConnectionList {
	viewport := viewport.New(layout.ConnectionListWidth, layout.BodyHeight)
	model := ConnectionList{
		connections:        connections,
		selectedConnection: 0,
		layout:             layout,
		viewport:           viewport,
	}
	model.viewport.SetContent(model.connectionsUI())
	return model
}

func (m ConnectionList) changeSelectedConnection() tea.Cmd {
	return func() tea.Msg {
		return SelectedConnectionMsg(m.connections[m.selectedConnection])
	}
}

func (m ConnectionList) connectionsUI() string {
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

func (m ConnectionList) Init() tea.Cmd {
	return nil
}

func (m ConnectionList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var viewPortCmd tea.Cmd
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
	case SelectedConnectionMsg:
		m.viewport.SetContent(m.connectionsUI())
	case LayoutUpdated:
		m.layout = utils.ConnectionManagerLayout(msg)
	}
	m.viewport, viewPortCmd = m.viewport.Update(msg)
	return m, tea.Batch(cmd, viewPortCmd)
}

func (m ConnectionList) View() string {
	info := fmt.Sprintf("Current rows: %d\nTotal rows: %d", m.viewport.VisibleLineCount(), m.viewport.TotalLineCount())
	return fmt.Sprintf("%s\n%s", m.viewport.View(), info)
}
