# App Store Scraper 项目全局解析文档

本文档面向想完整理解本项目的人，目标不是只讲“功能是什么”，而是讲“代码为什么这样写、每一段代码如何连接、关键行在做什么”。

文档采用四层结构：

1. 项目整体定位
2. 目录与模块关系
3. 请求链路与数据流
4. 按文件逐行解读

---

## 1. 项目整体定位

这个项目是一个纯 Go 标准库实现的 App Store 抓取库，核心职责是把 Apple 的几类数据源统一包装成一个高层 API：

- iTunes Lookup / Search JSON API
- App Store RSS Feed JSON API
- App Store 网页 HTML
- Search Hints XML 接口

项目的核心设计思路是：

- `Client` 负责 HTTP 请求和底层响应解析。
- `Scraper` 负责业务语义封装，把多个底层接口组合成对外 API。
- `types.go` 负责统一输入输出模型。
- `api_response.go` 负责对 Apple 原始响应建模。
- `example/main.go` 负责提供交互式演示入口。

---

## 2. 目录与模块关系

当前项目结构如下：

```text
app-store-scraper/
├── api_response.go
├── client.go
├── errors.go
├── go.mod
├── scraper.go
├── types.go
├── README.md
├── README_EN.md
├── docs/
│   ├── ARCHITECTURE.md
│   └── PROJECT_GLOBAL_LINE_BY_LINE_ANALYSIS.md
└── example/
    └── main.go
```

### 文件职责概览

- `go.mod`
  - 定义 Go 模块名 `github.com/yililith/app-store-scraper`。
- `errors.go`
  - 定义统一错误值，供整个库复用。
- `types.go`
  - 定义常量、输入参数结构体、输出结构体。
- `api_response.go`
  - 定义 Apple API 原始 JSON/RSS 结构体。
- `client.go`
  - 封装 HTTP 请求、参数拼装和底层响应反序列化。
- `scraper.go`
  - 实现业务 API，如 `App`、`Search`、`List`、`Reviews`、`Ratings`。
- `example/main.go`
  - 提供命令行测试菜单，方便手动验证功能。

---

## 3. 请求链路与数据流

### 3.1 最核心的调用链

以 `scraper.App()` 为例：

1. 调用方构造 `AppOptions`
2. `Scraper.App()` 做参数校验和默认值填充
3. `Scraper.App()` 构造查询参数
4. 调用 `Client.GetLookup()`
5. `Client.GetLookup()` 发起 HTTP GET
6. 把返回 JSON 反序列化到 `LookupResponse`
7. `Scraper.parseAppResult()` 把原始 `AppResult` 转成对外的 `AppInfo`
8. 如截图缺失，再回退到 HTML 页面抓图
9. 如要求评分分布，再额外调用 `Ratings()`

### 3.2 List / Similar / Developer 的共同模式

这三个接口的模式都是：

- 先拿到一组 app id
- 再批量调用 `lookup()` 把 id 转成 `AppInfo`

其中：

- `List()` 的 id 来自 RSS Feed
- `Similar()` 的 id 来自 App Store 页面 HTML
- `Developer()` 的 id 实际是开发者 id，也通过 lookup 路径查询

### 3.3 底层接口分工

- JSON API
  - `GetLookup()`
  - `GetRSSFeedTyped()`
  - `GetReviewFeed()`
- HTML 抓取
  - `GetHTML()`
  - `GetHTMLWithHeaders()`
- 原始文本
  - `GetRaw()`

---

## 4. 按文件逐行解读

以下内容按源码文件顺序展开。

---

## 4.1 `go.mod`

### 第 1 行

- `module github.com/yililith/app-store-scraper`
  - 定义模块导入路径。
  - 外部项目通过这个路径引入本库。

### 第 3 行

- `go 1.25.6`
  - 声明项目使用的 Go 版本。
  - 这里也意味着泛型语法、标准库行为都以该版本为准。

---

## 4.2 `errors.go`

### 第 5-6 行

- 注释说明这一段是统一错误定义。
  - 这能让外部调用者用 `switch err` 做稳定判断。

### 第 8-10 行

- `ErrInvalidParameter`
  - 表示输入参数不合法。
  - 例如 `AppOptions` 里同时没有 `ID` 和 `AppID` 时会返回它。

### 第 12-14 行

- `ErrNotFound`
  - 表示查询结果为空或目标资源不存在。
  - 属于业务层“没找到”，不是网络层错误。

### 第 16-18 行

- `ErrInvalidResponse`
  - 预留给响应格式异常的情况。
  - 当前项目里真正抛出的解析错误很多直接用 `fmt.Errorf("解析JSON失败: %w", err)` 包起来，而不是直接返回它。

### 第 20-22 行

- `ErrNetworkError`
  - 预定义网络错误值。
  - 目前代码没有系统性封装成这个错误，属于后续可继续统一的点。

### 第 24-26 行

- `ErrRateLimit`
  - 预定义速率限制错误。
  - 当前没有根据 HTTP 状态码专门映射到这里。

### 第 28-30 行

- `ErrServerError`
  - 预定义服务端错误。
  - 当前也没有在 `5xx` 状态时统一转换到这个错误。

### 小结

- `errors.go` 的价值主要是建立统一错误语义。
- 目前项目已经定义了这些错误，但还没完全做到“所有错误都统一落盘到这几种类型”。

---

## 4.3 `api_response.go`

这个文件的核心作用是“描述 Apple 返回的数据长什么样”，它是外部原始结构，不直接暴露给业务调用方。

### 第 3-6 行：`LookupResponse`

- 第 3 行定义顶层结构体。
- 第 4 行 `ResultCount int`
  - 对应 JSON 的 `resultCount`。
- 第 5 行 `Results []AppResult`
  - 真正的应用数据在这里。

这是 lookup/search 接口的统一顶层返回格式。

### 第 8-44 行：`AppResult`

这是最重要的原始响应结构体。

#### 第 9-12 行

- `TrackID`、`TrackName`、`BundleID`、`TrackViewURL`
  - 分别对应应用 id、名称、包名和页面链接。

#### 第 13-15 行

- `Description`、`ArtworkURL512`、`ArtworkURL100`
  - 应用描述和不同尺寸图标。

#### 第 16-19 行

- `Genres`、`GenreIDs`、`PrimaryGenreName`、`PrimaryGenreID`
  - 描述应用分类。
  - 这里的 `PrimaryGenreID` 已修正为 `int64`，因为 Apple 现在返回的是数字。

#### 第 20-27 行

- 内容分级、语言、包体积、最低系统、发布日期、当前版本发布日期、发布说明、版本号。
  - 这些是应用详情展示的主干字段。

#### 第 28-37 行

- `Price`、`Currency`
  - 价格相关。
- `ArtistID`、`ArtistName`、`ArtistViewURL`、`SellerURL`
  - 开发者相关。
- `AverageUserRating`、`UserRatingCount`
  - 总体评分及数量。
- `AverageUserRatingForCurrentVersion`、`UserRatingCountForCurrentVersion`
  - 当前版本评分及数量。

#### 第 38-43 行

- `ScreenshotURLs`、`IpadScreenshotURLs`、`AppletvScreenshotURLs`
  - 各端截图。
- `SupportedDevices`
  - 支持设备列表。
- `Kind`、`WrapperType`
  - 用于判断是不是软件类型。

### 第 46-52 行：RSS 顶层结构

- `RSSFeedResponse`
  - RSS 接口顶层对象。
- `RSSFeed`
  - `entry` 列表容器。

### 第 54-74 行：RSS 条目结构

- `RSSEntry`
  - 只保留了本项目当前会用到的字段：
    - `id`
    - `name`
    - `artist`

#### 第 60-66 行

- `RSSID` / `RSSIDAttributes`
  - Apple 把应用 id 嵌套在 `id.attributes.im:id` 中。

#### 第 68-74 行

- `RSSName` 与 `RSSArtist`
  - 两者都只是 `label` 字段包了一层。

### 第 76-129 行：评论 RSS 结构

这一段专门服务 `Reviews()`。

#### 第 76-82 行

- `ReviewFeedResponse` 和 `ReviewFeed`
  - 仍然是 `feed.entry` 的结构。

#### 第 84-92 行

- `ReviewEntry`
  - 包含评论 id、作者、评分、标题、内容、更新时间、评论对应版本。

#### 第 94-129 行

- 一系列小结构体本质上是在贴合 Apple RSS 的 XML/JSON 风格。
  - 比如 `ReviewImRating.Label`
  - 比如 `ReviewAuthor.Name.Label`

### 小结

- `api_response.go` 是“外部协议层模型”。
- 该文件尽量忠实表达 Apple 的原始字段，不做业务语义裁剪。
- 真正对外给使用者的模型在 `types.go`。

---

## 4.4 `types.go`

这个文件非常大，但逻辑很清晰：先定义枚举，再定义领域模型，再定义输入选项。

### 第 3-21 行：`Collection`

- 第 5 行定义 `Collection string`
  - 用字符串表示 App Store 集合名。

#### 第 7-20 行

- 这部分每一行都是一个 Apple RSS 集合常量。
  - `TopMac`
  - `TopFreeMac`
  - `TopGrossingMac`
  - `TopPaidMac`
  - `NewIOS`
  - `NewFreeIOS`
  - `NewPaidIOS`
  - `TopFreeIOS`
  - `TopFreeiPad`
  - `TopGrossingIOS`
  - `TopGrossingiPad`
  - `TopPaidIOS`
  - `TopPaidiPad`

这些值会被 `GetRSSFeedTyped()` 直接拼进 URL。

### 第 23-73 行：`Category`

- 第 25 行定义 `Category int`
  - 直接对应 Apple genre id。

#### 第 27-72 行

- 每一行一个分类 id。
- 如：
  - `Business = 6000`
  - `Finance = 6015`
  - `Games = 6014`
- 游戏子类继续延伸到 `7001-7019`。

这部分本质是“类型安全的数字枚举”：

- 优点是调用方不需要记住数字。
- 风险是 Apple 若调整 genre id，这里需要同步更新。

### 第 75-82 行：`Sort`

- `Sort string`
  - 评论排序方式。
- `Recent` 和 `Helpful`
  - 分别对应 Apple 接口需要的字符串值。

### 第 84-92 行：`Device`

- 目前设备常量定义了：
  - `All`
  - `iPad`
  - `Mac`

这部分当前在主流程里没有大量使用，更像扩展预留。

### 第 94-246 行：`Country`

这是整个文件最长的一段。

#### 第 96 行

- `type Country string`
  - 国家代码本质是字符串，但单独定义成类型后，接口签名更清晰。

#### 第 98-245 行

- 每一行定义一个国家/地区代码常量。
- 例如：
  - `CountryCN = "cn"`
  - `CountryTW = "tw"`
  - `CountryUS = "us"`
  - `CountryPH = "ph"`

逐行理解时，可以把这一大段看成一份静态映射表：

- 左边是更易用的 Go 常量名
- 右边是 Apple URL 和参数需要的国家代码

之所以占很多行，是因为项目把大部分支持国家都直接穷举出来了，这样 IDE 自动补全体验很好。

### 第 248-285 行：`AppInfo`

这是最重要的对外输出模型。

#### 第 250 行

- `type AppInfo struct`
  - 表示库最终返回给调用方的应用信息。

#### 第 251-260 行

- 基础标识信息：
  - `ID`
  - `AppID`
  - `Title`
  - `URL`
  - `Description`
  - `Icon`
  - `Genres`
  - `GenreIDs`
  - `PrimaryGenre`
  - `PrimaryGenreID`

#### 第 261-268 行

- 应用属性信息：
  - `ContentRating`
  - `Languages`
  - `Size`
  - `RequiredOsVersion`
  - `Released`
  - `Updated`
  - `ReleaseNotes`
  - `Version`

#### 第 269-284 行

- 商业与开发者信息：
  - `Price`
  - `Currency`
  - `Free`
  - `DeveloperID`
  - `Developer`
  - `DeveloperURL`
  - `DeveloperWebsite`
- 评分信息：
  - `Score`
  - `Reviews`
  - `CurrentVersionScore`
  - `CurrentVersionReviews`
- 多端截图与设备能力：
  - `ScreenshotURLs`
  - `IpadScreenshots`
  - `AppletvScreenshots`
  - `SupportedDevices`
- 附加扩展：
  - `Histogram`

这说明 `AppInfo` 是一个“归一化后的业务输出结构体”。

### 第 287-292 行：`SearchResult`

- 当前项目主流程里并没有大量直接使用它。
- 它代表“搜索返回应用列表”的通用模型。

### 第 294-305 行：`Review`

- 对外暴露的评论结构。
- 与 `ReviewEntry` 的差别在于：
  - `ReviewEntry` 保留 Apple 原始层级
  - `Review` 则已经变成平铺字段，便于消费

### 第 307-321 行：评分结构

#### `RatingHistogram`

- 用 `json:"1"` 到 `json:"5"` 映射星级计数。

#### `Ratings`

- 包含总评分数量和直方图。

### 第 323-329 行：`VersionHistory`

- 用于存储版本号、发布日期、发布说明。

### 第 331-336 行：`Suggestion`

- 搜索建议模型。
- `SecondTerm` 当前解析逻辑没有真正填充，属于预留字段。

### 第 338-342 行：`RequestOptions`

- 当前只支持一类扩展能力：自定义请求头。
- 这是整个项目可配置性的主要入口。

### 第 344-430 行：输入选项结构

这一段是对外 API 的输入模型定义。

#### `AppOptions`

- 支持按 `ID` 或 `AppID` 查询。
- 可附带 `Country`、`Lang`、`Ratings`。

#### `SearchOptions`

- 核心字段：
  - `Term`
  - `Country`
  - `Lang`
  - `Num`
  - `Page`
- `IdsOnly` 目前没有在搜索主流程中真正启用。

#### `ListOptions`

- 关键在于：
  - `Collection`
  - `Category`
  - `Country`
  - `Num`
  - `FullDetail`

#### `DeveloperOptions`

- 以 `DevID` 为主键查询开发者应用。

#### `ReviewsOptions`

- 支持 `ID` 或 `AppID`。
- 支持 `Sort` 和 `Page`。

#### `RatingsOptions`

- 最精简，只要 `ID` 和 `Country`。

#### `SimilarOptions`

- 与 `AppOptions` 类似，但用于“相似应用”。

#### `SuggestOptions`

- 只需要 `Term`。

#### `VersionHistoryOptions`

- 只需要 `ID` 和 `Country`。

### 小结

- `types.go` 解决的是“对外 API 该长什么样”。
- 它把大量 Apple 细节包装成更稳定的业务模型和参数模型。

---

## 4.5 `client.go`

这个文件是网络层核心。

### 第 13-19 行：`Client`

- 第 15 行定义 `Client struct`
  - 这是底层 HTTP 客户端封装。

字段解释：

- `httpClient`
  - 真正发请求的标准库客户端。
- `baseURL`
  - 统一 API 根地址，默认 `https://itunes.apple.com`。
- `userAgent`
  - 所有请求公用 UA。

### 第 21-31 行：`NewClient`

- 第 23 行构造一个默认客户端。
- 第 25-27 行设置 30 秒超时。
- 第 28 行设定 `baseURL`。
- 第 29 行设定一个浏览器风格 UA，避免 Apple 接口对默认 Go UA 的兼容性问题。

### 第 33-37 行：`SetTimeout`

- 对外暴露超时控制入口。
- 只改底层 `http.Client.Timeout`。

### 第 39-43 行：`SetUserAgent`

- 允许调用方替换请求头中的 UA。

### 第 45-94 行：`Get`

这个函数是一个通用 JSON GET 方法。

#### 第 49 行

- `func (c *Client) Get(path string, params map[string]string) (map[string]json.RawMessage, error)`
  - 输入是路径和查询参数。
  - 输出是“顶层 JSON 对象映射”。

#### 第 50 行

- `fullURL := c.baseURL + path`
  - 把相对路径拼成完整 URL。

#### 第 52-58 行

- 如果参数不为空，就创建 `url.Values` 并编码到查询串。
- 这里使用 `map[string]string`，避免了过去 `interface{}` 参数的不确定性。

#### 第 60-63 行

- 构造 GET 请求。
- 如果 URL 非法或其他问题，会在这里失败。

#### 第 65-67 行

- 设置默认请求头：
  - `User-Agent`
  - `Accept: application/json`
  - `Accept-Language`

#### 第 69-73 行

- 发请求并确保 `Body.Close()`。

#### 第 75-77 行

- 对 HTTP 状态码做最基本检查，只接受 `200 OK`。

#### 第 79-82 行

- 读取响应体字节。

#### 第 84-86 行

- 如果响应为空，或者看起来像 HTML 而不是 JSON，就直接报错。
- 这是一个很实用的防守式判断。

#### 第 88-91 行

- 反序列化成 `map[string]json.RawMessage`。
- 适合需要保留原始 JSON 片段但延迟解析的场景。

#### 第 93 行

- 返回解析结果。

说明：

- 当前主业务流基本不直接依赖这个通用方法，而更常用 `GetLookup()`。

### 第 96-138 行：`BuildParams`

这是参数构建辅助函数。

#### 第 98 行

- 使用泛型 `BuildParams[T any]`。
- 目的是统一处理多个 options 类型。

#### 第 99 行

- 初始化返回参数映射 `map[string]string`。

#### 第 101-135 行

- 使用类型分支，根据入参真实类型构造不同查询参数。

##### `SearchOptions`

- `term`
- `country`
- `lang`
- `limit`

##### `RatingsOptions`

- `id`
- `country`

##### `SimilarOptions`

- `id`
- `country`

#### 第 137 行

- 返回构建完成的参数表。

### 第 140-172 行：`GetHTML`

- 专用于抓 HTML 页面。
- 与 `Get()` 的差别在于：
  - `Accept` 头不同
  - 不做 JSON 解析
  - 直接返回字符串

### 第 174-211 行：`GetHTMLWithHeaders`

- 和 `GetHTML()` 几乎一样。
- 额外支持调用方传入自定义请求头。
- `Ratings()` 会用它注入 `X-Apple-Store-Front`。

### 第 213-259 行：`GetLookup`

这是目前最关键的 JSON 请求函数。

#### 第 214 行

- `func (c *Client) GetLookup(path string, params map[string]string) (*LookupResponse, error)`
  - 直接返回结构化 `LookupResponse`。

#### 第 215 行

- 把相对路径拼成完整 URL。

#### 第 217-223 行

- 若额外传了参数，则再编码到 URL。
- 这里加了 `params != nil` 判断，避免空 map / nil map 情况下逻辑混淆。

#### 第 225-258 行

- 与 `Get()` 的流程一致：
  - 发请求
  - 检查状态
  - 读 body
  - 检查异常 HTML
  - 反序列化

区别在于：

- 这里直接反序列化到 `LookupResponse`
- 调用方不用再自己解顶层字段

### 第 261-306 行：`GetRSSFeedTyped`

- 负责请求排行榜 RSS JSON。

关键点：

- 第 263-266 行
  - 如果没传 collection，回退到 `topfreeapplications`
- 第 268-274 行
  - URL 格式是 `/{country}/rss/{collection}/genre={category}/limit={limit}/json`
- 第 300-305 行
  - 直接解成 `RSSFeedResponse`

### 第 308-343 行：`GetReviewFeed`

- 专门请求评论 feed。
- 路径格式固定：
  - `/{country}/rss/customerreviews/page={page}/id={id}/sortby={sortby}/json`

### 第 345-377 行：`GetRaw`

- 返回原始字符串。
- `Suggest()` 会用它抓 XML 文本。

### 小结

- `client.go` 是底层网络访问层。
- 它负责把“请求是怎么发的”这件事从 `scraper.go` 中剥离出来。

---

## 4.6 `scraper.go`

这是项目的业务核心文件。

### 第 12-16 行：`Scraper`

- `Scraper` 只持有一个 `client *Client`。
- 这说明它本身是一个高层门面，而不是重网络实现。

### 第 18-31 行：构造函数

#### `NewScraper()`

- 创建默认客户端并挂到 `Scraper` 上。

#### `NewScraperWithClient()`

- 允许外部注入自定义客户端。
- 适合自定义超时、代理、请求头等场景。

### 第 33-100 行：`App`

这是“获取应用详情”的实现。

#### 第 36-39 行

- 如果 `ID` 和 `AppID` 都为空，直接返回 `ErrInvalidParameter`。

#### 第 41-43 行

- 国家为空时默认 `CountryUS`。

#### 第 45-53 行

- 统一把查询入口转为：
  - 路径 `/lookup`
  - 参数值 `idValue`

这里做了一层“兼容数字 id 或 bundle id”的抽象。

#### 第 55-62 行

- 构造查询参数：
  - `id`
  - `country`
  - `entity=software`
  - 可选 `lang`

#### 第 64-67 行

- 调用底层 `GetLookup()`。

#### 第 69-71 行

- `result.Results` 为空时视为未找到。

#### 第 73-76 行

- 如果结果不是软件类型，同样视为未找到。

#### 第 78 行

- 把原始 `AppResult` 转成 `AppInfo`。

#### 第 80-86 行

- 如果 lookup 没返回截图，就回退到 HTML 页面里抓图。
- 这是一个兜底逻辑，增强结果完整性。

#### 第 88-97 行

- 如果调用者要求附带评分分布，就继续调用 `Ratings()`。
- 这里对 `Ratings()` 错误做了“软失败”：
  - 出错不影响主接口成功
  - 只是缺少 `Histogram`

#### 第 99 行

- 返回最终 `AppInfo` 指针。

### 第 102-148 行：`Search`

#### 第 105-118 行

- 参数校验和默认值填充：
  - `Term` 必填
  - `Num` 默认 50
  - `Page` 默认 1
  - `Country` 默认 US

#### 第 120-123 行

- 通过 `BuildParams()` 构造基础参数，再手动补充：
  - `term`
  - `media=software`
  - `entity=software`

#### 第 125-128 行

- 调用 `/search` 接口。

#### 第 130-135 行

- 只保留 `item.Kind == "software"` 的结果。

#### 第 137-145 行

- 本地实现分页裁剪。
- 说明 Apple 搜索接口返回结果后，这里再按页数切片。

#### 第 147 行

- 返回当前页结果。

### 第 150-190 行：`List`

#### 第 153-162 行

- 设置默认：
  - `Num=50`
  - `Country=US`
  - `Collection=TopFreeIOS`

#### 第 164-167 行

- `limit` 上限被裁到 200。
- 这通常是为了贴合 Apple RSS 最大限制。

#### 第 169-172 行

- 调用 RSS Feed 接口。

#### 第 174-176 行

- 如果 feed 为空，返回空列表。

#### 第 178-183 行

- 从每个 RSS `entry` 中提取 `im:id` 并转成 `int64`。

#### 第 185-187 行

- 如果解析后一个 id 都没有，就返回空数组。

#### 第 189 行

- 调用 `lookup()` 批量补全详情。

说明：

- 当前实现无论 `FullDetail` 是否为真，都会走 `lookup()` 获取完整详情。
- 这意味着 `FullDetail` 字段存在，但实现上还没有真正分支控制。

### 第 192-206 行：`Developer`

#### 第 195-202 行

- 校验开发者 id。
- 填默认国家。

#### 第 204-205 行

- 把 `DevID` 包成 `[]int64` 后复用 `lookup()`。

这说明当前实现把开发者查询也套进 lookup 流程里。

### 第 208-253 行：`Reviews`

#### 第 211-227 行

- 参数校验：
  - `ID`/`AppID` 至少一个
  - `Page` 默认 1
  - 页码限定在 1-10
  - 默认国家 US
  - 默认排序 `Recent`

#### 第 229-239 行

- 如果只传了 `AppID`，先调用 `App()` 把 bundle id 转成数值 id。

#### 第 241-243 行

- 调用评论 feed 接口。

#### 第 246-250 行

- 从 `entry[1]` 开始遍历，而不是 `entry[0]`。
- 原因通常是 Apple 评论 feed 的第一项不是普通评论，而是应用元信息。

#### 第 252 行

- 返回平铺后的评论列表。

### 第 255-295 行：`Ratings`

#### 第 258-267 行

- 校验 app id。
- 填国家默认值。

#### 第 269-271 行

- 先根据国家算 `storeFront`。
- 再拼评分页面 URL。

#### 第 273-281 行

- 构造请求头。
- 关键是 `X-Apple-Store-Front`，这会影响 Apple 返回的商店上下文。
- 若有自定义 headers，会覆盖或追加进去。

#### 第 283-287 行

- 用带请求头的 HTML 请求方式抓页面。

#### 第 289-291 行

- 空 HTML 视为未找到。

#### 第 293-294 行

- 交给 `parseRatings()` 用正则解析评分分布。

### 第 297-335 行：`Similar`

#### 第 300-318 行

- 参数校验和根据 `AppID` 反查数值 id 的逻辑，与 `Reviews()` 类似。

#### 第 321-325 行

- 访问应用详情页 HTML。
- 抓页面失败时这里直接返回空数组而不是错误，属于“容错优先”策略。

#### 第 328-332 行

- 提取相似应用 id。
- 没提到就返回空数组。

#### 第 334 行

- 再复用 `lookup()` 批量转成 `AppInfo`。

### 第 337-357 行：`Suggest`

#### 第 340-344 行

- 只要求 `Term` 非空。

#### 第 346-347 行

- 用 `url.QueryEscape()` 转义关键词。

#### 第 349-352 行

- 请求 Apple hints 接口，返回 XML 文本。

#### 第 355-356 行

- 用 `parseSuggestXML()` 提取建议词。

### 第 359-384 行：`VersionHistory`

#### 第 362-370 行

- 校验 app id 并填默认国家。

#### 第 373-379 行

- 访问应用详情页 HTML。

#### 第 382-383 行

- 用 `parseVersionHistory()` 做正则提取。

### 第 386-432 行：`parseAppResult`

这是原始模型到业务模型的映射器。

#### 第 389-393 行

- 优先使用 `ArtworkURL512`，没有则回退到 `ArtworkURL100`。

#### 第 395-429 行

- 将 `AppResult` 逐字段映射到 `AppInfo`。
- 其中：
  - `Free` 不是接口原生字段，而是根据 `Price == 0` 推导出来的。

#### 第 431 行

- 返回 `AppInfo`。

### 第 434-450 行：`parseReviewEntry`

- 把带很多嵌套 `Label` 的 `ReviewEntry` 平铺成 `Review`。

### 第 452-460 行：`parseRatingStr`

- 小工具函数。
- 空字符串返回 0。
- 能转整数就返回，失败也返回 0。

### 第 463-470 行：`parseRSSEntry`

- 废弃函数。
- 当前只是占位，直接返回空 `AppInfo{}`。

这说明项目早期可能尝试过直接解析 RSS 条目成应用数据，后来改成统一通过 `lookup()` 补全详情。

### 第 472-477 行：`ScreenshotResult`

- 存储三类截图结果：
  - iPhone
  - iPad
  - Apple TV

### 第 479-499 行：`scrapeScreenshots`

- 请求 App Store HTML 页面。
- 调用 `extractScreenshots()` 三次分别抽三类截图。

### 第 501-542 行：`extractScreenshots`

这是 HTML 抽图正则逻辑。

#### 第 505-507 行

- 定位某类截图对应的 `<ul>` 容器。

#### 第 515-517 行

- 提取 `source` 标签里的 `srcset`。

#### 第 523-536 行

- 把 `srcset` 拆成多个候选 URL。
- 选出 URL 并做标准化替换。
- 用 `contains()` 去重。

### 第 544-552 行：`contains`

- 字符串切片去重辅助函数。

### 第 554-597 行：`parseRatings`

这部分用正则从 HTML 中提取评分统计。

#### 第 555-558 行

- 初始化 `Ratings`。

#### 第 560-567 行

- 找总评分数量。

#### 第 570-593 行

- 找各星级条目。
- 按出现顺序推导 5 星到 1 星。

### 第 599-628 行：`extractSimilarAppIds`

#### 第 603-605 行

- 先匹配所有 `/app/` 链接。

#### 第 613-624 行

- 再从链接里提取 `/id123456` 数字 id。
- 排除当前应用自身。
- 用 `containsInt64()` 去重。

### 第 630-638 行：`containsInt64`

- `int64` 版本的去重判断函数。

### 第 640-667 行：`lookup`

这是批量详情补全核心。

#### 第 642-646 行

- 把 `[]int64` 转成逗号分隔字符串。

#### 第 648-651 行

- 直接构造 `/lookup?...` 查询路径。
- 如果有语言参数，再补 `lang`。

#### 第 653-656 行

- 调用 `GetLookup()`。
- 这里传 `nil` 参数，表示路径里已经包含查询串。

#### 第 658-664 行

- 遍历结果，只保留 `software` 类型条目。

#### 第 666 行

- 返回应用详情列表。

### 第 669-693 行：`parseSuggestXML`

#### 第 675-677 行

- 用正则匹配所有 `<string>xxx</string>`。

#### 第 678-690 行

- 用 `seen` 去重。
- 每个唯一词封装成 `Suggestion`。

### 第 695-734 行：`storeId`

这是一张国家代码到 Apple storefront id 的映射表。

#### 第 697-728 行

- 每一组键值对都在声明：
  - `"us" -> 143441`
  - `"cn" -> 143465`
  - `"ph" -> 143474`

#### 第 730-733 行

- 查询不到时默认返回美国 storefront。

### 第 736-787 行：`parseVersionHistory`

这一段也是基于 HTML 正则解析。

#### 第 742-744 行

- 先匹配每一个版本卡片 `<article>`。

#### 第 752-757 行

- 提取版本号。

#### 第 760-765 行

- 提取发布日期。

#### 第 768-775 行

- 提取发布说明，并去掉 HTML 标签。

#### 第 777-783 行

- 只要版本号或日期不为空，就追加到结果。

#### 第 786 行

- 返回版本历史数组。

### 小结

- `scraper.go` 是“业务编排层”。
- 它把多个接口组合成可直接调用的高层方法。
- 同时也暴露出项目的主要技术路线：
  - JSON API 优先
  - HTML/正则作为补充

---

## 4.7 `example/main.go`

这个文件不是库本身的一部分，而是演示和手工测试入口。

### 第 14-57 行：`main`

#### 第 15 行

- 创建 `bufio.Scanner` 读取终端输入。

#### 第 17-56 行

- 一个无限循环的菜单系统：
  - 显示菜单
  - 读取用户输入
  - 按选项调用不同测试函数
  - 等待用户回车后回菜单

#### 第 26-50 行

- `switch choice` 控制所有测试入口。

### 第 59-74 行：`printMenu`

- 单纯打印 CLI 菜单。
- 这部分对库逻辑没有影响，只是交互壳。

### 第 76-101 行：`testAppDetail`

- 用 3 个固定测试用例验证：
  - 美国的 Candy Crush
  - 菲律宾的 MariBank
  - 中国的微信

它能验证：

- lookup 请求
- App 详情映射
- 多国家支持

### 第 103-131 行：`testSearch`

- 测试搜索接口。
- 使用不同关键词和不同国家。

### 第 133-157 行：`testList`

- 测试榜单接口。
- 当前只配了一个菲律宾财务榜单案例。

### 第 159-194 行：`testReviews`

- 测试评论抓取。
- 输出前三条评论标题。

### 第 196-225 行：`testRatings`

- 测试评分分布。
- 注意第 213-218 行在检查 `err` 前就读取了 `ratings.Histogram`，这在错误时理论上有空指针风险，是一个值得修补的小问题。

### 第 227-258 行：`testSimilar`

- 测试相似应用抓取。
- 输出前三个结果。

### 第 260-291 行：`testDeveloper`

- 测试开发者应用列表。
- 使用 OpenAI 和 Google LLC。

### 第 293-318 行：`testVersionHistory`

- 测试版本历史抓取。

### 第 320-390 行：`testAll`

- 顺序串联所有核心能力。
- 这是最接近“冒烟测试”的函数。

### 第 392-399 行：`printJSON`

- 泛型版本 JSON pretty print 工具。
- 当前主流程里不一定频繁使用，但对临时调试很有帮助。

### 小结

- `example/main.go` 的主要价值不是“实现业务”，而是“展示库怎么用”和“提供手工回归入口”。

---

## 5. 关键横切观察

### 5.1 项目的分层是清晰的

- 协议层：`api_response.go`
- 网络层：`client.go`
- 业务层：`scraper.go`
- 领域模型层：`types.go`
- 演示层：`example/main.go`

这是这个项目最值得保留的结构优点。

### 5.2 目前主要依赖正则解析 HTML

以下能力不是纯 JSON API，而是依赖 HTML 页面结构：

- `Ratings()`
- `Similar()`
- `VersionHistory()`
- `scrapeScreenshots()`

优点：

- 无第三方依赖
- 实现轻量

缺点：

- Apple 页面结构变化时，容易失效

### 5.3 `FullDetail` 字段尚未完全落地

`ListOptions.FullDetail` 已存在，但 `List()` 当前实际实现始终通过 `lookup()` 获取详情，因此还没有真正分成“快速模式”和“完整模式”。

### 5.4 错误体系定义了，但还不够统一

虽然 `errors.go` 中定义了多种错误，但当前代码中很多地方仍直接返回 `fmt.Errorf(...)` 包装错误，而不是统一映射到错误常量。

### 5.5 参数类型已经比早期版本更清晰

当前 `client.go` 中的查询参数已经是 `map[string]string`，不再依赖早期的动态类型，这使得：

- URL 拼装更稳定
- 调试更直观
- 接口语义更清晰

---

## 6. 推荐阅读顺序

如果你是第一次读这个项目，建议按以下顺序看源码：

1. `types.go`
   - 先理解有哪些输入输出模型
2. `api_response.go`
   - 再理解外部原始响应结构
3. `client.go`
   - 接着看底层 HTTP 请求怎么发
4. `scraper.go`
   - 再看高层业务 API 如何编排
5. `example/main.go`
   - 最后用示例程序验证理解

---

## 7. 一句话总结

这个项目本质上是一个“以 `Scraper` 为门面、以 `Client` 为网络层、以 `types/api_response` 为模型层”的 App Store 数据抓取库，优先使用官方 JSON 接口，在官方接口不足时再通过 HTML + 正则补齐信息。
