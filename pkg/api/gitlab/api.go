package gitlab

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	gitlab "github.com/xanzy/go-gitlab"
)

var (
	Git   *gitlab.Client
	token = config.Cfg.GetString("source.token")
	url   = config.Cfg.GetString("source.url")
)

func init() {
	var err error
	Git, err = gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		logger.Logger.Fatalf("Failed to create Gitlab client: %v", err)
	}
}

func GetProjects() ([]*gitlab.Project, error) {
	var Projects []*gitlab.Project
	page := 1
	for {
		projects, resp, err := Git.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    page,
			},
			Owned: gitlab.Bool(true),
		})
		if err != nil {
			logger.Logger.Fatalf("Failed to get Projects: %v", err)
		}
		Projects = append(Projects, projects...)
		if resp.NextPage == 0 {
			break
		}
		page++
	}
	return Projects, nil
}

// GetRelease 获取指定项目的release
func GetReleases(projectID int) (releases []*gitlab.Release, err error) {
	page := 1
	for {
		release, resp, err := Git.Releases.ListReleases(projectID, &gitlab.ListReleasesOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    page,
			},
			Sort: gitlab.String("asc"),
		})
		if err != nil {
			logger.Logger.Fatalf("Failed to get Releases: %v", err)
		}
		releases = append(releases, release...)
		if resp.NextPage == 0 {
			break
		}
		page++
	}
	return releases, nil
}

type ListUploadsRes struct {
	Id         int       `json:"id"`
	Size       int       `json:"size"`
	Filename   string    `json:"filename"`
	CreatedAt  time.Time `json:"created_at"`
	UploadedBy struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
	} `json:"uploaded_by"`
}

// ListUploads https://docs.gitlab.com/ee/api/project_markdown_uploads.html
func ListUploads(projectID string) (files map[string]int, err error) {
	u := fmt.Sprintf("%s/api/v4/projects/%s/uploads", url, projectID)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, u, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var uploads []ListUploadsRes
	err = json.Unmarshal(body, &uploads)
	if err != nil {
		return nil, err
	}
	files = make(map[string]int)
	for _, upload := range uploads {
		fileName := upload.Filename
		files[fileName] = upload.Id
	}
	return files, nil
}

func DownloadFile(projectID string, fileID int) (data []byte, err error) {
	u := fmt.Sprintf("%s/api/v4/projects/%s/uploads/%d", url, projectID, fileID)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, u, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
