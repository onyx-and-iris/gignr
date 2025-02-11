package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jasonuc/gignr/internal/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:     "create <template> [templates]...",
	Example: "gignr create gh:Go tt:clion jc:Rust",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Generate a .gitignore file using one or more templates",
	Long: `The create command generates a .gitignore file based on one or more templates of your choice.

Available templates are identified by prefixes:
  - tt: TopTal templates
  - gh: GitHub templates
  - ghg: GitHub Global templates
  - ghc: GitHub Community templates
  - [nickname]: User-added repositories
`,
	Run: func(cmd *cobra.Command, args []string) {
		var mergedContent strings.Builder
		templates.InitGitHubClient("")

		// Load user-added repositories
		userRepos := viper.GetStringMapString("repositories")

		for _, arg := range args {
			req := strings.SplitAfter(arg, ":")
			reqPrefix := strings.TrimSpace(req[0][:len(req[0])-1])
			templateName := strings.TrimSpace(req[1])

			// Define repo owner, repo name, and path for different sources
			var owner, repo, path string

			switch reqPrefix {
			case "tt":
				owner, repo, path = "toptal", "gitignore", "templates"
			case "gh":
				owner, repo, path = "github", "gitignore", ""
			case "ghc":
				owner, repo, path = "github", "gitignore", "community"
			case "ghg":
				owner, repo, path = "github", "gitignore", "Global"
			default:
				// Check if the prefix is a user-added repo
				if repoURL, exists := userRepos[reqPrefix]; exists {
					// Extract owner and repo from the URL
					splitURL := strings.Split(strings.TrimPrefix(repoURL, "https://github.com/"), "/")
					if len(splitURL) < 2 {
						fmt.Printf("Invalid repository format for %s: %s\n", reqPrefix, repoURL)
						continue
					}
					owner, repo, path = splitURL[0], splitURL[1], ""
				} else {
					fmt.Printf("Unknown template prefix: %s\n", reqPrefix)
					continue
				}
			}

			// Fetch available templates
			templateList, err := templates.FetchTemplates(owner, repo, path)
			if err != nil {
				fmt.Printf("Error fetching templates from %s: %v\n", reqPrefix, err)
				continue
			}

			// Find the requested template
			var downloadURL string
			for _, tmpl := range templateList {
				if strings.EqualFold(tmpl.Name, templateName+".gitignore") {
					downloadURL = tmpl.DownloadURL
					break
				}
			}

			if downloadURL == "" {
				fmt.Printf("Template %s not found in %s.\n", templateName, reqPrefix)
				continue
			}

			// Fetch the raw .gitignore content
			content, err := templates.FetchContent(downloadURL)
			if err != nil {
				fmt.Printf("Error fetching content for %s: %v\n", templateName, err)
				continue
			}

			// Merge content
			mergedContent.WriteString(fmt.Sprintf("\n##########  %s Template (%s)  ##########\n\n", strings.ToUpper(templateName), strings.ToUpper(reqPrefix)))
			mergedContent.Write(content)
			mergedContent.WriteString("\n\n")
		}

		// Write to .gitignore
		err := os.WriteFile(".gitignore", []byte(mergedContent.String()), 0644)
		if err != nil {
			fmt.Printf("Failed to write .gitignore file: %v\n", err)
			return
		}

		fmt.Println("Successfully created .gitignore!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
