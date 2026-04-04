package client

import (
	"golang.org/x/term"
	"os"

	adapters "app.lazygit/internal/adapters"
	editor "app.lazygit/internal/editor"
	explorer "app.lazygit/internal/explorer"
	utils "app.lazygit/internal/utils"
	viewer "app.lazygit/internal/viewer"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var MIN_WIDTH = 600
var MIN_HEIGHT = 400

type ConnectionContainerModel struct {
	explorer    tea.Model
	editor      tea.Model
	viewer      tea.Model
	active_view string
	layout      utils.ConnectionContainerLayout
}
type LayoutUpdated utils.ConnectionContainerLayout

func InitConnectionContainer(database adapters.Database) ConnectionContainerModel {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = MIN_WIDTH
		height = MIN_HEIGHT
	}

	layout := utils.CalculateConnectionContainerLayout(width, height)
	return ConnectionContainerModel{
		explorer:    explorer.InitExplorer(database, layout),
		editor:      editor.InitEditor(database, layout),
		viewer:      viewer.InitViewer(database, layout),
		active_view: "explorer",
		layout:      layout,
	}
}

func setLayout(width int, height int) tea.Cmd {
	return func() tea.Msg {
		return LayoutUpdated(utils.CalculateConnectionContainerLayout(width, height))
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
	var explorerCmd, editorCmd, viewerCmd, command tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyboardMsg(msg)
	case tea.WindowSizeMsg:
		command = setLayout(msg.Width, msg.Height)
	}

	m.explorer, explorerCmd = m.explorer.Update(msg)
	m.editor, editorCmd = m.editor.Update(msg)
	m.viewer, viewerCmd = m.viewer.Update(msg)

	return m, tea.Batch(command, explorerCmd, editorCmd, viewerCmd)
}

func (m ConnectionContainerModel) handleKeyboardMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "shift+tab":
		if m.active_view == "explorer" {
			m.active_view = "editor"
		} else if m.active_view == "editor" {
			m.active_view = "viewer"
		} else {
			m.active_view = "explorer"
		}
	}
	if m.active_view == "explorer" {
		m.explorer, cmd = m.explorer.Update(msg)
	} else if m.active_view == "editor" {
		m.editor, cmd = m.editor.Update(msg)
	} else {
		m.viewer, cmd = m.viewer.Update(msg)
	}
	return m, cmd
}

func (m ConnectionContainerModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top + 2,
		m.explorerView(),
		lipgloss.JoinVertical(lipgloss.Left,
			m.editorView(),
			m.viewerView(),
		),
	)
}

func (m ConnectionContainerModel) explorerView() string {
	style := lipgloss.
		NewStyle().
		Width(m.layout.ExplorerWidth - 6).
		Height(m.layout.ExplorerHeight - 6).
		Border(lipgloss.RoundedBorder()).
		Margin(1, 0, 0, 1)

	if m.active_view == "explorer" {
		style = style.BorderForeground(lipgloss.Color("205"))
	}

	return style.Render(m.explorer.View())
}

func (m ConnectionContainerModel) editorView() string {
	style := lipgloss.
		NewStyle().
		Width(m.layout.EditorWidth).
		Height(m.layout.EditorHeight - 4).
		Border(lipgloss.RoundedBorder()).
		Margin(1, 0, 0, 0)

	if m.active_view == "editor" {
		style = style.BorderForeground(lipgloss.Color("205"))
	}
	return style.Render(m.editor.View())
}

func (m ConnectionContainerModel) viewerView() string {
	style := lipgloss.
		NewStyle().
		Width(m.layout.ViewerWidth - 4).
		Height(m.layout.ViewerHeight).
		Border(lipgloss.RoundedBorder()).
		Margin(0, 0, 1, 0)

	if m.active_view == "viewer" {
		style = style.BorderForeground(lipgloss.Color("205"))
	}

	return style.Render(m.viewer.View())
}
