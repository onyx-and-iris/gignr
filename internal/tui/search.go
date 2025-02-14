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
	styles       *AppStyle
	width        int
	height       int
}

func newSearchModel() *SearchModel {
	styles := NewAppStyle(80, 24)

	ti := textinput.New()
	ti.Placeholder = "Type to search templates..."
	ti.Focus()
	ti.PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	ti.TextStyle = lipgloss.NewStyle().Foreground(textColor)
	ti.Width = styles.width - 4

	return &SearchModel{
		Tab:          newTabModel(styles),
		TextInput:    ti,
		TemplateList: newTemplateListModel(styles),
		styles:       styles,
		width:        styles.width,
		height:       styles.height,
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
	case tea.WindowSizeMsg:
		if _, cmd := m.Tab.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
		if _, cmd := m.TemplateList.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}

		m.width = msg.Width
		m.height = msg.Height
		m.styles.SetSize(msg.Width, msg.Height)

		m.TextInput.Width = msg.Width - 4

		if tab, ok := m.Tab.(*TabModel); ok {
			tab.SetStyles(m.styles)
		}
		if list, ok := m.TemplateList.(*TemplateListModel); ok {
			list.SetStyles(m.styles)
		}
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
	header := m.Tab.View()
	searchInput := m.styles.searchBox.Render(m.TextInput.View())
	templateList := m.TemplateList.View()
	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		searchInput,
		templateList,
	)
}

func RunSearch() error {
	if _, err := tea.NewProgram(newSearchModel()).Run(); err != nil {
		return err
	}
	return nil
}
