package conn_manager

import (
	"fmt"

	"app.lazygit/internal/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConnectionManager struct {
	width      int
	height     int
	inputs     []textinput.Model
	focusIndex int
}

func InitConnManager() ConnectionManager {
	hostInput := createHostInput()
	hostInput.Focus()

	inputs := []textinput.Model{hostInput, createPortInput(), createUserInput(), createPasswordInput()}
	return ConnectionManager{inputs: inputs[:], width: 80, height: 24, focusIndex: 0}
}

func (m ConnectionManager) Init() tea.Cmd {
	return nil
}

func (m ConnectionManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			s := msg.String()
			m.focusIndex = m.changeFocusIndex(s)
			cmds := m.changeFocusedInput()
			return m, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m ConnectionManager) View() string {
	view := fmt.Sprintf("Connection Manager\n\n%s\n\nPress q to quit.", m.renderInputs())

	container := utils.BorderedContainer().Width(m.width / 2).Height(m.height / 2).Render(view)
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		container,
	)
}

func (m ConnectionManager) changeFocusIndex(key string) int {
	if key == "tab" {
		return (m.focusIndex + 1) % 4
	} else if key == "shift+tab" {
		return (m.focusIndex - 1 + 4) % 4
	}
	return m.focusIndex
}

func (m ConnectionManager) changeFocusedInput() []tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.inputs {
		if i == m.focusIndex {
			cmds = append(cmds, m.inputs[i].Focus())
		} else {
			m.inputs[i].Blur()
		}
	}
	return cmds
}

func (m ConnectionManager) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m ConnectionManager) renderInputs() string {
	var result string
	for i, input := range m.inputs {
		if i == m.focusIndex {
			result += utils.FocusedTextInputStyle().Render(input.View())
		} else {
			result += utils.TextInputStyle().Render(input.View())
		}
		result += "\n"
	}
	return result
}
