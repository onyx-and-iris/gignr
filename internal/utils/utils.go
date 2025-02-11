package utils

import "strings"

// Extract owner and repo from a GitHub URL
func ExtractRepoDetails(url string) (string, string) {
	parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

// MatchRepoURL checks if a given path belongs to a user-added repository
func MatchRepoURL(repoURL, path string) bool {
	owner, repo := ExtractRepoDetails(repoURL)
	return strings.Contains(path, owner+"/"+repo)
}
