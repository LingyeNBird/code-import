package vcs

import (
	api "ccrctl/pkg/api/gitee"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"fmt"
	"strconv"
	"strings"
)

type GiteeVcs struct {
	httpURL  string
	RepoPath string
	RepoName string
	RepoType string
	Private  bool
	Desc     string
}

func (c *GiteeVcs) GetRepoPath() string {
	return c.RepoPath
}

func (c *GiteeVcs) GetRepoName() string {
	return c.RepoName
}

func (c *GiteeVcs) GetSubGroup() *SubGroup {
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

func (c *GiteeVcs) GetRepoType() string {
	return "Git"
}

func (c *GiteeVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.httpURL, c.GetUserName(), c.GetToken())
}

func (c *GiteeVcs) GetUserName() string {
	name, _ := api.GetUserName()
	return name
}

func (c *GiteeVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (c *GiteeVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (c *GiteeVcs) GetRepoPrivate() bool {
	return c.Private
}

func (c *GiteeVcs) GetReleases() (cnbReleases []Releases) {
	releases, err := api.GetReleases(c.RepoPath)
	if err != nil {
		panic(err)
	}
	for _, release := range releases {
		var assets []Asset
		// 私有仓库用户自定义上传的附件无法获取到，因此只在公开仓库获取附件
		if !c.Private {
			for index, asset := range release.Assets {
				// 跳过最后2个gitee自带的附件
				if index == len(release.Assets)-2 {
					break
				}
				assets = append(assets, Asset{
					Name: asset.Name,
					Url:  asset.BrowserDownloadUrl,
				})
			}
		}
		cnbReleases = append(cnbReleases, Releases{
			TagName:    release.TagName,
			Name:       release.Name,
			Body:       release.Body,
			Assets:     assets,
			Prerelease: release.Prerelease,
		})
	}
	return cnbReleases
}

func (c *GiteeVcs) GetProjectID() string {
	return strconv.Itoa(0)
}

func newGiteeRepo() ([]VCS, error) {
	repoList, err := api.GetRepoList()
	if err != nil {
		return nil, err
	}
	return GiteeCovertToVcs(repoList), nil
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
			Desc:     repo.Description,
		})
	}
	return VCS
}

func (c *GiteeVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	// 转换release描述中的附件链接为cnb附件链接
	attachments, images, exists := util.GiteeExtractAttachments(desc)
	if !exists {
		return nil, nil
	}

	var attachmentsList []Attachment

	// 统一处理附件和图片
	processAsset := func(name, url, assetType string) error {
		data, err := http_client.DownloadFromUrl(url)
		if err != nil {
			logger.Logger.Errorf("下载release资源失败 类型:%s 名称:%s URL:%s 错误:%v",
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

func (c *GiteeVcs) GetRepoDescription() string {
	return c.Desc
}

func (c *GiteeVcs) ListRepos() ([]VCS, error) {
	return newGiteeRepo()
}
