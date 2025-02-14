package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchModel struct {
	Tab          tea.Model
	TextInput    textinput.Model
	TemplateList tea.Model
	Keymap       any
}

func newSearchModel() *SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Type to search templates..."
	ti.Focus()
	ti.PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	ti.TextStyle = lipgloss.NewStyle().Foreground(textColor)

	return &SearchModel{
		Tab:          newTabModel(),
		TextInput:    ti,
		TemplateList: newTemplateListModel(),
	}
}

func (m *SearchModel) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		tea.EnterAltScreen,
		textinput.Blink,
	)
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
		case "up", "down", "enter", "home", "end", "pgup", "pgdown":
			m.TemplateList, cmd = m.TemplateList.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.TextInput, cmd = m.TextInput.Update(msg)
			cmds = append(cmds, cmd)

			if templateList, ok := m.TemplateList.(*TemplateListModel); ok {
				templateList.FilterTemplates(m.TextInput.Value())
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *SearchModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Tab.View(),
		searchBoxStyle.Render(m.TextInput.View()),
		m.TemplateList.View(),
	)
}

func RunSearch() error {
	if _, err := tea.NewProgram(newSearchModel()).Run(); err != nil {
		return err
	}
	return nil
}
