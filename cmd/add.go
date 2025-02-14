package cmd

import (
	"fmt"

	"github.com/jasonuc/gignr/internal/cache"
	"github.com/jasonuc/gignr/internal/tui"
	"github.com/jasonuc/gignr/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

		// Validate URL and nickname
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

		// Load repositories from Viper
		repos := viper.GetStringMapString("repositories")

		// If the nickname already exists, confirm before overwriting
		if _, exists := repos[nickname]; exists {
			prompt := fmt.Sprintf("The nickname '%s' already exists. Overwrite?\n%s → %s", nickname, nickname, repoURL)
			if !tui.RunConfirmation(prompt) {
				utils.PrintAlert("Operation canceled.")
				return
			}
		}

		// Save repository to Viper
		repos[nickname] = repoURL
		viper.Set("repositories", repos)

		cache.UpdateCacheNeedRefreshStatus(true)

		// Write changes to config file
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
