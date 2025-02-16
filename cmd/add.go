package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jasonuc/gignr/internal/cache"
	"github.com/jasonuc/gignr/internal/tui"
	"github.com/jasonuc/gignr/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GitHubContentResponse []struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var nickname string
var repoURL string

var addCmd = &cobra.Command{
	Use:     "add <url> -nickname <nickname>",
	Short:   "Add a custom GitHub repository with .gitignore templates",
	Long:    "Add a custom GitHub repository containing .gitignore templates. Once added, the repository will be used as a source for fetching templates.",
	Example: "gignr add https://github.com/jasonuc/gitignore -nickname jc",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL = args[0]

		if !utils.IsValidGitHubURL(repoURL) {
			utils.PrintError("Invalid GitHub URL. Must be in format: https://github.com/{user}/{repo}")
			return
		}
		if !utils.IsValidNickname(nickname) {
			utils.PrintError("Invalid nickname. Must be alphanumeric and contain no spaces.")
			return
		}
		if utils.IsReservedNickname(nickname) {
			utils.PrintError("Invalid nickname. Reserved names: gh, ghc, ghg, tt.")
			return
		}

		hasTemplates, err := validateGitignoreTemplates(repoURL)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to validate repository: %v", err))
			return
		}
		if !hasTemplates {
			utils.PrintError("Repository must contain named .gitignore files (e.g., python.gitignore, node.gitignore)")
			return
		}

		repos := viper.GetStringMapString("repositories")

		if _, exists := repos[nickname]; exists {
			prompt := fmt.Sprintf("The nickname '%s' already exists. Overwrite?\n%s → %s", nickname, nickname, repoURL)
			if !tui.RunConfirmation(prompt) {
				utils.PrintAlert("Operation canceled.")
				return
			}
		}

		repos[nickname] = repoURL
		viper.Set("repositories", repos)

		cache.UpdateCacheNeedRefreshStatus(true)

		if err := viper.WriteConfig(); err != nil {
			utils.PrintError(fmt.Sprintf("Unable to save repository: %v", err))
			return
		}

		utils.PrintSuccess(fmt.Sprintf("Added repository %s as %s\nUse with: gignr create %s:template-name", repoURL, nickname, nickname))
	},
}

func init() {
	addCmd.Flags().StringVarP(&nickname, "nickname", "n", "", "Nickname for the repository")
	addCmd.MarkFlagRequired("nickname")
	rootCmd.AddCommand(addCmd)
}

func validateGitignoreTemplates(repoURL string) (bool, error) {
	parts := strings.Split(strings.TrimSuffix(repoURL, "/"), "/")
	if len(parts) < 2 {
		return false, fmt.Errorf("invalid GitHub URL format")
	}
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return false, fmt.Errorf("failed to fetch repository contents: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch repository contents: %s", resp.Status)
	}

	var contents GitHubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return false, fmt.Errorf("failed to parse repository contents: %v", err)
	}

	var hasTemplates bool
	for _, item := range contents {
		if item.Type == "file" && strings.HasSuffix(item.Name, ".gitignore") && item.Name != ".gitignore" {
			hasTemplates = true
			break
		}
	}

	return hasTemplates, nil
}
