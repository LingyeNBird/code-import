package vcs

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"strings"
)

type CommonVcs struct {
	httpURL  string
	RepoPath string
	RepoName string
	RepoType string
}

func (c *CommonVcs) GetRepoPath() string {
	return c.RepoPath
}

func (c *CommonVcs) GetSubGroupName() string {
	parts := strings.Split(c.RepoPath, "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1] // 去掉仓库名
	}
	result := strings.Join(parts, "/")
	return result
}

func (c *CommonVcs) GetRepoName() string {
	return c.RepoName
}

func (c *CommonVcs) GetRepoType() string {
	return c.RepoType
}

func (c *CommonVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.httpURL, c.GetUserName(), c.GetToken())
}

func (c *CommonVcs) GetUserName() string {
	return config.Cfg.GetString("source.username")
}

func (c *CommonVcs) GetToken() string {
	return config.Cfg.GetString("source.password")
}

func (c *CommonVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath())
	if err != nil {
		return err
	}
	return nil
}

func (c *CommonVcs) GetRepoPrivate() bool {
	return true
}

func newCommonRepo() []VCS {
	VCS, _ := getRepos()
	return VCS
}

func getRepos() ([]VCS, error) {
	var VCS []VCS
	repos := config.Cfg.GetStringSlice("source.repo")
	for _, repo := range repos {
		VCS = append(VCS, &CommonVcs{
			httpURL:  config.Cfg.GetString("source.url") + "/" + repo + ".git",
			RepoPath: repo,
			RepoName: strings.Split(repo, "/")[len(strings.Split(repo, "/"))-1],
			RepoType: Git,
		})
	}
	logger.Logger.Debugw("获取到的仓库列表", "VCS", VCS)
	return VCS, nil
}
