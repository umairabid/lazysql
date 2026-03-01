package client

import (
	explorer "app.lazygit/internal/explorer"
	editor "app.lazygit/internal/editor"
	viewer "app.lazygit/internal/viewer"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionContainerModel struct {
	explorer explorer.ExplorerModel
	editor   editor.EditorModel
	viewer   viewer.ViewerModel
	active_view string
}

func InitConnectionContainer() ConnectionContainerModel {
	return ConnectionContainerModel{
		explorer: explorer.InitExplorer(),
		editor:   editor.InitEditor(),
		viewer:   viewer.InitViewer(),
	}
}

func (m ConnectionContainerModel) Init() tea.Cmd {
	return tea.Batch(
		m.explorer.Init(),
		m.editor.Init(),
		m.viewer.Init(),
	)
}

func (m ConnectionContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConnectionContainerModel) View() string {
	return "Welcome to LazyGit!\n\nPress q to quit."
}
