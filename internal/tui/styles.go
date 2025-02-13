package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Base colors
	primaryColor    = lipgloss.Color("#7D56F4")
	secondaryColor  = lipgloss.Color("#5A3DBF")
	backgroundColor = lipgloss.Color("#1A1B26")
	textColor       = lipgloss.Color("#C0CAF5")
	mutedTextColor  = lipgloss.Color("#565F89")
	successColor    = lipgloss.Color("#9ECE6A")
	highlightColor  = lipgloss.Color("#BB9AF7")

	// Common styles
	baseStyle = lipgloss.NewStyle().
			Background(backgroundColor).
			Foreground(textColor)

	// Title styles
	title = lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Foreground(primaryColor).
		Bold(true)

	// Tab styles
	tabSection = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	inactiveTabStyle = lipgloss.NewStyle().
				Bold(true).
				Padding(0, 1).
				Foreground(textColor)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(primaryColor).
			Foreground(backgroundColor)

	dividerStyle = lipgloss.NewStyle().
			SetString(" • ").
			Bold(true).
			Foreground(mutedTextColor)

	// Template list styles
	templateListStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1).
				Width(80)

	sourceHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				Padding(1).
				MarginBottom(1).
				Border(lipgloss.Border{
			Bottom: "─",
		}).
		BorderForeground(mutedTextColor)

	templateNameStyle = lipgloss.NewStyle().
				Foreground(textColor)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(backgroundColor).
				Background(primaryColor).
				Bold(true).
				Padding(0, 1)

	noTemplatesStyle = lipgloss.NewStyle().
				Foreground(mutedTextColor).
				Italic(true).
				Align(lipgloss.Center).
				Padding(2)

	checkboxStyle = lipgloss.NewStyle().
			Foreground(mutedTextColor)

	pointerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	hotkeysStyle = lipgloss.NewStyle().
			Foreground(mutedTextColor).
			Italic(true).
			Padding(0, 2).
			MarginTop(1)

	// Additional template list styles
	templateItemStyle = lipgloss.NewStyle().
				Padding(0, 1)

	checkmarkStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// Status styles
	statusStyle = lipgloss.NewStyle().
			Foreground(backgroundColor).
			Background(primaryColor).
			Bold(true).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7768E")).
			Bold(true).
			Padding(0, 1)

	// Progress styles
	progressStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginTop(1)
)
