package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	appstore "github.com/yililith/app-store-scraper"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMenu()
		fmt.Print("请输入选项 (1-11): ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			testAppDetail()
		case "2":
			testSearch()
		case "3":
			testList()
		case "4":
			testReviews()
		case "5":
			testRatings()
		case "6":
			testSimilar()
		case "7":
			testDeveloper()
		case "8":
			testVersionHistory()
		case "9":
			testAll()
		case "10", "q", "quit", "exit":
			fmt.Println("\n感谢使用！再见 👋")
			return
		default:
			fmt.Println("\n❌ 无效选项，请重新选择")
		}

		fmt.Println("\n按 Enter 键返回菜单...")
		if scanner.Scan() {
			scanner.Text()
		}
	}
}

func printMenu() {
	fmt.Println("\n╔════════════════════════════════════════════════╗")
	fmt.Println("║       App Store Scraper Go 测试菜单               ║")
	fmt.Println("╠════════════════════════════════════════════════╣")
	fmt.Println("║  1. 测试 App 应用详情                         ║")
	fmt.Println("║  2. 测试 Search 搜索功能                     ║")
	fmt.Println("║  3. 测试 List 排行榜功能                      ║")
	fmt.Println("║  4. 测试 Reviews 评论功能                      ║")
	fmt.Println("║  5. 测试 Ratings 评分分布                     ║")
	fmt.Println("║  6. 测试 Similar 相似应用                      ║")
	fmt.Println("║  7. 测试 Developer 开发者应用                ║")
	fmt.Println("║  8. 测试 VersionHistory 版本历史             ║")
	fmt.Println("║  9. 测试全部功能                               ║")
	fmt.Println("║ 10. 退出                                          ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
}

func testAppDetail() {
	fmt.Println("\n=== 测试 App 应用详情 ===\n")

	scraper := appstore.NewScraper()

	testCases := []struct {
		name    string
		options appstore.AppOptions
	}{
		{"Candy Crush (美国)", appstore.AppOptions{ID: 553834731, Country: appstore.CountryUS}},
		{"MariBank (菲律宾)", appstore.AppOptions{ID: 1592249158, Country: appstore.CountryPH}},
		{"微信 (中国)", appstore.AppOptions{ID: 414478124, Country: appstore.CountryCN}},
	}

	for i, tc := range testCases {
		fmt.Printf("%d. %s\n", i+1, tc.name)
		app, err := scraper.App(tc.options)
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
			continue
		}
		fmt.Printf("   ✅ 名称: %s\n", app.Title)
		fmt.Printf("   ✅ 开发者: %s\n", app.Developer)
		fmt.Printf("   ✅ 版本: %s\n", app.Version)
	}
}

func testSearch() {
	fmt.Println("\n=== 测试 Search 搜索功能 ===\n")

	scraper := appstore.NewScraper()

	testCases := []struct {
		name    string
		options appstore.SearchOptions
	}{
		{"搜索 TikTok (台湾)", appstore.SearchOptions{Term: "TikTok", Country: appstore.CountryTW, Num: 5}},
		{"搜索 Finance (菲律宾)", appstore.SearchOptions{Term: "finance banking", Country: appstore.CountryPH, Num: 5}},
		{"搜索 ChatGPT (美国)", appstore.SearchOptions{Term: "ChatGPT", Country: appstore.CountryUS, Num: 3}},
	}

	for i, tc := range testCases {
		fmt.Printf("%d. %s\n", i+1, tc.name)
		results, err := scraper.Search(tc.options)
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
			continue
		}
		fmt.Printf("   ✅ 找到 %d 个结果\n", len(results))
		for j, app := range results {
			if j < 3 {
				fmt.Printf("   %d. %s - %s\n", j+1, app.Title, app.Developer)
			}
		}
	}
}

func testList() {
	fmt.Println("\n=== 测试 List 排行榜功能 ===\n")

	scraper := appstore.NewScraper()

	testCases := []struct {
		name    string
		options appstore.ListOptions
	}{
		{"菲律宾免费财务榜单 (基本模式)", appstore.ListOptions{Collection: appstore.TopFreeIOS, Category: appstore.Finance, Country: appstore.CountryPH, Num: 10, FullDetail: false}},
		{"菲律宾免费财务榜单 (完整详情)", appstore.ListOptions{Collection: appstore.TopFreeIOS, Category: appstore.Finance, Country: appstore.CountryPH, Num: 5, FullDetail: true}},
		{"台湾免费游戏榜单", appstore.ListOptions{Collection: appstore.TopFreeIOS, Category: appstore.Games, Country: appstore.CountryTW, Num: 5, FullDetail: false}},
	}

	for i, tc := range testCases {
		fmt.Printf("%d. %s\n", i+1, tc.name)
		apps, err := scraper.List(tc.options)
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
			continue
		}
		fmt.Printf("   ✅ 获取到 %d 个应用\n", len(apps))
		for j, app := range apps {
			if j < 3 {
				fmt.Printf("   %d. %s - 评分: %.1f\n", j+1, app.Title, app.Score)
			}
		}
	}
}

func testReviews() {
	fmt.Println("\n=== 测试 Reviews 评论功能 ===\n")

	scraper := appstore.NewScraper()

	testCases := []struct {
		name    string
		options appstore.ReviewsOptions
	}{
		{"Candy Crush 评论 (美国)", appstore.ReviewsOptions{ID: 553834731, Country: appstore.CountryUS, Sort: appstore.Recent}},
		{"MariBank 评论 (菲律宾)", appstore.ReviewsOptions{ID: 1592249158, Country: appstore.CountryPH, Sort: appstore.Recent}},
	}

	for i, tc := range testCases {
		fmt.Printf("%d. %s\n", i+1, tc.name)
		reviews, err := scraper.Reviews(tc.options)
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
			continue
		}
		if len(reviews) == 0 {
			fmt.Printf("   ⚠️ 暂无评论\n")
		} else {
			fmt.Printf("   ✅ 获取到 %d 条评论\n", len(reviews))
			for j, review := range reviews {
				if j < 3 {
					title := review.Title
					if len(title) > 30 {
						title = title[:30] + "..."
					}
					fmt.Printf("   %d. [%d星] %s\n", j+1, review.Score, title)
				}
			}
		}
	}
}

func testRatings() {
	fmt.Println("\n=== 测试 Ratings 评分分布 ===\n")

	scraper := appstore.NewScraper()

	testCases := []struct {
		name    string
		id      int64
		country appstore.Country
	}{
		{"Candy Crush", 553834731, appstore.CountryUS},
		{"MariBank", 1592249158, appstore.CountryPH},
	}

	for i, tc := range testCases {
		fmt.Printf("%d. %s\n", i+1, tc.name)
		ratings, err := scraper.Ratings(appstore.RatingsOptions{ID: tc.id, Country: tc.country})
		fmt.Printf("   ✅ 评分分布\n")
		fmt.Printf("   	1星: %d\n", ratings.Histogram.OneStar)
		fmt.Printf("   	2星: %d\n", ratings.Histogram.TwoStars)
		fmt.Printf("   	3星: %d\n", ratings.Histogram.ThreeStars)
		fmt.Printf("   	4星: %d\n", ratings.Histogram.FourStars)
		fmt.Printf("   	5星: %d\n", ratings.Histogram.FiveStars)
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
			continue
		}
		fmt.Printf("   ✅ 总评分数: %d\n", ratings.Ratings)
	}
}

func testSimilar() {
	fmt.Println("\n=== 测试 Similar 相似应用 ===\n")

	scraper := appstore.NewScraper()

	testCases := []struct {
		name    string
		options appstore.SimilarOptions
	}{
		{"Candy Crush 相似应用", appstore.SimilarOptions{ID: 553834731, Country: appstore.CountryUS}},
		{"ChatGPT 相似应用", appstore.SimilarOptions{ID: 6448311069, Country: appstore.CountryUS}},
	}

	for i, tc := range testCases {
		fmt.Printf("%d. %s\n", i+1, tc.name)
		apps, err := scraper.Similar(tc.options)
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
			continue
		}
		if len(apps) == 0 {
			fmt.Printf("   ⚠️ 暂无相似应用\n")
		} else {
			fmt.Printf("   ✅ 找到 %d 个相似应用\n", len(apps))
			for j, app := range apps {
				if j < 3 {
					fmt.Printf("   %d. %s - %s\n", j+1, app.Title, app.Developer)
				}
			}
		}
	}
}

func testDeveloper() {
	fmt.Println("\n=== 测试 Developer 开发者应用 ===\n")

	scraper := appstore.NewScraper()

	testCases := []struct {
		name    string
		options appstore.DeveloperOptions
	}{
		{"OpenAI 开发者的应用", appstore.DeveloperOptions{DevID: 284882218, Country: appstore.CountryUS}},
		{"Google LLC 开发者的应用", appstore.DeveloperOptions{DevID: 284882215, Country: appstore.CountryUS}},
	}

	for i, tc := range testCases {
		fmt.Printf("%d. %s\n", i+1, tc.name)
		apps, err := scraper.Developer(tc.options)
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
			continue
		}
		if len(apps) == 0 {
			fmt.Printf("   ⚠️ 未找到该开发者的应用\n")
		} else {
			fmt.Printf("   ✅ 找到 %d 个应用\n", len(apps))
			for j, app := range apps {
				if j < 3 {
					fmt.Printf("   %d. %s (版本: %s)\n", j+1, app.Title, app.Version)
				}
			}
		}
	}
}

func testVersionHistory() {
	fmt.Println("\n=== 测试 VersionHistory 版本历史 ===\n")

	scraper := appstore.NewScraper()

	testCases := []struct {
		name    string
		options appstore.VersionHistoryOptions
	}{
		{"Candy Crush", appstore.VersionHistoryOptions{ID: 553834731, Country: appstore.CountryUS}},
		{"MariBank", appstore.VersionHistoryOptions{ID: 1592249158, Country: appstore.CountryPH}},
	}

	for i, tc := range testCases {
		fmt.Printf("%d. %s\n", i+1, tc.name)
		history, err := scraper.VersionHistory(tc.options)
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
			continue
		}
		if len(history) > 0 {
			fmt.Printf("   ✅ 版本数: %d\n", len(history))
			fmt.Printf("   ✅ 最新版本: %s\n", history[0].VersionDisplay)
		}
	}
}

func testAll() {
	fmt.Println("\n=== 测试全部功能 ===\n")

	scraper := appstore.NewScraper()

	fmt.Println("1. App 应用详情:")
	app, err := scraper.App(appstore.AppOptions{ID: 553834731, Country: appstore.CountryUS})
	if err != nil {
		fmt.Printf("   ❌ 错误: %v\n", err)
	} else {
		fmt.Printf("   ✅ %s - %s\n", app.Title, app.Developer)
	}

	fmt.Println("\n2. Search 搜索:")
	results, err := scraper.Search(appstore.SearchOptions{Term: "ChatGPT", Country: appstore.CountryUS, Num: 3})
	if err != nil {
		fmt.Printf("   ❌ 错误: %v\n", err)
	} else {
		fmt.Printf("   ✅ 找到 %d 个结果\n", len(results))
	}

	fmt.Println("\n3. List 排行榜:")
	apps, err := scraper.List(appstore.ListOptions{Collection: appstore.TopFreeIOS, Category: appstore.Games, Country: appstore.CountryTW, Num: 3})
	if err != nil {
		fmt.Printf("   ❌ 错误: %v\n", err)
	} else {
		fmt.Printf("   ✅ 获取到 %d 个应用\n", len(apps))
	}

	fmt.Println("\n4. Reviews 评论:")
	reviews, err := scraper.Reviews(appstore.ReviewsOptions{ID: 553834731, Country: appstore.CountryUS, Sort: appstore.Recent})
	if err != nil {
		fmt.Printf("   ❌ 错误: %v\n", err)
	} else {
		fmt.Printf("   ✅ 获取到 %d 条评论\n", len(reviews))
	}

	fmt.Println("\n5. Ratings 评分:")
	ratings, err := scraper.Ratings(appstore.RatingsOptions{ID: 553834731, Country: appstore.CountryUS})
	if err != nil {
		fmt.Printf("   ❌ 错误: %v\n", err)
	} else {
		fmt.Printf("   ✅ 总评分数: %d\n", ratings.Ratings)
	}

	fmt.Println("\n6. Similar 相似应用:")
	similar, err := scraper.Similar(appstore.SimilarOptions{ID: 553834731, Country: appstore.CountryUS})
	if err != nil {
		fmt.Printf("   ❌ 错误: %v\n", err)
	} else {
		fmt.Printf("   ✅ 找到 %d 个相似应用\n", len(similar))
	}

	fmt.Println("\n7. Developer 开发者应用:")
	devApps, err := scraper.Developer(appstore.DeveloperOptions{DevID: 284882218, Country: appstore.CountryUS})
	if err != nil {
		fmt.Printf("   ❌ 错误: %v\n", err)
	} else {
		fmt.Printf("   ✅ 找到 %d 个应用\n", len(devApps))
	}

	fmt.Println("\n8. VersionHistory 版本历史:")
	history, err := scraper.VersionHistory(appstore.VersionHistoryOptions{ID: 553834731, Country: appstore.CountryUS})
	if err != nil {
		fmt.Printf("   ❌ 错误: %v\n", err)
	} else {
		fmt.Printf("   ✅ 获取到 %d 个版本\n", len(history))
	}

	fmt.Println("\n✅ 全部功能测试完成！")
}

func printJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("JSON序列化失败: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}
