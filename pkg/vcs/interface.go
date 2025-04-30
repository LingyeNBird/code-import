package vcs

import (
	"ccrctl/pkg/config"
	"fmt"
)

var (
	allowIncompletePush = config.Cfg.GetBool("migrate.allow_incomplete_push")
)

type Releases struct {
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

type Attachment struct {
	Name     string
	Data     []byte
	Url      string
	Type     string
	RepoPath string
	Size     int
}

type SubGroup struct {
	Name   string
	Desc   string
	Remark string
}

type VCS interface {
	GetRepoPath() string
	GetRepoName() string
	GetSubGroup() *SubGroup
	GetRepoType() string
	GetCloneUrl() string
	GetUserName() string
	GetToken() string
	Clone() error
	GetRepoPrivate() bool
	GetReleases() []Releases
	GetProjectID() string
	GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error)
	GetRepoDescription() string
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
	case "cnb":
		return newCnbRepo(), nil
	default:
		return nil, fmt.Errorf("不支持的仓库平台: %s", sourceRepoPlatformName)
	}
}
