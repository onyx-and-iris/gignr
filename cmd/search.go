package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasonuc/gignr/internal/cache"
	"github.com/jasonuc/gignr/internal/templates"
	"github.com/jasonuc/gignr/internal/tui"
	"github.com/jasonuc/gignr/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search and browse templates interactively",
	Long: `The search command launches an interactive TUI (Terminal User Interface) to display available .gitignore templates.

Using the TUI, you can:
  - Browse all templates from different sources (GitHub, TopTal, etc.).
  - Filter templates by name or category.
  - Select templates to view details or mark them for use in other commands.
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Initialize GitHub client first
		templates.InitGitHubClient("")

		// Check if default caches exist and fetch if needed
		cacheDir := cache.GetCacheDir()

		// Check and fetch GitHub templates
		if _, err := os.Stat(filepath.Join(cacheDir, "github.json")); os.IsNotExist(err) {
			if fetchableTemplates, err := templates.FetchTemplates("github", "gitignore", ""); err == nil {
				templates.SaveTemplatesToCache("github.json", fetchableTemplates)
			}
		}

		// Check and fetch TopTal templates
		if _, err := os.Stat(filepath.Join(cacheDir, "toptal.json")); os.IsNotExist(err) {
			if fetchableTemplates, err := templates.FetchTemplates("toptal", "gitignore", "templates"); err == nil {
				templates.SaveTemplatesToCache("toptal.json", fetchableTemplates)
			}
		}

		// Check if cache needs refresh for custom repos
		if viper.GetBool("cache_needs_refresh") {
			// Get all custom repositories
			repos := viper.GetStringMapString("repositories")

			for nickname, repoURL := range repos {
				// Extract owner and repo from URL
				// URL format: https://github.com/owner/repo
				parts := strings.Split(repoURL, "/")
				if len(parts) < 5 {
					utils.PrintWarning(fmt.Sprintf("Invalid repository URL format for %s", nickname))
					continue
				}
				owner := parts[3]
				repo := parts[4]

				// Fetch fetchedTemplates from this repo
				if fetchedTemplates, err := templates.FetchTemplates(owner, repo, ""); err == nil {
					// Save to cache with nickname as identifier
					cacheFile := fmt.Sprintf("%s.json", nickname)
					templates.SaveTemplatesToCache(cacheFile, fetchedTemplates)
				} else {
					utils.PrintWarning(fmt.Sprintf("Failed to fetch templates from %s: %v", nickname, err))
				}
			}

			// Reset the refresh flag
			viper.Set("cache_needs_refresh", false)
			if err := viper.WriteConfig(); err != nil {
				utils.PrintWarning("Failed to update cache refresh status")
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.RunSearch(); err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
