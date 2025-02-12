package tui

import "github.com/charmbracelet/lipgloss"

var (
	title = lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#7D56F4")).Foreground(lipgloss.Color("#7D56F4"))
)

// TabModel
var (
	highlightColor   = lipgloss.Color("#7D56F4")
	tabSection       = lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#7D56F4"))
	inactiveTabStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	activeTabStyle   = lipgloss.NewStyle().Inherit(inactiveTabStyle).
				Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	dividerStyle = lipgloss.NewStyle().SetString(" • ").Bold(true)
)

// TemplatesListModel
var (
	templateListBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1, 2)

	templateSelectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("32")). // Green background for selected templates
				Foreground(lipgloss.Color("0")).
				Bold(true)

	highlightedPointerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")). // Pink for the pointer
				Bold(true)
)

var (
	sourceHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				MarginBottom(1)

	templateItemStyle = lipgloss.NewStyle().
				Padding(0, 2)

	templateNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF"))

	noTemplatesStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666")).
				Italic(true).
				Padding(1, 2)

	checkmarkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	pointerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)
)
