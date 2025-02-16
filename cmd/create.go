package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/jasonuc/gignr/internal/templates"
	"github.com/jasonuc/gignr/internal/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create <template> [templates]...",
	Example: "gignr create gh:Go tt:clion my-template",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Generate a .gitignore file using one or more templates",
	Long: `The create command generates a .gitignore file based on one or more templates of your choice.

Available templates are identified by prefixes:
  - tt: TopTal templates
  - gh: GitHub templates
  - ghg: GitHub Global templates
  - ghc: GitHub Community templates
  - (no prefix) → Fetch from local saved templates
`,
	Run: func(cmd *cobra.Command, args []string) {
		var mergedContent strings.Builder
		var hasErrors bool

		templates.InitGitHubClient("")
		repos := templates.LoadCustomRepositories()

		for _, arg := range args {
			content, err := processTemplate(arg, repos)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Error processing %s: %v", arg, err))
				hasErrors = true
				continue
			}

			addTemplateToContent(&mergedContent, arg, content)
		}

		if hasErrors {
			utils.PrintWarning("Some templates failed to process. .gitignore file will be incomplete.")
		}

		if err := os.WriteFile(".gitignore", []byte(mergedContent.String()), 0644); err != nil {
			utils.PrintError(fmt.Sprintf("Failed to write .gitignore file: %v", err))
			return
		}

		utils.PrintSuccess("Created .gitignore!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

type TemplateSource struct {
	owner string
	repo  string
	path  string
}

func resolveTemplateSource(prefix string, repos map[string]string) (*TemplateSource, error) {
	switch prefix {
	case "tt":
		return &TemplateSource{"toptal", "gitignore", "templates"}, nil
	case "gh":
		return &TemplateSource{"github", "gitignore", ""}, nil
	case "ghc":
		return &TemplateSource{"github", "gitignore", "community"}, nil
	case "ghg":
		return &TemplateSource{"github", "gitignore", "Global"}, nil
	default:
		// Check user-defined repos
		if repoURL, exists := repos[prefix]; exists {
			owner, repo, err := utils.ExtractRepoDetails(repoURL)
			if err != nil {
				return nil, fmt.Errorf("invalid repository URL for prefix %s", prefix)
			}
			return &TemplateSource{owner, repo, ""}, nil
		}
		return nil, fmt.Errorf("unknown template prefix or missing repository: %s", prefix)
	}
}

func findTemplate(templateName string, templates []templates.Template) (string, error) {
	// Try exact match first
	for _, tmpl := range templates {
		if tmpl.Name == templateName+".gitignore" {
			return tmpl.DownloadURL, nil
		}
	}

	// Try case-insensitive match
	for _, tmpl := range templates {
		if strings.EqualFold(tmpl.Name, templateName+".gitignore") {
			return tmpl.DownloadURL, nil
		}
	}

	return "", fmt.Errorf("template %s not found", templateName)
}

func createTemplateBox() box.Box {
	config := box.Config{Px: 1, Py: 1, Type: "", TitlePos: "Inside"}
	return box.Box{
		TopRight: "*", TopLeft: "*",
		BottomRight: "*", BottomLeft: "*",
		Horizontal: "-", Vertical: "|",
		Config: config,
	}
}

func addTemplateToContent(builder *strings.Builder, templateName string, content []byte) {
	boxTitle := fmt.Sprintf(" %s",
		strings.ToUpper(templateName))

	box := createTemplateBox()
	builder.WriteString(box.String("", boxTitle))
	builder.Write(content)
	builder.WriteString("\n\n")
}

func processTemplate(arg string, repos map[string]string) (content []byte, err error) {
	if strings.Contains(arg, ":") {
		// Handle remote templates (gh:, tt:, etc)
		parts := strings.SplitAfter(arg, ":")
		prefix := strings.TrimSpace(parts[0][:len(parts[0])-1])
		templateName := strings.TrimSpace(parts[1])

		src, err := resolveTemplateSource(prefix, repos)
		if err != nil {
			return nil, err
		}

		// Fetch available templates
		templateList, err := templates.FetchTemplates(src.owner, src.repo, src.path, prefix)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch templates from %s: %v", prefix, err)
		}

		// Find the specific template
		downloadURL, err := findTemplate(templateName, templateList)
		if err != nil {
			return nil, err
		}

		content, err = templates.GetTemplateContent(downloadURL)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch content: %v", err)
		}

		return content, nil
	}

	// Handle local templates
	content, err = templates.GetLocalTemplate(arg)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch local template: %v", err)
	}

	return content, nil
}
