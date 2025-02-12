package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SearchModel struct {
	Tab              tea.Model
	TextInput        any
	TemplatesDisplay any
	Keymap           any
}

func newSearchModel() *SearchModel {
	return &SearchModel{
		Tab:              newTabModel(),
		TextInput:        struct{}{},
		TemplatesDisplay: struct{}{},
		Keymap:           struct{}{},
	}
}

func (m *SearchModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "left", "right", "tab":
			var cmd tea.Cmd
			m.Tab, cmd = m.Tab.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *SearchModel) View() string {
	var view string

	view += m.Tab.View() + "\n"

	return view
}

func RunSearch() error {
	if _, err := tea.NewProgram(newSearchModel()).Run(); err != nil {
		return err
	}

	return nil
}
