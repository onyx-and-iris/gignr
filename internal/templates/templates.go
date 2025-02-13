package templates

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/go-github/v57/github"
	"github.com/jasonuc/gignr/internal/utils"
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

// FetchTemplates lists `.gitignore` templates, first checking the cache
func FetchTemplates(owner, repo, path string) ([]Template, error) {
	cacheFile := fmt.Sprintf("%s.json", owner) // Use owner name to determine cache file

	// Check the cache first
	cachedTemplates, err := LoadCachedTemplates(cacheFile)

	if err == nil {
		return cachedTemplates, nil
	}

	// If cache is missing or expired, fetch from GitHub
	ctx := context.Background()
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
				Source:      utils.DetectSource(contents.GetURL()),
			})
		}
	} else {
		// Case 2: If `dirContents` is a list of files/directories
		for _, content := range dirContents {
			if content.GetType() == "file" && strings.HasSuffix(content.GetName(), ".gitignore") {
				templates = append(templates, Template{
					Name:        content.GetName(),
					Path:        content.GetPath(),
					DownloadURL: content.GetDownloadURL(),
					Source:      utils.DetectSource(content.GetURL()),
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
	}

	// Save the fetched templates to the cache to avoid overloading the GitHub API
	SaveTemplatesToCache(cacheFile, templates)
	return templates, nil
}
