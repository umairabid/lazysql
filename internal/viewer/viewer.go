package viewer

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ViewerModel struct {
}

func InitViewer() ViewerModel {
	return ViewerModel{}
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
