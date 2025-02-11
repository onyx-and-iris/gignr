package templates

import (
	"github.com/spf13/viper"
)

// Load user-added repositories using urls from config
func LoadCustomRepositories() map[string]string {
	repos := viper.GetStringMapString("repositories")
	return repos
}
