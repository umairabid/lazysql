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

type LayoutUpdated utils.ConnectionContainerLayout
type ConnectionContainerModel struct {
	explorer    tea.Model
	editor      tea.Model
	viewer      tea.Model
	active_view string
	layout      utils.ConnectionContainerLayout
}

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

func (m ConnectionContainerModel) changeActiveView() tea.Cmd {
	return func() tea.Msg {
		var newActiveView string
		switch m.active_view {
		case "explorer":
			newActiveView = "editor"
		case "editor":
			newActiveView = "viewer"
		default:
			newActiveView = "explorer"
		}
		return utils.ActiveViewChanged(newActiveView)
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
	case LayoutUpdated:
		m.layout = utils.ConnectionContainerLayout(msg)
	case utils.ActiveViewChanged:
		m.active_view = string(msg)
	}

	m.explorer, explorerCmd = m.explorer.Update(msg)
	m.editor, editorCmd = m.editor.Update(msg)
	m.viewer, viewerCmd = m.viewer.Update(msg)

	return m, tea.Batch(command, explorerCmd, editorCmd, viewerCmd)
}

func (m ConnectionContainerModel) handleKeyboardMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd, activeViewCmd tea.Cmd
	switch msg.String() {
	case "shift+tab":
		cmd = m.changeActiveView()
	}
	if m.active_view == "explorer" {
		m.explorer, activeViewCmd = m.explorer.Update(msg)
	} else if m.active_view == "editor" {
		m.editor, activeViewCmd = m.editor.Update(msg)
	} else {
		m.viewer, activeViewCmd = m.viewer.Update(msg)
	}
	return m, tea.Batch(cmd, activeViewCmd)
}

func (m ConnectionContainerModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.explorer.View(),
		lipgloss.JoinVertical(lipgloss.Left,
			m.editor.View(),
			m.viewer.View(),
		),
	)
}
