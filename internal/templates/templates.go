package templates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// Type represents whether a GitHub object is a directory or file
type Type string

var (
	Directory Type = "dir"
	File      Type = "file"
)

// Template represents a .gitignore file
type Template struct {
	Name        string
	Path        string
	DownloadURL string
	Source      string
}

// HTTP client with timeout
var httpClient = &http.Client{Timeout: 10 * time.Second}

// GitHub client
var githubClient *github.Client

// Initialize GitHub Client (supports authentication if needed)
func InitGitHubClient(token string) {
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(context.Background(), ts)
		githubClient = github.NewClient(tc)
	} else {
		githubClient = github.NewClient(nil)
	}
}

// FetchTemplates lists `.gitignore` templates in a GitHub repository
func FetchTemplates(owner, repo, path string) ([]Template, error) {
	ctx := context.Background()

	// GetContents() can return a single file or a directory listing
	contents, dirContents, _, err := githubClient.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch templates: %w", err)
	}

	var templates []Template
	var wg sync.WaitGroup

	// Case 1: If `contents` is a single file
	if contents != nil {
		if strings.HasSuffix(contents.GetName(), ".gitignore") {
			templates = append(templates, Template{
				Name:        contents.GetName(),
				Path:        contents.GetPath(),
				DownloadURL: contents.GetDownloadURL(),
				Source:      detectSource(contents.GetPath()),
			})
		}
		return templates, nil
	}

	// Case 2: If `dirContents` is a list of files/directories
	for _, content := range dirContents {
		if content.GetType() == "file" && strings.HasSuffix(content.GetName(), ".gitignore") {
			templates = append(templates, Template{
				Name:        content.GetName(),
				Path:        content.GetPath(),
				DownloadURL: content.GetDownloadURL(),
				Source:      detectSource(content.GetPath()),
			})
		} else if content.GetType() == "dir" {
			// Recursively fetch templates from directories (Global, Community, etc.)
			wg.Add(1)
			go func(subPath string) {
				defer wg.Done()
				subTemplates, err := FetchTemplates(owner, repo, subPath)
				if err == nil {
					templates = append(templates, subTemplates...)
				}
			}(content.GetPath())
		}
	}
	wg.Wait()

	return templates, nil
}

// FetchContent gets the raw .gitignore file content
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

// detectSource determines where a template is from
func detectSource(path string) string {
	if strings.Contains(path, "community") {
		return "GitHub Community"
	}
	if strings.Contains(path, "Global") {
		return "GitHub Global"
	}
	return "GitHub"
}
