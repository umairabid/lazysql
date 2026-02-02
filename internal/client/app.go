package client

import (
	"fmt"

	conn_manager "app.lazygit/internal/conn_manager"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionFlag bool
type AppModel struct {
	checking_connection bool
	has_connection      ConnectionFlag
	conn_manager        conn_manager.ConnectionManager
}

func StartApp() {
	p := tea.NewProgram(initModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}

func initModel() AppModel {
	return AppModel{
		checking_connection: true,
	}
}

func (m AppModel) Init() tea.Cmd {
	return checkConnection()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.conn_manager = conn_manager.InitConnManager()
			m.conn_manager.Init()
		}
		return m, nil
	}
	return m, nil
}

func (m AppModel) View() string {
	if m.checking_connection {
		return "Checking connection...\n"
	} else {
		return m.conn_manager.View()
	}
}

func checkConnection() tea.Cmd {
	return func() tea.Msg {
		return ConnectionFlag(true)
	}
}
