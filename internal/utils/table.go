package utils

import (
	"github.com/charmbracelet/x/ansi"
	"slices"
	"strings"

	viewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

// escapeCell replaces control characters that would break the table layout.
func escapeCell(s string) string {
	r := strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
	)
	return r.Replace(s)
}

var MIN_COLUMN_WIDTH = 100

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
	viewport := viewport.New(width-2, height-2)
	var rows [][]string
	var cols []string

	if len(data) > 1 {
		rows = data[1:]
		cols = data[0]
	} else {
		rows = [][]string{}
		cols = []string{}
	}

	table := Table{
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
		columnWidths:        calculateColumnWidths(cols, rows),
	}
	content := table.renderColumns() + "\n" + table.renderRows()
	table.Viewport.SetContent(content)
	return table
}

func (t Table) Update(msg tea.Msg) (Table, tea.Cmd) {
	var viewportCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			t.SelectedColumn = (t.SelectedColumn + 1) % len(t.Columns)
			t.Viewport.ScrollRight(1)
		case "h":
			t.SelectedColumn = (t.SelectedColumn - 1 + len(t.Columns)) % len(t.Columns)
			t.Viewport.ScrollLeft(1)
		case "j":
			t.SelectedRow = (t.SelectedRow + 1) % len(t.Rows)
		case "k":
			t.SelectedRow = (t.SelectedRow - 1 + len(t.Rows)) % len(t.Rows)
		}
	}
	content := t.renderColumns() + "\n" + t.renderRows()
	t.Viewport.SetContent(content)
	t.Viewport, viewportCmd = t.Viewport.Update(msg)
	return t, viewportCmd
}

func (t Table) View() string {
	return t.Viewport.View()
}

func (t Table) HasData() bool {
	return len(t.Columns) > 0 && len(t.Rows) > 0
}

func calculateColumnWidths(cols []string, rows [][]string) []int {
	widths := make([]int, len(cols))
	for i, col := range cols {
		widths[i] = slices.Max([]int{widths[i], len(col) + 2})
	}
	for _, row := range rows {
		for j, cell := range row {
			widths[j] = slices.Min([]int{
				MIN_COLUMN_WIDTH,
				slices.Max([]int{widths[j], len(cell) + 2}),
			})
		}
	}
	return widths
}

func (t Table) renderColumns() string {
	var columns []string
	for i, col := range t.Columns {
		style := t.ColumnsStyle.
			Width(t.columnWidths[i]).
			Padding(0, 1, 0, 1)

		if i == t.SelectedColumn {
			style = style.Inherit(t.SelectedColumnStyle)
		}
		columns = append(columns, style.Render(col))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, columns...)
}

func (t Table) renderRows() string {
	var rows []string
	for i, row := range t.Rows {
		var columns []string
		for j, cell := range row {
			style := lipgloss.NewStyle().Width(t.columnWidths[j]).Padding(0, 1, 0, 1)
			if i == t.SelectedRow && j == t.SelectedColumn {
				style = style.Inherit(t.SelectedCellStyle)
			} else if i == t.SelectedRow {
				style = style.Inherit(t.SelectedRowStyle)
			} else if j == t.SelectedColumn {
				style = style.Inherit(t.SelectedColumnStyle)
			}
			columns = append(columns, style.Render(ansi.Truncate(escapeCell(cell), t.columnWidths[j]-2, "…")))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, columns...))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
