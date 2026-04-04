package editor

import (
	adapters "app.lazygit/internal/adapters"
	utils "app.lazygit/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type EditorModel struct {
	database adapters.Database
	layout   utils.ConnectionContainerLayout
	isActive bool
}

func InitEditor(database adapters.Database, layout utils.ConnectionContainerLayout) EditorModel {
	return EditorModel{
		database: database,
		layout:   layout,
		isActive: false,
	}
}

func (m EditorModel) Init() tea.Cmd {
	return nil
}

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case utils.ActiveViewChanged:
		m.isActive = string(msg) == "editor"
	case utils.LayoutUpdated:
		m.layout = utils.ConnectionContainerLayout(msg)
	}
	return m, cmd
}

func (m EditorModel) View() string {
	style := lipgloss.
		NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.layout.EditorWidth - 2).
		Height(m.layout.EditorHeight - 2)

	if m.isActive {
		style = style.BorderForeground(lipgloss.Color("205"))
	}
	return style.Render("Welcome to LazyGit Editor!\n\nPress q to quit.")
}
