package appstore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client HTTP客户端结构体
// 用于发送HTTP请求到App Store API
type Client struct {
	httpClient *http.Client // HTTP客户端实例
	baseURL    string       // API基础URL
	userAgent  string       // User-Agent头
}

// NewClient 创建并返回一个新的HTTP客户端
// 默认配置：30秒超时，使用合理的User-Agent
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // 30秒请求超时
		},
		baseURL:   "https://itunes.apple.com", // iTunes API基础URL
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

// SetTimeout 设置HTTP客户端的超时时间
// timeout: 超时时长
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// SetUserAgent 设置User-Agent头
// agent: User-Agent字符串
func (c *Client) SetUserAgent(agent string) {
	c.userAgent = agent
}

// Get 发送GET请求到指定URL
// path: 请求路径
// params: 查询参数
// 返回解析后的JSON数据和错误信息
func (c *Client) Get(path string, params map[string]string) (map[string]json.RawMessage, error) {
	fullURL := c.baseURL + path

	if params != nil && len(params) > 0 {
		queryParams := url.Values{}
		for key, value := range params {
			queryParams.Add(key, value)
		}
		fullURL += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if len(body) == 0 || len(body) > 0 && string(body[0]) == "<" {
		return nil, fmt.Errorf("响应内容为空或HTML格式，可能是API不可用")
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return result, nil
}

// BuildParams 构建API查询参数
// 将结构体字段转换为map，用于构建URL查询参数
func BuildParams[T any](opts T) map[string]string {
	params := make(map[string]string)

	switch opt := any(opts).(type) {
	case SearchOptions:
		params["term"] = opt.Term
		if opt.Country != "" {
			params["country"] = string(opt.Country)
		} else {
			params["country"] = string(CountryUS)
		}
		if opt.Lang != "" {
			params["lang"] = opt.Lang
		}
		if opt.Num > 0 {
			params["limit"] = strconv.Itoa(opt.Num)
		}

	case RatingsOptions:
		if opt.ID != 0 {
			params["id"] = strconv.FormatInt(opt.ID, 10)
		}
		if opt.Country != "" {
			params["country"] = string(opt.Country)
		} else {
			params["country"] = string(CountryUS)
		}

	case SimilarOptions:
		if opt.ID != 0 {
			params["id"] = strconv.FormatInt(opt.ID, 10)
		}
		if opt.Country != "" {
			params["country"] = string(opt.Country)
		} else {
			params["country"] = string(CountryUS)
		}
	}

	return params
}

// GetHTML 发送GET请求并返回HTML内容
func (c *Client) GetHTML(url string) (string, error) {
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取HTML内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取HTML失败: %w", err)
	}

	return string(body), nil
}

// GetHTMLWithHeaders 发送GET请求并返回HTML内容（带自定义请求头）
func (c *Client) GetHTMLWithHeaders(url string, headers map[string]string) (string, error) {
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置默认请求头
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 添加自定义请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取内容失败: %w", err)
	}

	return string(body), nil
}

// GetLookup 发送GET请求到lookup API并返回结构化的App结果
func (c *Client) GetLookup(path string, params map[string]string) (*LookupResponse, error) {
	fullURL := c.baseURL + path

	if params != nil && len(params) > 0 {
		queryParams := url.Values{}
		for key, value := range params {
			queryParams.Add(key, value)
		}
		fullURL += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if len(body) == 0 || len(body) > 0 && string(body[0]) == "<" {
		return nil, fmt.Errorf("响应内容为空或HTML格式，可能是API不可用")
	}

	var result LookupResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return &result, nil
}

// GetRSSFeedTyped 发送GET请求到RSS Feed并返回结构化的数据
func (c *Client) GetRSSFeedTyped(country string, collection Collection, category Category, limit int) (*RSSFeedResponse, error) {
	collectionStr := string(collection)
	if collectionStr == "" {
		collectionStr = "topfreeapplications"
	}

	url := fmt.Sprintf("%s/%s/rss/%s", c.baseURL, country, collectionStr)

	if category != 0 {
		url += fmt.Sprintf("/genre=%d", category)
	}

	url += fmt.Sprintf("/limit=%d/json", limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result RSSFeedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return &result, nil
}

// GetReviewFeed 获取评论Feed
func (c *Client) GetReviewFeed(country string, page, id int, sortby string) (*ReviewFeedResponse, error) {
	url := fmt.Sprintf("https://itunes.apple.com/%s/rss/customerreviews/page=%d/id=%d/sortby=%s/json",
		country, page, id, sortby)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ReviewFeedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return &result, nil
}

// GetRaw 发送GET请求并返回原始响应字符串
func (c *Client) GetRaw(url string) (string, error) {
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json,text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取内容失败: %w", err)
	}

	return string(body), nil
}
