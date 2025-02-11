package templates

import (
	"fmt"
	"time"

	"github.com/jasonuc/gignr/internal/cache"
)

// LoadCachedTemplates retrieves the list of available templates from cache
func LoadCachedTemplates(source string) ([]Template, error) {

	var cacheData struct {
		Updated   time.Time  `json:"updated"`
		Templates []Template `json:"templates"`
	}

	if err := cache.LoadCache(source, &cacheData); err != nil || cache.IsCacheExpired(cacheData.Updated) {
		return nil, fmt.Errorf("cache expired or missing")
	}

	return cacheData.Templates, nil
}

// SaveTemplatesToCache stores fetched templates in the cache
func SaveTemplatesToCache(cacheFile string, templates []Template) {
	cacheData := struct {
		Updated   time.Time  `json:"updated"`
		Templates []Template `json:"templates"`
	}{
		Updated:   time.Now(),
		Templates: templates,
	}

	cache.SaveCache(cacheFile, cacheData)
}

// LoadCachedTemplateContent retrieves the cached content of a specific template
func LoadCachedTemplateContent(url string) (string, error) {
	var cacheData = make(map[string]struct {
		Updated time.Time `json:"updated"`
		Content string    `json:"content"`
	})

	if err := cache.LoadCache("template-content.json", &cacheData); err != nil {
		return "", err
	}

	if entry, exists := cacheData[url]; exists && !cache.IsCacheExpired(entry.Updated) {
		return entry.Content, nil
	}

	return "", fmt.Errorf("content not found in cache")
}

// SaveTemplateContentToCache stores fetched `.gitignore` content in the cache
func SaveTemplateContentToCache(url, content string) {
	var cacheData = make(map[string]struct {
		Updated time.Time `json:"updated"`
		Content string    `json:"content"`
	})

	cache.LoadCache("template-content.json", &cacheData)

	cacheData[url] = struct {
		Updated time.Time `json:"updated"`
		Content string    `json:"content"`
	}{
		Updated: time.Now(),
		Content: content,
	}

	cache.SaveCache("template-content.json", cacheData)
}
