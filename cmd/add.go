package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Add a custom GitHub repository with .gitignore templates",
	Long: `Add a custom GitHub repository containing .gitignore templates. 
Once added, the repository will be used as a source for fetching templates.`,
	Example: "gignr add https://github.com/jasonuc/gitignore",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
