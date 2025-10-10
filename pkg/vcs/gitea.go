package vcs

import (
	api "ccrctl/pkg/api/gitea"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"fmt"
	"strconv"
	"strings"
)

// GiteaVcs Gitea VCS 实现
type GiteaVcs struct {
	httpURL  string
	RepoPath string
	RepoName string
	RepoType string
	Private  bool
	Desc     string
}

func (c *GiteaVcs) GetRepoPath() string {
	return c.RepoPath
}

func (c *GiteaVcs) GetRepoName() string {
	return c.RepoName
}

func (c *GiteaVcs) GetSubGroup() *SubGroup {
	parts := strings.Split(c.GetRepoPath(), "/")
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

func (c *GiteaVcs) GetRepoType() string {
	return "Git"
}

func (c *GiteaVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.httpURL, c.GetUserName(), c.GetToken())
}

func (c *GiteaVcs) GetUserName() string {
	name, _ := api.GetUserName()
	return name
}

func (c *GiteaVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (c *GiteaVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (c *GiteaVcs) GetRepoPrivate() bool {
	return c.Private
}

func (c *GiteaVcs) GetReleases() (cnbReleases []Releases) {
	releases, err := api.GetReleases(c.RepoPath)
	if err != nil {
		logger.Logger.Errorf("获取 Gitea Release 失败: %v", err)
		return nil
	}

	for _, release := range releases {
		var assets []Asset
		// 处理 Release 附件
		for _, asset := range release.Assets {
			assets = append(assets, Asset{
				Name: asset.Name,
				Url:  asset.BrowserDownloadUrl,
			})
		}

		cnbReleases = append(cnbReleases, Releases{
			TagName:    release.TagName,
			Name:       release.Name,
			Body:       release.Body,
			Assets:     assets,
			Prerelease: release.IsPrerelease,
		})
	}
	return cnbReleases
}

func (c *GiteaVcs) GetProjectID() string {
	return strconv.Itoa(0)
}

// newGiteaRepo 创建 Gitea 仓库列表
func newGiteaRepo() ([]VCS, error) {
	repoList, err := api.GetRepoList()
	if err != nil {
		return nil, err
	}
	return GiteaCovertToVcs(repoList), nil
}

// GiteaCovertToVcs 将 Gitea 仓库转换为 VCS 接口
func GiteaCovertToVcs(repoList []api.Repo) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &GiteaVcs{
			httpURL:  repo.CloneUrl,
			RepoPath: repo.FullName,
			RepoName: repo.Name,
			RepoType: Git,
			Private:  repo.Private,
			Desc:     repo.Description,
		})
	}
	return VCS
}

func (c *GiteaVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	// 转换 release 描述中的附件链接为 CNB 附件链接
	attachments, images, exists := util.GiteaExtractAttachments(desc)
	if !exists {
		return nil, nil
	}

	var attachmentsList []Attachment

	// 统一处理附件和图片
	processAsset := func(name, url, assetType string) error {
		data, err := http_client.DownloadFromUrl(url)
		if err != nil {
			logger.Logger.Errorf("下载 Release 资源失败 类型:%s 名称:%s URL:%s 错误:%v",
				assetType, name, url, err)
			return fmt.Errorf("下载%s资源'%s'失败: %w", assetType, name, err)
		}

		attachmentsList = append(attachmentsList, Attachment{
			Name:     name,
			Data:     data,
			Url:      url,
			Type:     assetType,
			RepoPath: repoPath,
			Size:     len(data),
		})
		return nil
	}

	// 处理普通附件
	for name, url := range attachments {
		if err := processAsset(name, url, "file"); err != nil {
			return nil, err
		}
	}

	// 处理图片附件
	for name, url := range images {
		if err := processAsset(name, url, "img"); err != nil {
			return nil, err
		}
	}

	return attachmentsList, nil
}

func (c *GiteaVcs) GetRepoDescription() string {
	return c.Desc
}

func (c *GiteaVcs) ListRepos() ([]VCS, error) {
	return newGiteaRepo()
}
