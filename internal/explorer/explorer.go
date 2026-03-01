package explorer

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ExplorerModel struct {
}

func InitExplorer() ExplorerModel {
	return ExplorerModel{}
}

func (m ExplorerModel) Init() tea.Cmd {
	return nil
}

func (m ExplorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ExplorerModel) View() string {
	return "Welcome to LazyGit Explorer!\n\nPress q to quit."
}
