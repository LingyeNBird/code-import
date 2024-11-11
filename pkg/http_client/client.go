package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client 是 OpenAPI 客户端的结构体
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient 创建一个新的 OpenAPI 客户端
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// Request 发送一个 HTTP 请求到 OpenAPI
func (c *Client) Request(method, endpoint string, token string, body interface{}) ([]byte, error) {
	// 将 body 转换为 JSON 格式
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest(method, c.BaseURL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) RequestV2(method, endpoint string, token string, body interface{}) ([]byte, http.Header, error) {
	// 将 body 转换为 JSON 格式
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}

	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest(method, c.BaseURL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, resp.Header, nil
}

func (c *Client) RequestV3(method, endpoint string, token string, body interface{}) ([]byte, http.Header, int, error) {
	// 将 body 转换为 JSON 格式
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, 0, err
	}

	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest(method, c.BaseURL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, 0, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, 0, err
	}
	return respBody, resp.Header, resp.StatusCode, nil
}

func (c *Client) Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
