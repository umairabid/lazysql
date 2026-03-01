package client

import (
	conn_manager "app.lazygit/internal/conn_manager"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	current_view tea.Model
}

func StartApp() {
	tea.NewProgram(initModel(), tea.WithAltScreen()).Run()
}

func initModel() AppModel {
	return AppModel{
		current_view: conn_manager.InitConnectionManager(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.current_view.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.current_view, cmd = m.current_view.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case conn_manager.ConnectedMsg:
		m.current_view = InitConnectionContainer()
		cmd = m.current_view.Init()
	}

	return m, cmd
}

func (m AppModel) View() string {
	return m.current_view.View()
}
