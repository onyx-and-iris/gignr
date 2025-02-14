package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type TabModel struct {
	currentTab int
	tabs       []string
	styles     *AppStyle
}

func newTabModel(styles *AppStyle) *TabModel {
	tabs := []string{"TopTal", "GitHub", "GitHub Community", "GitHub Global", "Others"}
	return &TabModel{
		currentTab: 0,
		tabs:       tabs,
		styles:     styles,
	}
}

func (m *TabModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m *TabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.currentTab > 0 {
				m.currentTab--
			} else {
				m.currentTab = len(m.tabs) - 1
			}
		case "right", "tab":
			if m.currentTab < len(m.tabs)-1 {
				m.currentTab++
			} else {
				m.currentTab = 0
			}
		}
	}

	return m, nil
}

func (m *TabModel) View() string {
	var tabContent strings.Builder

	for i, tab := range m.tabs {
		if i == m.currentTab {
			tabContent.WriteString(m.styles.activeTab.Render(tab))
		} else {
			tabContent.WriteString(m.styles.inactiveTab.Render(tab))
		}
		if i < len(m.tabs)-1 {
			tabContent.WriteString(m.styles.divider.Render())
		}
	}

	return m.styles.tabSection.Render(tabContent.String())
}

func (m *TabModel) SetStyles(styles *AppStyle) {
	m.styles = styles
}
