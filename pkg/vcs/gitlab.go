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
}

func (g *GitlabVcs) GetRepoPath() string {
	return g.RepoPath
}

func (g *GitlabVcs) GetSubGroupName() string {
	parts := strings.Split(g.RepoPath, "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1] // 去掉仓库名
	}
	result := strings.Join(parts, "/")
	return result
}

func (g *GitlabVcs) GetRepoName() string {
	return g.RepoName
}

func (g *GitlabVcs) GetRepoType() string {
	return g.RepoType
}

func (g *GitlabVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(g.httpURL, GitlabUserName, g.GetToken())
}

func (g *GitlabVcs) GetUserName() string {
	return GitlabUserName
}

func (g *GitlabVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (g *GitlabVcs) Clone() error {
	err := git.Clone(g.GetCloneUrl(), g.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (g *GitlabVcs) GetRepoPrivate() bool {
	switch g.Private {
	case "private":
		return true
	case "internal":
		return true
	}
	return false
}

func (g *GitlabVcs) GetReleases() (cnbReleases []Releases) {
	gitlabReleases, err := api.GetReleases(g.ProjectId)
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

func (g *GitlabVcs) GetProjectID() string {
	return strconv.Itoa(g.ProjectId)
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
		})
	}
	return VCS
}

func (g *GitlabVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
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
