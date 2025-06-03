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
	Desc      string
}

func (c *GithubVcs) GetRepoPath() string {
	return c.RepoPath
}

func (c *GithubVcs) GetSubGroup() *SubGroup {
	parts := strings.Split(c.RepoPath, "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1] // 去掉仓库名
	}
	result := strings.Join(parts, "/")
	return &SubGroup{
		Name:   result,
		Desc:   "",
		Remark: "",
	}
}

func (c *GithubVcs) GetRepoName() string {
	return c.RepoName
}

func (c *GithubVcs) GetRepoType() string {
	return c.RepoType
}

func (c *GithubVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.httpURL, c.GetUserName(), c.GetToken())
}

func (c *GithubVcs) GetUserName() string {
	return api.GetUserName()
}

func (c *GithubVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (c *GithubVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (c *GithubVcs) GetRepoPrivate() bool {
	return c.Private
}

func (c *GithubVcs) GetReleases() (cnbReleases []Releases) {
	parts := strings.Split(c.RepoPath, "/")
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
		if !c.Private {
			for _, asset := range githubRelease.Assets {
				assetName := ""
				if asset.Name != nil {
					assetName = *asset.Name
				}
				assetURL := ""
				if asset.BrowserDownloadURL != nil {
					assetURL = *asset.BrowserDownloadURL
				}
				assets = append(assets, Asset{
					Name: assetName,
					Url:  assetURL,
				})
			}
		}
		tagName := ""
		if githubRelease.TagName != nil {
			tagName = *githubRelease.TagName
		}
		name := ""
		if githubRelease.Name != nil {
			name = *githubRelease.Name
		}
		body := ""
		if githubRelease.Body != nil {
			body = *githubRelease.Body
		}
		makeLatest := ""
		if githubRelease.MakeLatest != nil {
			makeLatest = *githubRelease.MakeLatest
		}
		cnbReleases = append(cnbReleases, Releases{
			TagName:    tagName,
			Name:       name,
			Body:       body,
			Assets:     assets,
			Prerelease: githubRelease.Prerelease != nil && *githubRelease.Prerelease,
			Draft:      githubRelease.Draft != nil && *githubRelease.Draft,
			MakeLatest: makeLatest,
		})
	}
	return cnbReleases
}

func (c *GithubVcs) GetProjectID() string {
	return strconv.Itoa(c.ProjectId)
}

func newGithubRepo() ([]VCS, error) {
	repoList, err := api.GetRepos()
	if err != nil {
		return nil, err
	}
	return GithubCovertToVcs(repoList), nil
}

func GithubCovertToVcs(repoList []*github.Repository) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		desc := ""
		if repo.Description != nil {
			desc = *repo.Description
		}
		VCS = append(VCS, &GithubVcs{
			httpURL:   *repo.CloneURL,
			RepoPath:  *repo.FullName,
			RepoName:  *repo.Name,
			RepoType:  Git,
			Private:   *repo.Private,
			ProjectId: int(*repo.ID),
			Desc:      desc,
		})
	}
	return VCS
}

// GetReleaseAttachments Github release 描述里的普通附件需要鉴权，且未提供相关openAPI,因此无法迁移
func (c *GithubVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}

func (c *GithubVcs) GetRepoDescription() string {
	return c.Desc
}

func (c *GithubVcs) ListRepos() ([]VCS, error) {
	return newGithubRepo()
}
