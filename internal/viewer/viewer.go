package viewer

import (
	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewerModel struct {
	database adapters.Database
	content  string
}

func InitViewer(database adapters.Database) ViewerModel {
	return ViewerModel{
		database: database,
		content:  "",
	}
}

func (m ViewerModel) Init() tea.Cmd {
	return nil
}

func (m ViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case utils.ViewerTableData:
		m.content = "There is new content"
	}
	return m, cmd
}

func (m ViewerModel) View() string {
	if m.content != "" {
		return m.content
	} else {
		return "Welcome to LazyGit Viewer!\n\nPress q to quit."
	}
}
