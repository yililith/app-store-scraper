# App Store Scraper 函数详解文档

## 目录

- [scraper.go 函数详解](#scrapergo-函数详解)
- [client.go 函数详解](#clientgo-函数详解)
- [types.go 类型定义](#typesgo-类型定义)
- [errors.go 错误定义](#errorsgo-错误定义)

---

## scraper.go 函数详解

### Scraper 结构体

```go
type Scraper struct {
    client *Client // HTTP客户端
}
```

**说明**：核心爬虫结构体，包含一个 HTTP 客户端实例，用于发送所有 App Store API 请求。

---

### NewScraper

```go
func NewScraper() *Scraper
```

**功能**：创建并返回一个新的 App Store 爬虫实例，使用默认配置。

**返回值**：`*Scraper` - 新的爬虫实例

**示例**：
```go
scraper := appstore.NewScraper()
```

---

### NewScraperWithClient

```go
func NewScraperWithClient(client *Client) *Scraper
```

**功能**：使用指定的 HTTP 客户端创建爬虫实例。

**参数**：
- `client` (*Client) - 自定义的 HTTP 客户端实例

**返回值**：`*Scraper` - 新的爬虫实例

**示例**：
```go
client := &appstore.Client{
    HTTPClient: &http.Client{
        Timeout: 60 * time.Second,
    },
}
scraper := appstore.NewScraperWithClient(client)
```

---

### App

```go
func (s *Scraper) App(opts AppOptions) (*AppInfo, error)
```

**功能**：根据应用 ID 或 Bundle ID 获取应用的详细信息。

**参数**：
- `opts` (AppOptions) - 获取选项，包含以下字段：
  - `ID` (int64) - 应用 ID，与 AppID 二选一
  - `AppID` (string) - Bundle ID，与 ID 二选一
  - `Country` (Country) - 国家/地区代码，默认为 us
  - `Lang` (string) - 语言代码，用于指定返回数据的语言
  - `Ratings` (bool) - 是否包含评分直方图数据
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `*AppInfo` - 应用详情结构体指针
- `error` - 错误信息

**API 调用**：
1. 调用 `/lookup` API 获取基本信息
2. 如果 API 未返回截图，爬取 App Store 页面提取截图
3. 如果 `Ratings=true`，额外调用 Ratings API 获取评分分布

**示例**：
```go
app, err := scraper.App(appstore.AppOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
    Ratings: true,
})
```

---

### Search

```go
func (s *Scraper) Search(opts SearchOptions) ([]AppInfo, error)
```

**功能**：根据关键词搜索 App Store 中的应用程序。

**参数**：
- `opts` (SearchOptions) - 搜索选项，包含以下字段：
  - `Term` (string) - 搜索关键词（必填）
  - `Country` (Country) - 国家/地区代码
  - `Lang` (string) - 语言代码
  - `Num` (int) - 返回结果数量，默认 50，最大 200
  - `Page` (int) - 页码，默认 1
  - `IdsOnly` (bool) - 是否只返回 ID（当前未使用）
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `[]AppInfo` - 应用列表
- `error` - 错误信息

**API 调用**：调用 `/search` API

**示例**：
```go
results, err := scraper.Search(appstore.SearchOptions{
    Term:    "minecraft",
    Country: appstore.CountryUS,
    Num:     10,
    Page:    1,
})
```

---

### List

```go
func (s *Scraper) List(opts ListOptions) ([]AppInfo, error)
```

**功能**：获取应用排行榜数据，支持多种集合类型和类别筛选。

**参数**：
- `opts` (ListOptions) - 列表选项，包含以下字段：
  - `Collection` (Collection) - 应用集合类型
  - `Category` (Category) - 应用类别
  - `Country` (Country) - 国家/地区代码
  - `Lang` (string) - 语言代码
  - `Num` (int) - 返回结果数量
  - `FullDetail` (bool) - 是否获取完整详情，默认 false
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `[]AppInfo` - 应用列表
- `error` - 错误信息

**API 调用**：
1. 调用 RSS Feed API 获取应用 ID 列表
2. 如果 `FullDetail=true`，并发调用 `/lookup` API 获取每个应用的详细信息

**示例**：
```go
apps, err := scraper.List(appstore.ListOptions{
    Collection: appstore.TopFreeIOS,
    Category:   appstore.Games,
    Country:    appstore.CountryTW,
    Num:        10,
    FullDetail: true,
})
```

---

### Developer

```go
func (s *Scraper) Developer(opts DeveloperOptions) ([]AppInfo, error)
```

**功能**：获取指定开发者的所有应用。

**参数**：
- `opts` (DeveloperOptions) - 开发者选项，包含以下字段：
  - `DevID` (int64) - 开发者 ID（必填）
  - `Country` (Country) - 国家/地区代码
  - `Lang` (string) - 语言代码
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `[]AppInfo` - 应用列表
- `error` - 错误信息

**API 调用**：调用 `/lookup?artistId={DevID}` API

**示例**：
```go
apps, err := scraper.Developer(appstore.DeveloperOptions{
    DevID:   284882218, // OpenAI
    Country: appstore.CountryUS,
})
```

---

### Reviews

```go
func (s *Scraper) Reviews(opts ReviewsOptions) ([]Review, error)
```

**功能**：获取应用的用户评论列表。

**参数**：
- `opts` (ReviewsOptions) - 评论选项，包含以下字段：
  - `ID` (int64) - 应用 ID
  - `AppID` (string) - Bundle ID，与 ID 二选一
  - `Country` (Country) - 国家/地区代码
  - `Sort` (Sort) - 排序方式：Recent（最新）或 Helpful（最有用）
  - `Page` (int) - 页码，1-10
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `[]Review` - 评论列表
- `error` - 错误信息

**API 调用**：调用 RSS Feed 评论 API

**注意**：
- 第一个条目通常是应用元数据，会被跳过
- 页码范围为 1-10

**示例**：
```go
reviews, err := scraper.Reviews(appstore.ReviewsOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
    Sort:    appstore.Recent,
    Page:    1,
})
```

---

### Ratings

```go
func (s *Scraper) Ratings(opts RatingsOptions) (*Ratings, error)
```

**功能**：获取应用的评分分布直方图数据。

**参数**：
- `opts` (RatingsOptions) - 评分选项，包含以下字段：
  - `ID` (int64) - 应用 ID（必填）
  - `Country` (Country) - 国家/地区代码
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `*Ratings` - 评分分布结构体指针
- `error` - 错误信息

**API 调用**：
1. 根据国家代码查找对应的 Store ID
2. 爬取 App Store 评论页面
3. 使用正则表达式从 HTML 中提取评分分布

**返回数据结构**：
```go
type Ratings struct {
    Ratings   int             // 总评分数量
    Histogram RatingHistogram  // 评分直方图
}

type RatingHistogram struct {
    OneStar    int // 1星评价数量
    TwoStars   int // 2星评价数量
    ThreeStars int // 3星评价数量
    FourStars  int // 4星评价数量
    FiveStars  int // 5星评价数量
}
```

**示例**：
```go
ratings, err := scraper.Ratings(appstore.RatingsOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})
fmt.Printf("总评分数: %d\n", ratings.Ratings)
fmt.Printf("5星: %d\n", ratings.Histogram.FiveStars)
```

---

### Similar

```go
func (s *Scraper) Similar(opts SimilarOptions) ([]AppInfo, error)
```

**功能**：获取与指定应用相似的其他应用。

**参数**：
- `opts` (SimilarOptions) - 相似应用选项，包含以下字段：
  - `ID` (int64) - 应用 ID
  - `AppID` (string) - Bundle ID，与 ID 二选一
  - `Country` (Country) - 国家/地区代码
  - `Lang` (string) - 语言代码
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `[]AppInfo` - 相似应用列表
- `error` - 错误信息

**API 调用**：
1. 爬取 App Store 应用页面
2. 使用正则表达式提取页面中相关的应用链接
3. 调用 `/lookup` API 获取相似应用的详细信息

**示例**：
```go
apps, err := scraper.Similar(appstore.SimilarOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})
```

---

### Suggest

```go
func (s *Scraper) Suggest(opts SuggestOptions) ([]Suggestion, error)
```

**功能**：获取搜索关键词的建议列表。

**参数**：
- `opts` (SuggestOptions) - 搜索建议选项，包含以下字段：
  - `Term` (string) - 搜索关键词（必填）
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `[]Suggestion` - 搜索建议列表
- `error` - 错误信息

**API 调用**：调用 iTunes 搜索建议 API，返回 XML 格式数据

**示例**：
```go
suggestions, err := scraper.Suggest(appstore.SuggestOptions{
    Term: "chat",
})
```

---

### VersionHistory

```go
func (s *Scraper) VersionHistory(opts VersionHistoryOptions) ([]VersionHistory, error)
```

**功能**：获取应用的版本更新历史。

**参数**：
- `opts` (VersionHistoryOptions) - 版本历史选项，包含以下字段：
  - `ID` (int64) - 应用 ID（必填）
  - `Country` (Country) - 国家/地区代码
  - `RequestOptions` (*RequestOptions) - 自定义请求选项

**返回值**：
- `[]VersionHistory` - 版本历史列表
- `error` - 错误信息

**API 调用**：
1. 爬取 App Store 应用页面
2. 使用正则表达式解析版本历史条目（`<article>` 标签）
3. 提取版本号、发布日期、更新说明

**返回数据结构**：
```go
type VersionHistory struct {
    VersionDisplay string // 版本号
    ReleaseDate    string // 发布日期
    ReleaseNotes   string // 发布说明
}
```

**示例**：
```go
history, err := scraper.VersionHistory(appstore.VersionHistoryOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})
for _, version := range history {
    fmt.Printf("版本 %s (%s)\n", version.VersionDisplay, version.ReleaseDate)
}
```

---

## client.go 函数详解

### Client 结构体

```go
type Client struct {
    httpClient *http.Client  // HTTP客户端实例
    baseURL    string        // API基础URL
    userAgent  string        // User-Agent头
}
```

**说明**：HTTP 客户端结构体，用于发送请求到 App Store API。

---

### NewClient

```go
func NewClient() *Client
```

**功能**：创建并返回一个新的 HTTP 客户端。

**默认配置**：
- 超时时间：30 秒
- User-Agent：Chrome 浏览器标识

**返回值**：`*Client` - 新的客户端实例

---

### SetTimeout

```go
func (c *Client) SetTimeout(timeout time.Duration)
```

**功能**：设置 HTTP 客户端的超时时间。

**参数**：
- `timeout` (time.Duration) - 超时时长

**示例**：
```go
client := appstore.NewClient()
client.SetTimeout(60 * time.Second)
```

---

### SetUserAgent

```go
func (c *Client) SetUserAgent(agent string)
```

**功能**：设置 User-Agent 请求头。

**参数**：
- `agent` (string) - User-Agent 字符串

---

### Get

```go
func (c *Client) Get(path string, params map[string]interface{}) (map[string]interface{}, error)
```

**功能**：发送 GET 请求到指定 URL，返回解析后的 JSON 数据。

**参数**：
- `path` (string) - 请求路径
- `params` (map[string]interface{}) - 查询参数

**返回值**：
- `map[string]interface{}` - 解析后的 JSON 数据
- `error` - 错误信息

**请求头设置**：
```
User-Agent: Chrome
Accept: application/json
Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
```

---

### GetRaw

```go
func (c *Client) GetRaw(url string) (string, error)
```

**功能**：发送 GET 请求并返回原始响应字符串。

**参数**：
- `url` (string) - 完整 URL

**返回值**：
- `string` - 原始响应内容
- `error` - 错误信息

**适用场景**：获取 XML、HTML 或原始 JSON 数据

---

### GetHTML

```go
func (c *Client) GetHTML(url string) (string, error)
```

**功能**：发送 GET 请求并返回 HTML 内容。

**参数**：
- `url` (string) - 完整 URL

**返回值**：
- `string` - HTML 内容
- `error` - 错误信息

---

### GetHTMLWithHeaders

```go
func (c *Client) GetHTMLWithHeaders(url string, headers map[string]string) (string, error)
```

**功能**：发送 GET 请求并返回 HTML 内容（带自定义请求头）。

**参数**：
- `url` (string) - 完整 URL
- `headers` (map[string]string) - 自定义请求头

**返回值**：
- `string` - HTML 内容
- `error` - 错误信息

**用途**：用于评分 API，需要设置 `X-Apple-Store-Front` 请求头

---

### GetRSSFeed

```go
func (c *Client) GetRSSFeed(country string, collection Collection, category Category, limit int) (map[string]interface{}, error)
```

**功能**：发送 GET 请求到 RSS Feed 并返回解析后的数据。

**参数**：
- `country` (string) - 国家代码
- `collection` (Collection) - 集合类型
- `category` (Category) - 类别 ID
- `limit` (int) - 返回数量

**返回值**：
- `map[string]interface{}` - 解析后的 RSS Feed 数据
- `error` - 错误信息

---

### BuildParams

```go
func BuildParams(opts interface{}) map[string]interface{}
```

**功能**：将结构体选项转换为 API 查询参数 map。

**参数**：
- `opts` (interface{}) - 选项结构体（SearchOptions、RatingsOptions 等）

**返回值**：map[string]interface{} - 查询参数字典

**支持类型**：
- SearchOptions
- RatingsOptions
- SimilarOptions

---

### 辅助解析函数

以下函数用于将 `interface{}` 类型安全转换为具体类型：

#### parseInterfaceArray

```go
func parseInterfaceArray(data interface{}) []string
```

**功能**：将 interface{} 数组转换为字符串数组。

**支持类型**：`[]interface{}`、`[]string`

#### parseInterfaceInt

```go
func parseInterfaceInt(data interface{}) int
```

**功能**：将 interface{} 转换为 int。

**支持类型**：`float64`、`int`、`string`

#### parseInterfaceInt64

```go
func parseInterfaceInt64(data interface{}) int64
```

**功能**：将 interface{} 转换为 int64。

**支持类型**：`float64`、`int`、`int64`、`string`

#### parseInterfaceFloat64

```go
func parseInterfaceFloat64(data interface{}) float64
```

**功能**：将 interface{} 转换为 float64。

**支持类型**：`float64`、`int`、`int64`、`string`

#### parseInterfaceBool

```go
func parseInterfaceBool(data interface{}) bool
```

**功能**：将 interface{} 转换为 bool。

**支持类型**：`bool`、`string`（"true"/"1"）、`int`

#### parseInterfaceString

```go
func parseInterfaceString(data interface{}) string
```

**功能**：将 interface{} 转换为 string。

**支持类型**：`string`、`float64`、`int`、`int64`

---

## types.go 类型定义

### 常量类型

#### Collection

```go
type Collection string
```

**定义**：App Store 中的应用集合类型。

| 常量 | 值 | 说明 |
|------|------|------|
| TopFreeIOS | "topfreeapplications" | iOS 免费应用排行 |
| TopGrossingIOS | "topgrossingapplications" | iOS 畅销应用排行 |
| TopPaidIOS | "toppaidapplications" | iOS 付费应用排行 |
| NewIOS | "newapplications" | iOS 新上架应用 |
| NewFreeIOS | "newfreeapplications" | iOS 新上架免费应用 |
| NewPaidIOS | "newpaidapplications" | iOS 新上架付费应用 |
| TopFreeiPad | "topfreeipadapplications" | iPad 免费应用排行 |
| TopGrossingiPad | "topgrossingipadapplications" | iPad 畅销应用排行 |
| TopPaidiPad | "toppaidipadapplications" | iPad 付费应用排行 |
| TopMac | "topmacapps" | Mac 应用总榜 |
| TopFreeMac | "topfreemacapps" | Mac 免费应用排行 |
| TopGrossingMac | "topgrossingmacapps" | Mac 畅销应用排行 |
| TopPaidMac | "toppaidmacapps" | Mac 付费应用排行 |

#### Category

```go
type Category int
```

**定义**：App Store 应用类别，使用官方 Genre ID。

| 常量 | 值 | 说明 |
|------|------|------|
| Games | 6014 | 游戏 |
| Finance | 6015 | 财务 |
| Business | 6000 | 商务 |
| Education | 6017 | 教育 |
| Entertainment | 6016 | 娱乐 |
| HealthAndFitness | 6013 | 健康与健身 |
| Lifestyle | 6012 | 生活 |
| Music | 6011 | 音乐 |
| News | 6009 | 新闻 |
| PhotoAndVideo | 6008 | 摄影与视频 |
| Productivity | 6007 | 效率 |
| Shopping | 6024 | 购物 |
| SocialNetworking | 6005 | 社交 |
| Sports | 6004 | 体育 |
| Travel | 6003 | 旅游 |
| Utilities | 6002 | 工具 |
| Weather | 6001 | 天气 |

**游戏子类别**：GamesAction (7001)、GamesAdventure (7002)、GamesArcade (7003) 等

#### Country

```go
type Country string
```

**定义**：App Store 国家/地区代码，支持 175+ 国家。

| 常量 | 说明 |
|------|------|
| CountryUS | 美国（默认） |
| CountryCN | 中国 |
| CountryJP | 日本 |
| CountryKR | 韩国 |
| CountryUK | 英国 |
| CountryDE | 德国 |
| CountryFR | 法国 |
| CountryTW | 台湾 |
| CountryHK | 香港 |
| CountrySG | 新加坡 |

#### Sort

```go
type Sort string
```

**定义**：评论排序方式。

| 常量 | 值 | 说明 |
|------|------|------|
| Recent | "mostRecent" | 按最新排序 |
| Helpful | "mostHelpful" | 按最有用排序 |

---

### 数据结构

#### AppInfo

```go
type AppInfo struct {
    ID                    int64             // 应用ID
    AppID                 string            // Bundle标识符
    Title                 string            // 应用名称
    URL                   string            // App Store链接
    Description           string            // 应用描述
    Icon                  string            // 图标URL
    Genres                []string          // 类别列表
    GenreIDs              []string          // 类别ID列表
    PrimaryGenre          string            // 主要类别
    PrimaryGenreID        string            // 主要类别ID
    ContentRating         string            // 内容评级
    Languages             []string          // 支持的语言
    Size                  string            // 文件大小
    RequiredOsVersion     string            // 最低系统版本要求
    Released              string            // 发布日期
    Updated               string            // 更新日期
    ReleaseNotes          string            // 发布说明
    Version               string            // 当前版本号
    Price                 float64           // 价格
    Currency              string            // 货币类型
    Free                  bool              // 是否免费
    DeveloperID           int64             // 开发者ID
    Developer             string            // 开发者名称
    DeveloperURL          string            // 开发者iTunes页面链接
    DeveloperWebsite      string            // 开发者网站
    Score                 float64           // 当前版本平均评分
    Reviews               int               // 当前版本评分总数
    CurrentVersionScore   float64           // 当前版本平均评分
    CurrentVersionReviews  int               // 当前版本评分总数
    ScreenshotURLs        []string          // iPhone/iPod截图URL列表
    IpadScreenshots       []string          // iPad截图URL列表
    AppletvScreenshots    []string          // Apple TV截图URL列表
    SupportedDevices      []string          // 支持的设备
    Histogram             *RatingHistogram  // 评分直方图
}
```

#### Review

```go
type Review struct {
    ID            string  // 评论ID
    UserName      string  // 用户名
    UserReviewURL string  // 用户评论页面URL
    Version       string  // 应用版本
    Score         int     // 评分（1-5）
    Title         string  // 评论标题
    Text          string  // 评论内容
    Updated       string  // 更新日期
}
```

#### Ratings

```go
type Ratings struct {
    Ratings   int              // 总评分数量
    Histogram RatingHistogram  // 评分直方图
}

type RatingHistogram struct {
    OneStar    int  // 1星评价数量
    TwoStars   int  // 2星评价数量
    ThreeStars int  // 3星评价数量
    FourStars  int  // 4星评价数量
    FiveStars  int  // 5星评价数量
}
```

#### VersionHistory

```go
type VersionHistory struct {
    VersionDisplay string  // 版本号
    ReleaseDate    string  // 发布日期
    ReleaseNotes   string  // 发布说明
}
```

#### Suggestion

```go
type Suggestion struct {
    Term       string  // 搜索关键词
    SecondTerm string  // 第二个关键词
}
```

---

### 选项结构体

| 结构体 | 用途 |
|--------|------|
| AppOptions | App 方法选项 |
| SearchOptions | Search 方法选项 |
| ListOptions | List 方法选项 |
| DeveloperOptions | Developer 方法选项 |
| ReviewsOptions | Reviews 方法选项 |
| RatingsOptions | Ratings 方法选项 |
| SimilarOptions | Similar 方法选项 |
| SuggestOptions | Suggest 方法选项 |
| VersionHistoryOptions | VersionHistory 方法选项 |
| RequestOptions | 自定义请求选项（请求头） |

---

## errors.go 错误定义

```go
var ErrInvalidParameter = errors.New("无效参数：必须提供必要的参数")
var ErrNotFound = errors.New("未找到：请求的资源不存在")
var ErrInvalidResponse = errors.New("无效响应：无法解析API返回的数据")
var ErrNetworkError = errors.New("网络错误：无法连接到App Store")
var ErrRateLimit = errors.New("速率限制：请稍后重试")
var ErrServerError = errors.New("服务器错误：App Store服务器出现问题")
```

| 错误 | 触发场景 |
|------|----------|
| ErrInvalidParameter | 缺少必要参数或参数无效 |
| ErrNotFound | 应用/资源在 App Store 中不存在 |
| ErrInvalidResponse | API 返回的数据格式无法解析 |
| ErrNetworkError | 网络连接失败 |
| ErrRateLimit | 请求过于频繁被限制 |
| ErrServerError | Apple 服务器错误 |

**错误处理示例**：
```go
app, err := scraper.App(appstore.AppOptions{ID: 553834731})
if err != nil {
    switch err {
    case appstore.ErrInvalidParameter:
        fmt.Println("参数错误")
    case appstore.ErrNotFound:
        fmt.Println("未找到该应用")
    case appstore.ErrInvalidResponse:
        fmt.Println("服务器响应错误")
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
    return
}
```
