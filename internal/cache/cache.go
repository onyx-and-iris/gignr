package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type TemplateContentCache struct {
	Updated time.Time `json:"updated"`
	Content string    `json:"content"`
}

func GetCacheDir() string {
	configDir := filepath.Dir(viper.ConfigFileUsed())
	if configDir == "." || configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configDir = filepath.Join(home, ".config", "gignr")
	}
	cacheDir := filepath.Join(configDir, "cache")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		panic(err)
	}

	return cacheDir
}

func LoadCache(fileName string, target interface{}) error {
	cachePath := filepath.Join(GetCacheDir(), fileName)
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func SaveCache(fileName string, data interface{}) error {
	cachePath := filepath.Join(GetCacheDir(), fileName)
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath, content, 0644)
}

// IsCacheExpired checks if a cache entry is expired (TTL: 2 weeks)
func IsCacheExpired(updatedTime time.Time) bool {
	return time.Since(updatedTime) > 14*24*time.Hour
}
