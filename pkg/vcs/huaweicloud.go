// Package vcs 华为云CodeArts VCS实现
package vcs

import (
	"ccrctl/pkg/api/huaweicloud"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
	"fmt"
	"path"
	"strings"
)

const (
	huaweiUserName = "private-token"
)

// HuaweiCloudVcs 华为云CodeArts VCS实现
type HuaweiCloudVcs struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CloneURL    string `json:"clone_url"`
	SSHURL      string `json:"ssh_url"`
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	RepoPath    string `json:"repo_path"`
	SubGroup    string `json:"sub_group"`
	RepoType    string `json:"repo_type"`
	UserName    string `json:"user_name"`
	Private     bool   `json:"private"`
}

// ListRepos 获取华为云CodeArts仓库列表
func (c *HuaweiCloudVcs) ListRepos() ([]VCS, error) {
	return newHuaweiCloudRepo()
}

// GetCloneUrl 获取克隆URL
func (c *HuaweiCloudVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.CloneURL, huaweiUserName, c.GetToken())
}

// GetRepoName 获取仓库名称
func (c *HuaweiCloudVcs) GetRepoName() string {
	return c.Name
}

// GetDescription 获取仓库描述
func (c *HuaweiCloudVcs) GetDescription() string {
	return c.Description
}

// GetRepoPath 获取仓库路径
func (c *HuaweiCloudVcs) GetRepoPath() string {
	return c.RepoPath
}

// GetSubGroup 获取子组信息
func (c *HuaweiCloudVcs) GetSubGroup() *SubGroup {
	if c.SubGroup == "" {
		return nil
	}
	return &SubGroup{
		Name: c.SubGroup,
	}
}

// GetRepoType 获取仓库类型
func (c *HuaweiCloudVcs) GetRepoType() string {
	return Git
}

// GetUserName 获取用户名
func (c *HuaweiCloudVcs) GetUserName() string {
	return "huawei"
}

// GetRepoPrivate 获取仓库是否私有
func (c *HuaweiCloudVcs) GetRepoPrivate() bool {
	return c.Private
}

// GetRepoDescription 获取仓库描述（别名）
func (c *HuaweiCloudVcs) GetRepoDescription() string {
	return c.Description
}

// GetReleases 获取发布信息
func (c *HuaweiCloudVcs) GetReleases() []Releases {
	// 华为云CodeArts暂不支持Release功能，返回空列表
	return []Releases{}
}

// GetReleaseAttachments 获取发布附件
func (c *HuaweiCloudVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	// 华为云CodeArts暂不支持Release附件功能，返回空列表
	return []Attachment{}, nil
}

// GetToken 获取访问令牌
func (c *HuaweiCloudVcs) GetToken() string {
	// 华为云使用AK/SK认证，这里返回空字符串
	return config.Cfg.GetString("source.token")
}

// Clone 克隆仓库（VCS接口要求的方法）
func (c *HuaweiCloudVcs) Clone() error {
	// 这里可以实现具体的克隆逻辑
	// 目前返回 nil 表示成功
	return git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
}

func (c *HuaweiCloudVcs) GetProjectID() string {
	return c.ProjectID
}

// newHuaweiCloudRepo 创建华为云CodeArts仓库客户端并获取仓库列表
func newHuaweiCloudRepo() ([]VCS, error) {
	// 调用华为云API获取仓库列表
	response, err := huaweicloud.GetRepositories()
	if err != nil {
		return nil, fmt.Errorf("获取华为云CodeArts仓库列表失败: %w", err)
	}

	// 检查响应是否为空
	if response == nil {
		return []VCS{}, nil
	}

	var repos []VCS
	projectsMap, err := huaweicloud.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("获取华为云CodeArts项目列表失败: %w", err)
	}

	if response != nil {
		for _, repo := range response {
			var isPrivate bool
			// 0 私有 20 公开
			if *repo.VisibilityLevel == 20 {
				isPrivate = false
			} else {
				isPrivate = true
			}
			projectDisplayName := projectsMap[*repo.ProjectUuid]
			subGroup := path.Join(projectDisplayName, covertSubGroup(*repo.GroupName))
			vcsRepo := &HuaweiCloudVcs{
				Name:        *repo.RepositoryName,
				CloneURL:    *repo.HttpsUrl,
				Private:     isPrivate,
				RepoPath:    path.Join(subGroup, *repo.RepositoryName),
				SubGroup:    subGroup,
				ProjectName: projectDisplayName,
			}
			repos = append(repos, vcsRepo)
		}
	}

	return repos, nil
}

// covertSubGroup 适配代码组结构，去掉路径中最左边的部分
// e.g GroupName: "57fb765fa45f4a5b97073fc0c0df5a79/test" -> "test"
// e.g GroupName: "57fb765fa45f4a5b97073fc0c0df5a79/test/test1" -> "test/test1"
func covertSubGroup(groupName string) string {
	if groupName == "" {
		return ""
	}

	// 查找第一个 '/' 的位置
	index := strings.Index(groupName, "/")
	if index == -1 {
		// 如果没有找到 '/'，返回空字符串
		return ""
	}

	// 返回第一个 '/' 之后的部分
	return groupName[index+1:]
}
