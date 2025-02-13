package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonuc/gignr/internal/cache"
	"github.com/jasonuc/gignr/internal/templates"
	"github.com/spf13/viper"
)

type templateSrc string

const (
	TopTal          templateSrc = "TopTal"
	GitHub          templateSrc = "GitHub"
	GitHubCommunity templateSrc = "GitHub Community"
	GitHubGlobal    templateSrc = "GitHub Global"
	Others          templateSrc = "Others"
)

type TemplateEntry struct {
	Name     string
	Selected bool
	Source   string
}

type SourceData struct {
	Templates    []TemplateEntry
	CurrentIndex int
}

type CachedTemplates struct {
	Sources map[string]*SourceData
}

type sourceChangeMsg struct {
	NewSource templateSrc
}

type TemplateListModel struct {
	Templates    CachedTemplates
	ActiveSource templateSrc
	viewport     viewport.Model
	pageSize     int
	currentPage  int
	width        int
	height       int
	focused      bool
}

func newTemplateListModel() *TemplateListModel {
	model := &TemplateListModel{
		ActiveSource: TopTal,
		Templates:    CachedTemplates{Sources: make(map[string]*SourceData)},
		pageSize:     10,
		currentPage:  0,
		viewport:     viewport.New(80, 20),
		width:        80,
		height:       20,
		focused:      true,
	}

	// Initialize all sources
	sources := []templateSrc{TopTal, GitHub, GitHubCommunity, GitHubGlobal, Others}
	for _, source := range sources {
		model.Templates.Sources[string(source)] = &SourceData{
			Templates:    make([]TemplateEntry, 0),
			CurrentIndex: 0,
		}
	}

	// Load and distribute templates
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

	// Load local templates
	if localTemplates, err := listLocalTemplates(); err == nil {
		othersData := model.Templates.Sources[string(Others)]
		for _, file := range localTemplates {
			othersData.Templates = append(othersData.Templates, TemplateEntry{
				Name:     file,
				Selected: false,
			})
		}
	}

	return model
}

func (m *TemplateListModel) Init() tea.Cmd {
	return nil
}

func (m *TemplateListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	sourceData, exists := m.Templates.Sources[string(m.ActiveSource)]
	if !exists {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 6 // Account for header and footer
		m.viewport.Width = m.width
		m.viewport.Height = m.height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if sourceData.CurrentIndex > 0 {
				sourceData.CurrentIndex--
				m.ensureVisibleItem()
			}
		case "down":
			if sourceData.CurrentIndex < len(sourceData.Templates)-1 {
				sourceData.CurrentIndex++
				m.ensureVisibleItem()
			}
		case "enter":
			if len(sourceData.Templates) > 0 {
				sourceData.Templates[sourceData.CurrentIndex].Selected =
					!sourceData.Templates[sourceData.CurrentIndex].Selected
			}
		case "home":
			sourceData.CurrentIndex = 0
			m.ensureVisibleItem()
		case "end":
			sourceData.CurrentIndex = len(sourceData.Templates) - 1
			m.ensureVisibleItem()
		case "pgup":
			sourceData.CurrentIndex -= m.pageSize
			if sourceData.CurrentIndex < 0 {
				sourceData.CurrentIndex = 0
			}
			m.ensureVisibleItem()
		case "pgdown":
			sourceData.CurrentIndex += m.pageSize
			if sourceData.CurrentIndex >= len(sourceData.Templates) {
				sourceData.CurrentIndex = len(sourceData.Templates) - 1
			}
			m.ensureVisibleItem()
		}

	case sourceChangeMsg:
		m.ActiveSource = msg.NewSource
		m.currentPage = 0
		m.viewport.GotoTop()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
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

func (m *TemplateListModel) View() string {
	sourceData, exists := m.Templates.Sources[string(m.ActiveSource)]
	if !exists || len(sourceData.Templates) == 0 {
		return templateListStyle.Render(
			noTemplatesStyle.Render("No templates found in this source"))
	}

	// Build simple progress indicator
	currentIdx := sourceData.CurrentIndex + 1
	total := len(sourceData.Templates)
	progress := progressStyle.Render(fmt.Sprintf("(%d/%d)", currentIdx, total))

	// Build template list
	var content strings.Builder
	for i, template := range sourceData.Templates {
		templateLine := m.renderTemplateItem(template, i == sourceData.CurrentIndex)
		content.WriteString(templateLine + "\n")
	}

	// Set viewport content
	m.viewport.SetContent(content.String())

	// Build footer with hotkeys
	footer := hotkeysStyle.Render("↑/↓: Navigate • Enter: Select • PgUp/PgDn: Page • Home/End: Jump • Shift+C: Copy Command")

	// Combine all elements
	listContent := templateListStyle.Render(m.viewport.View())

	// Combine everything with the progress indicator right-aligned
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(m.width).Align(lipgloss.Right).Render(progress),
		listContent,
		footer,
	)
}

func (m *TemplateListModel) renderTemplateItem(template TemplateEntry, isCurrent bool) string {
	var builder strings.Builder

	// Cursor indicator
	if isCurrent {
		builder.WriteString(pointerStyle.Render("→ "))
	} else {
		builder.WriteString("  ")
	}

	// Checkbox
	if template.Selected {
		builder.WriteString(checkboxStyle.Render("[✓] "))
	} else {
		builder.WriteString(checkboxStyle.Render("[ ] "))
	}

	// Template name
	name := template.Name
	if isCurrent {
		name = selectedItemStyle.Render(name)
	} else {
		name = templateNameStyle.Render(name)
	}
	builder.WriteString(name)

	return builder.String()
}

func mapTemplateSource(filename, source string) string {
	// First check for known filenames
	lcFilename := strings.ToLower(filename)
	if strings.Contains(lcFilename, "toptal") {
		return string(TopTal)
	}

	// Then check custom repositories (the filename would be the nickname)
	repos := viper.GetStringMapString("repositories")
	if _, exists := repos[strings.TrimSuffix(filename, ".json")]; exists {
		return string(Others)
	}

	// Finally check GitHub sources
	switch source {
	case "GitHub Community":
		return string(GitHubCommunity)
	case "GitHub Global":
		return string(GitHubGlobal)
	case "GitHub":
		return string(GitHub)
	default:
		// If file is from a custom repo (checked above) or local template
		return string(Others)
	}
}

func listCacheFiles() (map[string]string, error) {
	cacheDir := cache.GetCacheDir()
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}

	sourceMap := make(map[string]string)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			sourceMap[strings.TrimSuffix(file.Name(), ".json")] = file.Name()
		}
	}
	return sourceMap, nil
}

func listLocalTemplates() ([]string, error) {
	templatesDir := filepath.Join(os.Getenv("HOME"), ".config/gignr/templates")
	files, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".gitignore") {
			templates = append(templates, file.Name())
		}
	}
	return templates, nil
}

// GetSelectedTemplates returns all selected templates across all sources
func (m *TemplateListModel) GetSelectedTemplates() []TemplateEntry {
	selected := make([]TemplateEntry, 0)
	for _, sourceData := range m.Templates.Sources {
		for _, template := range sourceData.Templates {
			if template.Selected {
				selected = append(selected, template)
			}
		}
	}
	return selected
}

func (m *SearchModel) HandleSave() {
	templateList, ok := m.TemplateList.(*TemplateListModel)
	if !ok {
		return
	}

	// Get all selected templates
	var parts []string
	for sourceName, sourceData := range templateList.Templates.Sources {
		for _, template := range sourceData.Templates {
			if !template.Selected {
				continue
			}

			// Map source to prefix
			var prefix string
			switch templateSrc(sourceName) {
			case GitHub:
				prefix = "gh:"
			case GitHubCommunity:
				prefix = "ghc:"
			case GitHubGlobal:
				prefix = "ghg:"
			case TopTal:
				prefix = "tt:"
			case Others:
				storagePath := viper.GetString("templates.storage_path")

				// Check if this template exists in storage path
				localPath := filepath.Join(storagePath, template.Name)
				if _, err := os.Stat(localPath); err == nil {
					// It's a local template
					prefix = ""
				} else {
					// It's from a custom repo - get the nickname from the repository map
					repos := viper.GetStringMapString("repositories")
					for nickname := range repos {
						// Check if this template came from nickname's cache file
						cacheFile := fmt.Sprintf("%s.json", nickname)
						if cachedTemplates, err := templates.LoadCachedTemplates(cacheFile); err == nil {
							// Look for this template in the cache
							for _, ct := range cachedTemplates {
								if ct.Name == template.Name {
									prefix = nickname + ":"
									break
								}
							}
						}
						if prefix != "" {
							break
						}
					}
				}
			}

			templateName := strings.TrimSuffix(template.Name, ".gitignore")
			if prefix == "" {
				parts = append(parts, templateName)
			} else {
				parts = append(parts, prefix+templateName)
			}
		}
	}

	command := fmt.Sprintf("gignr create %s", strings.Join(parts, " "))

	if err := clipboard.WriteAll(command); err != nil {
		return
	}
}
