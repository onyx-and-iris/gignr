package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonuc/gignr/internal/templates"
)

type TemplateListModel struct {
	Templates         CachedTemplates
	ActiveSource      templateSrc
	viewport          viewport.Model
	pageSize          int
	currentPage       int
	filteredTemplates []TemplateEntry
	filterText        string
	styles            *AppStyle
}

func newTemplateListModel(styles *AppStyle) *TemplateListModel {
	viewportWidth := styles.width - 4
	viewportHeight := styles.templateList.GetHeight()

	model := &TemplateListModel{
		ActiveSource: TopTal,
		Templates:    CachedTemplates{Sources: make(map[string]*SourceData)},
		pageSize:     10,
		currentPage:  0,
		viewport:     viewport.New(viewportWidth, viewportHeight),
		filterText:   "",
		styles:       styles,
	}

	model.viewport.KeyMap.PageDown.SetEnabled(false)

	sources := []templateSrc{TopTal, GitHub, GitHubCommunity, GitHubGlobal, Others}
	for _, source := range sources {
		model.Templates.Sources[string(source)] = &SourceData{
			Templates:    make([]TemplateEntry, 0),
			CurrentIndex: 0,
		}
	}

	if cacheFiles, err := listCacheFiles(); err == nil {
		for filename, cacheFile := range cacheFiles {
			if templates, err := templates.LoadCachedTemplates(cacheFile); err == nil {
				for _, template := range templates {
					sourceKey := mapTemplateSource(filename, template.Source)
					if sourceData, exists := model.Templates.Sources[sourceKey]; exists {
						sourceData.Templates = append(sourceData.Templates, TemplateEntry{
							Name:     template.Name,
							Selected: false,
						})
					}
				}
			}
		}
	}

	if localTemplates, err := listLocalTemplates(); err == nil {
		othersData := model.Templates.Sources[string(Others)]
		for _, file := range localTemplates {
			othersData.Templates = append(othersData.Templates, TemplateEntry{
				Name:     file,
				Selected: false,
			})
		}
	}

	if sourceData := model.Templates.Sources[string(model.ActiveSource)]; sourceData != nil {
		model.filteredTemplates = sourceData.Templates
	}

	return model
}

func (m *TemplateListModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m *TemplateListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	sourceData, exists := m.Templates.Sources[string(m.ActiveSource)]
	if !exists {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if sourceData.CurrentIndex > 0 {
				sourceData.CurrentIndex--
				m.ensureVisibleItem()
			}
		case "down":
			if sourceData.CurrentIndex < len(m.filteredTemplates)-1 {
				sourceData.CurrentIndex++
				m.ensureVisibleItem()
			}
		case "enter", " ":
			if len(m.filteredTemplates) > 0 {
				idx := sourceData.CurrentIndex
				name := m.filteredTemplates[idx].Name

				for i, t := range sourceData.Templates {
					if t.Name == name {
						sourceData.Templates[i].Selected = !sourceData.Templates[i].Selected
						newSelectedState := sourceData.Templates[i].Selected

						for j, ft := range m.filteredTemplates {
							if ft.Name == name {
								m.filteredTemplates[j].Selected = newSelectedState
								break
							}
						}
						break
					}
				}
			}
		}
	case sourceChangeMsg:
		m.ActiveSource = msg.NewSource
		m.currentPage = 0
		m.viewport.GotoTop()
		m.FilterTemplates(m.filterText)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *TemplateListModel) View() string {
	sourceData := m.Templates.Sources[string(m.ActiveSource)]
	if sourceData == nil || len(m.filteredTemplates) == 0 {
		message := "No templates found"
		if m.filterText != "" {
			message = "No matching templates found"
		}
		return m.styles.templateList.Render(
			m.styles.noTemplates.Render(message))
	}

	currentIdx := sourceData.CurrentIndex + 1
	total := len(m.filteredTemplates)
	progress := fmt.Sprintf("(%d/%d)", currentIdx, total)

	var content strings.Builder
	for i, template := range m.filteredTemplates {
		templateLine := m.renderTemplateItem(template, i == sourceData.CurrentIndex)
		content.WriteString(templateLine + "\n")
	}

	m.viewport.SetContent(content.String())

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.progress.Render(progress),
		m.styles.templateList.Render(m.viewport.View()),
		m.styles.hotkeys.Render("↑/↓: Navigate • Enter/Space: Select • Shift+C: Copy Command"),
	)

	return mainContent
}

func (m *TemplateListModel) ensureVisibleItem() {
	sourceData := m.Templates.Sources[string(m.ActiveSource)]
	itemHeight := 1
	viewportStart := m.viewport.YOffset
	viewportEnd := viewportStart + m.viewport.Height

	itemTop := sourceData.CurrentIndex * itemHeight
	itemBottom := itemTop + itemHeight

	if itemTop < viewportStart {
		m.viewport.YOffset = itemTop
	} else if itemBottom > viewportEnd {
		m.viewport.YOffset = itemBottom - m.viewport.Height
	}
}

func (m *TemplateListModel) SetStyles(styles *AppStyle) {
	m.styles = styles
	m.viewport.Width = styles.width - 4
	m.viewport.Height = styles.templateList.GetHeight()
}
