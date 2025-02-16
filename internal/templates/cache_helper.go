package templates

import (
	"fmt"
	"time"

	"github.com/jasonuc/gignr/internal/cache"
)

type TemplatesCache struct {
	Updated   time.Time  `json:"updated"`
	Templates []Template `json:"templates"`
}

// LoadCachedTemplates retrieves the list of available templates from cache
func LoadCachedTemplates(source string) ([]Template, error) {

	var cacheData TemplatesCache

	if err := cache.LoadCache(source, &cacheData); err != nil || cache.IsCacheExpired(cacheData.Updated) {
		return nil, fmt.Errorf("cache expired or missing")
	}

	return cacheData.Templates, nil
}

// SaveTemplatesToCache stores fetched templates in the cache
func SaveTemplatesToCache(cacheFile string, newTemplates []Template) {
	var existingCache TemplatesCache
	cache.LoadCache(cacheFile, &existingCache)

	if !cache.IsCacheExpired(existingCache.Updated) {
		existing := make(map[string]Template)
		for _, t := range existingCache.Templates {
			existing[t.Path] = t
		}

		for _, t := range newTemplates {
			existing[t.Path] = t
		}

		var merged []Template
		for _, t := range existing {
			merged = append(merged, t)
		}
		newTemplates = merged
	}

	cacheData := TemplatesCache{
		Updated:   time.Now(),
		Templates: newTemplates,
	}

	cache.SaveCache(cacheFile, cacheData)
}

// LoadCachedTemplateContent retrieves the cached content of a specific template
func LoadCachedTemplateContent(url string) ([]byte, error) {
	var cacheData = make(map[string]cache.TemplateContentCache)

	if err := cache.LoadCache("template-content.json", &cacheData); err == nil {
		if entry, exists := cacheData[url]; exists && !cache.IsCacheExpired(entry.Updated) {
			return []byte(entry.Content), nil
		}
	}

	return nil, fmt.Errorf("content not found in cache")
}

// SaveTemplateContentToCache stores fetched `.gitignore` content in the cache
func SaveTemplateContentToCache(url, content string) {
	var cacheData = make(map[string]cache.TemplateContentCache)

	cache.LoadCache("template-content.json", &cacheData)

	cacheData[url] = cache.TemplateContentCache{
		Updated: time.Now(),
		Content: content,
	}

	cache.SaveCache("template-content.json", cacheData)
}
