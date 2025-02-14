package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
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

func (m *TemplateListModel) renderTemplateItem(template TemplateEntry, isCurrent bool) string {
	var b strings.Builder
	prefix := "  "
	if isCurrent {
		prefix = m.styles.pointer.Render("→ ")
	}

	checkbox := "[ ] "
	if template.Selected {
		checkbox = "[✓] "
	}

	name := strings.TrimSuffix(template.Name, ".gitignore")
	if isCurrent {
		name = m.styles.selectedItem.Render(name)
	} else {
		name = m.styles.templateName.Render(name)
	}

	b.WriteString(prefix)
	b.WriteString(m.styles.checkbox.Render(checkbox))
	b.WriteString(name)

	return b.String()
}

func mapTemplateSource(filename, source string) string {
	lcFilename := strings.ToLower(filename)
	if strings.Contains(lcFilename, "toptal") {
		return string(TopTal)
	}

	repos := viper.GetStringMapString("repositories")
	if _, exists := repos[strings.TrimSuffix(filename, ".json")]; exists {
		return string(Others)
	}

	switch source {
	case "GitHub Community":
		return string(GitHubCommunity)
	case "GitHub Global":
		return string(GitHubGlobal)
	case "GitHub":
		return string(GitHub)
	default:
		return string(Others)
	}
}

func listCacheFiles() (map[string]string, error) {
	files, err := os.ReadDir(cache.GetCacheDir())
	if err != nil {
		return nil, err
	}

	sources := make(map[string]string)
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".json") {
			sources[strings.TrimSuffix(name, ".json")] = name
		}
	}
	return sources, nil
}

func listLocalTemplates() ([]string, error) {
	files, err := os.ReadDir(filepath.Join(os.Getenv("HOME"), ".config/gignr/templates"))
	if err != nil {
		return nil, err
	}

	templates := make([]string, 0, len(files))
	for _, f := range files {
		if name := f.Name(); strings.HasSuffix(name, ".gitignore") {
			templates = append(templates, name)
		}
	}
	return templates, nil
}

func (m *TemplateListModel) GetSelectedTemplates() []TemplateEntry {
	var selected []TemplateEntry
	for _, src := range m.Templates.Sources {
		for _, tmpl := range src.Templates {
			if tmpl.Selected {
				selected = append(selected, tmpl)
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

	command := fmt.Sprintf("gignr create %s", strings.Join(buildTemplateParts(templateList), " "))
	clipboard.WriteAll(command)
}

func buildTemplateParts(templateList *TemplateListModel) []string {
	var parts []string
	for src, data := range templateList.Templates.Sources {
		for _, tmpl := range data.Templates {
			if !tmpl.Selected {
				continue
			}
			if part := formatTemplatePart(templateSrc(src), tmpl); part != "" {
				parts = append(parts, part)
			}
		}
	}
	return parts
}

func formatTemplatePart(source templateSrc, tmpl TemplateEntry) string {
	prefix := getSourcePrefix(source, tmpl)
	name := strings.TrimSuffix(tmpl.Name, ".gitignore")
	if prefix == "" {
		return name
	}
	return prefix + name
}

func getSourcePrefix(source templateSrc, tmpl TemplateEntry) string {
	switch source {
	case GitHub:
		return "gh:"
	case GitHubCommunity:
		return "ghc:"
	case GitHubGlobal:
		return "ghg:"
	case TopTal:
		return "tt:"
	case Others:
		return getOthersPrefix(tmpl)
	default:
		return ""
	}
}

func getOthersPrefix(tmpl TemplateEntry) string {
	storagePath := viper.GetString("templates.storage_path")
	if _, err := os.Stat(filepath.Join(storagePath, tmpl.Name)); err == nil {
		return ""
	}

	repos := viper.GetStringMapString("repositories")
	for nick := range repos {
		cached, err := templates.LoadCachedTemplates(fmt.Sprintf("%s.json", nick))
		if err != nil {
			continue
		}
		for _, t := range cached {
			if t.Name == tmpl.Name {
				return nick + ":"
			}
		}
	}
	return ""
}

func (m *TemplateListModel) FilterTemplates(searchText string) {
	m.filterText = searchText
	src := m.Templates.Sources[string(m.ActiveSource)]
	if src == nil {
		return
	}

	if searchText == "" {
		m.filteredTemplates = src.Templates
		return
	}

	searchLower := strings.ToLower(searchText)
	filtered := make([]TemplateEntry, 0, len(src.Templates))
	for _, tmpl := range src.Templates {
		name := strings.ToLower(strings.TrimSuffix(tmpl.Name, ".gitignore"))
		if strings.Contains(name, searchLower) {
			filtered = append(filtered, tmpl)
		}
	}

	m.filteredTemplates = filtered
	if src.CurrentIndex >= len(filtered) {
		src.CurrentIndex = 0
	}
}
