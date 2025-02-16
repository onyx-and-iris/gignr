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

type Type string

const (
	TypeDirectory Type = "dir"
	TypeFile      Type = "file"
)

type Template struct {
	Name        string
	Path        string
	DownloadURL string
	Source      string
}

var githubClient *github.Client

func InitGitHubClient(token string) {
	if token == "" {
		githubClient = github.NewClient(nil)
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	githubClient = github.NewClient(tc)
}

func FetchTemplates(owner, repo, path, sourceID string) ([]Template, error) {
	if templates, err := LoadCachedTemplates(getCacheFileName(owner, sourceID)); err == nil {
		return templates, nil
	}

	templates, err := fetchFromGitHub(owner, repo, path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch templates: %w", err)
	}

	SaveTemplatesToCache(getCacheFileName(owner, sourceID), templates)
	return templates, nil
}

func getCacheFileName(owner, _ string) string {
	switch owner {
	case "github":
		return "github.json"
	case "toptal":
		return "toptal.json"
	default:
		return fmt.Sprintf("%s.json", owner)
	}
}

func fetchFromGitHub(owner, repo, path string) ([]Template, error) {
	ctx := context.Background()
	contents, dirContents, _, err := githubClient.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, err
	}

	var templates []Template
	if contents != nil {
		templates = handleSingleFile(contents)
	} else {
		templates = handleDirectory(owner, repo, dirContents)
	}

	return templates, nil
}

func handleSingleFile(content *github.RepositoryContent) []Template {
	if !strings.HasSuffix(content.GetName(), ".gitignore") {
		return nil
	}

	return []Template{{
		Name:        content.GetName(),
		Path:        content.GetPath(),
		DownloadURL: content.GetDownloadURL(),
		Source:      utils.DetectSource(content.GetURL()),
	}}
}

func handleDirectory(owner, repo string, contents []*github.RepositoryContent) []Template {
	var templates []Template
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, content := range contents {
		if content.GetType() == "file" {
			if tmpl := handleSingleFile(content); tmpl != nil {
				templates = append(templates, tmpl...)
			}
			continue
		}

		if content.GetType() == "dir" {
			wg.Add(1)
			go func(subPath string) {
				defer wg.Done()
				subTemplates, err := fetchFromGitHub(owner, repo, subPath)
				if err != nil {
					return
				}
				mu.Lock()
				templates = append(templates, subTemplates...)
				mu.Unlock()
			}(content.GetPath())
		}
	}

	wg.Wait()
	return templates
}
