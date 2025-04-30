package vcs

import (
	"ccrctl/pkg/api/coding"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"strconv"
)

const (
	CodingUserName = "coding"
)

type CodingVcs struct {
	httpURL      string
	RepoPath     string
	SubGroupName string
	RepoName     string
	RepoType     string
	Private      bool
	Desc         string
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
		logger.Logger.Errorf(err.Error())
		panic(err)
	}
	return &SubGroup{
		Name:   c.SubGroupName,
		Desc:   project.Description,
		Remark: project.DisplayName,
	}
}

func (c *CodingVcs) GetRepoType() string {
	return c.RepoType
}

func (c *CodingVcs) GetCloneUrl() string {
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
	return nil
}

func (c *CodingVcs) GetProjectID() string {
	return strconv.Itoa(0)
}

func newCodingRepo() []VCS {
	repoList, err := coding.GetDepotList(config.Cfg.GetString("migrate.type"))
	if err != nil {
		logger.Logger.Errorf(err.Error())
		panic(err)
	}
	return CodingCovertToVcs(repoList)
}

func CodingCovertToVcs(repoList []coding.Depots) []VCS {
	var VCS []VCS
	for _, repo := range repoList {
		VCS = append(VCS, &CodingVcs{
			httpURL:      repo.HttpsUrl,
			RepoPath:     repo.GetRepoPath(),
			SubGroupName: repo.ProjectName,
			RepoName:     repo.Name,
			RepoType:     repo.RepoType,
			Private:      !repo.IsShared,
			Desc:         repo.Description,
		})
	}
	return VCS
}

func (c *CodingVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}

func (c *CodingVcs) GetRepoDescription() string {
	return c.Desc
}
