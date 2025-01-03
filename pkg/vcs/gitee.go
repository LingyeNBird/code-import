package vcs

import (
	api "ccrctl/pkg/api/gitee"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
	"strconv"
	"strings"
)

type GiteeVcs struct {
	httpURL  string
	RepoPath string
	RepoName string
	RepoType string
	Private  bool
}

func (r *GiteeVcs) GetRepoPath() string {
	return r.RepoPath
}

func (r *GiteeVcs) GetRepoName() string {
	return r.RepoName
}

func (r *GiteeVcs) GetSubGroupName() string {
	parts := strings.Split(r.GetRepoPath(), "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1] // 去掉仓库名
	}
	result := strings.Join(parts, "/")
	return result
}

func (r *GiteeVcs) GetRepoType() string {
	return "Git"
}

func (r *GiteeVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(r.httpURL, r.GetUserName(), r.GetToken())
}

func (r *GiteeVcs) GetUserName() string {
	name, _ := api.GetUserName()
	return name
}

func (r *GiteeVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (r *GiteeVcs) Clone() error {
	err := git.Clone(r.GetCloneUrl(), r.GetRepoPath())
	if err != nil {
		return err
	}
	return nil
}

func (r *GiteeVcs) GetRepoPrivate() bool {
	return r.Private
}

func (r *GiteeVcs) GetReleases() []releases {
	return nil
}

func (r *GiteeVcs) GetProjectID() string {
	return strconv.Itoa(0)
}

func newGiteeRepo() []VCS {
	repoList, err := api.GetRepoList()
	if err != nil {
		panic(err)
	}
	return GiteeCovertToVcs(repoList)
}

func GiteeCovertToVcs(repoList []api.Repo) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &GiteeVcs{
			httpURL:  repo.HtmlUrl,
			RepoPath: repo.FullName,
			RepoName: repo.Name,
			RepoType: Git,
			Private:  repo.Private,
		})
	}
	return VCS
}
