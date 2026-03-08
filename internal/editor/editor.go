package editor

import (
	tea "github.com/charmbracelet/bubbletea"
	adapters "app.lazygit/internal/adapters"
)

type EditorModel struct {
	database adapters.Database
}

func InitEditor(database adapters.Database) EditorModel {
	return EditorModel{
		database: database,
	}
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
