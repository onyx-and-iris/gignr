package templates

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Type string

var (
	Directory Type = "dir"
	File      Type = "file"
)

type Template struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadUrl string `json:"download_url"`
	URL         string `json:"url"`
	Type        Type   `json:"type"`
	Source      string
	Content     []byte
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func FetchTemplates(url string) ([]Template, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	var templates []Template
	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return templates, nil
}

func FetchContent(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

func ProcessTemplates(url string, results chan<- Template, wg *sync.WaitGroup) {
	defer wg.Done()

	templates, err := FetchTemplates(url)
	if err != nil {
		log.Println(err)
		return
	}

	for _, tmpl := range templates {
		if strings.HasSuffix(tmpl.Name, ".gitignore") {
			tmpl.Source = detectSource(tmpl)
			results <- tmpl
		} else if tmpl.Type == Directory {
			wg.Add(1)
			go ProcessTemplates(tmpl.URL, results, wg)
		}
	}
}

func detectSource(tmpl Template) string {
	if strings.Contains(tmpl.URL, "github/gitignore") {
		if strings.Contains(tmpl.Path, "community") {
			return "GitHub Community"
		}
		if strings.Contains(tmpl.Path, "Global") {
			return "GitHub Global"
		}
		return "GitHub"
	} else if strings.Contains(tmpl.URL, "toptal/gitignore") {
		return "TopTal"
	}
	return "Unknown"
}
