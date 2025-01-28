package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create <template> [templates]...",
	Example: "gignr create gh:Go",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Generate a .gitignore file using a selected template",
	Long: `The create command generates a .gitignore file based on a template of your choice.

Available templates are identified by prefixes:
  - tt: TopTal templates
  - gh: GitHub templates
  - ghg: GitHub Global templates
  - ghc: GitHub Community templates
`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
