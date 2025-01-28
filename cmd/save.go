package cmd

import (
	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save -n <name>",
	Short: "Save the .gitignore file in your current directory to your templates",
	Long: `The save command allows you to save the current .gitignore file in your directory 
as a custom template under "My Templates".

You can specify a name for the template using the -n flag, which will be used to reference this template later.`,
	Args: cobra.MatchAll(cobra.NoArgs),
	Run:  func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
