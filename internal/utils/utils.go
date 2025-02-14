package utils

import (
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func ExtractRepoDetails(url string) (owner, repo string, ok bool) {
	if strings.Contains(url, "https://api.github.com/repos/") {
		url = strings.Replace(url, "https://api.github.com/repos/", "", 1)
	}
	if strings.Contains(url, "https://github.com/") {
		url = strings.Replace(url, "https://github.com/", "", 1)
	}

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", false
	}

	owner = parts[0]
	repo = parts[1]

	repo = strings.TrimSuffix(repo, ".git")

	if strings.Contains(url, "blob/") || strings.Contains(parts[len(parts)-1], ".") {
		return "", "", false
	}
	if owner == "" || repo == "" {
		return "", "", false
	}
	return owner, repo, true
}

// MatchRepoURL checks if a given path belongs to a user-added repository
func MatchRepoURL(repoURL, path string) bool {
	owner, repo, ok := ExtractRepoDetails(repoURL)
	if !ok {
		return false
	}
	return strings.Contains(path, owner+"/"+repo)
}

func DetectSource(path string) string {
	owner, repo, ok := ExtractRepoDetails(path)
	if !ok {
		return "Unknown"
	}

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

	repos := viper.GetStringMapString("repositories")

	for nickname, repoURL := range repos {
		if MatchRepoURL(repoURL, path) {
			return "Custom (" + nickname + ")"
		}
	}

	return "Unknown"
}

func IsValidGitHubURL(url string) bool {
	if len(url) > 19 && url[:19] == "https://github.com/" {
		resp, err := http.Get(url)
		if err != nil {
			PrintWarning("Warning: Could not verify repository (network error). Assuming invalid.")
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == http.StatusOK
	}
	return false
}

func IsReservedNickname(nickname string) bool {
	reserved := map[string]bool{"gh": true, "ghc": true, "ghg": true, "tt": true}
	return reserved[strings.ToLower(nickname)]
}

func IsValidNickname(nickname string) bool {
	if nickname == "" {
		return false
	}
	if strings.Contains(nickname, " ") {
		return false
	}
	return true
}
