package http_client

import (
	"bytes"
	"ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

const (
	giteeApiPath = "/api/v5"
	giteeApiHost = "https://gitee.com"
)

var (
	CnbURL = config.Cfg.GetString("cnb.url")
)

// Client 是 OpenAPI 客户端的结构体
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	Limiter    *rate.Limiter
}

// NewClientV2 创建一个新的 OpenAPI 客户端
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		Limiter:    rate.NewLimiter(rate.Every(time.Second), 10),
	}
}

// NewClientV2 NewClient 创建一个新的 OpenAPI 客户端
func NewClientV2() *Client {
	return &Client{
		BaseURL:    config.ConvertToApiURL(CnbURL),
		HTTPClient: &http.Client{},
		Token:      config.Cfg.GetString("cnb.token"),
		Limiter:    rate.NewLimiter(rate.Every(time.Second), 10),
	}
}

// NewCNBClient NewClientV3 创建一个新的 OpenAPI 客户端
func NewCNBClient() *Client {
	return &Client{
		BaseURL:    config.ConvertToApiURL(config.Cfg.GetString("source.url")),
		HTTPClient: &http.Client{},
		Token:      config.Cfg.GetString("source.token"),
		Limiter:    rate.NewLimiter(rate.Every(time.Second), 10),
	}
}

func NewGiteeClient() *Client {
	return &Client{
		BaseURL:    giteeApiHost + giteeApiPath,
		HTTPClient: &http.Client{},
		Token:      config.Cfg.GetString("source.token"),
		Limiter:    rate.NewLimiter(rate.Every(time.Second), 1),
	}
}

// Request 发送一个 HTTP 请求到 OpenAPI
func (c *Client) Request(method, endpoint string, token string, body interface{}) ([]byte, error) {
	defer logger.Logger.Debugw("Request", "body", body, "reqPath", endpoint, "url", c.BaseURL+endpoint)
	if err := c.Limiter.Wait(context.Background()); err != nil {
		return nil, err
	}
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
	if err := c.Limiter.Wait(context.Background()); err != nil {
		return nil, nil, err
	}
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
	if err := c.Limiter.Wait(context.Background()); err != nil {
		return nil, nil, 0, err
	}
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

func (c *Client) RequestV4(method, endpoint string, body interface{}) ([]byte, http.Header, int, error) {
	if err := c.Limiter.Wait(context.Background()); err != nil {
		return nil, nil, 0, err
	}
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
	req.Header.Set("Authorization", "Bearer "+c.Token)
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
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, resp.Header, resp.StatusCode, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, resp.Header, resp.StatusCode, nil
}

func (c *Client) RequestWithURL(method, url string, body interface{}) ([]byte, http.Header, int, error) {
	if err := c.Limiter.Wait(context.Background()); err != nil {
		return nil, nil, 0, err
	}
	// 将 body 转换为 JSON 格式
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, 0, err
	}

	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, 0, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
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
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, resp.Header, resp.StatusCode, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, resp.Header, resp.StatusCode, nil
}

func (c *Client) GiteeClient(method, endpoint string, body interface{}) ([]byte, http.Header, int, error) {
	if err := c.Limiter.Wait(context.Background()); err != nil {
		return nil, nil, 0, err
	}
	// 将 body 转换为 JSON 格式
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, 0, err
	}

	fullUrl := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, 0, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

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

func (c *Client) GiteeRequest(method, endpoint string, body interface{}, values url.Values) ([]byte, http.Header, int, error) {
	logger.Logger.Debugf("start gitee request %s", endpoint)
	if err := c.Limiter.Wait(context.Background()); err != nil {
		return nil, nil, 0, err
	}
	// 将 body 转换为 JSON 格式
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, 0, err
	}

	values.Add("access_token", c.Token)
	endpoint = fmt.Sprintf("%s?%s", endpoint, values.Encode())
	fullUrl := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	logger.Logger.Debugf("gitee request url %s", fullUrl)
	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, 0, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

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

// SendUploadRequest 发送上传请求
func (c *Client) SendUploadRequest(url, contentType string, body *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return []byte(""), fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	newClient := &http.Client{}

	resp, err := newClient.Do(req)
	if err != nil {
		return []byte(""), fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) UploadData(url string, data []byte) (err error) {
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	newClient := &http.Client{}

	resp, err := newClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
