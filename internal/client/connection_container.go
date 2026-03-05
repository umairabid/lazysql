package client

import (
	"golang.org/x/term"
	"os"

	editor "app.lazygit/internal/editor"
	explorer "app.lazygit/internal/explorer"
	viewer "app.lazygit/internal/viewer"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var MIN_WIDTH = 600
var MIN_HEIGHT = 400

type ConnectionContainerModel struct {
	explorer    explorer.ExplorerModel
	editor      editor.EditorModel
	viewer      viewer.ViewerModel
	width       int
	height      int
	active_view string
}

func InitConnectionContainer() ConnectionContainerModel {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = MIN_WIDTH
		height = MIN_HEIGHT
	}

	return ConnectionContainerModel{
		explorer: explorer.InitExplorer(),
		editor:   editor.InitEditor(),
		viewer:   viewer.InitViewer(),
		width:    width,
		height:   height,
	}
}

func (m ConnectionContainerModel) Init() tea.Cmd {
	return tea.Batch(
		m.explorer.Init(),
		m.editor.Init(),
		m.viewer.Init(),
	)
}

func (m ConnectionContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConnectionContainerModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.explorerView(),
		lipgloss.JoinVertical(lipgloss.Left,
			m.editorView(),
			m.viewerView(),
		),
	)
}

func (m ConnectionContainerModel) explorerView() string {
	explorerWidth := (m.width / 3) - 8
	explorerHeight := m.height - 4
	return lipgloss.
		NewStyle().
		Width(explorerWidth).
		Height(explorerHeight).
		Border(lipgloss.RoundedBorder()).
		Margin(1, 0, 0, 1).
		Render(m.explorer.View())
}

func (m ConnectionContainerModel) editorView() string {
	editorWidth := m.width - (m.width / 3) + 3
	editorHeight := (m.height / 2) - 8
	return lipgloss.
		NewStyle().
		Width(editorWidth).
		Height(editorHeight).
		Border(lipgloss.RoundedBorder()).
		Margin(1, 0, 0, 0).
		Render(m.editor.View())
}

func (m ConnectionContainerModel) viewerView() string {
	viewerWidth := m.width - (m.width / 3) + 3
	viewerHeight := (m.height / 2) + 3
	return lipgloss.
		NewStyle().
		Width(viewerWidth).
		Height(viewerHeight).
		Border(lipgloss.RoundedBorder()).
		Margin(0, 0, 1, 0).
		Render(m.viewer.View())
}
