package vcs

import (
	api "ccrctl/pkg/api/aliyun"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
	"strings"
)

type AliyunVcs struct {
	httpURL           string
	PathWithNamespace string
	RepoName          string
	RepoType          string
	Private           string
	Desc              string
}

func (c *AliyunVcs) GetRepoPath() string {
	return c.PathWithNamespace
}

func (c *AliyunVcs) GetSubGroup() *SubGroup {
	parts := strings.Split(c.PathWithNamespace, "/")
	if len(parts) > 0 {
		parts = parts[1 : len(parts)-1] // 去掉仓库名
	}
	result := strings.Join(parts, "/")
	return &SubGroup{
		Name:   result,
		Desc:   "",
		Remark: "",
	}
}

func (c *AliyunVcs) GetRepoName() string {
	return c.RepoName
}

func (c *AliyunVcs) GetRepoType() string {
	return c.RepoType
}

func (c *AliyunVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.httpURL, c.GetUserName(), c.GetToken()) + ".git"
}

func (c *AliyunVcs) GetUserName() string {
	return "aliyun"
}

func (c *AliyunVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (c *AliyunVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (c *AliyunVcs) GetRepoPrivate() bool {
	return true
}

func (c *AliyunVcs) GetReleases() (cnbReleases []Releases) {
	return nil
}

func (c *AliyunVcs) GetProjectID() string {
	return ""
}

func newAliyunRepo() ([]VCS, error) {
	repoList, err := api.GetAllRepositories()
	if err != nil {
		return nil, err
	}
	return aliyunCovertToVcs(repoList), nil
}

func aliyunCovertToVcs(repoList []api.Repository) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &AliyunVcs{
			httpURL:           repo.WebUrl,
			PathWithNamespace: repo.PathWithNamespace,
			RepoName:          repo.Path,
			RepoType:          Git,
			Private:           repo.Visibility,
			Desc:              repo.Description,
		})
	}
	return VCS
}

func (c *AliyunVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}

func (c *AliyunVcs) GetRepoDescription() string {
	return c.Desc
}

func (c *AliyunVcs) ListRepos() ([]VCS, error) {
	return newAliyunRepo()
}
