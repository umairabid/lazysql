package utils

import "github.com/charmbracelet/lipgloss"

func BorderedContainer() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Foreground(lipgloss.Color("5"))
}
