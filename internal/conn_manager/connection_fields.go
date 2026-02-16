package conn_manager

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func createTextInput(value string) textinput.Model {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.SetValue(value)
	return ti
}

func createHostInput(value string) textinput.Model {
	ti := createTextInput(value)
	ti.Placeholder = "Host"
	return ti
}

func createPortInput(value string) textinput.Model {
	ti := createTextInput(value)
	ti.Placeholder = "Port"
	return ti
}

func createUserInput(value string) textinput.Model {
	ti := createTextInput(value)
	ti.Placeholder = "User"
	return ti
}

func createPasswordInput(value string) textinput.Model {
	ti := createTextInput(value)
	ti.Placeholder = "Password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	return ti
}
