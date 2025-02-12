package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jasonuc/gignr/internal/cache"
	"github.com/jasonuc/gignr/internal/templates"
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
	pageSize     int
	currentPage  int
}

func newTemplateListModel() *TemplateListModel {
	model := &TemplateListModel{
		ActiveSource: GitHub,
		Templates:    CachedTemplates{Sources: make(map[string]*SourceData)},
		pageSize:     10,
		currentPage:  0,
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		sourceData, exists := m.Templates.Sources[string(m.ActiveSource)]
		if !exists {
			return m, nil
		}

		switch msg.String() {
		case "up":
			if sourceData.CurrentIndex > 0 {
				sourceData.CurrentIndex--
				m.currentPage = sourceData.CurrentIndex / m.pageSize
			}
		case "down":
			if sourceData.CurrentIndex < len(sourceData.Templates)-1 {
				sourceData.CurrentIndex++
				m.currentPage = sourceData.CurrentIndex / m.pageSize
			}
		case "enter":
			if len(sourceData.Templates) > 0 {
				sourceData.Templates[sourceData.CurrentIndex].Selected =
					!sourceData.Templates[sourceData.CurrentIndex].Selected
			}
		}
	case sourceChangeMsg:
		m.ActiveSource = msg.NewSource
		m.currentPage = 0
	}
	return m, nil
}

func (m *TemplateListModel) View() string {
	sourceData, exists := m.Templates.Sources[string(m.ActiveSource)]
	if !exists || len(sourceData.Templates) == 0 {
		return "No templates found\n"
	}

	startIdx := m.currentPage * m.pageSize
	endIdx := startIdx + m.pageSize
	if endIdx > len(sourceData.Templates) {
		endIdx = len(sourceData.Templates)
	}

	var view strings.Builder
	view.WriteString(fmt.Sprintf("%s Templates (Page %d/%d)\n",
		m.ActiveSource,
		m.currentPage+1,
		(len(sourceData.Templates)+m.pageSize-1)/m.pageSize))

	for i := startIdx; i < endIdx; i++ {
		template := sourceData.Templates[i]
		cursor := " "
		if i == sourceData.CurrentIndex {
			cursor = "→"
		}
		checkmark := " "
		if template.Selected {
			checkmark = "✓"
		}
		view.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checkmark, template.Name))
	}

	return view.String()
}

func mapTemplateSource(filename, source string) string {
	lcFilename := strings.ToLower(filename)
	if strings.Contains(lcFilename, "toptal") {
		return string(TopTal)
	}

	switch source {
	case "GitHub Community":
		return string(GitHubCommunity)
	case "GitHub Global":
		return string(GitHubGlobal)
	case "GitHub":
		return string(GitHub)
	default:
		return string(GitHub)
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

func (m *SearchModel) HandleSave() {
	// templateList, ok := m.TemplateList.(*TemplateListModel)
	// if !ok {
	// 	return
	// }

	// selectedTemplates := templateList.GetSelectedTemplates()
	// // TODO: do something with the selected templates
	// for _, template := range selectedTemplates {
	// }
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
