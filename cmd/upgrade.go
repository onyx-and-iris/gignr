package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/jasonuc/gignr/internal/utils"
	"github.com/jasonuc/gignr/internal/version"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "Installs the latest version of the CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.FromContext(cmd.Context())
		if !info.IsOutdated {
			utils.PrintAlert("gignr is already up to date.")
			return
		}

		// install the latest version
		command := exec.Command("go", "install", "github.com/jasonuc/gignr@latest")
		_, err := command.Output()
		cobra.CheckErr(err)

		// Get the new version info
		command = exec.Command("gignr", "--version")
		b, err := command.Output()
		cobra.CheckErr(err)
		re := regexp.MustCompile(`v\d+\.\d+\.\d+`)
		version := re.FindString(string(b))
		utils.PrintSuccess(fmt.Sprintf("gignr upgraded to %s\n", version))
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
