package gitee

import "C"
import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	apiPath     = "/api/v5"
	host        = "https://gitee.com"
	getRepoList = "/user/repos"
	getUser     = "/user"
)

type Repo struct {
	Id        int    `json:"id"`
	FullName  string `json:"full_name"`
	HumanName string `json:"human_name"`
	Url       string `json:"url"`
	Namespace struct {
		Id      int    `json:"id"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Path    string `json:"path"`
		HtmlUrl string `json:"html_url"`
	} `json:"namespace"`
	Path  string `json:"path"`
	Name  string `json:"name"`
	Owner struct {
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarUrl         string `json:"avatar_url"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		Remark            string `json:"remark"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
	} `json:"owner"`
	Assigner struct {
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarUrl         string `json:"avatar_url"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		Remark            string `json:"remark"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
	} `json:"assigner"`
	Description         string      `json:"description"`
	Private             bool        `json:"private"`
	Public              bool        `json:"public"`
	Internal            bool        `json:"internal"`
	Fork                bool        `json:"fork"`
	HtmlUrl             string      `json:"html_url"`
	SshUrl              string      `json:"ssh_url"`
	ForksUrl            string      `json:"forks_url"`
	KeysUrl             string      `json:"keys_url"`
	CollaboratorsUrl    string      `json:"collaborators_url"`
	HooksUrl            string      `json:"hooks_url"`
	BranchesUrl         string      `json:"branches_url"`
	TagsUrl             string      `json:"tags_url"`
	BlobsUrl            string      `json:"blobs_url"`
	StargazersUrl       string      `json:"stargazers_url"`
	ContributorsUrl     string      `json:"contributors_url"`
	CommitsUrl          string      `json:"commits_url"`
	CommentsUrl         string      `json:"comments_url"`
	IssueCommentUrl     string      `json:"issue_comment_url"`
	IssuesUrl           string      `json:"issues_url"`
	PullsUrl            string      `json:"pulls_url"`
	MilestonesUrl       string      `json:"milestones_url"`
	NotificationsUrl    string      `json:"notifications_url"`
	LabelsUrl           string      `json:"labels_url"`
	ReleasesUrl         string      `json:"releases_url"`
	Recommend           bool        `json:"recommend"`
	Gvp                 bool        `json:"gvp"`
	Homepage            interface{} `json:"homepage"`
	Language            string      `json:"language"`
	ForksCount          int         `json:"forks_count"`
	StargazersCount     int         `json:"stargazers_count"`
	WatchersCount       int         `json:"watchers_count"`
	DefaultBranch       string      `json:"default_branch"`
	OpenIssuesCount     int         `json:"open_issues_count"`
	HasIssues           bool        `json:"has_issues"`
	HasWiki             bool        `json:"has_wiki"`
	IssueComment        bool        `json:"issue_comment"`
	CanComment          bool        `json:"can_comment"`
	PullRequestsEnabled bool        `json:"pull_requests_enabled"`
	HasPage             bool        `json:"has_page"`
	License             interface{} `json:"license"`
	Outsourced          bool        `json:"outsourced"`
	ProjectCreator      string      `json:"project_creator"`
	Members             []string    `json:"members"`
	PushedAt            time.Time   `json:"pushed_at"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	Parent              interface{} `json:"parent"`
	Paas                interface{} `json:"paas"`
	Stared              bool        `json:"stared"`
	Watched             bool        `json:"watched"`
	Permission          struct {
		Pull  bool `json:"pull"`
		Push  bool `json:"push"`
		Admin bool `json:"admin"`
	} `json:"permission"`
	Relation            string        `json:"relation"`
	AssigneesNumber     int           `json:"assignees_number"`
	TestersNumber       int           `json:"testers_number"`
	Assignee            []interface{} `json:"assignee"`
	Testers             []interface{} `json:"testers"`
	Status              string        `json:"status"`
	Programs            []interface{} `json:"programs"`
	Enterprise          interface{}   `json:"enterprise"`
	ProjectLabels       []interface{} `json:"project_labels"`
	IssueTemplateSource string        `json:"issue_template_source"`
}

type ErrorResp struct {
	Message string `json:"message"`
}

type User struct {
	Id                int         `json:"id"`
	Login             string      `json:"login"`
	Name              string      `json:"name"`
	AvatarUrl         string      `json:"avatar_url"`
	Url               string      `json:"url"`
	HtmlUrl           string      `json:"html_url"`
	Remark            string      `json:"remark"`
	FollowersUrl      string      `json:"followers_url"`
	FollowingUrl      string      `json:"following_url"`
	GistsUrl          string      `json:"gists_url"`
	StarredUrl        string      `json:"starred_url"`
	SubscriptionsUrl  string      `json:"subscriptions_url"`
	OrganizationsUrl  string      `json:"organizations_url"`
	ReposUrl          string      `json:"repos_url"`
	EventsUrl         string      `json:"events_url"`
	ReceivedEventsUrl string      `json:"received_events_url"`
	Type              string      `json:"type"`
	Blog              string      `json:"blog"`
	Weibo             interface{} `json:"weibo"`
	Bio               string      `json:"bio"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	Stared            int         `json:"stared"`
	Watched           int         `json:"watched"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	Email             interface{} `json:"email"`
}

func GetRepoListFetchPage(page string) ([]Repo, http.Header, error) {
	queryParams := url.Values{}
	queryParams.Add("access_token", config.Cfg.GetString("source.token"))
	queryParams.Add("affiliation", "admin")
	queryParams.Add("sort", "full_name")
	queryParams.Add("per_page", "100")
	queryParams.Add("page", page)
	endPoint := apiPath + getRepoList + "?" + queryParams.Encode()
	c := http_client.NewClient(host)
	resp, header, respCode, err := c.GiteeClient(http.MethodGet, endPoint, nil)
	if err != nil {
		logger.Logger.Error("Failed to get repo list", err)
		return nil, nil, err
	}
	if respCode != http.StatusOK {
		var e ErrorResp
		err = c.Unmarshal(resp, &e)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, fmt.Errorf("failed to get repo list: %s", e.Message)
	}
	var repoList []Repo
	err = c.Unmarshal(resp, &repoList)
	if err != nil {
		return nil, nil, err
	}
	return repoList, header, err
}

func GetRepoList() ([]Repo, error) {
	var repoList []Repo
	page := 1 // 将page初始化为整数
	for {
		list, header, err := GetRepoListFetchPage(strconv.Itoa(page)) // 使用strconv.Itoa将整数转换为字符串
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			break
		}
		repoList = append(repoList, list...)
		if strconv.Itoa(page) == header.Get("total_page") {
			break
		}
		page++ // 直接递增整数
	}
	return repoList, nil
}

func GetUserName() (name string, err error) {
	c := http_client.NewClient(host)
	queryParams := url.Values{}
	queryParams.Add("access_token", config.Cfg.GetString("source.token"))
	endPoint := apiPath + getUser + "?" + queryParams.Encode()
	resp, _, respCode, err := c.GiteeClient(http.MethodGet, endPoint, nil)
	if err != nil {
		logger.Logger.Error("Failed to get username", err)
		return name, err
	}
	if respCode != http.StatusOK {
		var e ErrorResp
		err = c.Unmarshal(resp, &e)
		if err != nil {
			return name, err
		}
		return name, fmt.Errorf("failed to get username: %s", e.Message)
	}
	var data User
	err = c.Unmarshal(resp, &data)
	if err != nil {
		return name, err
	}
	name = data.Login
	return name, err
}
