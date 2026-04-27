package appstore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client HTTP客户端结构体
// 用于发送HTTP请求到App Store API
type Client struct {
	httpClient *http.Client      // HTTP客户端实例
	baseURL    string            // API基础URL
	userAgent  string            // User-Agent头
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
func (c *Client) Get(path string, params map[string]interface{}) (map[string]interface{}, error) {
	// 构建完整URL
	fullURL := c.baseURL + path

	// 添加查询参数
	if len(params) > 0 {
		queryParams := url.Values{}
		for key, value := range params {
			queryParams.Add(key, fmt.Sprintf("%v", value))
		}
		fullURL += "?" + queryParams.Encode()
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查响应内容是否为空或HTML
	if len(body) == 0 || len(body) > 0 && string(body[0]) == "<" {
		return nil, fmt.Errorf("响应内容为空或HTML格式，可能是API不可用")
	}

	// 解析JSON响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return result, nil
}

// min 返回两个整数中的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetList 发送GET请求并返回JSON数组
// path: 请求路径
// params: 查询参数
// 返回解析后的JSON数组和错误信息
func (c *Client) GetList(path string, params map[string]interface{}) ([]map[string]interface{}, error) {
	// 构建完整URL
	fullURL := c.baseURL + path

	// 添加查询参数
	if len(params) > 0 {
		queryParams := url.Values{}
		for key, value := range params {
			queryParams.Add(key, fmt.Sprintf("%v", value))
		}
		fullURL += "?" + queryParams.Encode()
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return result, nil
}

// BuildParams 构建API查询参数
// 将结构体字段转换为map，用于构建URL查询参数
func BuildParams(opts interface{}) map[string]interface{} {
	params := make(map[string]interface{})

	switch opt := opts.(type) {
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

// parseInterfaceArray 将interface{}数组转换为字符串数组
// data: 需要转换的数据
// 返回转换后的字符串数组
func parseInterfaceArray(data interface{}) []string {
	result := []string{}
	if data == nil {
		return result
	}

	switch v := data.(type) {
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
	case []string:
		result = v
	}
	return result
}

// parseInterfaceInt 将interface{}转换为int
// data: 需要转换的数据
// 返回转换后的整数值
func parseInterfaceInt(data interface{}) int {
	if data == nil {
		return 0
	}

	switch v := data.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		result, _ := strconv.Atoi(v)
		return result
	}
	return 0
}

// parseInterfaceInt64 将interface{}转换为int64
// data: 需要转换的数据
// 返回转换后的64位整数值
func parseInterfaceInt64(data interface{}) int64 {
	if data == nil {
		return 0
	}

	switch v := data.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	case string:
		result, _ := strconv.ParseInt(v, 10, 64)
		return result
	}
	return 0
}

// parseInterfaceFloat64 将interface{}转换为float64
// data: 需要转换的数据
// 返回转换后的浮点数值
func parseInterfaceFloat64(data interface{}) float64 {
	if data == nil {
		return 0.0
	}

	switch v := data.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		result, _ := strconv.ParseFloat(v, 64)
		return result
	}
	return 0.0
}

// parseInterfaceBool 将interface{}转换为bool
// data: 需要转换的数据
// 返回转换后的布尔值
func parseInterfaceBool(data interface{}) bool {
	if data == nil {
		return false
	}

	switch v := data.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true" || v == "1"
	case int:
		return v != 0
	}
	return false
}

// parseInterfaceString 将interface{}转换为string
// data: 需要转换的数据
// 返回转换后的字符串值
func parseInterfaceString(data interface{}) string {
	if data == nil {
		return ""
	}

	switch v := data.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	}
	return ""
}

// GetRSSFeed 发送GET请求到RSS Feed并返回解析后的数据
// country: 国家代码
// collection: 集合类型（字符串）
// category: 类别ID（可选，0表示不限制）
// limit: 返回数量
// 返回解析后的RSS Feed数据和错误信息
func (c *Client) GetRSSFeed(country string, collection Collection, category Category, limit int) (map[string]interface{}, error) {
	// 构建URL路径
	// 格式: https://itunes.apple.com/{country}/rss/{collection}/limit={limit}/json
	// 如果有category: https://itunes.apple.com/{country}/rss/{collection}/genre={category}/limit={limit}/json

	// Collection 现在是字符串类型，直接使用
	collectionStr := string(collection)
	if collectionStr == "" {
		collectionStr = "topfreeapplications" // 默认值
	}

	// 构建URL
	url := fmt.Sprintf("%s/%s/rss/%s", c.baseURL, country, collectionStr)

	// 如果指定了类别，添加genre参数
	if category != 0 {
		url += fmt.Sprintf("/genre=%d", category)
	}

	// 添加limit参数
	url += fmt.Sprintf("/limit=%d/json", limit)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return result, nil
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
