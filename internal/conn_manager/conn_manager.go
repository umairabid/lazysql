package conn_manager

import (
	"fmt"
	"os"
	"slices"

	"app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type ConnectionManager struct {
	width  int
	height int
	list   tea.Model
	form   tea.Model
}

func InitConnectionManager() ConnectionManager {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = 80
		height = 24
	}
	return ConnectionManager{
		width:  width,
		height: height,
		list:   InitConnectionList(),
		form:   InitConnForm(),
	}
}

var MIN_WIDTH = 80
var MIN_HEIGHT = 24

func (m ConnectionManager) Init() tea.Cmd {
	return tea.Batch(m.list.Init(), m.form.Init())
}

func (m ConnectionManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var listCmd, formCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.list, listCmd = m.list.Update(msg)
	m.form, formCmd = m.form.Update(msg)
	cmd := tea.Batch(listCmd, formCmd)
	return m, cmd
}

func (m ConnectionManager) View() string {
	widths := []int{MIN_WIDTH, m.width / 3}
	heights := []int{MIN_HEIGHT, m.height / 3}
	width := slices.Max(widths)
	height := slices.Max(heights)
	header := utils.BottomBorder().Width(width).Padding(1).Render("Connection Manager")
	footer := utils.TopBorder().Width(width).Padding(1).Render("Press 'q' to quit")
	bodyHeight := height - (lipgloss.Height(header) + lipgloss.Height(footer))
	listView := utils.RightBorder().Width(width / 2).Height(bodyHeight).Render(m.list.View())
	formView := lipgloss.NewStyle().Width(width/2).Height(bodyHeight).Padding(1, 2).Render(m.form.View())
	listAndFormView := lipgloss.JoinHorizontal(lipgloss.Top, listView, formView)
	body := lipgloss.NewStyle().Height(bodyHeight - 3).Render(listAndFormView)
	container := utils.Border().Width(width).Height(height).Render(
		fmt.Sprintf("%s\n%s\n%s", header, body, footer),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, container)
}
