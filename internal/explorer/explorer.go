package explorer

import (
	"fmt"
	"strings"

	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type DatabasesLoadedError string
type DatabasesLoaded []string

type ExplorerModel struct {
	database          adapters.Database
	databaseLoadError string
	explorerList      utils.ExplorerList
}

func InitExplorer(database adapters.Database) ExplorerModel {
	list := utils.ExplorerList{}
	list.Initialize()
	return ExplorerModel{
		database:          database,
		databaseLoadError: "",
		explorerList:      list,
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
	case DatabasesLoadedError:
		m.databaseLoadError = string(msg)
	case DatabasesLoaded:
		m.databaseLoadError = ""
		m.explorerList = m.createDatabaseList([]string(msg))
	}
	return m, cmd
}

func (m ExplorerModel) View() string {
	return m.ListNode(m.explorerList.Root, 0)
}

func (m ExplorerModel) handleKeyboardActions(msg tea.Msg) (ExplorerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			m.explorerList.Expand(
				[]utils.ExplorerNode{
					utils.ExplorerNode{Title: "users", Type: "tables"},
					utils.ExplorerNode{Title: "accounts", Type: "tables"},
					utils.ExplorerNode{Title: "pages", Type: "tables"},
					utils.ExplorerNode{Title: "transactions", Type: "tables"},
					utils.ExplorerNode{Title: "categories", Type: "tables"},
					utils.ExplorerNode{Title: "tokens", Type: "tables"},
				},
			)
		case "h":
			m.explorerList.Contract()
		case "j":
			m.explorerList.MoveDown()
		case "k":
			m.explorerList.MoveUp()
		}
	}
	return m, cmd
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
