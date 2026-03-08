package explorer

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	adapters "app.lazygit/internal/adapters"
)

type DatabasesLoadedError string
type DatabasesLoaded []string

type ExplorerModel struct {
	database adapters.Database
}

func InitExplorer(database adapters.Database) ExplorerModel {
	return ExplorerModel{
		database: database,
	}
}

func (m ExplorerModel) loadDatabases() tea.Cmd {
	return func() tea.Msg {
		databases, err := m.database.GetDatabases()
		if err != nil {
			return DatabasesLoadedError(fmt.Sprintf("Failed to load databases: %v", err))
		}
		return DatabasesLoaded(databases)
	}
}

func (m ExplorerModel) Init() tea.Cmd {
	return m.loadDatabases()
}

func (m ExplorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ExplorerModel) View() string {
	return "Welcome to LazyGit Explorer!\n\nPress q to quit."
}
