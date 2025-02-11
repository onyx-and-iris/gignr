package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var saveCmd = &cobra.Command{
	Use:     "save <name>",
	Short:   "Save the current .gitignore file locally",
	Long:    `Save the current .gitignore file in your configured templates directory.`,
	Example: `gignr save my-template`,
	Args:    cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		setDefaultStoragePath()
	},

	Run: func(cmd *cobra.Command, args []string) {
		saveName := args[0]
		gitignorePath := ".gitignore"

		if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
			fmt.Println("No .gitignore file found in the current directory.")
			return
		}

		if !isValidTemplateName(saveName) {
			fmt.Println("Invalid name. Template names can only contain letters, numbers, dashes, and underscores.")
			return
		}

		if err := saveLocally(saveName, gitignorePath); err != nil {
			fmt.Println("Error saving locally:", err)
		} else {
			fmt.Println("Successfully saved locally as:", saveName)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}

// Ensures a default storage path is set in config.yaml
func setDefaultStoragePath() {
	if !viper.IsSet("templates.storage_path") {
		defaultPath := filepath.Join(os.Getenv("HOME"), ".config/gignr/templates")
		viper.Set("templates.storage_path", defaultPath)
		viper.WriteConfig()
	}
}

// Validates template names (only letters, numbers, dashes, and underscores)
func isValidTemplateName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched
}

// Saves .gitignore locally under the configured directory
func saveLocally(name, path string) error {
	configDir := viper.GetString("templates.storage_path")

	savePath := filepath.Join(configDir, name+".gitignore")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return os.WriteFile(savePath, content, 0644)
}
