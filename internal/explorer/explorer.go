package explorer

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	adapters "app.lazygit/internal/adapters"
)

type DatabasesLoadedError string
type DatabasesLoaded []ExplorerNodeModel

type ExplorerModel struct {
	database adapters.Database
	databaseLoadError string
	databaseList tea.Model
}

func InitExplorer(database adapters.Database) ExplorerModel {
	return ExplorerModel{
		database: database,
		databaseLoadError: "",
		databaseList: ExplorerNodeModel{Title: "", Expanded: true, Type: "root", Selected: true},
	}
}

func (m ExplorerModel) loadDatabases() tea.Cmd {
	return func() tea.Msg {
		databases, err := m.database.GetDatabases()
		if err != nil {
			return DatabasesLoadedError(fmt.Sprintf("Failed to load databases: %v", err))
		}
		var nodes []ExplorerNodeModel
		for _, db := range databases {
			node := m.databaseList.(ExplorerNodeModel)
			nodes = append(nodes, ExplorerNodeModel{Title: db, Type: "database", Parent: &node})
		}
		nodes[0].Selected = true
		return DatabasesLoaded(nodes)
	}
}

func (m ExplorerModel) Init() tea.Cmd {
	return tea.Batch(m.loadDatabases(), m.databaseList.Init())
}

func (m ExplorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd, databaseListCmd tea.Cmd
	switch msg := msg.(type) {
	case DatabasesLoadedError:
		m.databaseLoadError = string(msg)
	case DatabasesLoaded:
		m.databaseLoadError = ""
		var nodeCmds []tea.Cmd
		for _, node := range msg {
			nodeCmds = append(nodeCmds, node.Init())
		}
		dbList := m.databaseList.(ExplorerNodeModel)
		nodeCmds = append(nodeCmds, dbList.setNodes(msg))
		cmd = tea.Batch(nodeCmds...)
	}
	m.databaseList, databaseListCmd = m.databaseList.Update(msg)
	return m, tea.Batch(cmd, databaseListCmd)
}

func (m ExplorerModel) View() string {
	return fmt.Sprintf("%s", m.databaseList.View())
}
