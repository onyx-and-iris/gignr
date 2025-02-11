package templates

import (
	"fmt"

	"github.com/jasonuc/gignr/internal/utils"
	"github.com/spf13/viper"
)

func FetchUserTemplates() []Template {
	userRepos := viper.GetStringMapString("repositories")
	var userTemplates []Template

	for nickname, repoURL := range userRepos {
		owner, repo := utils.ExtractRepoDetails(repoURL)
		repoTemplates, err := FetchTemplates(owner, repo, "")
		if err != nil {
			continue
		}
		for _, tmpl := range repoTemplates {
			tmpl.Name = fmt.Sprintf("%s:%s", nickname, tmpl.Name)
			userTemplates = append(userTemplates, tmpl)
		}
	}

	return userTemplates
}
