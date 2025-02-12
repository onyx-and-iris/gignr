package templates

import (
	"io"
	"net/http"
)

// GetTemplateContent retrieves a `.gitignore` template (cached or fresh)
func GetTemplateContent(url string) ([]byte, error) {
	if cacheContent, err := LoadCachedTemplateContent(url); err != nil {
		return cacheContent, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	SaveTemplateContentToCache(url, string(body))
	return body, nil
}
