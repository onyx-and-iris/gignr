package utils

import (
	"fmt"
	"strings"
)

func ConvertToAPIURL(repoURL string) (string, error) {
	if !strings.HasPrefix(repoURL, "https://github.com/") {
		return "", fmt.Errorf("invalid GitHub URL: must be in the format 'https://github.com/{username}/{repoName}'")
	}

	apiURL := strings.Replace(repoURL, "https://github.com/", "https://api.github.com/repos/", 1) + "/contents"

	return apiURL, nil
}
