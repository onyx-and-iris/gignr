package tui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor    = lipgloss.Color("#7D56F4")
	backgroundColor = lipgloss.Color("#1A1B26")
	textColor       = lipgloss.Color("#C0CAF5")
	mutedTextColor  = lipgloss.Color("#565F89")
)

type AppStyle struct {
	width  int
	height int

	tabSection   lipgloss.Style
	inactiveTab  lipgloss.Style
	activeTab    lipgloss.Style
	divider      lipgloss.Style
	templateList lipgloss.Style
	templateName lipgloss.Style
	selectedItem lipgloss.Style
	noTemplates  lipgloss.Style
	checkbox     lipgloss.Style
	pointer      lipgloss.Style
	hotkeys      lipgloss.Style
	progress     lipgloss.Style
	searchBox    lipgloss.Style
}

func NewAppStyle(width, height int) *AppStyle {
	s := &AppStyle{
		width:  width,
		height: height,
	}
	s.refresh()
	return s
}

func (s *AppStyle) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.refresh()
}

func (s *AppStyle) refresh() {
	contentWidth := s.width - 4

	templateListHeight := s.height - 15

	s.tabSection = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Width(contentWidth)

	s.inactiveTab = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(textColor)

	s.activeTab = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(primaryColor).
		Foreground(backgroundColor)

	s.divider = lipgloss.NewStyle().
		SetString(" • ").
		Bold(true).
		Foreground(mutedTextColor)

	s.templateList = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(contentWidth).
		Height(templateListHeight)

	s.templateName = lipgloss.NewStyle().
		Foreground(textColor)

	s.selectedItem = lipgloss.NewStyle().
		Foreground(backgroundColor).
		Background(primaryColor).
		Bold(true).
		Padding(0, 1)

	s.noTemplates = lipgloss.NewStyle().
		Foreground(mutedTextColor).
		Italic(true).
		Align(lipgloss.Center).
		Padding(2).
		Width(contentWidth)

	s.checkbox = lipgloss.NewStyle().
		Foreground(mutedTextColor)

	s.pointer = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)

	s.hotkeys = lipgloss.NewStyle().
		Foreground(mutedTextColor).
		Italic(true).
		Padding(0, 2).
		Margin(1, 0)

	s.progress = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).Align(lipgloss.Right).PaddingRight(2)

	s.searchBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1).
		Width(contentWidth)
}
