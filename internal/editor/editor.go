package editor

import (
	tea "github.com/charmbracelet/bubbletea"
)

type EditorModel struct {
}

func InitEditor() EditorModel {
	return EditorModel{}
}

func (m EditorModel) Init() tea.Cmd {
	return nil
}

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m EditorModel) View() string {
	return "Welcome to LazyGit Editor!\n\nPress q to quit."
}
