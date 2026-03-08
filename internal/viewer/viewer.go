package viewer

import (
	tea "github.com/charmbracelet/bubbletea"
	adapters "app.lazygit/internal/adapters"
)

type ViewerModel struct {
	database adapters.Database
}

func InitViewer(database adapters.Database) ViewerModel {
	return ViewerModel{
		database: database,
	}
}

func (m ViewerModel) Init() tea.Cmd {
	return nil
}

func (m ViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ViewerModel) View() string {
	return "Welcome to LazyGit Viewer!\n\nPress q to quit."
}
