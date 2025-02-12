package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TabModel struct {
	currentTab int
	tabs       []string
}

func newTabModel() *TabModel {
	tabs := []string{"TopTal", "GitHub", "GitHub Community", "GitHub Global", "Others"}

	return &TabModel{
		currentTab: 0,
		tabs:       tabs,
	}
}

func (m *TabModel) Init() tea.Cmd {
	return nil
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
	var view string

	for i, tab := range m.tabs {
		if i == m.currentTab {
			view += activeTabStyle.Render(tab)
		} else {
			view += inactiveTabStyle.Render(tab)
		}

		if i < len(m.tabs)-1 {
			view += dividerStyle.Render()
		}
	}

	return tabSection.Render(view)
}
