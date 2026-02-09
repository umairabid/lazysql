package utils

import "github.com/charmbracelet/lipgloss"

func Border() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())
}

func BottomBorder() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false)
}

func TopBorder() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false)
}


func RightBorder() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false)
}

func FocusedTextInputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(30).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)
}

func TextInputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(30).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)
}
