package vcs

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"fmt"
	"strconv"
	"strings"
)

type CommonVcs struct {
	httpURL  string
	RepoPath string
	RepoName string
	RepoType string
}

func (c *CommonVcs) GetRepoPath() string {
	return c.RepoPath
}

func (c *CommonVcs) GetSubGroup() *SubGroup {
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

func (c *CommonVcs) GetRepoName() string {
	return c.RepoName
}

func (c *CommonVcs) GetRepoType() string {
	return c.RepoType
}

func (c *CommonVcs) GetCloneUrl() string {
	ssh := config.Cfg.GetBool("migrate.ssh")
	if ssh {
		// 将 HTTP URL 转换为 SSH 格式
		// 示例转换: http://example.com/group/repo.git -> git@example.com:group/repo.git
		sshURL := strings.Replace(c.httpURL, "http://", "", 1)
		sshURL = strings.Replace(sshURL, "https://", "", 1)
		sshURL = strings.Replace(sshURL, "/", ":", 1)
		return fmt.Sprintf("git@%s", sshURL)
	}
	// 默认保持 HTTP 协议克隆
	return util.ConvertUrlWithAuth(c.httpURL, c.GetUserName(), c.GetToken())
}

func (c *CommonVcs) GetUserName() string {
	return config.Cfg.GetString("source.username")
}

func (c *CommonVcs) GetToken() string {
	return config.Cfg.GetString("source.password")
}

func (c *CommonVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommonVcs) GetRepoPrivate() bool {
	return true
}

func (c *CommonVcs) GetReleases() (cnbReleases []Releases) {
	return nil
}

func (c *CommonVcs) GetProjectID() string {
	return strconv.Itoa(0)
}

func newCommonRepo() ([]VCS, error) {
	VCS, err := getRepos()
	if err != nil {
		return nil, err
	}
	return VCS, nil
}

func getRepos() ([]VCS, error) {
	var VCS []VCS
	repos := config.Cfg.GetStringSlice("source.repo")
	for _, repo := range repos {
		VCS = append(VCS, &CommonVcs{
			httpURL:  config.Cfg.GetString("source.url") + "/" + repo + ".git",
			RepoPath: repo,
			RepoName: strings.Split(repo, "/")[len(strings.Split(repo, "/"))-1],
			RepoType: Git,
		})
	}
	logger.Logger.Debugw("获取到的仓库列表", "VCS", VCS)
	return VCS, nil
}

func (c *CommonVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}

func (c *CommonVcs) GetRepoDescription() string {
	return ""
}

func (c *CommonVcs) ListRepos() ([]VCS, error) {
	return newCommonRepo()
}
