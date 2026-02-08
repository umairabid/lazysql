package utils

import "github.com/charmbracelet/lipgloss"

func BorderedContainer() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Foreground(lipgloss.Color("5"))
}

func FocusedTextInputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(30).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)
}

func TextInputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(30).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)
}
