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
	databaseLoadError string
	databaseList tea.Model
}

func InitExplorer(database adapters.Database) ExplorerModel {
	return ExplorerModel{
		database: database,
		databaseLoadError: "",
		databaseList: ExplorerNodeModel{Title: "Loading...", Type: "root", Parent: nil},
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

func (m ExplorerModel) createDatabaseList(databases []string) ExplorerNodeModel {
		model := ExplorerNodeModel{Title: "", Type: "root", Parent: nil}
		var nodes []ExplorerNodeModel
		for _, db := range databases {
			nodes = append(nodes, ExplorerNodeModel{Title: db, Type: "database", Parent: &model})
		}
		nodes[0].Selected = true
		model.Children = nodes
		return model
}

func (m ExplorerModel) Init() tea.Cmd {
	return m.loadDatabases()
}

func (m ExplorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd, databaseListCmd tea.Cmd
	switch msg := msg.(type) {
	case DatabasesLoadedError:
		m.databaseLoadError = string(msg)
	case DatabasesLoaded:
		m.databaseLoadError = ""
		m.databaseList = m.createDatabaseList(msg)
		var nodeCmds []tea.Cmd
		for _, node := range m.databaseList.(ExplorerNodeModel).Children {
			nodeCmds = append(nodeCmds, node.Init())
		}
		nodeCmds = append(nodeCmds, m.databaseList.Init())
		cmd = tea.Batch(nodeCmds...)
	}
	m.databaseList, databaseListCmd = m.databaseList.Update(msg)
	return m, tea.Batch(cmd, databaseListCmd)
}

func (m ExplorerModel) View() string {
	return fmt.Sprintf("%s", m.databaseList.View())
}
