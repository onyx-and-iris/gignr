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
	cacheDir := filepath.Join(configDir, "cache")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		panic(err)
	}

	return cacheDir
}

// LoadCache reads a JSON cache file and unmarshals it into a given structure
func LoadCache(fileName string, target interface{}) error {
	cachePath := filepath.Join(GetCacheDir(), fileName)
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// SaveCache writes a given structure into a JSON cache file
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
