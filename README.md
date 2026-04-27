# App Store Scraper for Go

## 🌐 语言切换 | Language Switch

📖 **您正在阅读中文版本** | [English Version](README_EN.md)

---

## 目录

- [简介](#简介)
- [功能特点](#功能特点)
- [安装](#安装)
- [快速开始](#快速开始)
- [核心功能](#核心功能)
- [API 参考](#api-参考)
- [常量定义](#常量定义)
- [示例代码](#示例代码)
- [错误处理](#错误处理)
- [性能优化](#性能优化)
- [常见问题](#常见问题)
- [许可证](#许可证)

## 简介

App Store Scraper 是一个用 Go 语言编写的现代化库，专门用于从 iTunes/Mac App Store 抓取应用程序数据。该库提供了完整的类型安全保证，无需任何外部依赖，支持全球 100+ 个国家和地区的 App Store 数据获取。

本项目设计初衷是为了解决在 Go 语言生态中缺乏可靠的 App Store 数据抓取工具的问题。通过封装 Apple 官方的 iTunes Search API 和 RSS Feed API，我们能够获取应用的详细信息、用户评论、评分数据、排行榜信息、开发者应用列表等丰富的数据。

### 借鉴并参考
+ JavaScript: [facundoolano/app-store-scraper](https://github.com/facundoolano/app-store-scraper)
+ TypeScript: [plahteenlahti/app-store-scraper](https://github.com/plahteenlahti/app-store-scraper)

## 功能特点

### 完整类型安全

所有 API 方法都经过精心设计，返回明确的结构体类型。在编译时就能确保数据结构的正确性，避免了运行时类型断言的风险。使用本库时，IDE 能够提供完整的代码补全和类型检查功能。

### 无外部依赖

本库仅使用 Go 标准库实现，不依赖任何第三方包。这意味着您无需担心依赖冲突、版本兼容性问题或额外的安全风险。导入即可使用，部署简单便捷。

### 轻量级、易于使用

API 设计简洁直观，通过链式调用和合理的默认参数设置，可以仅用几行代码就完成复杂的数据抓取任务。每个方法都有清晰的职责边界，易于理解和使用。

### 多地区支持

支持全球 175+ 个国家和地区的 App Store，包括美国、中国（含港澳台）、日本、韩国、英国、德国、法国、俄罗斯、东南亚各国等。每个国家都有对应的常量定义，方便调用。

### 丰富的 API 方法

提供 10+ 个核心方法，涵盖应用获取、搜索、排行榜、评论、评分、相似应用、开发者应用、版本历史、搜索建议等完整的数据获取能力，满足各种业务场景需求。

### 性能优化

在需要获取多个应用详情时，库内部使用 goroutine 并发请求，显著提升数据获取效率。同时提供 FullDetail 选项，允许在快速模式和完整模式之间灵活切换。

### 评分直方图支持

除了平均评分外，还能获取完整的评分分布数据，包括 1 星到 5 星各等级的评价数量，帮助您更全面地了解应用的用户反馈情况。

## 安装

> 📝 **注意**：安装说明将在包正式发布到 Go 官方模块管理后再补充。

```bash
go get github.com/yililith/app-store-scraper
```

### 前置要求

- Go 1.25.6 或更高版本

### 模块引入

安装完成后，在您的 Go 项目中引入：

```go
import appstore "github.com/yililith/app-store-scraper"
```

## 快速开始

### 基本使用流程

使用本库抓取 App Store 数据的基本流程非常简单：

1. 导入 github.com/yililith/app-store-scraper 包
2. 创建 Scraper 实例
3. 调用相应的 API 方法获取数据
4. 处理返回的结果或错误

### 运行示例程序

项目中包含了一个交互式的测试菜单程序，可以让您快速体验各项功能：

```bash
# 进入示例目录
cd example

# 运行交互式测试菜单
go run main.go
```

### 交互式测试菜单

运行程序后会显示一个交互式菜单，您可以输入数字选择要测试的功能：

```
╔════════════════════════════════════════════════╗
║       App Store Scraper Go 测试菜单               ║
╠════════════════════════════════════════════════╣
║  1. 测试 App 应用详情                         ║
║  2. 测试 Search 搜索功能                     ║
║  3. 测试 List 排行榜功能                      ║
║  4. 测试 Reviews 评论功能                      ║
║  5. 测试 Ratings 评分分布                     ║
║  6. 测试 Similar 相似应用                      ║
║  7. 测试 Developer 开发者应用                ║
║  8. 测试 VersionHistory 版本历史             ║
║  9. 测试全部功能                               ║
║ 10. 退出                                          ║
╚════════════════════════════════════════════════╝
```

## 核心功能

### 1. 获取应用详情

通过应用 ID 或 Bundle ID 获取应用的完整信息，包括名称、开发者、价格、评分、描述、截图等。适合用于展示应用详情页面或进行应用数据统计分析。

### 2. 搜索应用

根据关键词搜索 App Store 中的应用程序。支持分页查询，可以获取搜索结果中的任意数量的应用信息。适合实现应用搜索功能或关键词热度分析。

### 3. 获取排行榜

获取各类别的应用排行榜数据，支持免费榜、付费榜、畅销榜等不同维度的排行。支持按应用类别筛选，可以获取特定分类下的排行榜信息。

### 4. 获取评论

获取应用的用户评论列表，支持按最新或最有用排序。包含评论标题、内容、评分、用户信息等完整数据。适合用于情感分析或用户反馈收集。

### 5. 评分分布

获取应用的评分分布直方图数据，包括各星级评价的数量分布。可以帮助您了解应用的整体用户满意度分布情况。

### 6. 相似应用

获取与目标应用相似的其他应用推荐。基于 App Store 的相似应用算法，返回可能感兴趣的应用列表。

### 7. 开发者应用

根据开发者 ID 获取该开发者发布的所有应用列表。适合用于开发者信息展示或竞品分析。

### 8. 版本历史

获取应用的完整版本更新历史，包括每个版本的版本号、发布日期和更新说明。适合用于版本追踪或应用更新分析。

### 9. 搜索建议

获取搜索关键词的建议列表，帮助用户输入更准确的搜索词或发现相关应用。

## API 参考

### Scraper 结构体

Scraper 是本库的核心结构体，提供所有 App Store 数据抓取功能。

```go
// 创建默认的 Scraper 实例
scraper := appstore.NewScraper()

// 使用自定义 HTTP 客户端创建 Scraper
scraper := appstore.NewScraperWithClient(client)
```

### App - 获取应用详情

根据应用 ID 或 Bundle ID 获取应用的详细信息。

```go
func (s *Scraper) App(opts AppOptions) (*AppInfo, error)
```

**参数说明**：

- `ID`：应用 ID（数字），与 AppID 二选一
- `AppID`：Bundle ID（字符串），与 ID 二选一
- `Country`：国家/地区代码，默认为美国
- `Lang`：语言代码，用于指定返回数据的语言
- `Ratings`：是否包含评分直方图数据

**返回值**：

- `*AppInfo`：应用详情结构体指针
- `error`：错误信息

**使用示例**：

```go
// 通过应用 ID 获取
app, err := scraper.App(appstore.AppOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})

// 通过 Bundle ID 获取并包含评分
app, err := scraper.App(appstore.AppOptions{
    AppID:   "com.king.candycrushsaga",
    Country: appstore.CountryUS,
    Ratings: true,
})

// 指定语言和地区
app, err := scraper.App(appstore.AppOptions{
    ID:      414478124,
    Country: appstore.CountryCN,
    Lang:    "zh_cn",
})
```

### Search - 搜索应用

根据关键词搜索应用程序。

```go
func (s *Scraper) Search(opts SearchOptions) ([]AppInfo, error)
```

**参数说明**：

- `Term`：搜索关键词（必填）
- `Country`：国家/地区代码
- `Lang`：语言代码
- `Num`：返回结果数量，默认 50，最大 200
- `Page`：页码，默认 1

**使用示例**：

```go
// 基本搜索
results, err := scraper.Search(appstore.SearchOptions{
    Term:    "minecraft",
    Country: appstore.CountryUS,
    Num:     10,
})

// 分页搜索
results, err := scraper.Search(appstore.SearchOptions{
    Term:    "finance banking",
    Country: appstore.CountryTW,
    Num:     10,
    Page:    2,
})
```

### List - 获取排行榜

获取应用排行榜数据，支持多种集合类型和类别筛选。

```go
func (s *Scraper) List(opts ListOptions) ([]AppInfo, error)
```

**参数说明**：

- `Collection`：应用集合类型，如免费榜、付费榜等
- `Category`：应用类别，如游戏、财务等
- `Country`：国家/地区代码
- `Num`：返回结果数量
- `FullDetail`：是否获取完整详情，true 会额外调用 API 获取详细信息

**使用示例**：

```go
// 获取免费财务应用排行榜（快速模式）
apps, err := scraper.List(appstore.ListOptions{
    Collection: appstore.TopFreeIOS,
    Category:   appstore.Finance,
    Country:    appstore.CountryPH,
    Num:        10,
    FullDetail: false,
})

// 获取免费游戏排行榜（完整详情）
apps, err := scraper.List(appstore.ListOptions{
    Collection: appstore.TopFreeIOS,
    Category:   appstore.Games,
    Country:    appstore.CountryTW,
    Num:        5,
    FullDetail: true,
})
```

### Reviews - 获取评论

获取应用的用户评论列表。

```go
func (s *Scraper) Reviews(opts ReviewsOptions) ([]Review, error)
```

**参数说明**：

- `ID`：应用 ID
- `AppID`：Bundle ID，与 ID 二选一
- `Country`：国家/地区代码
- `Sort`：排序方式，Recent（最新）或 Helpful（最有用）
- `Page`：页码，1-10

**使用示例**：

```go
reviews, err := scraper.Reviews(appstore.ReviewsOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
    Sort:    appstore.Recent,
    Page:    1,
})
```

### Ratings - 评分分布

获取应用的评分分布直方图。

```go
func (s *Scraper) Ratings(opts RatingsOptions) (*Ratings, error)
```

**参数说明**：

- `ID`：应用 ID（必填）
- `Country`：国家/地区代码

**使用示例**：

```go
ratings, err := scraper.Ratings(appstore.RatingsOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})

// 访问评分数据
fmt.Printf("总评分数: %d\n", ratings.Ratings)
fmt.Printf("5星: %d\n", ratings.Histogram.FiveStars)
fmt.Printf("4星: %d\n", ratings.Histogram.FourStars)
fmt.Printf("3星: %d\n", ratings.Histogram.ThreeStars)
fmt.Printf("2星: %d\n", ratings.Histogram.TwoStars)
fmt.Printf("1星: %d\n", ratings.Histogram.OneStar)
```

### Similar - 相似应用

获取与指定应用相似的其他应用。

```go
func (s *Scraper) Similar(opts SimilarOptions) ([]AppInfo, error)
```

**参数说明**：

- `ID`：应用 ID
- `AppID`：Bundle ID，与 ID 二选一
- `Country`：国家/地区代码

**使用示例**：

```go
apps, err := scraper.Similar(appstore.SimilarOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})
```

### Developer - 开发者应用

获取指定开发者的所有应用。

```go
func (s *Scraper) Developer(opts DeveloperOptions) ([]AppInfo, error)
```

**参数说明**：

- `DevID`：开发者 ID（必填）
- `Country`：国家/地区代码

**使用示例**：

```go
// 获取 OpenAI 开发者的应用
apps, err := scraper.Developer(appstore.DeveloperOptions{
    DevID:   284882218,
    Country: appstore.CountryUS,
})

// 获取 Google LLC 开发者的应用
apps, err := scraper.Developer(appstore.DeveloperOptions{
    DevID:   284882215,
    Country: appstore.CountryUS,
})
```

### VersionHistory - 版本历史

获取应用的版本更新历史。

```go
func (s *Scraper) VersionHistory(opts VersionHistoryOptions) ([]VersionHistory, error)
```

**参数说明**：

- `ID`：应用 ID（必填）
- `Country`：国家/地区代码

**使用示例**：

```go
history, err := scraper.VersionHistory(appstore.VersionHistoryOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})

for i, version := range history {
    fmt.Printf("版本 %s (%s)\n", version.VersionDisplay, version.ReleaseDate)
    if version.ReleaseNotes != "" {
        fmt.Printf("更新说明: %s\n", version.ReleaseNotes)
    }
}
```

### Suggest - 搜索建议

获取搜索关键词的建议。

```go
func (s *Scraper) Suggest(opts SuggestOptions) ([]Suggestion, error)
```

**参数说明**：

- `Term`：搜索关键词（必填）

**使用示例**：

```go
suggestions, err := scraper.Suggest(appstore.SuggestOptions{
    Term: "chat",
})
```

## 常量定义

### Collection - 应用集合

定义 App Store 中的各种应用集合类型。

```go
// iOS 应用集合
appstore.TopFreeIOS      // iOS 免费应用排行
appstore.TopGrossingIOS  // iOS 畅销应用排行
appstore.TopPaidIOS      // iOS 付费应用排行
appstore.NewFreeIOS     // iOS 新上架免费应用
appstore.NewPaidIOS     // iOS 新上架付费应用
appstore.NewIOS         // iOS 新上架应用

// iPad 应用集合
appstore.TopFreeiPad     // iPad 免费应用排行
appstore.TopGrossingiPad // iPad 畅销应用排行
appstore.TopPaidiPad     // iPad 付费应用排行

// Mac 应用集合
appstore.TopMac          // Mac 应用总榜
appstore.TopFreeMac      // Mac 免费应用排行
appstore.TopGrossingMac  // Mac 畅销应用排行
appstore.TopPaidMac      // Mac 付费应用排行
```

### Category - 应用类别

定义 App Store 中的应用类别，使用官方 Genre ID。

```go
// 主流类别
appstore.Games               // 游戏
appstore.Finance             // 财务
appstore.Business            // 商务
appstore.Education           // 教育
appstore.Entertainment      // 娱乐
appstore.HealthAndFitness    // 健康与健身
appstore.Lifestyle           // 生活
appstore.Music              // 音乐
appstore.News               // 新闻
appstore.PhotoAndVideo       // 摄影与视频
appstore.Productivity        // 效率
appstore.Shopping           // 购物
appstore.SocialNetworking    // 社交
appstore.Sports             // 体育
appstore.Travel             // 旅游
appstore.Utilities          // 工具
appstore.Weather           // 天气

// 游戏子类别
appstore.GamesAction        // 游戏-动作
appstore.GamesAdventure     // 游戏-冒险
appstore.GamesArcade        // 游戏-街机
appstore.GamesBoard         // 游戏-棋牌
appstore.GamesCard          // 游戏-卡牌
appstore.GamesCasino        // 游戏-娱乐场
appstore.GamesEducational   // 游戏-教育
appstore.GamesFamily        // 游戏-家庭
appstore.GamesKids          // 游戏-儿童
appstore.GamesMusic         // 游戏-音乐
appstore.GamesPuzzle        // 游戏-益智
appstore.GamesRacing        // 游戏-赛车
appstore.GamesRolePlaying   // 游戏-角色扮演
appstore.GamesSimulation    // 游戏-模拟
appstore.GamesSports        // 游戏-体育
appstore.GamesStrategy      // 游戏-策略
appstore.GamesTrivia        // 游戏-问答
appstore.GamesWord          // 游戏-文字
```

### Country - 国家/地区代码

定义支持的 App Store 国家/地区代码（部分列表）：

```go
// 亚洲
appstore.CountryCN  // 中国
appstore.CountryTW  // 台湾
appstore.CountryHK  // 香港
appstore.CountryMO  // 澳门
appstore.CountryJP  // 日本
appstore.CountryKR  // 韩国
appstore.CountrySG  // 新加坡
appstore.CountryTH  // 泰国
appstore.CountryMY  // 马来西亚
appstore.CountryID  // 印度尼西亚
appstore.CountryPH  // 菲律宾
appstore.CountryVN  // 越南
appstore.CountryIN  // 印度
appstore.CountryPK  // 巴基斯坦

// 欧洲
appstore.CountryUK  // 英国
appstore.CountryDE  // 德国
appstore.CountryFR  // 法国
appstore.CountryIT  // 意大利
appstore.CountryES  // 西班牙
appstore.CountryRU  // 俄罗斯
appstore.CountryNL  // 荷兰
appstore.CountryBE  // 比利时
appstore.CountryCH  // 瑞士
appstore.CountrySE  // 瑞典
appstore.CountryNO  // 挪威
appstore.CountryDK  // 丹麦
appstore.CountryFI  // 芬兰
appstore.CountryPL  // 波兰

// 北美
appstore.CountryUS  // 美国
appstore.CountryCA  // 加拿大
appstore.CountryMX  // 墨西哥

// 南美
appstore.CountryBR  // 巴西
appstore.CountryAR  // 阿根廷
appstore.CountryCL  // 智利
appstore.CountryCO  // 哥伦比亚
appstore.CountryPE  // 秘鲁

// 大洋洲
appstore.CountryAU  // 澳大利亚
appstore.CountryNZ  // 新西兰

// 中东/非洲
appstore.CountryAE  // 阿联酋
appstore.CountrySA  // 沙特阿拉伯
appstore.CountryIL  // 以色列
appstore.CountryZA  // 南非
```

### Sort - 排序方式

定义评论排序方式：

```go
appstore.Recent  // 按最新排序
appstore.Helpful // 按最有用排序
```

## 示例代码

### 完整使用示例

以下是一个完整的示例，展示如何获取应用信息、搜索应用和获取评论：

```go
package main

import (
    "fmt"
    appstore "github.com/yililith/app-store-scraper"
)

func main() {
    scraper := appstore.NewScraper()

    // 1. 获取应用详情
    app, err := scraper.App(appstore.AppOptions{
        ID:      553834731,
        Country: appstore.CountryUS,
    })
    if err != nil {
        fmt.Printf("获取应用详情失败: %v\n", err)
        return
    }
    fmt.Printf("应用名称: %s\n", app.Title)
    fmt.Printf("开发者: %s\n", app.Developer)
    fmt.Printf("版本: %s\n", app.Version)
    fmt.Printf("评分: %.1f (%d 条评价)\n", app.Score, app.Reviews)

    // 2. 搜索应用
    results, err := scraper.Search(appstore.SearchOptions{
        Term:    "productivity",
        Country: appstore.CountryUS,
        Num:     5,
    })
    if err != nil {
        fmt.Printf("搜索失败: %v\n", err)
        return
    }
    fmt.Printf("找到 %d 个相关应用\n", len(results))

    // 3. 获取评论
    reviews, err := scraper.Reviews(appstore.ReviewsOptions{
        ID:      553834731,
        Country: appstore.CountryUS,
        Sort:    appstore.Recent,
    })
    if err != nil {
        fmt.Printf("获取评论失败: %v\n", err)
        return
    }
    fmt.Printf("最新评论 (%d 条):\n", len(reviews))
    for i, review := range reviews {
        if i >= 3 {
            break
        }
        fmt.Printf("  [%d星] %s\n", review.Score, review.Title)
    }

    // 4. 获取评分分布
    ratings, err := scraper.Ratings(appstore.RatingsOptions{
        ID:      553834731,
        Country: appstore.CountryUS,
    })
    if err != nil {
        fmt.Printf("获取评分失败: %v\n", err)
        return
    }
    fmt.Printf("评分分布: %d 个评分\n", ratings.Ratings)
}
```

### 测试数据

项目中使用的测试应用数据：

| 应用名称 | App ID | 国家/地区 | 说明 |
|---------|--------|----------|------|
| Candy Crush | 553834731 | 美国 | 热门消除游戏 |
| MariBank | 1592249158 | 菲律宾 | 数字银行应用 |
| ChatGPT | 6448311069 | 美国 | OpenAI 的 AI 助手 |
| 微信 | 414478124 | 中国 | 腾讯社交应用 |

## 错误处理

本库定义了以下错误类型，建议在调用 API 后进行错误检查：

```go
// 参数错误
appstore.ErrInvalidParameter // 缺少必要参数或参数无效

// 数据错误
appstore.ErrNotFound          // 未找到请求的数据
appstore.ErrInvalidResponse   // 服务器响应格式错误

// 网络错误
// 标准 net/http 错误
```

**错误处理示例**：

```go
app, err := scraper.App(appstore.AppOptions{ID: 553834731})
if err != nil {
    switch err {
    case appstore.ErrInvalidParameter:
        fmt.Println("参数错误：请检查 App ID 或 Bundle ID")
    case appstore.ErrNotFound:
        fmt.Println("未找到该应用")
    case appstore.ErrInvalidResponse:
        fmt.Println("服务器响应错误")
    default:
        fmt.Printf("网络或其他错误: %v\n", err)
    }
    return
}
```

## 性能优化

### FullDetail 选项

List 方法支持两种数据获取模式：

**快速模式（FullDetail=false）**：
- 仅调用 RSS Feed API
- 返回应用的基本信息（ID、名称、图标、排名等）
- 响应速度快，适合需要大量数据的场景

**完整模式（FullDetail=true）**：
- 先调用 RSS Feed API 获取应用 ID 列表
- 再并发调用 Lookup API 获取每个应用的详细信息
- 返回完整应用数据，包括描述、截图等
- 适合需要展示详细信息的场景

### 并发优化

当使用 FullDetail=true 模式时，库内部会使用 goroutine 并发请求各应用的详细信息，显著提升数据获取效率。并发请求数根据目标应用数量动态调整，同时保证返回结果的顺序与排行榜顺序一致。

### 请求限制

请注意，Apple 的 iTunes API 可能有隐式的速率限制。建议：

- 避免短时间内发送大量请求
- 对于需要获取大量数据的场景，使用快速模式
- 在生产环境中添加适当的请求间隔

## 常见问题

### Q: 如何获取应用的 Bundle ID？

A: Bundle ID 是应用在 App Store 中的唯一标识符，可以在应用页面的 URL 中找到。例如，微信的 App Store URL 是 `https://apps.apple.com/cn/app/id414478124`，其中 `id414478124` 后面的数字部分（414478124）就是应用 ID。

### Q: 如何获取开发者的 ID？

A: 开发者 ID 可以通过访问开发者的 App Store 页面获取。例如，OpenAI 的开发者页面 URL 是 `https://apps.apple.com/us/developer/openai/id284882218`，其中 `id` 后面的数字（284882218）就是开发者 ID。

### Q: 为什么某些应用在特定国家找不到？

A: 部分应用可能只在特定国家/地区的 App Store 上架。请确保使用正确的国家代码，或者尝试使用应用的上架国家进行查询。

### Q: 如何处理请求超时？

A: Scraper 使用默认的 HTTP 客户端设置。如需自定义超时，可以创建自定义的 Client 并设置 Timeout：

```go
client := &appstore.Client{
    HTTPClient: &http.Client{
        Timeout: 30 * time.Second,
    },
}
scraper := appstore.NewScraperWithClient(client)
```

### 开发环境

- Go 1.25.6 或更高版本
- 代码格式化工具：gofmt
- 建议使用 VS Code 或 GoLand 等 IDE

### 代码规范

- 遵循 Go 官方代码规范
- 所有公共函数必须有文档注释
- 添加必要的单元测试

## 许可证

本项目采用 MIT 许可证开源。详见 [LICENSE](LICENSE) 文件。
