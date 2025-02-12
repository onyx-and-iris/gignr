package tui

import "github.com/charmbracelet/lipgloss"

// TabModel
var (
	highlightColor   = lipgloss.Color("#7D56F4")
	tabSection       = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#7D56F4"))
	inactiveTabStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	activeTabStyle   = lipgloss.NewStyle().Inherit(inactiveTabStyle).
				Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	dividerStyle = lipgloss.NewStyle().SetString(" • ").Bold(true)
)
