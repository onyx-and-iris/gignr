package cmd

import (
	"github.com/spf13/cobra"
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
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
