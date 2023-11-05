package ui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	BaseBorder lipgloss.Style
	Title      lipgloss.Style
	Error      lipgloss.Style
	Modal      lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		BaseBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true).
			Padding(0, 1),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Padding(0, 1),
		Modal: lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Width(50).Padding(1, 2),
	}
}
