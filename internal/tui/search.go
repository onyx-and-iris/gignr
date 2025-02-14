package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SearchModel struct {
	Tab          tea.Model
	TextInput    any
	TemplateList tea.Model
	Keymap       any
}

func newSearchModel() *SearchModel {
	return &SearchModel{
		Tab:          newTabModel(),
		TextInput:    struct{}{},
		TemplateList: newTemplateListModel(),
		Keymap:       struct{}{},
	}
}

func (m *SearchModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "C":
			m.HandleSave()
		case "left", "right", "tab":
			m.Tab, cmd = m.Tab.Update(msg)
			cmds = append(cmds, cmd)

			tab, ok := m.Tab.(*TabModel)
			if !ok {
				return m, tea.Quit
			}
			newSource := templateSrc(tab.tabs[tab.currentTab])
			m.TemplateList, cmd = m.TemplateList.Update(sourceChangeMsg{newSource})
			cmds = append(cmds, cmd)
		case "up", "down", "enter":
			m.TemplateList, cmd = m.TemplateList.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *SearchModel) View() string {
	var view string

	view += m.Tab.View() + "\n"
	view += m.TemplateList.View() + "\n"

	return view
}

func RunSearch() error {
	if _, err := tea.NewProgram(newSearchModel()).Run(); err != nil {
		return err
	}
	return nil
}
