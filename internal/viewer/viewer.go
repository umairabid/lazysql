package viewer

import (
	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	paddingStyle    = lipgloss.NewStyle().Padding(1)
	selectedStyle   = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("230"))
	borderColor     = lipgloss.Color("240")
	borderStyle     = lipgloss.NewStyle().Foreground(borderColor)
)

type ViewerModel struct {
	database    adapters.Database
	headers     []string
	rows        [][]string
	selectedRow int
}

func InitViewer(database adapters.Database) ViewerModel {
	return ViewerModel{database: database}
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
	for i := range colWidths {
		colWidths[i] += 2 // 1 space padding each side
	}

	sep := func(left, mid, join, right string) string {
		var s strings.Builder
		s.WriteString(borderStyle.Render(left))
		for i, w := range colWidths {
			s.WriteString(borderStyle.Render(strings.Repeat("─", w)))
			if i < len(colWidths)-1 {
				s.WriteString(borderStyle.Render(join))
			}
		}
		s.WriteString(borderStyle.Render(right))
		return s.String()
	}

	renderRow := func(cells []string, selected bool) string {
		var s strings.Builder
		s.WriteString(borderStyle.Render("│"))
		for i, w := range colWidths {
			cell := ""
			if i < len(cells) {
				cell = cells[i]
			}
			content := " " + cell + strings.Repeat(" ", w-len(cell)-1)
			if selected {
				s.WriteString(selectedStyle.Render(content))
			} else {
				s.WriteString(content)
			}
			s.WriteString(borderStyle.Render("│"))
		}
		return s.String()
	}

	var sb strings.Builder
	sb.WriteString(sep("┌", "─", "┬", "┐") + "\n")
	sb.WriteString(renderRow(headers, false) + "\n")
	sb.WriteString(sep("├", "─", "┼", "┤") + "\n")
	for i, row := range rows {
		sb.WriteString(renderRow(row, i == selectedRow) + "\n")
		if i < len(rows)-1 {
			sb.WriteString(sep("├", "─", "┼", "┤") + "\n")
		}
	}
	sb.WriteString(sep("└", "─", "┴", "┘"))

	return sb.String()
}

func (m ViewerModel) View() string {
	if len(m.headers) > 0 {
		return paddingStyle.Render(renderTable(m.headers, m.rows, m.selectedRow))
	}
	return "Welcome to LazyGit Viewer!\n\nPress q to quit."
}
