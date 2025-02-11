package templates

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetLocalTemplate(name string) ([]byte, error) {
	storagePath := viper.GetString("templates.storage_path")
	if storagePath == "" {
		storagePath = filepath.Join(os.Getenv("HOME"), ".config/gignr/templates")
	}

	templatePath := filepath.Join(storagePath, name+".gitignore")
	return os.ReadFile(templatePath)
}
