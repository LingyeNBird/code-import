package vcs

import (
	api "ccrctl/pkg/api/aliyun"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
	"strings"

	devops20210625 "github.com/alibabacloud-go/devops-20210625/v5/client"
)

type AliyunVcs struct {
	httpURL           string
	PathWithNamespace string
	RepoName          string
	RepoType          string
	Private           string
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
	return config.Cfg.GetString("source.username")
}

func (c *AliyunVcs) GetToken() string {
	return config.Cfg.GetString("source.password")
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

func newAliyunRepo() []VCS {
	repoList := api.ListRepository(config.Cfg.GetString("source.organizationId"))
	return aliyunCovertToVcs(repoList)
}

func aliyunCovertToVcs(repoList []*devops20210625.ListRepositoriesResponseBodyResult) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &AliyunVcs{
			httpURL:           *repo.WebUrl,
			PathWithNamespace: *repo.PathWithNamespace,
			RepoName:          *repo.Name,
			RepoType:          Git,
			Private:           *repo.VisibilityLevel,
		})
	}
	return VCS
}

func (c *AliyunVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}

func (c *AliyunVcs) GetRepoDescription() string {
	return ""
}
