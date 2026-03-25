package explorer

import (
	"fmt"
	"strings"

	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type DatabaseError string
type DatabasesLoaded []string
type TablesLoaded []string
type TableLoaded []string

type ExplorerModel struct {
	database      adapters.Database
	databaseError string
	explorerList  utils.ExplorerList
}

func InitExplorer(database adapters.Database) ExplorerModel {
	list := utils.ExplorerList{}
	list.Initialize()
	return ExplorerModel{
		database:      database,
		databaseError: "",
		explorerList:  list,
	}
}

func (m ExplorerModel) loadDatabases() tea.Cmd {
	return func() tea.Msg {
		databases, err := m.database.GetDatabases()
		if err != nil {
			return DatabaseError(fmt.Sprintf("Failed to load databases: %v", err))
		}
		return DatabasesLoaded(databases)
	}
}

func (m ExplorerModel) expandSelectedNode() tea.Cmd {
	return func() tea.Msg {
		if m.explorerList.Selected.Type == "database" {
			database := m.explorerList.Selected.Title
			tables, err := m.database.GetTables(database)
			if err != nil {
				return DatabaseError(fmt.Sprintf("Failed to load tables for database %s: %v", database, err))
			}
			return TablesLoaded(tables)
		} else if m.explorerList.Selected.Type == "table" {
			return TableLoaded([]string{"data", "schema", "indexes"})
		} else if m.explorerList.Selected.Type == "table_item" {
			tableItem := m.explorerList.Selected.Title
			tableName := m.explorerList.Selected.Parent.Title
			database := m.explorerList.Selected.Parent.Parent.Title
			itemData, err := m.database.GetTableItem(database, tableName, tableItem)
			if err != nil {
				return DatabaseError(fmt.Sprintf("Failed to load data for table item %s: %v", tableItem, err))
			}
			return TableLoaded(itemData)
		} else {
			return nil
		}
	}
}

func (m ExplorerModel) createDatabaseList(databases []string) utils.ExplorerList {
	var nodes []utils.ExplorerNode
	for _, db := range databases {
		nodes = append(nodes, utils.ExplorerNode{Title: db, Type: "database"})
	}
	m.explorerList.Expand(nodes)
	return m.explorerList
}

func (m ExplorerModel) Init() tea.Cmd {
	return m.loadDatabases()
}

func (m ExplorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m, cmd = m.handleKeyboardActions(msg)
	switch msg := msg.(type) {
	case DatabaseError:
		m.databaseError = string(msg)
	case DatabasesLoaded:
		m.databaseError = ""
		m.explorerList = m.createDatabaseList([]string(msg))
	case TablesLoaded:
		m.databaseError = ""
		var nodes []utils.ExplorerNode
		for _, table := range msg {
			nodes = append(nodes, utils.ExplorerNode{Title: table, Type: "table"})
		}
		m.explorerList.Expand(nodes)
	case TableLoaded:
		m.databaseError = ""
		var nodes []utils.ExplorerNode
		for _, item := range msg {
			nodes = append(nodes, utils.ExplorerNode{Title: item, Type: "table_item"})
		}
		m.explorerList.Expand(nodes)
	}
	return m, cmd
}

func (m ExplorerModel) handleKeyboardActions(msg tea.Msg) (ExplorerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			cmd = m.expandSelectedNode()
		case "h":
			if m.explorerList.Selected.Expanded {
				m.explorerList.Contract()
			} else {
				m.explorerList.ContractParent()
			}
		case "j":
			m.explorerList.MoveDown()
		case "k":
			m.explorerList.MoveUp()
		}
	}
	return m, cmd
}

func (m ExplorerModel) View() string {
	return fmt.Sprintf("%s\n%s", m.ListNode(m.explorerList.Root, 0), m.databaseError)
}

func (m ExplorerModel) ListNode(node *utils.ExplorerNode, indent int) string {
	prefix := strings.Repeat("  ", indent)
	var newIndent int
	var result string
	var icon string
	var style lipgloss.Style

	if node.Type == "root" {
		result = ""
		newIndent = indent
	} else {
		if node.Expanded {
			icon = "[-] "
		} else if node.Type == "table_item" {
			icon = ""
		} else {
			icon = "[+] "
		}

		if m.explorerList.Selected == node {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
		} else {
			style = lipgloss.NewStyle()
		}
		title := style.Render(fmt.Sprintf("%s%s", icon, node.Title))
		result = fmt.Sprintf("%s%s\n", prefix, title)
		newIndent = indent + 1
	}
	if node.Expanded {
		for i := range node.Children {
			result += m.ListNode(&node.Children[i], newIndent)
		}
	}
	return result
}
