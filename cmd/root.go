package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gignr",
	Short: "Effortlessly manage and generate .gitignore files",
	Long: fmt.Sprintf(`gignr is a CLI tool designed to help you fetch, manage, and customize .gitignore templates 
from popular repositories. Simplify your project setup with ease.

Created by github.com/jasonuc.
Visit %v for more information.`, color.New(color.FgBlue).Sprint("https://github.com/jasonuc/gignr")),
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gignr.yaml)")
}
