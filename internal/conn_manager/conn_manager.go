package conn_manager

import (
	"fmt"

	"app.lazygit/internal/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConnectionManager struct {
	width  int
	height int
	inputs []textinput.Model
}

func InitConnManager() ConnectionManager {
	hostInput := createHostInput()
	hostInput.Focus()

	inputs := []textinput.Model{hostInput, createPortInput(), createUserInput(), createPasswordInput()}
	return ConnectionManager{inputs: inputs[:], width: 80, height: 24}
}

func (m ConnectionManager) Init() tea.Cmd {
	return nil
}

func (m ConnectionManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m ConnectionManager) View() string {
	view := fmt.Sprintf("Connection Manager\n\n%s\n\nPress q to quit.", renderInputs(m.inputs))

	container := utils.BorderedContainer().Width(m.width / 2).Height(m.height / 2).Render(view)
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		container,
	)
}

func renderInputs(inputs []textinput.Model) string {
	var result string
	for _, input := range inputs {
		result += input.View() + "\n"
	}
	return result
}
