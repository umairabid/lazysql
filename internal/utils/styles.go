package utils

import "github.com/charmbracelet/lipgloss"

func CenteredContainer(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Margin(0, 0).
		Padding(1, 2).
		Align(lipgloss.Center, lipgloss.Center)
	}
