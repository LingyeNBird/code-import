package aliyun

import (
	"ccrctl/pkg/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	AliyunEndpoint  = "https://openapi-rdc.aliyuncs.com"
	defaultPageSize = 100
)

// Repository 表示代码仓库信息
type Repository struct {
	AccessLevel       int    `json:"accessLevel"`
	Archived          bool   `json:"archived"`
	CreatedAt         string `json:"createdAt"`
	CreatorId         int    `json:"creatorId"`
	DemoProject       bool   `json:"demoProject"`
	Encrypted         bool   `json:"encrypted"`
	Id                int    `json:"id"`
	LastActivityAt    string `json:"lastActivityAt"`
	Name              string `json:"name"`
	NameWithNamespace string `json:"nameWithNamespace"`
	NamespaceId       int    `json:"namespaceId"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"pathWithNamespace"`
	StarCount         int    `json:"starCount"`
	Starred           bool   `json:"starred"`
	UpdatedAt         string `json:"updatedAt"`
	Visibility        string `json:"visibility"`
	WebUrl            string `json:"webUrl"`
	Description string `json:"description"`
}

// RepositoryListResponse 表示获取仓库列表的响应
type RepositoryListResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total    int          `json:"total"`
		Page     int          `json:"page"`
		PageSize int          `json:"pageSize"`
		List     []Repository `json:"list"`
	} `json:"data"`
}

// GetRepositories 获取仓库列表
// page: 页码，从1开始
func GetRepositories(page int) ([]Repository, int, error) {
	organizationID := config.Cfg.GetString("source.organizationId")
	token := config.Cfg.GetString("source.token")

	url := fmt.Sprintf("%s/oapi/v1/codeup/organizations/%s/repositories?page=%d&perPage=%d",
		AliyunEndpoint, organizationID, page, defaultPageSize)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Add("x-yunxiao-token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, 0, fmt.Errorf("请求失败: %s", string(body))
	}

	var result []Repository
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, fmt.Errorf("解析响应失败: %w", err)
	}

	totalPages := 0
	if totalPagesStr := resp.Header.Get("X-Total-Pages"); totalPagesStr != "" {
		fmt.Sscanf(totalPagesStr, "%d", &totalPages)
	}

	return result, totalPages, nil
}

// GetAllRepositories 获取所有仓库列表（自动处理分页）
func GetAllRepositories() ([]Repository, error) {
	var allRepos []Repository
	page := 1

	for {
		repos, totalPages, err := GetRepositories(page)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if page >= totalPages {
			break
		}
		page++
	}
	return allRepos, nil
}
