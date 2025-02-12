package utils

import (
	"strings"

	"github.com/spf13/viper"
)

// Extract owner and repo from a GitHub URL
func ExtractRepoDetails(url string) (string, string) {
	if strings.Contains(url, "https://api.github.com/repos/") {
		url = strings.Replace(url, "https://api.github.com/repos/", "", 1)
	}
	if strings.Contains(url, "https://github.com/") {
		url = strings.Replace(url, "https://github.com/", "", 1)
	}

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", ""
	}

	owner := parts[0]
	repo := parts[1]

	return owner, repo
}

// MatchRepoURL checks if a given path belongs to a user-added repository
func MatchRepoURL(repoURL, path string) bool {
	owner, repo := ExtractRepoDetails(repoURL)
	return strings.Contains(path, owner+"/"+repo)
}

func DetectSource(path string) string {
	owner, repo := ExtractRepoDetails(path)

	fullRepoPath := owner + "/" + repo

	if fullRepoPath == "github/gitignore" {
		if strings.Contains(path, "/community/") {
			return "GitHub Community"
		}
		if strings.Contains(path, "/Global/") {
			return "GitHub Global"
		}
		return "GitHub"
	}

	if fullRepoPath == "toptal/gitignore" {
		return "TopTal"
	}

	// Check user-added repositories
	repos := viper.GetStringMapString("repositories")

	for nickname, repoURL := range repos {
		if MatchRepoURL(repoURL, path) {
			return "Custom (" + nickname + ")"
		}
	}

	return "Unknown"
}
