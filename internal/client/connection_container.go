package client

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionContainerModel struct {
}

func InitConnectionContainer() ConnectionContainerModel {
	return ConnectionContainerModel{}
}

func (m ConnectionContainerModel) Init() tea.Cmd {
	return nil
}

func (m ConnectionContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConnectionContainerModel) View() string {
	return "Welcome to LazyGit!\n\nPress q to quit."
}
