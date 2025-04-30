package vcs

import (
	api "ccrctl/pkg/api/gitlab"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"strconv"
	"strings"

	"github.com/xanzy/go-gitlab"
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
	Desc      string
}

func (c *GitlabVcs) GetRepoPath() string {
	return c.RepoPath
}

func (c *GitlabVcs) GetSubGroup() *SubGroup {
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

func (c *GitlabVcs) GetRepoName() string {
	return c.RepoName
}

func (c *GitlabVcs) GetRepoType() string {
	return c.RepoType
}

func (c *GitlabVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.httpURL, GitlabUserName, c.GetToken())
}

func (c *GitlabVcs) GetUserName() string {
	return GitlabUserName
}

func (c *GitlabVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (c *GitlabVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (c *GitlabVcs) GetRepoPrivate() bool {
	switch c.Private {
	case "private":
		return true
	case "internal":
		return true
	}
	return false
}

func (c *GitlabVcs) GetReleases() (cnbReleases []Releases) {
	gitlabReleases, err := api.GetReleases(c.ProjectId)
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
		cnbReleases = append(cnbReleases, Releases{
			TagName:    gitlabRelease.TagName,
			Name:       gitlabRelease.Name,
			Body:       gitlabRelease.Description,
			Assets:     assets,
			Prerelease: gitlabRelease.UpcomingRelease,
		})
	}
	return cnbReleases
}

func (c *GitlabVcs) GetProjectID() string {
	return strconv.Itoa(c.ProjectId)
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
			Desc:      repo.Description,
		})
	}
	return VCS
}

func (c *GitlabVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	// 转换release描述中的附件链接为cnb附件链接
	attachments, images, exists := util.ExtractAttachments(desc)
	if !exists {
		return nil, nil
	}
	var attachmentsList []Attachment
	for attachmentName, attachmentUrl := range attachments {
		uploadFiles, err := api.ListUploads(projectID)
		if err != nil {
			return nil, err
		}
		fileID, ok := uploadFiles[attachmentName]
		if !ok {
			logger.Logger.Warnf("%s 附件 %s 不存在", repoPath, attachmentName)
			continue
		}
		data, err := api.DownloadFile(projectID, fileID)
		if err != nil {
			logger.Logger.Errorf("%s 下载release asset %s 失败: %s", attachmentUrl, attachmentName, err)
			return nil, err
		}
		attachmentsList = append(attachmentsList, Attachment{
			Name:     attachmentName,
			Data:     data,
			Url:      attachmentUrl,
			Type:     "file",
			RepoPath: repoPath,
			Size:     len(data),
		})
	}
	for imageName, imageUrl := range images {
		uploadFiles, err := api.ListUploads(projectID)
		if err != nil {
			return nil, err
		}
		fileID, ok := uploadFiles[imageName]
		if !ok {
			logger.Logger.Warnf("%s 附件 %s 不存在", repoPath, imageName)
			continue
		}
		data, err := api.DownloadFile(projectID, fileID)
		if err != nil {
			logger.Logger.Errorf("%s 下载release asset %s 失败: %s", imageUrl, imageName, err)
			return nil, err
		}
		attachmentsList = append(attachmentsList, Attachment{
			Name:     imageName,
			Data:     data,
			Url:      imageUrl,
			Type:     "img",
			RepoPath: repoPath,
			Size:     len(data),
		})
	}
	return attachmentsList, nil

}

func (c *GitlabVcs) GetRepoDescription() string {
	return c.Desc
}
