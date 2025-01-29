package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "gignr",
	Short: "Effortlessly manage and generate .gitignore files",
	Long: fmt.Sprintf(`gignr is a CLI tool designed to help you fetch, manage, and customize .gitignore templates 
from popular repositories. Simplify your project setup with ease.

Created by github.com/jasonuc.
Visit %v for more information.`, color.New(color.FgBlue).Sprint("https://github.com/jasonuc/gignr")),
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("aaa", "vvv")
		viper.WriteConfig()
	},
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
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configDir := home + "/.config/gignr"

	if err := os.MkdirAll(configDir, 0755); err != nil {
		cobra.CheckErr(err)
	}

	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		// If the file does not exist, create a default one
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			configFile := configDir + "/config.yaml"
			if err := viper.WriteConfigAs(configFile); err != nil {
				cobra.CheckErr(err)
			}

			return
		}
		cobra.CheckErr(err)
	}
}
