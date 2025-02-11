package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jasonuc/gignr/internal/tui"
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

		// Validate URL
		if !isValidGitHubURL(repoURL) {
			fmt.Println("Invalid GitHub URL. Must be in format: https://github.com/{user}/{repo}")
			return
		}

		// Validate nickname
		if isReservedNickname(nickname) {
			fmt.Println("Invalid nickname. Reserved names: gh, ghc, ghg, tt.")
			return
		}

		// Load repositories from Viper
		repos := viper.GetStringMapString("repositories")

		// If the nickname already exists, confirm before overwriting
		if _, exists := repos[nickname]; exists {
			prompt := fmt.Sprintf("The nickname '%s' already exists. Overwrite?\n%s → %s", nickname, nickname, repoURL)
			if !tui.RunConfirmation(prompt) {
				fmt.Println("Operation canceled.")
				return
			}
		}

		// Save repository to Viper
		repos[nickname] = repoURL
		viper.Set("repositories", repos)

		// Write changes to config file
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error saving repository:", err)
			return
		}

		fmt.Printf("Successfully added repository %s as %s\n", repoURL, nickname)
	},
}

func init() {
	addCmd.Flags().StringVarP(&nickname, "nickname", "n", "", "Nickname for the repository")
	addCmd.MarkFlagRequired("nickname")
	rootCmd.AddCommand(addCmd)
}

// isValidGitHubURL checks if a URL is a valid GitHub repo URL
func isValidGitHubURL(url string) bool {
	if len(url) > 19 && url[:19] == "https://github.com/" {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Warning: Could not verify repository (network error). Assuming invalid.")
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == http.StatusOK
	}
	return false
}

// isReservedNickname ensures nickname is not a reserved name
func isReservedNickname(nickname string) bool {
	reserved := map[string]bool{"gh": true, "ghc": true, "ghg": true, "tt": true}
	return reserved[strings.ToLower(nickname)]
}
