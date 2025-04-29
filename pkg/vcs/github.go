package vcs

import (
	api "ccrctl/pkg/api/github"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
	"strconv"
	"strings"

	"github.com/google/go-github/v66/github"
)

type GithubVcs struct {
	httpURL   string
	RepoPath  string
	RepoName  string
	RepoType  string
	Private   bool
	ProjectId int
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

func (g *GithubVcs) GetReleases() (cnbReleases []Releases) {
	parts := strings.Split(g.RepoPath, "/")
	if len(parts) != 2 {
		return nil
	}
	owner := parts[0]
	repo := parts[1]

	githubReleases, err := api.GetReleases(owner, repo)
	if err != nil {
		panic(err)
	}
	for _, githubRelease := range githubReleases {
		var assets []Asset
		// 私有仓库用户自定义上传的附件无法获取到，因此只在公开仓库获取附件
		if !g.Private {
			for _, asset := range githubRelease.Assets {
				assets = append(assets, Asset{
					Name: *asset.Name,
					Url:  *asset.BrowserDownloadURL,
				})
			}
		}
		cnbReleases = append(cnbReleases, Releases{
			TagName:    *githubRelease.TagName,
			Name:       *githubRelease.Name,
			Body:       *githubRelease.Body,
			Assets:     assets,
			Prerelease: *githubRelease.Prerelease,
			Draft:      *githubRelease.Draft,
			MakeLatest: *githubRelease.MakeLatest,
		})
	}
	return cnbReleases
}

func (g *GithubVcs) GetProjectID() string {
	return strconv.Itoa(g.ProjectId)
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
			httpURL:   *repo.CloneURL,
			RepoPath:  *repo.FullName,
			RepoName:  *repo.Name,
			RepoType:  Git,
			Private:   *repo.Private,
			ProjectId: int(*repo.ID),
		})
	}
	return VCS
}

// GetReleaseAttachments Github release 描述里的普通附件需要鉴权，且未提供相关openAPI,因此无法迁移
func (g *GithubVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}
