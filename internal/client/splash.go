package client

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SplashModel struct {
}

func InitSplash() SplashModel {
	return SplashModel{}
}

func (m SplashModel) Init() tea.Cmd {
	return nil
}

func (m SplashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SplashModel) View() string {
	return "Welcome to LazyGit!\n\nPress q to quit."
}
