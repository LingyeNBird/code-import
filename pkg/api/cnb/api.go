package cnb

import (
	"ccrctl/pkg/http_client"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	c = http_client.NewCNBClient()
)

type Repos struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Freeze          bool      `json:"freeze"`
	Status          int       `json:"status"`
	VisibilityLevel string    `json:"visibility_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Description     string    `json:"description"`
	Site            string    `json:"site"`
	Topics          string    `json:"topics"`
	License         string    `json:"license"`
	DisplayModule   struct {
		Activity     bool `json:"activity"`
		Contributors bool `json:"contributors"`
		Release      bool `json:"release"`
	} `json:"display_module"`
	StarCount            int         `json:"star_count"`
	ForkCount            int         `json:"fork_count"`
	MarkCount            int         `json:"mark_count"`
	LastUpdatedAt        interface{} `json:"last_updated_at"`
	Language             string      `json:"language"`
	WebUrl               string      `json:"web_url"`
	Path                 string      `json:"path"`
	Tags                 interface{} `json:"tags"`
	OpenIssueCount       int         `json:"open_issue_count"`
	OpenPullRequestCount int         `json:"open_pull_request_count"`
	LastUpdateUsername   string      `json:"last_update_username"`
	LastUpdateNickname   string      `json:"last_update_nickname"`
}

func GetUserRepoFetchPage(page int) (repos []Repos, totalRow, pageSize int, err error) {
	endpoint := fmt.Sprintf("/user/repos?page=%d&page_size=100&desc=false", page)
	resp, header, _, err := c.RequestV4(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	err = c.Unmarshal(resp, &repos)
	if err != nil {
		return nil, 0, 0, err
	}
	totalRow, err = strconv.Atoi(header.Get("x-cnb-total"))
	if err != nil {
		return nil, 0, 0, err
	}
	pageSize, err = strconv.Atoi(header.Get("x-cnb-page-size"))
	if err != nil {
		return nil, 0, 0, err
	}
	return repos, totalRow, pageSize, nil
}

func GetUserRepos() (Repos []Repos, err error) {
	page := 1
	for {
		repos, totalRow, pageSize, err := GetUserRepoFetchPage(page)
		if err != nil {
			return nil, err
		}
		Repos = append(Repos, repos...)
		if page*pageSize >= totalRow {
			break
		}
		page++
	}
	return Repos, nil
}

func GetReposByGroupFetchPage(groupName string, page int) (repos []Repos, totalRow, pageSize int, err error) {
	endpoint := fmt.Sprintf("/%s/-/repos?page=%d&page_size=100&desc=false", groupName, page)
	resp, header, _, err := c.RequestV4(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	err = c.Unmarshal(resp, &repos)
	if err != nil {
		return nil, 0, 0, err
	}
	totalRow, err = strconv.Atoi(header.Get("x-cnb-total"))
	if err != nil {
		return nil, 0, 0, err
	}
	pageSize, err = strconv.Atoi(header.Get("x-cnb-page-size"))
	if err != nil {
		return nil, 0, 0, err
	}
	return repos, totalRow, pageSize, nil
}

func GetReposByGroup(group string) ([]Repos, error) {
	var Repos []Repos
	page := 1
	for {
		apiRepos, totalRow, pageSize, err := GetReposByGroupFetchPage(group, page)
		if err != nil {
			return nil, err
		}
		for _, repo := range apiRepos {
			Repos = append(Repos, repo)
		}
		if page*pageSize >= totalRow {
			break
		}
		page++
	}
	return Repos, nil
}
