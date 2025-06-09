package gongfeng

import (
	"ccrctl/pkg/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	DefaultPerPage = 100 // 每页最大数量

	// 项目可见度级别
	VisibilityLevelPrivate  = 0  // 私有项目，必须显式授予每个用户访问权限
	VisibilityLevelInternal = 10 // 内部公开项目，任何登录用户都可以访问
)

// Project 工蜂项目结构
type Project struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Path            string `json:"path"`
	PathWithNS      string `json:"path_with_namespace"`
	Description     string `json:"description"`
	VisibilityLevel int    `json:"visibility_level"` // 0: 私有, 10: 内部公开
	WebURL          string `json:"web_url"`
	HTTPURL         string `json:"http_url_to_repo"`
}

// IsPrivate 判断项目是否为私有
// visibility_level = 0: 私有项目，必须显式授予每个用户访问权限
// visibility_level = 10: 内部公开项目，任何登录用户都可以访问
func (p *Project) IsPrivate() bool {
	return p.VisibilityLevel == VisibilityLevelPrivate || p.VisibilityLevel == VisibilityLevelInternal
}

// GetProjects 获取工蜂平台的所有项目列表
func GetProjects() ([]Project, error) {
	url := config.Cfg.GetString("source.url")
	token := config.Cfg.GetString("source.token")

	var allProjects []Project
	page := 1

	for {
		projects, totalPages, err := getProjectsByPage(url, token, page)
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)

		if page >= totalPages {
			break
		}
		page++
	}

	return allProjects, nil
}

// getProjectsByPage 获取指定页的项目列表
func getProjectsByPage(baseURL, token string, page int) ([]Project, int, error) {
	// 构建并发送请求
	resp, err := sendProjectsRequest(baseURL, token, page)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// 获取总页数
	totalPages, err := getTotalPages(resp)
	if err != nil {
		return nil, 0, err
	}

	// 解析项目列表
	projects, err := parseProjectsResponse(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return projects, totalPages, nil
}

// sendProjectsRequest 发送获取项目列表的请求
func sendProjectsRequest(baseURL, token string, page int) (*http.Response, error) {
	apiURL := fmt.Sprintf("%s/api/v3/projects?private_token=%s&page=%d&per_page=%d&owned=true",
		baseURL, token, page, DefaultPerPage)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("API 请求失败，状态码: %d", resp.StatusCode)
	}

	return resp, nil
}

// getTotalPages 从响应头中获取总页数
func getTotalPages(resp *http.Response) (int, error) {
	totalPages, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
	if err != nil {
		return 0, fmt.Errorf("获取总页数失败: %v", err)
	}
	return totalPages, nil
}

// parseProjectsResponse 解析响应体中的项目列表
func parseProjectsResponse(body io.Reader) ([]Project, error) {
	var projects []Project
	if err := json.NewDecoder(body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	return projects, nil
}
