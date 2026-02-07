package conn_manager

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func createTextInput() textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ti.CharLimit = 256
	return ti
}

func createHostInput() textinput.Model {
	ti := createTextInput()
	ti.Placeholder = "Host"
	return ti
}

func createPortInput() textinput.Model {
	ti := createTextInput()
	ti.Placeholder = "Port"
	return ti
}

func createUserInput() textinput.Model {
	ti := createTextInput()
	ti.Placeholder = "User"
	return ti
}

func createPasswordInput() textinput.Model {
	ti := createTextInput()
	ti.Placeholder = "Password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	return ti
}
