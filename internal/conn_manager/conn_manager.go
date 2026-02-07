package conn_manager

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"app.lazygit/internal/utils"
)

type ConnectionManager struct {
	inputs []textinput.Model
}

func InitConnManager() ConnectionManager {
	hostInput := createHostInput()
	hostInput.Focus()

	inputs := []textinput.Model{hostInput, createPortInput(), createUserInput(), createPasswordInput()}
	return ConnectionManager{inputs: inputs[:]}
}

func (m ConnectionManager) Init() tea.Cmd {
	return nil
}

func (m ConnectionManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConnectionManager) View() string {
	view := fmt.Sprintf("Connection Manager\n\n%s\n\nPress q to quit.", renderInputs(m.inputs))

	return utils.CenteredContainer(50, 20).Render(view)
}

func renderInputs(inputs []textinput.Model) string {
	var result string
	for _, input := range inputs {
		result += input.View() + "\n"
	}
	return result
}
