package vcs

import (
	api "ccrctl/pkg/api/github"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
	"github.com/google/go-github/v66/github"
	"strconv"
	"strings"
)

type GithubVcs struct {
	httpURL  string
	RepoPath string
	RepoName string
	RepoType string
	Private  bool
}

func (g *GithubVcs) GetRepoPath() string {
	return g.RepoPath
}

func (g *GithubVcs) GetSubGroupName() string {
	parts := strings.Split(g.RepoPath, "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1] // 去掉仓库名
	}
	result := strings.Join(parts, "/")
	return result
}

func (g *GithubVcs) GetRepoName() string {
	return g.RepoName
}

func (g *GithubVcs) GetRepoType() string {
	return g.RepoType
}

func (g *GithubVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(g.httpURL, g.GetUserName(), g.GetToken())
}

func (g *GithubVcs) GetUserName() string {
	return api.GetUserName()
}

func (g *GithubVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (g *GithubVcs) Clone() error {
	err := git.Clone(g.GetCloneUrl(), g.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (g *GithubVcs) GetRepoPrivate() bool {
	return g.Private
}

func (g *GithubVcs) GetReleases() (cnbReleases []releases) {
	return nil
}

func (g *GithubVcs) GetProjectID() string {
	return strconv.Itoa(0)
}

func newGithubRepo() []VCS {
	repoList, err := api.GetRepos()
	if err != nil {
		panic(err)
	}
	return GithubCovertToVcs(repoList)
}

func GithubCovertToVcs(repoList []*github.Repository) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &GithubVcs{
			httpURL:  *repo.CloneURL,
			RepoPath: *repo.FullName,
			RepoName: *repo.Name,
			RepoType: Git,
			Private:  *repo.Private,
		})
	}
	return VCS
}
