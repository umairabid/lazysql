package conn_manager

import (
	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type ConnectionForm struct {
	inputs     []textinput.Model
	mode       string
	focusIndex int
	layout     utils.ConnectionManagerLayout
}

func InitConnForm(connection adapters.DbConnection, layout utils.ConnectionManagerLayout) ConnectionForm {
	inputs := []textinput.Model{
		createDriverInput(connection.Driver),
		createNameInput(connection.Name),
		createHostInput(connection.Host),
		createPortInput(connection.Port),
		createUserInput(connection.Username),
		createPasswordInput(connection.Password),
		createUrlInput(connection.Url),
		createCommandInput(connection.Command),
	}
	return ConnectionForm{
		inputs:     inputs[:],
		focusIndex: -1,
		layout:     layout,
		mode:       "url",
	}
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
		conn := adapters.DbConnection(msg)
		m.inputs[0].SetValue(conn.Host)
		m.inputs[1].SetValue(conn.Port)
		m.inputs[2].SetValue(conn.Username)
		m.inputs[3].SetValue(conn.Password)
	case EditConnectionMsg:
		canEdit := bool(msg)
		if canEdit {
			m.focusIndex = 0
		} else {
			m.focusIndex = -1
		}
		cmds := m.changeFocusedInput()
		return m, tea.Batch(cmds...)
	case LayoutUpdated:
		m.layout = utils.ConnectionManagerLayout(msg)
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m ConnectionForm) View() string {
	result := m.renderConnectionInfoFields()
	if m.mode == "credentials" {
		result = append(result, m.renderCredentialsFields()...)
	}
	if m.mode == "url" {
		result = append(result, m.renderFieldsForIndexes([]int{6})...)
	}
	if m.mode == "command" {
		result = append(result, m.renderFieldsForIndexes([]int{7})...)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		Width(m.layout.ConnectionFormWidth).
		Height(m.layout.BodyHeight).Render(lipgloss.JoinVertical(lipgloss.Left, result...))
}

func (m ConnectionForm) renderConnectionInfoFields() []string {
	indexes := []int{0, 1}
	return m.renderFieldsForIndexes(indexes)
}

func (m ConnectionForm) renderCredentialsFields() []string {
	indexes := []int{2, 3, 4, 5}
	return m.renderFieldsForIndexes(indexes)
}

func (m ConnectionForm) renderFieldsForIndexes(indexes []int) []string {
	var result []string
	for _, index := range indexes {
		var inputView string
		if index == m.focusIndex {
			inputView = utils.FocusedTextInputStyle().Render(m.inputs[index].View())
		} else {
			inputView = utils.TextInputStyle().Render(m.inputs[index].View())
		}
		result = append(result, inputView)
	}
	return result
}

func (m ConnectionForm) changeFocusIndex(key string) int {
	if key == "tab" {
		return (m.focusIndex + 1) % len(m.inputs)
	} else if key == "shift+tab" {
		return (m.focusIndex - 1 + len(m.inputs)) % len(m.inputs)
	}
	return m.focusIndex
}
