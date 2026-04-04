package viewer

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"

	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type ViewerModel struct {
	database    adapters.Database
	headers     []string
	rows        [][]string
	selectedRow int
	viewport    viewport.Model
	layout      utils.ConnectionContainerLayout
	isActive    bool
}

func InitViewer(database adapters.Database, layout utils.ConnectionContainerLayout) ViewerModel {
	return ViewerModel{
		database: database,
		layout:   layout,
		isActive: false,
	}
}

func (m ViewerModel) Init() tea.Cmd { return nil }

func (m ViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case utils.ViewerTableData:
		if len(msg) > 0 {
			m.headers = msg[0]
			m.rows = msg[1:]
			m.selectedRow = 0
		}
	case utils.ActiveViewChanged:
		m.isActive = string(msg) == "viewer"
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.selectedRow < len(m.rows)-1 {
				m.selectedRow++
			}
		case "k", "up":
			if m.selectedRow > 0 {
				m.selectedRow--
			}
		}
	}
	return m, nil
}

func renderTable(headers []string, rows [][]string, selectedRow int) string {
	// Calculate column widths: max of header length and widest cell per column
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	cols := make([]table.Column, len(headers))
	for i, h := range headers {
		cols[i] = table.Column{Title: h, Width: colWidths[i]}
	}

	tableRows := make([]table.Row, len(rows))
	for i, r := range rows {
		tableRows[i] = table.Row(r)
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(tableRows),
		table.WithFocused(true),
	)
	t.SetCursor(selectedRow)

	return t.View()
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
	if len(m.headers) > 0 {
		content = renderTable(m.headers, m.rows, m.selectedRow)
	} else {
		content = "Welcome to LazyGit Viewer!\n\nPress q to quit."
	}
	return style.Render(content)
}
