package appstore

import "errors"

// 错误定义
// 预定义应用程序中可能遇到的常见错误

// ErrInvalidParameter 无效参数错误
// 当传递给函数的参数不符合要求时返回
var ErrInvalidParameter = errors.New("无效参数：必须提供必要的参数")

// ErrNotFound 未找到错误
// 当请求的资源在App Store中不存在时返回
var ErrNotFound = errors.New("未找到：请求的资源不存在")

// ErrInvalidResponse 无效响应错误
// 当API返回的数据格式不符合预期时返回
var ErrInvalidResponse = errors.New("无效响应：无法解析API返回的数据")

// ErrNetworkError 网络错误
// 当网络请求失败时返回
var ErrNetworkError = errors.New("网络错误：无法连接到App Store")

// ErrRateLimit 速率限制错误
// 当请求过于频繁被限制时返回
var ErrRateLimit = errors.New("速率限制：请稍后重试")

// ErrServerError 服务器错误
// 当App Store服务器返回错误时返回
var ErrServerError = errors.New("服务器错误：App Store服务器出现问题")
