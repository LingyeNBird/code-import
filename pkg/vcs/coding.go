package vcs

import (
	"ccrctl/pkg/coding"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/util"
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
}

func (r *CodingVcs) GetRepoPath() string {
	return r.RepoPath
}

func (r *CodingVcs) GetRepoName() string {
	return r.RepoName
}

func (r *CodingVcs) GetSubGroupName() string {
	return r.SubGroupName
}

func (r *CodingVcs) GetRepoType() string {
	return r.RepoType
}

func (r *CodingVcs) GetCloneUrl() string {
	return util.ConvertUrlWithAuth(r.httpURL, CodingUserName, r.GetToken())
}

func (r *CodingVcs) GetUserName() string {
	return CodingUserName
}

func (r *CodingVcs) Clone() error {
	err := git.Clone(r.GetCloneUrl(), r.GetRepoPath())
	if err != nil {
		return err
	}
	return nil
}

func (r *CodingVcs) GetToken() string {
	return config.Cfg.GetString("source.token")
}

func (r *CodingVcs) GetRepoPrivate() bool {
	return true
}

func newCodingRepo() []VCS {
	repoList, err := coding.GetDepotList(config.Cfg.GetString("migrate.type"))
	if err != nil {
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
		})
	}
	return VCS
}
