package vcs

import (
	api "ccrctl/pkg/api/cnb"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"strconv"
	"strings"
)

const (
	CNBUserName = "cnb"
)

func newCnbRepo() []VCS {
	sourceGroup := config.Cfg.GetString("source.group")
	if sourceGroup != "" {
		repos, err := api.GetReposByGroup(config.Cfg.GetString("source.group"))
		if err != nil {
			logger.Logger.Errorf("Failed to fetch repos from CNB: %v", err)
			panic(err)
		}
		return CNBCovertToVcs(repos)
	} else {
		repos, err := api.GetUserRepos()
		if err != nil {
			logger.Logger.Errorf("Failed to fetch repos from CNB: %v", err)
			panic(err)
		}
		return CNBCovertToVcs(repos)
	}
}

func CNBCovertToVcs(repoList []api.Repos) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &CNBVcs{
			httpURL:  repo.WebUrl,
			RepoPath: repo.Path,
			RepoName: repo.Name,
			Private:  isPrivate(repo.VisibilityLevel),
			Desc:     repo.Description,
		})
	}
	return VCS
}

func isPrivate(visibilityLevel string) bool {
	switch visibilityLevel {
	case "Private":
		return true
	case "Public":
		return false
	}
	return true
}

type CNBVcs struct {
	httpURL  string
	RepoPath string
	RepoName string
	Private  bool
	Desc     string
}

func (c *CNBVcs) GetRepoPath() string {
	parts := strings.Split(c.RepoPath, "/")
	parts = parts[1:] // 去掉根组织名
	result := strings.Join(parts, "/")
	return result
}

func (c *CNBVcs) GetRepoName() string {
	return c.RepoName
}

func (c *CNBVcs) GetSubGroup() *SubGroup {
	var result string
	parts := strings.Split(c.GetRepoPath(), "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1] // 去掉仓库名
		result = strings.Join(parts, "/")
	} else {
		result = ""
	}
	return &SubGroup{
		Name:   result,
		Desc:   "",
		Remark: "",
	}
}

func (c *CNBVcs) GetRepoType() string {
	return Git
}

func (c *CNBVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(c.httpURL, CNBUserName, c.GetToken())
}

func (c *CNBVcs) GetUserName() string {
	return CNBUserName
}

func (c *CNBVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (c *CNBVcs) Clone() error {
	err := git.Clone(c.GetCloneUrl(), c.GetRepoPath(), allowIncompletePush)
	if err != nil {
		return err
	}
	return nil
}

func (c *CNBVcs) GetRepoPrivate() bool {
	return c.Private
}

func (c *CNBVcs) GetReleases() []Releases {
	return nil
}

func (c *CNBVcs) GetProjectID() string {
	return strconv.Itoa(0)
}

func (c *CNBVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}

func (c *CNBVcs) GetRepoDescription() string {
	return c.Desc
}
