package vcs

import (
	api "ccrctl/pkg/api/gitlab"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
	"github.com/xanzy/go-gitlab"
	"strconv"
	"strings"
)

const (
	GitlabUserName = "gitlab"
	Git            = "git"
)

type GitlabVcs struct {
	httpURL   string
	RepoPath  string
	RepoName  string
	RepoType  string
	Private   string
	ProjectId int
}

func (g *GitlabVcs) GetRepoPath() string {
	return g.RepoPath
}

func (g *GitlabVcs) GetSubGroupName() string {
	parts := strings.Split(g.RepoPath, "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1] // 去掉仓库名
	}
	result := strings.Join(parts, "/")
	return result
}

func (g *GitlabVcs) GetRepoName() string {
	return g.RepoName
}

func (g *GitlabVcs) GetRepoType() string {
	return g.RepoType
}

func (g *GitlabVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(g.httpURL, GitlabUserName, g.GetToken())
}

func (g *GitlabVcs) GetUserName() string {
	return GitlabUserName
}

func (g *GitlabVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (g *GitlabVcs) Clone() error {
	err := git.Clone(g.GetCloneUrl(), g.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (g *GitlabVcs) GetRepoPrivate() bool {
	switch g.Private {
	case "private":
		return true
	case "internal":
		return true
	}
	return false
}

func (g *GitlabVcs) GetReleases() (cnbReleases []releases) {
	gitlabReleases, err := api.GetReleases(g.ProjectId)
	if err != nil {
		panic(err)
	}
	for _, gitlabRelease := range gitlabReleases {
		var assets []Asset
		for _, link := range gitlabRelease.Assets.Links {
			assets = append(assets, Asset{
				Name: link.Name,
				Url:  link.URL,
			})
		}
		cnbReleases = append(cnbReleases, releases{
			TagName:    gitlabRelease.TagName,
			Name:       gitlabRelease.Name,
			Body:       gitlabRelease.Description,
			Assets:     assets,
			Prerelease: gitlabRelease.UpcomingRelease,
		})
	}
	return cnbReleases
}

func (g *GitlabVcs) GetProjectID() string {
	return strconv.Itoa(g.ProjectId)
}

func newGitlabRepo() []VCS {
	repoList, err := api.GetProjects()
	if err != nil {
		panic(err)
	}
	return GitlabCovertToVcs(repoList)
}

func GitlabCovertToVcs(repoList []*gitlab.Project) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &GitlabVcs{
			httpURL:   repo.HTTPURLToRepo,
			RepoPath:  repo.PathWithNamespace,
			RepoName:  repo.Name,
			RepoType:  Git,
			Private:   string(repo.Visibility),
			ProjectId: repo.ID,
		})
	}
	return VCS
}
