package gitea

import (
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	apiPath     = "/api/v1"
	getRepoList = "/user/repos"
	getUser     = "/user"
	getReleases = "/repos/%s/releases"
)

var (
	c = http_client.NewGiteaClient()
)

// Repo Gitea 仓库结构体
type Repo struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Private       bool   `json:"private"`
	Fork          bool   `json:"fork"`
	HtmlUrl       string `json:"html_url"`
	SshUrl        string `json:"ssh_url"`
	CloneUrl      string `json:"clone_url"`
	Website       string `json:"website"`
	Language      string `json:"language"`
	DefaultBranch string `json:"default_branch"`
	Owner         struct {
		Id        int    `json:"id"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
		Username  string `json:"username"`
	} `json:"owner"`
	Permissions struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
	HasIssues                 bool      `json:"has_issues"`
	HasWiki                   bool      `json:"has_wiki"`
	HasPullRequests           bool      `json:"has_pull_requests"`
	HasProjects               bool      `json:"has_projects"`
	IgnoreWhitespaceConflicts bool      `json:"ignore_whitespace_conflicts"`
	AllowMergeCommits         bool      `json:"allow_merge_commits"`
	AllowRebase               bool      `json:"allow_rebase"`
	AllowRebaseMerge          bool      `json:"allow_rebase_explicit"`
	AllowSquashMerge          bool      `json:"allow_squash_merge"`
	AvatarUrl                 string    `json:"avatar_url"`
	Internal                  bool      `json:"internal"`
	MirrorInterval            string    `json:"mirror_interval"`
	Size                      int       `json:"size"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// User Gitea 用户结构体
type User struct {
	Id        int       `json:"id"`
	Login     string    `json:"login"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatar_url"`
	Username  string    `json:"username"`
	Language  string    `json:"language"`
	IsAdmin   bool      `json:"is_admin"`
	LastLogin time.Time `json:"last_login"`
	Created   time.Time `json:"created"`
}

// Release Gitea Release 结构体
type Release struct {
	Id              int       `json:"id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Url             string    `json:"url"`
	HtmlUrl         string    `json:"html_url"`
	TarballUrl      string    `json:"tarball_url"`
	ZipballUrl      string    `json:"zipball_url"`
	IsDraft         bool      `json:"draft"`
	IsPrerelease    bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          struct {
		Id        int    `json:"id"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
		Username  string `json:"username"`
	} `json:"author"`
	Assets []struct {
		Id                 int       `json:"id"`
		Name               string    `json:"name"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		Uuid               string    `json:"uuid"`
		BrowserDownloadUrl string    `json:"browser_download_url"`
	} `json:"assets"`
}

// ErrorResp Gitea 错误响应结构体
type ErrorResp struct {
	Message string `json:"message"`
	Url     string `json:"url"`
}

// GetRepoListFetchPage 分页获取仓库列表
func GetRepoListFetchPage(page string) ([]Repo, http.Header, error) {
	queryParams := url.Values{}
	queryParams.Add("page", page)
	queryParams.Add("limit", "50") // Gitea 默认每页最多50个

	endpoint := fmt.Sprintf("%s?%s", getRepoList, queryParams.Encode())
	resp, header, respCode, err := c.GiteaRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logger.Logger.Errorf("获取仓库列表失败: %v", err)
		return nil, nil, err
	}

	if respCode != http.StatusOK {
		var e ErrorResp
		err = c.Unmarshal(resp, &e)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, fmt.Errorf("获取仓库列表失败: %s", e.Message)
	}

	var repoList []Repo
	err = c.Unmarshal(resp, &repoList)
	if err != nil {
		return nil, nil, err
	}

	return repoList, header, nil
}

// GetRepoList 获取所有仓库列表
func GetRepoList() ([]Repo, error) {
	var repoList []Repo
	page := 1

	for {
		list, _, err := GetRepoListFetchPage(strconv.Itoa(page))
		if err != nil {
			return nil, err
		}

		if len(list) == 0 {
			break
		}

		repoList = append(repoList, list...)

		// 如果返回的数量少于每页限制，说明已经是最后一页
		if len(list) < 50 {
			break
		}

		page++
	}

	return repoList, nil
}

// GetUserName 获取当前用户名
func GetUserName() (string, error) {
	resp, _, respCode, err := c.GiteaRequest(http.MethodGet, getUser, nil)
	if err != nil {
		logger.Logger.Errorf("获取用户信息失败: %v", err)
		return "", err
	}

	if respCode != http.StatusOK {
		var e ErrorResp
		err = c.Unmarshal(resp, &e)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("获取用户信息失败: %s", e.Message)
	}

	var user User
	err = c.Unmarshal(resp, &user)
	if err != nil {
		return "", err
	}

	return user.Login, nil
}

// GetReleasesFetchPage 分页获取 Release 列表
func GetReleasesFetchPage(repoPath string, pageInt int) ([]Release, error) {
	page := strconv.Itoa(pageInt)
	queryParams := url.Values{}
	queryParams.Add("page", page)
	queryParams.Add("limit", "50")

	endpoint := fmt.Sprintf(getReleases, repoPath)
	endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	resp, _, respCode, err := c.GiteaRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logger.Logger.Errorf("获取 Release 列表失败: %v", err)
		return nil, err
	}

	if respCode != http.StatusOK {
		var e ErrorResp
		err = c.Unmarshal(resp, &e)
		if err != nil {
			logger.Logger.Errorf("解析 Release 错误响应失败: %v", err)
			return nil, err
		}
		logger.Logger.Errorf("获取 Release 列表失败: %s", e.Message)
		return nil, fmt.Errorf("获取 Release 列表失败: %s", e.Message)
	}

	var releases []Release
	err = c.Unmarshal(resp, &releases)
	if err != nil {
		logger.Logger.Errorf("解析 Release 列表失败: %v", err)
		return nil, err
	}

	return releases, nil
}

// GetReleases 获取所有 Release 列表
func GetReleases(repoPath string) ([]Release, error) {
	page := 1
	var releases []Release

	for {
		data, err := GetReleasesFetchPage(repoPath, page)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			break
		}

		releases = append(releases, data...)

		// 如果返回的数量少于每页限制，说明已经是最后一页
		if len(data) < 50 {
			break
		}

		page++
	}

	return releases, nil
}
