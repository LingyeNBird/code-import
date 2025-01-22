package vcs

import (
	"ccrctl/pkg/config"
	"fmt"
)

var (
	allowIncompletePush = config.Cfg.GetBool("migrate.allow_incomplete_push")
)

type releases struct {
	Body            string  `json:"body"`
	Draft           bool    `json:"draft"`
	MakeLatest      string  `json:"make_latest"`
	Name            string  `json:"name"`
	Prerelease      bool    `json:"prerelease"`
	TagName         string  `json:"tag_name"`
	TargetCommitish string  `json:"target_commitish"`
	Assets          []Asset `json:"assets"`
}

type Asset struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type VCS interface {
	GetRepoPath() string
	GetRepoName() string
	GetSubGroupName() string
	GetRepoType() string
	GetCloneUrl() string
	GetUserName() string
	GetToken() string
	Clone() error
	GetRepoPrivate() bool
	GetReleases() []releases
	GetProjectID() string
}

func NewVcs(sourceRepoPlatformName string) ([]VCS, error) {
	switch sourceRepoPlatformName {
	case "coding":
		return newCodingRepo(), nil
	case "gitlab":
		return newGitlabRepo(), nil
	case "github":
		return newGithubRepo(), nil
	case "gitee":
		return newGiteeRepo(), nil
	case "common":
		return newCommonRepo(), nil
	case "aliyun":
		return newAliyunRepo(), nil
	default:
		return nil, fmt.Errorf("不支持的仓库平台: %s", sourceRepoPlatformName)
	}
}
