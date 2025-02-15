package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jasonuc/gignr/internal/cache"
	"github.com/jasonuc/gignr/internal/templates"
	"github.com/jasonuc/gignr/internal/tui"
	"github.com/jasonuc/gignr/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search and browse templates interactively",
	Long: `Search and browse .gitignore templates using an interactive TUI.
Navigate between sources, filter templates, and select them for use.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		templates.InitGitHubClient("")
		ensureDefaultCaches()
		refreshCustomRepos()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.RunSearch(); err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func ensureDefaultCaches() {
	cacheDir := cache.GetCacheDir()
	defaultSources := map[string]struct {
		owner string
		repo  string
		path  string
	}{
		"github.json": {"github", "gitignore", ""},
		"toptal.json": {"toptal", "gitignore", "templates"},
	}

	for filename, source := range defaultSources {
		cachePath := filepath.Join(cacheDir, filename)
		if _, err := os.Stat(cachePath); os.IsNotExist(err) {
			fetchAndSaveTemplates(source.owner, source.repo, source.path, "", filename)
		}
	}
}

func refreshCustomRepos() {
	if !viper.GetBool("cache_needs_refresh") {
		return
	}

	repos := viper.GetStringMapString("repositories")
	for nickname, repoURL := range repos {
		owner, repo, err := utils.ExtractRepoDetails(repoURL)
		if err != nil {
			utils.PrintWarning(fmt.Sprintf("Invalid repository URL format for %s", nickname))
			continue
		}

		cacheFile := fmt.Sprintf("%s.json", nickname)
		fetchAndSaveTemplates(owner, repo, "", nickname, cacheFile)
	}

	cache.UpdateCacheNeedRefreshStatus(false)
}

func fetchAndSaveTemplates(owner, repo, path, sourceID, cacheFile string) {
	fetchedTemplates, err := templates.FetchTemplates(owner, repo, path, sourceID)
	if err != nil {
		utils.PrintWarning(fmt.Sprintf("Failed to fetch templates from %s: %v", owner, err))
		return
	}
	templates.SaveTemplatesToCache(cacheFile, fetchedTemplates)
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
