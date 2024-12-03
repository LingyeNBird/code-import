package vcs

import (
	"fmt"
)

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
	default:
		return nil, fmt.Errorf("不支持的仓库平台: %s", sourceRepoPlatformName)
	}
}
