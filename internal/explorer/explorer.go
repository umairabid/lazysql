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
	databaseNames []string
	databaseLoadError string
	databaseList ExplorerNodeModel
}

func InitExplorer(database adapters.Database) ExplorerModel {
	return ExplorerModel{
		database: database,
		databaseNames: []string{},
		databaseLoadError: "",
		databaseList: ExplorerNodeModel{Title: ""},
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case DatabasesLoadedError:
		m.databaseLoadError = string(msg)
	case DatabasesLoaded:
		m.databaseLoadError = ""
		m.databaseNames = []string(msg)
		cmd = m.databaseList.setNodes(m.databaseNames)
	}

	return m, cmd
}

func (m ExplorerModel) View() string {

}
