package vcs

import (
	"ccrctl/pkg/api/gitee"
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
		vcs, err := gitee.GetRepoList()
		if err != nil {
			return nil, fmt.Errorf("获取gitee仓库列表失败")
		}
		return vcs, nil
	case "common":
		return newCommonRepo(), nil
	default:
		return nil, fmt.Errorf("不支持的仓库平台: %s", sourceRepoPlatformName)
	}
}
