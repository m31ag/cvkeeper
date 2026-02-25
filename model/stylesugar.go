package model

import "github.com/charmbracelet/lipgloss"

func Bold() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}
func Faint() lipgloss.Style {
	return lipgloss.NewStyle().Faint(true)
}

func RoundedBorder() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
}
