package utils

import (
	"github.com/charmbracelet/x/ansi"
	viewport "github.com/charmbracelet/bubbles/viewport"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type Table struct {
	Columns             []string
	Rows                [][]string
	Width               int
	Height              int
	SelectedRow         int
	SelectedColumn      int
	ColumnsStyle        lipgloss.Style
	SelectedRowStyle    lipgloss.Style
	SelectedColumnStyle lipgloss.Style
	SelectedCellStyle   lipgloss.Style
	Viewport            viewport.Model

	columnWidths []int
}

func InitTable(data [][]string, width int, height int) Table {
	viewport := viewport.New(width, height)
	var rows [][]string
	var cols []string
	if len(data) > 1 {
		rows = data[1:]
		cols = data[0]
	} else {
		rows = [][]string{}
		cols = []string{}
	}

	

	return Table{
		Columns:             cols,
		Rows:                rows,
		Width:               width,
		Height:              height,
		SelectedRow:         0,
		SelectedColumn:      0,
		ColumnsStyle:        lipgloss.NewStyle().Bold(true),
		SelectedRowStyle:    lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("229")),
		SelectedColumnStyle: lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("229")),
		SelectedCellStyle:   lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("229")),
		Viewport:            viewport,
		columnWidths:        calculateColumnWidths(cols),
	}
}

func calculateColumnWidths(columns []string) []int {
	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = len(col) + 2
	}
	return widths
}

func (t Table) renderColumns() string {
	var result string
	for i, col := range t.Columns {
		result += t.ColumnsStyle.Padding(0, 1, 0, 1).Width(t.columnWidths[i]).Render(col)
	}
	return result
}

func (t Table) renderRows() string {
	var result string
	for _, row := range t.Rows {
		for j, cell := range row {
			style := lipgloss.NewStyle()
			style = style.Width(t.columnWidths[j]).Padding(0, 1, 0, 1)
			result += style.Render(ansi.Truncate(cell, t.columnWidths[j] - 2, "…"))
		}
		result += "\n"
	}
	return result
}

func (t Table) View() string {
	content := t.renderColumns() + "\n" + t.renderRows()
	t.Viewport.SetContent(content)
	return t.Viewport.View()
}

func (t Table) HasData() bool {
	return len(t.Columns) > 0 && len(t.Rows) > 0
}
