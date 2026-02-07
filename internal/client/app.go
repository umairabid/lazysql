package client

import (
	conn_manager "app.lazygit/internal/conn_manager"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionFlag bool
type AppModel struct {
	checking_connection bool
	has_connection      ConnectionFlag
	current_view        tea.Model
}

func StartApp() {
	tea.NewProgram(initModel(), tea.WithAltScreen()).Run()
}

func initModel() AppModel {
	return AppModel{
		checking_connection: true,
		current_view:        SplashModel{},
	}
}

func (m AppModel) Init() tea.Cmd {
	return checkConnection()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.current_view, cmd = m.current_view.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case ConnectionFlag:
		m.checking_connection = false
		m.has_connection = msg
		if m.has_connection {
			m.current_view = conn_manager.InitConnManager()
			m.current_view.Init()
		}
		return m, nil
	}
	return m, cmd
}

func (m AppModel) View() string {
	return m.current_view.View()
}

func checkConnection() tea.Cmd {
	return func() tea.Msg {
		return ConnectionFlag(true)
	}
}
