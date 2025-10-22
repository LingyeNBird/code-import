package vcs

import (
	"ccrctl/pkg/api/coding"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	CodingUserName = "coding"
)

type CodingVcs struct {
	httpURL      string
	sshURL       string
	RepoPath     string
	SubGroupName string
	RepoName     string
	RepoType     string
	Private      bool
	Desc         string
	id           int
}

func (c *CodingVcs) GetRepoPath() string {
	return c.RepoPath
}

func (c *CodingVcs) GetRepoName() string {
	return c.RepoName
}

func (c *CodingVcs) GetSubGroup() *SubGroup {
	project, err := coding.GetProjectByName(config.Cfg.GetString("source.url"), c.GetToken(), c.SubGroupName)
	if err != nil {
		logger.Logger.Errorf("获取项目信息失败: %v", err)
		os.Exit(1)
	}
	// 初始化描述和备注字段
	var desc, remark string
	// 如果配置了映射 Coding 项目描述，则使用项目的描述作为子组织描述
	if config.Cfg.GetBool("migrate.map_coding_description") {
		desc = strings.TrimSpace(project.Description)
	}
	// 如果配置了映射 Coding 项目显示名称，则使用项目的显示名称作为子组织别名
	if config.Cfg.GetBool("migrate.map_coding_display_name") {
		remark = strings.TrimSpace(project.DisplayName)
	}
	return &SubGroup{
		Name:   c.SubGroupName,
		Desc:   desc,
		Remark: remark,
	}
}

func (c *CodingVcs) GetRepoType() string {
	return c.RepoType
}

func (c *CodingVcs) GetCloneUrl() string {
	ssh := config.Cfg.GetBool("migrate.ssh")
	if ssh {
		return c.sshURL
	}
	return util.ConvertUrlWithAuth(c.httpURL, CodingUserName, c.GetToken())
}

func (c *CodingVcs) GetUserName() string {
	return CodingUserName
}

func (c *CodingVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (c *CodingVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (c *CodingVcs) GetRepoPrivate() bool {
	return c.Private
}

func (c *CodingVcs) GetReleases() []Releases {
	codingReleases, err := coding.GetReleasesList(c.id)
	if err != nil {
		logger.Logger.Errorf(err.Error())
		panic(err)
	}
	return covertToReleases(codingReleases)
}

func covertToReleases(codingReleases []coding.Releases) []Releases {
	var releases []Releases
	for _, codingRelease := range codingReleases {
		var assets []Asset
		if codingRelease.ReleaseAttachment != nil && len(codingRelease.ReleaseAttachment) > 0 {
			for _, attachment := range codingRelease.ReleaseAttachment {
				assets = append(assets, Asset{
					Name: attachment.AttachmentName,
					Url:  attachment.AttachmentDownloadUrl,
				})
			}
		}
		imgURLs := util.GetReleaseAttachmentMarkdownURL(codingRelease.Body)
		body := codingRelease.Body
		if len(imgURLs) > 0 && len(codingRelease.ImageDownloadUrl) > 0 {
			// 假设imgURLs和ImageDownloadUrl顺序对应，直接按索引替换
			for i, imgURL := range imgURLs {
				if i < len(codingRelease.ImageDownloadUrl) {
					body = strings.ReplaceAll(body, imgURL, codingRelease.ImageDownloadUrl[i])
				}
			}
		}
		releases = append(releases, Releases{
			TagName:    codingRelease.TagName,
			Name:       codingRelease.Title,
			Body:       body,
			Assets:     assets,
			Prerelease: codingRelease.Pre,
		})
	}
	return releases
}

func (c *CodingVcs) GetProjectID() string {
	return strconv.Itoa(0)
}

func newCodingRepo() ([]VCS, error) {
	repoList, err := coding.GetDepotList(config.Cfg.GetString("migrate.type"))
	if err != nil {
		return nil, err
	}
	return CodingCovertToVcs(repoList), nil
}

func CodingCovertToVcs(repoList []coding.Depots) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &CodingVcs{
			httpURL:      repo.HttpsUrl,
			sshURL:       repo.SshUrl,
			RepoPath:     repo.GetRepoPath(),
			SubGroupName: repo.ProjectName,
			RepoName:     repo.Name,
			RepoType:     repo.RepoType,
			Private:      !repo.IsShared,
			Desc:         repo.Description,
			id:           repo.Id,
		})
	}
	return VCS
}

func (c *CodingVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	// 转换release描述中的附件链接为cnb附件链接
	images, exists := util.CodingExtractAttachments(desc)
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

	// 处理图片附件
	for name, url := range images {
		if err := processAsset(name, url, "img"); err != nil {
			return nil, err
		}
	}

	return attachmentsList, nil
}

func (c *CodingVcs) GetRepoDescription() string {
	return c.Desc
}

func (c *CodingVcs) ListRepos() ([]VCS, error) {
	return newCodingRepo()
}
