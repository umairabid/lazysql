package conn_manager

import (
	"app.lazygit/internal/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionForm struct {
	inputs     []textinput.Model
	focusIndex int
}

func InitConnForm(connection Connection) ConnectionForm {
	inputs := []textinput.Model{
		createHostInput(connection.host),
		createPortInput(connection.port),
		createUserInput(connection.username),
		createPasswordInput(connection.password),
	}
	return ConnectionForm{inputs: inputs[:], focusIndex: -1}
}

func (m ConnectionForm) Init() tea.Cmd {
	return nil
}

func (m ConnectionForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			if m.focusIndex != -1 {

				s := msg.String()
				m.focusIndex = m.changeFocusIndex(s)
				cmds := m.changeFocusedInput()
				return m, tea.Batch(cmds...)
			}
		}
	case SelectedConnectionMsg:
		conn := Connection(msg)
		m.inputs[0].SetValue(conn.host)
		m.inputs[1].SetValue(conn.port)
		m.inputs[2].SetValue(conn.username)
		m.inputs[3].SetValue(conn.password)
	case EditConnectionMsg:
		canEdit := bool(msg)
		if canEdit {
			m.focusIndex = 0
		} else {
			m.focusIndex = -1
		}
		cmds := m.changeFocusedInput()
		return m, tea.Batch(cmds...)
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m ConnectionForm) View() string {
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

func (m ConnectionForm) changeFocusIndex(key string) int {
	if key == "tab" {
		return (m.focusIndex + 1) % 4
	} else if key == "shift+tab" {
		return (m.focusIndex - 1 + 4) % 4
	}
	return m.focusIndex
}

func (m ConnectionForm) changeFocusedInput() []tea.Cmd {
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

func (m ConnectionForm) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
