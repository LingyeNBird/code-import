package gitlab

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	"github.com/xanzy/go-gitlab"
)

var (
	Git *gitlab.Client
)

func init() {
	var err error
	Git, err = gitlab.NewClient(config.Cfg.GetString("source.token"), gitlab.WithBaseURL(config.Cfg.GetString("source.url")))
	if err != nil {
		logger.Logger.Fatalf("Failed to create Gitlab client: %v", err)
	}
}

func GetProjects() ([]*gitlab.Project, error) {
	var Projects []*gitlab.Project
	page := 1
	for {
		projects, resp, err := Git.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    page,
			},
			Owned: gitlab.Bool(true),
		})
		if err != nil {
			logger.Logger.Fatalf("Failed to get Projects: %v", err)
		}
		Projects = append(Projects, projects...)
		if resp.NextPage == 0 {
			break
		}
		page++
	}
	return Projects, nil
}
