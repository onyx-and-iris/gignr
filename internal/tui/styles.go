package tui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor    = lipgloss.Color("#7D56F4")
	backgroundColor = lipgloss.Color("#1A1B26")
	textColor       = lipgloss.Color("#C0CAF5")
	mutedTextColor  = lipgloss.Color("#565F89")

	title = lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Foreground(primaryColor).
		Bold(true)

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

	templateListStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1).
				Width(80)

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

	progressStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginTop(1)

	searchBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)
)
