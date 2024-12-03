package github

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	"context"
	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

var client *github.Client

func init() {
	if config.Cfg.GetString("source.platform") == "github" {
		var err error
		// 创建一个上下文
		ctx := context.Background()

		// 替换为你的GitHub访问令牌
		token := config.Cfg.GetString("source.token")

		// 创建一个OAuth2客户端
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		// 创建一个GitHub客户端
		client = github.NewClient(tc)
		if err != nil {
			logger.Logger.Fatalf("Failed to create Github client: %v", err)
		}
	}
}

func GetRepos() ([]*github.Repository, error) {
	// 创建一个上下文
	ctx := context.Background()

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10}, // 每页显示30个仓库
	}
	var allRepos []*github.Repository

	for {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			logger.Logger.Fatalf("Failed to list repositories: %v", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func GetUserName() string {
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		logger.Logger.Fatalf("Failed to get github user: %v", err)
	}
	return user.GetName()
}
