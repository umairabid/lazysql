package viewer

import (
	"fmt"

	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type ViewerModel struct {
	database    adapters.Database
	selectedRow int
	table	 utils.Table
	layout      utils.ConnectionContainerLayout
	isActive    bool
}

func InitViewer(database adapters.Database, layout utils.ConnectionContainerLayout) ViewerModel {
	return ViewerModel{
		database: database,
		layout:   layout,
		isActive: false,
		table:    utils.InitTable([][]string{}, layout.ViewerWidth, layout.ViewerHeight),
	}
}

func createTableFromData(data [][]string, layout utils.ConnectionContainerLayout) utils.Table {
	return utils.InitTable(data, layout.ViewerWidth, layout.ViewerHeight)
}

func (m ViewerModel) Init() tea.Cmd { return nil }

func (m ViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var viewPortCmd tea.Cmd
	switch msg := msg.(type) {
	case utils.ViewerTableData:
		m.table = createTableFromData(msg, m.layout)
	case utils.ActiveViewChanged:
		m.isActive = string(msg) == "viewer"
	case utils.LayoutUpdated:
		m.layout = utils.ConnectionContainerLayout(msg)
	}
	return m, viewPortCmd
}

func (m ViewerModel) View() string {
	style := lipgloss.
		NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.layout.ViewerWidth - 2).
		Height(m.layout.ViewerHeight - 2)

	if m.isActive {
		style = style.BorderForeground(lipgloss.Color("205"))
	}
	var content string
	if m.table.HasData() {
		content = fmt.Sprintf("%s\n", m.table.View())
	} else {
		content = "Welcome to LazyGit Viewer!\n\nPress q to quit."
	}
	return style.Render(content)
}
