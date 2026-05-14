package appstore

// Collection 应用集合常量
// 定义App Store中的各种应用集合（使用字符串ID）
type Collection string

const (
	TopMac          Collection = "topmacapps"                  // Mac应用总榜
	TopFreeMac      Collection = "topfreemacapps"              // Mac免费应用排行
	TopGrossingMac  Collection = "topgrossingmacapps"          // Mac畅销应用排行
	TopPaidMac      Collection = "toppaidmacapps"              // Mac付费应用排行
	NewIOS          Collection = "newapplications"             // iOS新上架应用
	NewFreeIOS      Collection = "newfreeapplications"         // iOS新上架免费应用
	NewPaidIOS      Collection = "newpaidapplications"         // iOS新上架付费应用
	TopFreeIOS      Collection = "topfreeapplications"         // iOS免费应用排行
	TopFreeiPad     Collection = "topfreeipadapplications"     // iPad免费应用排行
	TopGrossingIOS  Collection = "topgrossingapplications"     // iOS畅销应用排行
	TopGrossingiPad Collection = "topgrossingipadapplications" // iPad畅销应用排行
	TopPaidIOS      Collection = "toppaidapplications"         // iOS付费应用排行
	TopPaidiPad     Collection = "toppaidipadapplications"     // iPad付费应用排行
)

// Category 应用类别常量
// 定义App Store中的各种应用类别（使用App Store官方Genre ID）
type Category int

const (
	Books                  Category = 6018 // 图书
	Business               Category = 6000 // 商务
	Catalogs               Category = 6022 // 商品目录
	Education              Category = 6017 // 教育
	Entertainment          Category = 6016 // 娱乐
	Finance                Category = 6015 // 财务 ← 修正：之前错误定义为24
	FoodAndDrink           Category = 6023 // 美食与餐饮
	Games                  Category = 6014 // 游戏 ← 修正：之前错误定义为1
	GamesAction            Category = 7001 // 游戏-动作
	GamesAdventure         Category = 7002 // 游戏-冒险
	GamesArcade            Category = 7003 // 游戏-街机
	GamesBoard             Category = 7004 // 游戏-棋牌
	GamesCard              Category = 7005 // 游戏-卡牌
	GamesCasino            Category = 7006 // 游戏-娱乐场
	GamesDice              Category = 7007 // 游戏-骰子
	GamesEducational       Category = 7008 // 游戏-教育
	GamesFamily            Category = 7009 // 游戏-家庭
	GamesKids              Category = 7010 // 游戏-儿童
	GamesMusic             Category = 7011 // 游戏-音乐
	GamesPuzzle            Category = 7012 // 游戏-益智
	GamesRacing            Category = 7013 // 游戏-赛车
	GamesRolePlaying       Category = 7014 // 游戏-角色扮演
	GamesSimulation        Category = 7015 // 游戏-模拟
	GamesSports            Category = 7016 // 游戏-体育
	GamesStrategy          Category = 7017 // 游戏-策略
	GamesTrivia            Category = 7018 // 游戏-问答
	GamesWord              Category = 7019 // 游戏-文字
	HealthAndFitness       Category = 6013 // 健康与健身
	Lifestyle              Category = 6012 // 生活
	Medical                Category = 6020 // 医疗
	Music                  Category = 6011 // 音乐
	Navigation             Category = 6010 // 导航
	News                   Category = 6009 // 新闻
	PhotoAndVideo          Category = 6008 // 摄影与视频
	Productivity           Category = 6007 // 效率
	Reference              Category = 6006 // 参考
	Shopping               Category = 6024 // 购物
	SocialNetworking       Category = 6005 // 社交
	Sports                 Category = 6004 // 体育
	Stickers               Category = 6025 // 表情贴纸
	Travel                 Category = 6003 // 旅游
	Utilities              Category = 6002 // 工具
	Weather                Category = 6001 // 天气
	DeveloperTools         Category = 0    // 开发者工具（通用）
	MagazinesAndNewspapers Category = 6021 // 杂志与报刊
)

// Sort 排序选项常量
// 定义评论和列表的排序方式
type Sort string

const (
	Recent  Sort = "mostRecent"  // 按最新排序
	Helpful Sort = "mostHelpful" // 按最有用排序
)

// Device 设备类型常量
// 定义支持的设备类型
type Device string

const (
	All  Device = "software"     // 所有软件
	iPad Device = "iPadSoftware" // iPad
	Mac  Device = "macSoftware"  // Mac
)

// Country 国家/地区代码
// 定义支持的App Store国家代码
type Country string

const (
	CountryDZ Country = "dz" // 阿尔及利亚
	CountryAO Country = "ao" // 安哥拉
	CountryAI Country = "ai" // 安圭拉
	CountryAG Country = "ag" // 安提瓜和巴布达
	CountryAR Country = "ar" // 阿根廷
	CountryAM Country = "am" // 亚美尼亚
	CountryAU Country = "au" // 澳大利亚
	CountryAT Country = "at" // 奥地利
	CountryAZ Country = "az" // 阿塞拜疆
	CountryBS Country = "bs" // 巴哈马
	CountryBH Country = "bh" // 巴林
	CountryBB Country = "bb" // 巴巴多斯
	CountryBY Country = "by" // 白俄罗斯
	CountryBE Country = "be" // 比利时
	CountryBZ Country = "bz" // 伯利兹
	CountryBJ Country = "bj" // 贝宁
	CountryBM Country = "bm" // 百慕大
	CountryBO Country = "bo" // 玻利维亚
	CountryBW Country = "bw" // 博茨瓦纳
	CountryBR Country = "br" // 巴西
	CountryVG Country = "vg" // 英属维尔京群岛
	CountryBN Country = "bn" // 文莱
	CountryBG Country = "bg" // 保加利亚
	CountryBF Country = "bf" // 布基纳法索
	CountryCA Country = "ca" // 加拿大
	CountryKY Country = "ky" // 开曼群岛
	CountryTD Country = "td" // 乍得
	CountryCL Country = "cl" // 智利
	CountryCN Country = "cn" // 中国
	CountryCO Country = "co" // 哥伦比亚
	CountryCR Country = "cr" // 哥斯达黎加
	CountryCI Country = "ci" // 科特迪瓦
	CountryHR Country = "hr" // 克罗地亚
	CountryCY Country = "cy" // 塞浦路斯
	CountryCZ Country = "cz" // 捷克
	CountryDK Country = "dk" // 丹麦
	CountryDM Country = "dm" // 多米尼克
	CountryDO Country = "do" // 多米尼加共和国
	CountryEC Country = "ec" // 厄瓜多尔
	CountryEG Country = "eg" // 埃及
	CountrySV Country = "sv" // 萨尔瓦多
	CountryEE Country = "ee" // 爱沙尼亚
	CountryFJ Country = "fj" // 斐济
	CountryFI Country = "fi" // 芬兰
	CountryFR Country = "fr" // 法国
	CountryGM Country = "gm" // 冈比亚
	CountryDE Country = "de" // 德国
	CountryGH Country = "gh" // 加纳
	CountryGR Country = "gr" // 希腊
	CountryGD Country = "gd" // 格林纳达
	CountryGT Country = "gt" // 危地马拉
	CountryGW Country = "gw" // 几内亚比绍
	CountryGY Country = "gy" // 圭亚那
	CountryHN Country = "hn" // 洪都拉斯
	CountryHK Country = "hk" // 香港
	CountryHU Country = "hu" // 匈牙利
	CountryIS Country = "is" // 冰岛
	CountryIN Country = "in" // 印度
	CountryID Country = "id" // 印度尼西亚
	CountryIE Country = "ie" // 爱尔兰
	CountryIL Country = "il" // 以色列
	CountryIT Country = "it" // 意大利
	CountryJM Country = "jm" // 牙买加
	CountryJP Country = "jp" // 日本
	CountryJO Country = "jo" // 约旦
	CountryKZ Country = "kz" // 哈萨克斯坦
	CountryKE Country = "ke" // 肯尼亚
	CountryKR Country = "kr" // 韩国
	CountryKW Country = "kw" // 科威特
	CountryLV Country = "lv" // 拉脱维亚
	CountryLB Country = "lb" // 黎巴嫩
	CountryLT Country = "lt" // 立陶宛
	CountryLU Country = "lu" // 卢森堡
	CountryMO Country = "mo" // 澳门
	CountryMK Country = "mk" // 北马其顿
	CountryMG Country = "mg" // 马达加斯加
	CountryMW Country = "mw" // 马拉维
	CountryMY Country = "my" // 马来西亚
	CountryML Country = "ml" // 马里
	CountryMT Country = "mt" // 马耳他
	CountryMR Country = "mr" // 毛里塔尼亚
	CountryMU Country = "mu" // 毛里求斯
	CountryMX Country = "mx" // 墨西哥
	CountryFM Country = "fm" // 密克罗尼西亚联邦
	CountryMD Country = "md" // 摩尔多瓦
	CountryMN Country = "mn" // 蒙古
	CountryMS Country = "ms" // 蒙特塞拉特
	CountryMZ Country = "mz" // 莫桑比克
	CountryNA Country = "na" // 纳米比亚
	CountryNP Country = "np" // 尼泊尔
	CountryNL Country = "nl" // 荷兰
	CountryNZ Country = "nz" // 新西兰
	CountryNI Country = "ni" // 尼加拉瓜
	CountryNE Country = "ne" // 尼日尔
	CountryNG Country = "ng" // 尼日利亚
	CountryNO Country = "no" // 挪威
	CountryOM Country = "om" // 阿曼
	CountryPK Country = "pk" // 巴基斯坦
	CountryPW Country = "pw" // 帕劳
	CountryPA Country = "pa" // 巴拿马
	CountryPG Country = "pg" // 巴布亚新几内亚
	CountryPY Country = "py" // 巴拉圭
	CountryPE Country = "pe" // 秘鲁
	CountryPH Country = "ph" // 菲律宾 ← 你需要的国家！
	CountryPL Country = "pl" // 波兰
	CountryPT Country = "pt" // 葡萄牙
	CountryQA Country = "qa" // 卡塔尔
	CountryRO Country = "ro" // 罗马尼亚
	CountryRU Country = "ru" // 俄罗斯
	CountryKN Country = "kn" // 圣基茨和尼维斯
	CountryLC Country = "lc" // 圣卢西亚
	CountryVC Country = "vc" // 圣文森特和格林纳丁斯
	CountryST Country = "st" // 圣多美和普林西比
	CountrySA Country = "sa" // 沙特阿拉伯
	CountrySN Country = "sn" // 塞内加尔
	CountrySC Country = "sc" // 塞舌尔
	CountrySL Country = "sl" // 塞拉利昂
	CountrySG Country = "sg" // 新加坡
	CountrySK Country = "sk" // 斯洛伐克
	CountrySI Country = "si" // 斯洛文尼亚
	CountrySB Country = "sb" // 所罗门群岛
	CountryZA Country = "za" // 南非
	CountryES Country = "es" // 西班牙
	CountryLK Country = "lk" // 斯里兰卡
	CountrySR Country = "sr" // 苏里南
	CountrySZ Country = "sz" // 斯威士兰
	CountrySE Country = "se" // 瑞典
	CountryCH Country = "ch" // 瑞士
	CountryTW Country = "tw" // 台湾
	CountryTJ Country = "tj" // 塔吉克斯坦
	CountryTZ Country = "tz" // 坦桑尼亚
	CountryTH Country = "th" // 泰国
	CountryTN Country = "tn" // 突尼斯
	CountryTR Country = "tr" // 土耳其
	CountryTM Country = "tm" // 土库曼斯坦
	CountryTC Country = "tc" // 特克斯和凯科斯群岛
	CountryUG Country = "ug" // 乌干达
	CountryUA Country = "ua" // 乌克兰
	CountryAE Country = "ae" // 阿联酋
	CountryUK Country = "gb" // 英国
	CountryUS Country = "us" // 美国 ← 默认
	CountryUY Country = "uy" // 乌拉圭
	CountryUZ Country = "uz" // 乌兹别克斯坦
	CountryVE Country = "ve" // 委内瑞拉
	CountryVN Country = "vn" // 越南
	CountryYE Country = "ye" // 也门
	CountryZW Country = "zw" // 津巴布韦
)

// AppInfo 应用详细信息结构体
// 包含应用的所有基本信息，如名称、开发者、价格、评分等
type AppInfo struct {
	ID                    int64            `json:"id"`                    // 应用ID
	AppID                 string           `json:"appId"`                 // Bundle标识符
	Title                 string           `json:"title"`                 // 应用名称
	URL                   string           `json:"url"`                   // App Store链接
	Description           string           `json:"description"`           // 应用描述
	Icon                  string           `json:"icon"`                  // 图标URL
	Genres                []string         `json:"genres"`                // 类别列表
	GenreIDs              []string         `json:"genreIds"`              // 类别ID列表
	PrimaryGenre          string           `json:"primaryGenre"`          // 主要类别
	PrimaryGenreID        int64    `json:"primaryGenreId"`        // 主要类别ID
	ContentRating         string           `json:"contentRating"`         // 内容评级
	Languages             []string         `json:"languages"`             // 支持的语言
	Size                  string           `json:"size"`                  // 文件大小
	RequiredOsVersion     string           `json:"requiredOsVersion"`     // 最低系统版本要求
	Released              string           `json:"released"`              // 发布日期
	Updated               string           `json:"updated"`               // 更新日期
	ReleaseNotes          string           `json:"releaseNotes"`          // 发布说明
	Version               string           `json:"version"`               // 当前版本号
	Price                 float64          `json:"price"`                 // 价格
	Currency              string           `json:"currency"`              // 货币类型
	Free                  bool             `json:"free"`                  // 是否免费
	DeveloperID           int64            `json:"developerId"`           // 开发者ID
	Developer             string           `json:"developer"`             // 开发者名称
	DeveloperURL          string           `json:"developerUrl"`          // 开发者iTunes页面链接
	DeveloperWebsite      string           `json:"developerWebsite"`      // 开发者网站
	Score                 float64          `json:"score"`                 // 当前版本平均评分
	Reviews               int              `json:"reviews"`               // 当前版本评分总数
	CurrentVersionScore   float64          `json:"currentVersionScore"`   // 当前版本平均评分（与Score相同）
	CurrentVersionReviews int              `json:"currentVersionReviews"` // 当前版本评分总数
	ScreenshotURLs        []string         `json:"screenshotUrls"`        // iPhone/iPod截图URL列表
	IpadScreenshots       []string         `json:"ipadScreenshots"`       // iPad截图URL列表
	AppletvScreenshots    []string         `json:"appletvScreenshots"`    // Apple TV截图URL列表
	SupportedDevices      []string         `json:"supportedDevices"`      // 支持的设备
	Histogram             *RatingHistogram `json:"histogram,omitempty"`   // 评分直方图
}

// SearchResult 搜索结果结构体
// 用于存储搜索返回的应用列表信息
type SearchResult struct {
	ResultCount int       `json:"resultCount"` // 搜索结果数量
	Results     []AppInfo `json:"results"`     // 应用列表
}

// Review 评论结构体
// 包含用户评论的详细信息
type Review struct {
	ID            string `json:"id"`       // 评论ID
	UserName      string `json:"userName"` // 用户名
	UserReviewURL string `json:"userUrl"`  // 用户评论页面URL
	Version       string `json:"version"`  // 应用版本
	Score         int    `json:"score"`    // 评分（1-5）
	Title         string `json:"title"`    // 评论标题
	Text          string `json:"text"`     // 评论内容
	Updated       string `json:"updated"`  // 更新日期
}

// RatingHistogram 评分分布直方图
// 用于存储各星级评分的分布情况
type RatingHistogram struct {
	OneStar    int `json:"1"` // 1星评价数量
	TwoStars   int `json:"2"` // 2星评价数量
	ThreeStars int `json:"3"` // 3星评价数量
	FourStars  int `json:"4"` // 4星评价数量
	FiveStars  int `json:"5"` // 5星评价数量
}

// Ratings 评分响应结构体
type Ratings struct {
	Ratings   int             `json:"ratings"`   // 总评分数量
	Histogram RatingHistogram `json:"histogram"` // 评分直方图
}

// VersionHistory 版本历史记录
// 包含应用的历史版本信息
type VersionHistory struct {
	VersionDisplay string `json:"versionDisplay"` // 版本号
	ReleaseDate    string `json:"releaseDate"`    // 发布日期
	ReleaseNotes   string `json:"releaseNotes"`   // 发布说明
}

// Suggestion 搜索建议结构体
// 用于存储搜索建议返回的关键词
type Suggestion struct {
	Term       string `json:"term"`       // 搜索关键词
	SecondTerm string `json:"secondTerm"` // 第二个关键词
}

// RequestOptions 请求选项
// 用于配置HTTP请求的行为
type RequestOptions struct {
	Headers map[string]string // 自定义请求头
}

// AppOptions 获取应用详情的选项参数
// 用于配置app函数的行为
type AppOptions struct {
	ID             int64           // 应用ID（与AppID二选一）
	AppID          string          // Bundle ID（与ID二选一）
	Country        Country         // 国家/地区代码，默认为us
	Lang           string          // 语言代码
	Ratings        bool            // 是否包含评分直方图
	RequestOptions *RequestOptions // 自定义请求选项
}

// SearchOptions 搜索选项参数
// 用于配置search函数的行为
type SearchOptions struct {
	Term           string          // 搜索关键词
	Country        Country         // 国家/地区代码
	Lang           string          // 语言代码
	Num            int             // 返回结果数量，默认50
	Page           int             // 页码，默认1
	IdsOnly        bool            // 是否只返回ID
	RequestOptions *RequestOptions // 自定义请求选项
}

// ListOptions 列表选项参数
// 用于配置list函数的行为
type ListOptions struct {
	Collection     Collection      // 应用集合类型
	Category       Category        // 应用类别
	Country        Country         // 国家/地区代码
	Lang           string          // 语言代码
	Num            int             // 返回结果数量
	FullDetail     bool            // 是否获取完整详情，默认false
	RequestOptions *RequestOptions // 自定义请求选项
}

// DeveloperOptions 开发者选项参数
// 用于配置developer函数的行为
type DeveloperOptions struct {
	DevID          int64           // 开发者ID（必需）
	Country        Country         // 国家/地区代码
	Lang           string          // 语言代码
	RequestOptions *RequestOptions // 自定义请求选项
}

// ReviewsOptions 评论选项参数
// 用于配置reviews函数的行为
type ReviewsOptions struct {
	ID             int64           // 应用ID
	AppID          string          // Bundle ID
	Country        Country         // 国家/地区代码
	Sort           Sort            // 排序方式
	Page           int             // 页码
	RequestOptions *RequestOptions // 自定义请求选项
}

// RatingsOptions 评分选项参数
// 用于配置ratings函数的行为
type RatingsOptions struct {
	ID             int64           // 应用ID
	Country        Country         // 国家/地区代码
	RequestOptions *RequestOptions // 自定义请求选项
}

// SimilarOptions 相似应用选项参数
// 用于配置similar函数的行为
type SimilarOptions struct {
	ID             int64           // 应用ID
	AppID          string          // Bundle ID
	Country        Country         // 国家/地区代码
	Lang           string          // 语言代码
	RequestOptions *RequestOptions // 自定义请求选项
}

// SuggestOptions 搜索建议选项参数
// 用于配置suggest函数的行为
type SuggestOptions struct {
	Term           string          // 搜索关键词
	RequestOptions *RequestOptions // 自定义请求选项
}

// VersionHistoryOptions 版本历史选项参数
// 用于配置versionHistory函数的行为
type VersionHistoryOptions struct {
	ID             int64           // 应用ID
	Country        Country         // 国家/地区代码
	RequestOptions *RequestOptions // 自定义请求选项
}
