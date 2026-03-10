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
	databaseList ExplorerNodeModel
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
		fmt.Println("Loaded databases:", databases)
		if err != nil {
			return DatabasesLoadedError(fmt.Sprintf("Failed to load databases: %v", err))
		}
		var nodes []ExplorerNodeModel
		for _, db := range databases {
			nodes = append(nodes, ExplorerNodeModel{Title: db, Type: "database", Parent: &m.databaseList})
		}
		fmt.Println(m.databaseList)
		fmt.Println("Created nodes:", nodes)
		return DatabasesLoaded(nodes)
	}
}

func (m ExplorerModel) Init() tea.Cmd {
	return tea.Batch(m.loadDatabases(), m.databaseList.Init())
}

func (m ExplorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Println("hello there")
	var cmd tea.Cmd
	m.databaseList, databaseListCmd = m.databaseList.Update(msg)
	switch msg := msg.(type) {
	case DatabasesLoadedError:
		m.databaseLoadError = string(msg)
	case DatabasesLoaded:
		m.databaseLoadError = ""
		var nodeCmds []tea.Cmd
		for _, node := range msg {
			nodeCmds = append(nodeCmds, node.Init())
		}
		nodeCmds = append(nodeCmds, m.databaseList.setNodes(msg))
		cmd = tea.Batch(nodeCmds...)
	}
	return m, tea.Batch(cmd, databaseListCmd)
}

func (m ExplorerModel) View() string {
	return fmt.Sprintf("Database Explorer\n\n%s\n\n%s", m.databaseList.View(), m.databaseList.Selected)
}
