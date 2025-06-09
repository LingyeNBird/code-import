package vcs

import (
	"ccrctl/pkg/api/gongfeng"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
	"strconv"
	"strings"
)

const (
	GongfengUserName = "private"
)

type GongfengVcs struct {
	httpURL   string
	RepoPath  string
	RepoName  string
	RepoType  string
	Private   bool
	ProjectId int
	Desc      string
}

// GetRepoPath 返回仓库路径
func (c *GongfengVcs) GetRepoPath() string {
	return c.RepoPath
}

// GetRepoName 返回仓库名称
func (c *GongfengVcs) GetRepoName() string {
	return c.RepoName
}

// GetSubGroup 返回子组信息
func (c *GongfengVcs) GetSubGroup() *SubGroup {
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

// GetRepoType 返回仓库类型
func (c *GongfengVcs) GetRepoType() string {
	return c.RepoType
}

// GetCloneUrl 返回克隆 URL
func (c *GongfengVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.httpURL, GongfengUserName, c.GetToken())
}

// GetUserName 返回用户名
func (c *GongfengVcs) GetUserName() string {
	return GongfengUserName
}

// GetToken 返回访问令牌
func (c *GongfengVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

// Clone 克隆仓库
func (c *GongfengVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

// GetRepoPrivate 返回仓库是否私有
func (c *GongfengVcs) GetRepoPrivate() bool {
	return c.Private
}

// GetReleases 返回仓库的发布列表（暂不实现）
func (c *GongfengVcs) GetReleases() []Releases {
	return nil
}

// GetProjectID 返回项目 ID
func (c *GongfengVcs) GetProjectID() string {
	return strconv.Itoa(c.ProjectId)
}

// GetReleaseAttachments 返回发布附件（暂不实现）
func (c *GongfengVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}

// GetRepoDescription 返回仓库描述
func (c *GongfengVcs) GetRepoDescription() string {
	return c.Desc
}

// ListRepos 列出所有仓库
func (c *GongfengVcs) ListRepos() ([]VCS, error) {
	return newGongfengRepo()
}

// newGongfengRepo 创建工蜂仓库实例
func newGongfengRepo() ([]VCS, error) {
	projects, err := gongfeng.GetProjects()
	if err != nil {
		return nil, err
	}
	return GongfengConvertToVcs(projects), nil
}

// GongfengConvertToVcs 将工蜂项目转换为 VCS 接口
func GongfengConvertToVcs(projects []gongfeng.Project) []VCS {
	var vcsList []VCS
	for _, project := range projects {
		vcsList = append(vcsList, &GongfengVcs{
			httpURL:   project.HTTPURL,
			RepoPath:  project.PathWithNS,
			RepoName:  project.Name,
			RepoType:  Git,
			Private:   project.IsPrivate(),
			ProjectId: project.ID,
			Desc:      project.Description,
		})
	}
	return vcsList
}
