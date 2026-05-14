package appstore

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Scraper App Store爬虫结构体
// 提供所有App Store数据抓取功能
type Scraper struct {
	client *Client // HTTP客户端
}

// NewScraper 创建并返回一个新的App Store爬虫实例
func NewScraper() *Scraper {
	return &Scraper{
		client: NewClient(),
	}
}

// NewScraperWithClient 使用指定的HTTP客户端创建爬虫
// client: HTTP客户端实例
func NewScraperWithClient(client *Client) *Scraper {
	return &Scraper{
		client: client,
	}
}

// App 根据应用ID或Bundle ID获取应用详细信息
// opts: 获取选项，包含ID或AppID
// 返回应用信息结构体和错误信息
func (s *Scraper) App(opts AppOptions) (*AppInfo, error) {
	if opts.ID == 0 && opts.AppID == "" {
		return nil, ErrInvalidParameter
	}

	if opts.Country == "" {
		opts.Country = CountryUS
	}

	var path string
	var idValue string
	if opts.ID != 0 {
		path = "/lookup"
		idValue = strconv.FormatInt(opts.ID, 10)
	} else {
		path = "/lookup"
		idValue = opts.AppID
	}

	params := map[string]string{
		"id":      idValue,
		"country": string(opts.Country),
		"entity":  "software",
	}
	if opts.Lang != "" {
		params["lang"] = opts.Lang
	}

	result, err := s.client.GetLookup(path, params)
	if err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, ErrNotFound
	}

	appResult := result.Results[0]
	if appResult.Kind != "software" && appResult.WrapperType != "software" {
		return nil, ErrNotFound
	}

	app := s.parseAppResult(appResult)

	hasScreenshots := len(app.ScreenshotURLs) > 0 || len(app.IpadScreenshots) > 0 || len(app.AppletvScreenshots) > 0
	if !hasScreenshots {
		screenshots := s.scrapeScreenshots(app.ID, string(opts.Country), opts.RequestOptions)
		app.ScreenshotURLs = screenshots.ScreenshotURLs
		app.IpadScreenshots = screenshots.IpadScreenshots
		app.AppletvScreenshots = screenshots.AppletvScreenshots
	}

	if opts.Ratings {
		ratingsData, err := s.Ratings(RatingsOptions{
			ID:             app.ID,
			Country:        opts.Country,
			RequestOptions: opts.RequestOptions,
		})
		if err == nil {
			app.Histogram = &ratingsData.Histogram
		}
	}

	return &app, nil
}

// Search 根据关键词搜索应用
// opts: 搜索选项，包含关键词、分页等参数
// 返回应用列表或ID列表和错误信息
func (s *Scraper) Search(opts SearchOptions) ([]AppInfo, error) {
	if opts.Term == "" {
		return nil, ErrInvalidParameter
	}

	if opts.Num == 0 {
		opts.Num = 50
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.Country == "" {
		opts.Country = CountryUS
	}

	params := BuildParams(opts)
	params["term"] = opts.Term
	params["media"] = "software"
	params["entity"] = "software"

	result, err := s.client.GetLookup("/search", params)
	if err != nil {
		return nil, err
	}

	var apps []AppInfo
	for _, item := range result.Results {
		if item.Kind == "software" {
			app := s.parseAppResult(item)
			apps = append(apps, app)
		}
	}

	start := (opts.Page - 1) * opts.Num
	end := start + opts.Num
	if start > len(apps) {
		return []AppInfo{}, nil
	}
	if end > len(apps) {
		end = len(apps)
	}

	return apps[start:end], nil
}

// List 获取应用列表（排行、新品等）
// opts: 列表选项，包含集合类型、分类等参数
// 返回应用列表和错误信息
func (s *Scraper) List(opts ListOptions) ([]AppInfo, error) {
	if opts.Num == 0 {
		opts.Num = 50
	}
	if opts.Country == "" {
		opts.Country = CountryUS
	}
	if opts.Collection == "" {
		opts.Collection = TopFreeIOS
	}

	limit := opts.Num
	if limit > 200 {
		limit = 200
	}

	feed, err := s.client.GetRSSFeedTyped(string(opts.Country), opts.Collection, opts.Category, limit)
	if err != nil {
		return nil, err
	}

	if len(feed.Feed.Entry) == 0 {
		return []AppInfo{}, nil
	}

	appIDs := make([]int64, 0, len(feed.Feed.Entry))
	for _, entry := range feed.Feed.Entry {
		if id, err := strconv.ParseInt(entry.ID.Attributes.ImID, 10, 64); err == nil {
			appIDs = append(appIDs, id)
		}
	}

	if len(appIDs) == 0 {
		return []AppInfo{}, nil
	}

	return s.lookup(appIDs, string(opts.Country), opts.Lang, opts.RequestOptions)
}

// Developer 获取指定开发者的所有应用
// opts: 开发者选项，包含开发者ID
// 返回应用列表和错误信息
func (s *Scraper) Developer(opts DeveloperOptions) ([]AppInfo, error) {
	if opts.DevID == 0 {
		return nil, ErrInvalidParameter
	}

	if opts.Country == "" {
		opts.Country = CountryUS
	}

	ids := []int64{opts.DevID}
	return s.lookup(ids, string(opts.Country), opts.Lang, opts.RequestOptions)
}

// Reviews 获取应用的评论列表
// opts: 评论选项，包含应用ID、分页、排序等参数
// 返回评论列表和错误信息
func (s *Scraper) Reviews(opts ReviewsOptions) ([]Review, error) {
	if opts.ID == 0 && opts.AppID == "" {
		return nil, ErrInvalidParameter
	}

	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.Page < 1 || opts.Page > 10 {
		return nil, fmt.Errorf("页码必须在1到10之间")
	}
	if opts.Country == "" {
		opts.Country = CountryUS
	}
	if opts.Sort == "" {
		opts.Sort = Recent
	}

	id := opts.ID
	if opts.AppID != "" && opts.ID == 0 {
		app, err := s.App(AppOptions{
			AppID:   opts.AppID,
			Country: opts.Country,
		})
		if err != nil {
			return nil, err
		}
		id = app.ID
	}

	feed, err := s.client.GetReviewFeed(string(opts.Country), opts.Page, int(id), string(opts.Sort))
	if err != nil {
		return nil, err
	}

	var reviews []Review
	for i := 1; i < len(feed.Feed.Entry); i++ {
		review := s.parseReviewEntry(feed.Feed.Entry[i])
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// Ratings 获取应用的评分分布直方图
// opts: 评分选项，包含应用ID和国家参数
// 返回评分分布和错误信息
func (s *Scraper) Ratings(opts RatingsOptions) (*Ratings, error) {
	// 验证必填参数
	if opts.ID == 0 {
		return nil, ErrInvalidParameter
	}

	// 设置默认值
	if opts.Country == "" {
		opts.Country = CountryUS
	}

	// 获取storeId
	storeFront := storeId(string(opts.Country))
	url := fmt.Sprintf("https://itunes.apple.com/%s/customer-reviews/id%d?displayable-kind=11", opts.Country, opts.ID)

	// 构建请求头
	headers := map[string]string{
		"X-Apple-Store-Front": fmt.Sprintf("%d,12", storeFront),
	}
	if opts.RequestOptions != nil && opts.RequestOptions.Headers != nil {
		for k, v := range opts.RequestOptions.Headers {
			headers[k] = v
		}
	}

	// 发送请求
	html, err := s.client.GetHTMLWithHeaders(url, headers)
	if err != nil {
		return nil, err
	}

	if html == "" {
		return nil, ErrNotFound
	}

	// 解析HTML
	return s.parseRatings(html)
}

// Similar 获取相似应用推荐
// opts: 相似应用选项，包含应用ID等参数
// 返回相似应用列表和错误信息
func (s *Scraper) Similar(opts SimilarOptions) ([]AppInfo, error) {
	if opts.ID == 0 && opts.AppID == "" {
		return nil, ErrInvalidParameter
	}

	if opts.Country == "" {
		opts.Country = CountryUS
	}

	id := opts.ID
	if opts.AppID != "" && opts.ID == 0 {
		app, err := s.App(AppOptions{
			AppID:   opts.AppID,
			Country: opts.Country,
		})
		if err != nil {
			return nil, err
		}
		id = app.ID
	}

	url := fmt.Sprintf("https://apps.apple.com/%s/app/id%d", opts.Country, id)

	html, err := s.client.GetHTML(url)
	if err != nil {
		return []AppInfo{}, nil
	}

	similarIds := s.extractSimilarAppIds(html, id)

	if len(similarIds) == 0 {
		return []AppInfo{}, nil
	}

	return s.lookup(similarIds, string(opts.Country), opts.Lang, opts.RequestOptions)
}

// Suggest 获取搜索建议
// opts: 搜索建议选项，包含关键词等参数
// 返回搜索建议列表和错误信息
func (s *Scraper) Suggest(opts SuggestOptions) ([]Suggestion, error) {
	// 验证必填参数
	if opts.Term == "" {
		return nil, ErrInvalidParameter
	}

	// 构建URL
	url := fmt.Sprintf("https://search.itunes.apple.com/WebObjects/MZSearchHints.woa/wa/hints?clientApplication=Software&term=%s", url.QueryEscape(opts.Term))

	// 发送请求
	xmlData, err := s.client.GetRaw(url)
	if err != nil {
		return nil, err
	}

	// 解析XML
	return s.parseSuggestXML(xmlData)
}

// VersionHistory 获取应用版本历史
// opts: 版本历史选项，包含应用ID等参数
// 返回版本历史列表和错误信息
func (s *Scraper) VersionHistory(opts VersionHistoryOptions) ([]VersionHistory, error) {
	// 验证必填参数
	if opts.ID == 0 {
		return nil, ErrInvalidParameter
	}

	// 设置默认值
	if opts.Country == "" {
		opts.Country = CountryUS
	}

	// 构建URL
	url := fmt.Sprintf("https://apps.apple.com/%s/app/id%d", opts.Country, opts.ID)

	// 发送请求
	html, err := s.client.GetHTML(url)
	if err != nil {
		return nil, err
	}

	// 解析HTML内容，提取版本历史
	return s.parseVersionHistory(html)
}

// parseAppResult 解析应用信息
// result: API响应的AppResult结构
// 返回解析后的AppInfo结构体
func (s *Scraper) parseAppResult(result AppResult) AppInfo {
	icon := result.ArtworkURL512
	if icon == "" {
		icon = result.ArtworkURL100
	}

	app := AppInfo{
		ID:                    result.TrackID,
		AppID:                 result.BundleID,
		Title:                 result.TrackName,
		URL:                   result.TrackViewURL,
		Description:           result.Description,
		Icon:                  icon,
		Genres:                result.Genres,
		GenreIDs:              result.GenreIDs,
		PrimaryGenre:          result.PrimaryGenreName,
		PrimaryGenreID:        result.PrimaryGenreID,
		ContentRating:         result.ContentAdvisoryRating,
		Languages:             result.LanguageCodesISO2A,
		Size:                  result.FileSizeBytes,
		RequiredOsVersion:     result.MinimumOsVersion,
		Released:              result.ReleaseDate,
		Updated:               result.CurrentVersionReleaseDate,
		ReleaseNotes:          result.ReleaseNotes,
		Version:               result.Version,
		Price:                 result.Price,
		Currency:              result.Currency,
		Free:                  result.Price == 0,
		DeveloperID:           result.ArtistID,
		Developer:             result.ArtistName,
		DeveloperURL:          result.ArtistViewURL,
		DeveloperWebsite:      result.SellerURL,
		Score:                 result.AverageUserRating,
		Reviews:               result.UserRatingCount,
		CurrentVersionScore:   result.AverageUserRatingForCurrentVersion,
		CurrentVersionReviews: result.UserRatingCountForCurrentVersion,
		ScreenshotURLs:        result.ScreenshotURLs,
		IpadScreenshots:       result.IpadScreenshotURLs,
		AppletvScreenshots:    result.AppletvScreenshotURLs,
		SupportedDevices:      result.SupportedDevices,
	}

	return app
}

// parseReviewEntry 解析评论信息
// entry: RSS Feed中的ReviewEntry结构
// 返回解析后的Review结构体
func (s *Scraper) parseReviewEntry(entry ReviewEntry) Review {
	review := Review{
		ID:            entry.ID.Label,
		UserName:      entry.Author.Name.Label,
		UserReviewURL: entry.Author.URI.Label,
		Version:       entry.ImVersion.Label,
		Score:         parseRatingStr(entry.ImRating.Label),
		Title:         entry.Title.Label,
		Text:          entry.Content.Label,
		Updated:       entry.Updated.Label,
	}

	return review
}

// parseRatingStr 解析评分为整数
func parseRatingStr(value string) int {
	if value == "" {
		return 0
	}
	if rating, err := strconv.Atoi(value); err == nil {
		return rating
	}
	return 0
}

// parseRSSEntry 解析RSS Feed条目中的基本信息
// 用于快速获取列表，避免每个应用都调用一次API
// entry: RSS feed中的entry数据
// 返回解析后的AppInfo结构体（基本信息）
// 已废弃，使用 lookup 函数代替
func (s *Scraper) parseRSSEntry(entry json.RawMessage) AppInfo {
	return AppInfo{}
}

// ScreenshotResult 截图爬取结果
type ScreenshotResult struct {
	ScreenshotURLs     []string
	IpadScreenshots    []string
	AppletvScreenshots []string
}

// scrapeScreenshots 从App Store页面爬取截图
func (s *Scraper) scrapeScreenshots(appID int64, country string, opts *RequestOptions) ScreenshotResult {
	result := ScreenshotResult{
		ScreenshotURLs:     []string{},
		IpadScreenshots:    []string{},
		AppletvScreenshots: []string{},
	}

	url := fmt.Sprintf("https://apps.apple.com/%s/app/id%d", country, appID)
	html, err := s.client.GetHTML(url)
	if err != nil {
		return result
	}

	// 提取截图URL
	result.ScreenshotURLs = extractScreenshots(html, "ScreenshotPhone")
	result.IpadScreenshots = extractScreenshots(html, "ScreenshotPad")
	result.AppletvScreenshots = extractScreenshots(html, "ScreenshotAppleTv")

	return result
}

// extractScreenshots 从HTML中提取截图URL
func extractScreenshots(html string, screenshotType string) []string {
	var screenshots []string

	// 查找对应的截图容器
	pattern := regexp.MustCompile(fmt.Sprintf(`<ul[^>]*class=["'][^"']*grid-type-%s[^"']*["'][^>]*>[\s\S]*?</ul>`, screenshotType))
	matches := pattern.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) < 1 {
			continue
		}
		content := match[0]

		// 提取source标签的srcset属性
		srcsetPattern := regexp.MustCompile(`<source[^>]*type=["']image/webp["'][^>]*srcset=["']([^"']+)["'][^>]*>`)
		srcsetMatches := srcsetPattern.FindAllStringSubmatch(content, -1)

		for _, srcsetMatch := range srcsetMatches {
			if len(srcsetMatch) < 2 {
				continue
			}
			srcset := srcsetMatch[1]
			// 解析srcset
			entries := strings.Split(srcset, ",")
			for _, entry := range entries {
				parts := strings.TrimSpace(entry)
				urlParts := strings.Split(parts, " ")
				if len(urlParts) > 0 {
					url := urlParts[0]
					// 标准化URL
					url = regexp.MustCompile(`(/\d+x\d+bb(-\d+)?\.(webp|jpg|jpeg|png))$`).ReplaceAllString(url, "/392x696bb.png")
					if !contains(screenshots, url) {
						screenshots = append(screenshots, url)
					}
				}
			}
		}
	}

	return screenshots
}

// contains 检查切片是否包含某个元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// parseRatings 解析评分HTML
func (s *Scraper) parseRatings(html string) (*Ratings, error) {
	ratings := &Ratings{
		Histogram: RatingHistogram{},
	}

	// 提取总评分数量
	ratingCountPattern := regexp.MustCompile(`<span[^>]*class=["'][^"']*rating-count[^"']*["'][^>]*>[\s\S]*?(\d+)[\s\S]*?</span>`)
	ratingCountMatch := ratingCountPattern.FindStringSubmatch(html)
	if len(ratingCountMatch) > 1 {
		if count, err := strconv.Atoi(ratingCountMatch[1]); err == nil {
			ratings.Ratings = count
		}
	}

	// 提取各星级评分数量
	// 格式通常是: <span class="total">123</span>
	totalPattern := regexp.MustCompile(`<span[^>]*class=["']total["'][^>]*>(\d+)</span>`)
	totalMatches := totalPattern.FindAllStringSubmatch(html, -1)

	// 从5星到1星
	for i, match := range totalMatches {
		if len(match) < 2 {
			continue
		}
		if count, err := strconv.Atoi(match[1]); err == nil {
			star := 5 - i
			switch star {
			case 5:
				ratings.Histogram.FiveStars = count
			case 4:
				ratings.Histogram.FourStars = count
			case 3:
				ratings.Histogram.ThreeStars = count
			case 2:
				ratings.Histogram.TwoStars = count
			case 1:
				ratings.Histogram.OneStar = count
			}
		}
	}

	return ratings, nil
}

// extractSimilarAppIds 从HTML中提取相似应用的ID
func (s *Scraper) extractSimilarAppIds(html string, currentAppID int64) []int64 {
	var ids []int64

	// 查找所有 /app/ 链接
	linkPattern := regexp.MustCompile(`href=["']([^"']*\/app\/[^"']*)["']`)
	matches := linkPattern.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		href := match[1]

		// 提取ID
		idPattern := regexp.MustCompile(`\/id(\d+)`)
		idMatch := idPattern.FindStringSubmatch(href)
		if len(idMatch) < 2 {
			continue
		}
		if id, err := strconv.ParseInt(idMatch[1], 10, 64); err == nil {
			// 避免重复和当前应用
			if id != currentAppID && !containsInt64(ids, id) {
				ids = append(ids, id)
			}
		}
	}

	return ids
}

// containsInt64 检查切片是否包含某个int64
func containsInt64(slice []int64, item int64) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// lookup 通过lookup API获取应用信息
func (s *Scraper) lookup(ids []int64, country, lang string, opts *RequestOptions) ([]AppInfo, error) {
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.FormatInt(id, 10)
	}
	idsString := strings.Join(idStrs, ",")

	url := fmt.Sprintf("/lookup?id=%s&country=%s&entity=software", idsString, country)
	if lang != "" {
		url += "&lang=" + lang
	}

	result, err := s.client.GetLookup(url, nil)
	if err != nil {
		return nil, err
	}

	var apps []AppInfo
	for _, item := range result.Results {
		if item.Kind == "software" || item.WrapperType == "software" {
			app := s.parseAppResult(item)
			apps = append(apps, app)
		}
	}

	return apps, nil
}

// parseSuggestXML 解析建议XML响应
func (s *Scraper) parseSuggestXML(xmlData string) ([]Suggestion, error) {
	suggestions := []Suggestion{}

	// 简单的XML解析
	// 查找 <string> 标签
	stringPattern := regexp.MustCompile(`<string[^>]*>([^<]+)</string>`)
	matches := stringPattern.FindAllStringSubmatch(xmlData, -1)

	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		term := strings.TrimSpace(match[1])
		if term != "" && !seen[term] {
			suggestions = append(suggestions, Suggestion{
				Term: term,
			})
			seen[term] = true
		}
	}

	return suggestions, nil
}

// storeId 获取国家代码对应的Apple Store ID
func storeId(country string) int {
	storeIds := map[string]int{
		"dz": 143563, "ao": 143564, "ai": 143538, "ag": 143540, "ar": 143505,
		"am": 143524, "au": 143460, "at": 143445, "az": 143568, "bs": 143539,
		"bh": 143559, "bb": 143541, "by": 143565, "be": 143446, "bz": 143555,
		"bj": 143576, "bm": 143542, "bo": 143556, "bw": 143525, "br": 143503,
		"vg": 143543, "bn": 143560, "bg": 143526, "bf": 143578, "ca": 143455,
		"ky": 143544, "td": 143581, "cl": 143483, "cn": 143465, "co": 143501,
		"cr": 143495, "ci": 143527, "hr": 143494, "cy": 143557, "cz": 143489,
		"dk": 143458, "dm": 143545, "do": 143508, "ec": 143509, "eg": 143516,
		"sv": 143506, "ee": 143518, "fj": 143583, "fi": 143447, "fr": 143442,
		"gm": 143584, "de": 143443, "gh": 143573, "gr": 143448, "gd": 143546,
		"gt": 143504, "gw": 143585, "gy": 143553, "hn": 143510, "hk": 143463,
		"hu": 143482, "is": 143558, "in": 143467, "id": 143476, "ie": 143449,
		"il": 143491, "it": 143450, "jm": 143511, "jp": 143462, "jo": 143528,
		"kz": 143517, "ke": 143529, "kr": 143466, "kw": 143493, "lv": 143519,
		"lb": 143497, "lt": 143520, "lu": 143451, "mo": 143515, "mk": 143530,
		"mg": 143531, "mw": 143589, "my": 143473, "ml": 143532, "mt": 143521,
		"mr": 143590, "mu": 143533, "mx": 143468, "fm": 143591, "md": 143523,
		"mn": 143592, "ms": 143547, "mz": 143593, "na": 143594, "np": 143484,
		"nl": 143452, "nz": 143461, "ni": 143512, "ne": 143534, "ng": 143561,
		"no": 143457, "om": 143562, "pk": 143477, "pw": 143595, "pa": 143485,
		"pg": 143597, "py": 143513, "pe": 143507, "ph": 143474, "pl": 143478,
		"pt": 143453, "qa": 143498, "ro": 143487, "ru": 143469, "kn": 143548,
		"lc": 143549, "vc": 143550, "st": 143598, "sa": 143479, "sn": 143535,
		"sc": 143599, "sl": 143600, "sg": 143464, "sk": 143496, "si": 143499,
		"sb": 143601, "za": 143472, "es": 143454, "lk": 143486, "sr": 143554,
		"sz": 143602, "se": 143456, "ch": 143459, "tw": 143470, "tj": 143603,
		"tz": 143572, "th": 143475, "tn": 143536, "tr": 143480, "tm": 143604,
		"tc": 143551, "ug": 143537, "ua": 143492, "ae": 143481, "gb": 143444,
		"us": 143441, "uy": 143514, "uz": 143566, "ve": 143502, "vn": 143471,
		"ye": 143571, "zw": 143605,
	}

	if id, ok := storeIds[strings.ToLower(country)]; ok {
		return id
	}
	return 143441 // 默认US
}

// parseVersionHistory 解析版本历史HTML
func (s *Scraper) parseVersionHistory(html string) ([]VersionHistory, error) {
	var versions []VersionHistory

	// 查找版本历史条目
	// 格式: <article class="svelte-13339ih">
	articlePattern := regexp.MustCompile(`<article[^>]*class=["'][^"']*svelte-13339ih[^"']*["'][^>]*>[\s\S]*?</article>`)
	articles := articlePattern.FindAllStringSubmatch(html, -1)

	for _, article := range articles {
		if len(article) < 1 {
			continue
		}
		content := article[0]

		// 提取版本号
		versionPattern := regexp.MustCompile(`<h4[^>]*class=["'][^"']*svelte-13339ih[^"']*["'][^>]*>([^<]+)</h4>`)
		versionMatch := versionPattern.FindStringSubmatch(content)
		versionDisplay := ""
		if len(versionMatch) > 1 {
			versionDisplay = strings.TrimSpace(versionMatch[1])
		}

		// 提取发布日期
		datePattern := regexp.MustCompile(`<time[^>]*datetime=["']([^"']+)["'][^>]*>`)
		dateMatch := datePattern.FindStringSubmatch(content)
		releaseDate := ""
		if len(dateMatch) > 1 {
			releaseDate = strings.TrimSpace(dateMatch[1])
		}

		// 提取发布说明
		notesPattern := regexp.MustCompile(`<p[^>]*class=["'][^"']*svelte-13339ih[^"']*["'][^>]*>([\s\S]*?)</p>`)
		notesMatch := notesPattern.FindStringSubmatch(content)
		releaseNotes := ""
		if len(notesMatch) > 1 {
			// 去除HTML标签
			releaseNotes = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(notesMatch[1], "")
			releaseNotes = strings.TrimSpace(releaseNotes)
		}

		if versionDisplay != "" || releaseDate != "" {
			versions = append(versions, VersionHistory{
				VersionDisplay: versionDisplay,
				ReleaseDate:    releaseDate,
				ReleaseNotes:   releaseNotes,
			})
		}
	}

	return versions, nil
}
