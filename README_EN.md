# App Store Scraper for Go

## 🌐 语言切换 | Language Switch

📖 **You are reading English version** | [中文版本](./README.md)

---

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Features](#core-features)
- [API Reference](#api-reference)
- [Constants](#constants)
- [Code Examples](#code-examples)
- [Error Handling](#error-handling)
- [Performance Optimization](#performance-optimization)
- [FAQ](#faq)
- [License](#license)

## Introduction

App Store Scraper is a modern Go library designed specifically for scraping application data from the iTunes/Mac App Store. This library provides complete type safety guarantees, requires no external dependencies, and supports data retrieval from over 100 countries and regions worldwide.

The library was created to address the lack of reliable App Store data scraping tools in the Go ecosystem. By wrapping Apple's official iTunes Search API and RSS Feed API, we can retrieve comprehensive data including app details, user reviews, rating statistics, ranking information, developer app lists, and more.

### Learn from and refer to
+ JavaScript: [facundoolano/app-store-scraper](https://github.com/facundoolano/app-store-scraper)
+ TypeScript: [plahteenlahti/app-store-scraper](https://github.com/plahteenlahti/app-store-scraper)

## Features

### Complete Type Safety

All API methods are carefully designed to return explicit struct types. Data structure correctness is ensured at compile time, avoiding the risks of runtime type assertions. When using this library, your IDE provides complete code completion and type checking functionality.

### No External Dependencies

This library is implemented using only Go standard libraries, with no third-party dependencies. This means you don't need to worry about dependency conflicts, version compatibility issues, or additional security risks. Import and use immediately.

### Lightweight and Easy to Use

The API is designed to be clean and intuitive. Through method chaining and sensible default parameters, complex data scraping tasks can be completed with just a few lines of code. Each method has a clear responsibility boundary, making it easy to understand and use.

### Multi-Region Support

Supports App Stores from 175+ countries and regions worldwide, including the United States, China (including Hong Kong, Macau, and Taiwan), Japan, South Korea, United Kingdom, Germany, France, Russia, Southeast Asian countries, and more. Each country has a corresponding constant definition for easy calling.

### Rich API Methods

Provides 10+ core methods covering complete data retrieval capabilities such as app fetching, search, rankings, reviews, ratings, similar apps, developer apps, version history, and search suggestions, meeting various business scenario requirements.

### Performance Optimization

When retrieving multiple app details, the library internally uses goroutines for concurrent requests, significantly improving data retrieval efficiency. The FullDetail option is also provided to allow flexible switching between fast mode and complete mode.

### Rating Histogram Support

In addition to average ratings, complete rating distribution data can be retrieved, including the number of reviews from 1-star to 5-star ratings, helping you get a more comprehensive understanding of the app's user feedback.

## Installation
```bash
go get github.com/yililith/app-store-scraper
```

### Prerequisites

- Go 1.25.6 or higher

### Module Import

After installation, import in your Go project:

```go
import appstore "github.com/yililith/app-store-scraper"
```

## Quick Start

### Basic Usage Flow

Using this library to scrape App Store data is straightforward:

1. Import the app-store-scraper package
2. Create a Scraper instance
3. Call the corresponding API method to get data
4. Handle the returned result or error

### Running the Example Program

The project includes an interactive test menu program for you to quickly experience all features:

```bash
# Navigate to the example directory
cd example

# Run the interactive test menu
go run main.go
```

### Interactive Test Menu

After running the program, an interactive menu will be displayed. You can enter a number to select the feature to test:

```
╔════════════════════════════════════════════════╗
║       App Store Scraper Go Test Menu            ║
╠════════════════════════════════════════════════╣
║  1. Test App Details                           ║
║  2. Test Search                                ║
║  3. Test List/Rankings                         ║
║  4. Test Reviews                               ║
║  5. Test Ratings                               ║
║  6. Test Similar Apps                          ║
║  7. Test Developer Apps                        ║
║  8. Test Version History                       ║
║  9. Test All Features                          ║
║ 10. Exit                                        ║
╚════════════════════════════════════════════════╝
```

## Core Features

### 1. Get App Details

Retrieve complete app information including name, developer, price, rating, description, and screenshots by app ID or Bundle ID. Suitable for displaying app detail pages or conducting app data statistical analysis.

### 2. Search Apps

Search for applications in the App Store based on keywords. Supports paginated queries and can retrieve any number of apps from search results. Suitable for implementing app search functionality or keyword trend analysis.

### 3. Get Rankings

Retrieve app ranking data for various categories, supporting different dimensions such as free, paid, and top-grossing charts. Supports filtering by app category to get ranking information for specific categories.

### 4. Get Reviews

Retrieve user review lists for apps, supporting sorting by most recent or most helpful. Includes complete data such as review titles, content, ratings, and user information. Suitable for sentiment analysis or user feedback collection.

### 5. Rating Distribution

Retrieve the app's rating distribution histogram data, including the distribution of ratings across different star levels. Helps you understand the overall user satisfaction distribution of an app.

### 6. Similar Apps

Retrieve other app recommendations similar to the target app. Based on App Store's similar apps algorithm, returns a list of apps users might be interested in.

### 7. Developer Apps

Retrieve all apps published by a specified developer by developer ID. Suitable for developer information display or competitive analysis.

### 8. Version History

Retrieve complete version update history for an app, including version number, release date, and update notes for each version. Suitable for version tracking or app update analysis.

### 9. Search Suggestions

Retrieve suggestions for search keywords to help users enter more accurate search terms or discover related apps.

## API Reference

### Scraper Struct

Scraper is the core struct of this library, providing all App Store data scraping functionality.

```go
// Create a default Scraper instance
scraper := appstore.NewScraper()

// Create a Scraper with a custom HTTP client
scraper := appstore.NewScraperWithClient(client)
```

### App - Get App Details

Retrieve detailed app information by app ID or Bundle ID.

```go
func (s *Scraper) App(opts AppOptions) (*AppInfo, error)
```

**Parameters**:

- `ID`: App ID (number), mutually exclusive with AppID
- `AppID`: Bundle ID (string), mutually exclusive with ID
- `Country`: Country/region code, defaults to US
- `Lang`: Language code for specifying the language of returned data
- `Ratings`: Whether to include rating histogram data

**Returns**:

- `*AppInfo`: Pointer to the app detail struct
- `error`: Error information

**Usage Examples**:

```go
// Get by app ID
app, err := scraper.App(appstore.AppOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})

// Get by Bundle ID and include ratings
app, err := scraper.App(appstore.AppOptions{
    AppID:   "com.king.candycrushsaga",
    Country: appstore.CountryUS,
    Ratings: true,
})

// Specify language and region
app, err := scraper.App(appstore.AppOptions{
    ID:      414478124,
    Country: appstore.CountryCN,
    Lang:    "zh_cn",
})
```

### Search - Search Apps

Search for applications by keywords.

```go
func (s *Scraper) Search(opts SearchOptions) ([]AppInfo, error)
```

**Parameters**:

- `Term`: Search keyword (required)
- `Country`: Country/region code
- `Lang`: Language code
- `Num`: Number of results to return, default 50, max 200
- `Page`: Page number, default 1

**Usage Examples**:

```go
// Basic search
results, err := scraper.Search(appstore.SearchOptions{
    Term:    "minecraft",
    Country: appstore.CountryUS,
    Num:     10,
})

// Paginated search
results, err := scraper.Search(appstore.SearchOptions{
    Term:    "finance banking",
    Country: appstore.CountryTW,
    Num:     10,
    Page:    2,
})
```

### List - Get Rankings

Retrieve app ranking data, supporting multiple collection types and category filters.

```go
func (s *Scraper) List(opts ListOptions) ([]AppInfo, error)
```

**Parameters**:

- `Collection`: App collection type, such as free chart, paid chart, etc.
- `Category`: App category, such as games, finance, etc.
- `Country`: Country/region code
- `Num`: Number of results to return
- `FullDetail`: Whether to get complete details; true will make additional API calls

**Usage Examples**:

```go
// Get free finance apps ranking (fast mode)
apps, err := scraper.List(appstore.ListOptions{
    Collection: appstore.TopFreeIOS,
    Category:   appstore.Finance,
    Country:    appstore.CountryPH,
    Num:        10,
    FullDetail: false,
})

// Get free games ranking (full details)
apps, err := scraper.List(appstore.ListOptions{
    Collection: appstore.TopFreeIOS,
    Category:   appstore.Games,
    Country:    appstore.CountryTW,
    Num:        5,
    FullDetail: true,
})
```

### Reviews - Get Reviews

Retrieve user review lists for an app.

```go
func (s *Scraper) Reviews(opts ReviewsOptions) ([]Review, error)
```

**Parameters**:

- `ID`: App ID
- `AppID`: Bundle ID, mutually exclusive with ID
- `Country`: Country/region code
- `Sort`: Sort order, Recent or Helpful
- `Page`: Page number, 1-10

**Usage Examples**:

```go
reviews, err := scraper.Reviews(appstore.ReviewsOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
    Sort:    appstore.Recent,
    Page:    1,
})
```

### Ratings - Rating Distribution

Retrieve the app's rating distribution histogram.

```go
func (s *Scraper) Ratings(opts RatingsOptions) (*Ratings, error)
```

**Parameters**:

- `ID`: App ID (required)
- `Country`: Country/region code

**Usage Examples**:

```go
ratings, err := scraper.Ratings(appstore.RatingsOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})

// Access rating data
fmt.Printf("Total ratings: %d\n", ratings.Ratings)
fmt.Printf("5 stars: %d\n", ratings.Histogram.FiveStars)
fmt.Printf("4 stars: %d\n", ratings.Histogram.FourStars)
fmt.Printf("3 stars: %d\n", ratings.Histogram.ThreeStars)
fmt.Printf("2 stars: %d\n", ratings.Histogram.TwoStars)
fmt.Printf("1 star: %d\n", ratings.Histogram.OneStar)
```

### Similar - Similar Apps

Get other apps similar to the specified app.

```go
func (s *Scraper) Similar(opts SimilarOptions) ([]AppInfo, error)
```

**Parameters**:

- `ID`: App ID
- `AppID`: Bundle ID, mutually exclusive with ID
- `Country`: Country/region code

**Usage Examples**:

```go
apps, err := scraper.Similar(appstore.SimilarOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})
```

### Developer - Developer Apps

Retrieve all apps by a specified developer.

```go
func (s *Scraper) Developer(opts DeveloperOptions) ([]AppInfo, error)
```

**Parameters**:

- `DevID`: Developer ID (required)
- `Country`: Country/region code

**Usage Examples**:

```go
// Get OpenAI's apps
apps, err := scraper.Developer(appstore.DeveloperOptions{
    DevID:   284882218,
    Country: appstore.CountryUS,
})

// Get Google LLC's apps
apps, err := scraper.Developer(appstore.DeveloperOptions{
    DevID:   284882215,
    Country: appstore.CountryUS,
})
```

### VersionHistory - Version History

Retrieve the app's version update history.

```go
func (s *Scraper) VersionHistory(opts VersionHistoryOptions) ([]VersionHistory, error)
```

**Parameters**:

- `ID`: App ID (required)
- `Country`: Country/region code

**Usage Examples**:

```go
history, err := scraper.VersionHistory(appstore.VersionHistoryOptions{
    ID:      553834731,
    Country: appstore.CountryUS,
})

for i, version := range history {
    fmt.Printf("Version %s (%s)\n", version.VersionDisplay, version.ReleaseDate)
    if version.ReleaseNotes != "" {
        fmt.Printf("Release notes: %s\n", version.ReleaseNotes)
    }
}
```

### Suggest - Search Suggestions

Get suggestions for search keywords.

```go
func (s *Scraper) Suggest(opts SuggestOptions) ([]Suggestion, error)
```

**Parameters**:

- `Term`: Search keyword (required)

**Usage Examples**:

```go
suggestions, err := scraper.Suggest(appstore.SuggestOptions{
    Term: "chat",
})
```

## Constants

### Collection - App Collections

Defines various app collection types in the App Store.

```go
// iOS App Collections
appstore.TopFreeIOS      // iOS Top Free Apps
appstore.TopGrossingIOS  // iOS Top Grossing Apps
appstore.TopPaidIOS      // iOS Top Paid Apps
appstore.NewFreeIOS     // iOS New Free Apps
appstore.NewPaidIOS     // iOS New Paid Apps
appstore.NewIOS         // iOS New Apps

// iPad App Collections
appstore.TopFreeiPad     // iPad Top Free Apps
appstore.TopGrossingiPad // iPad Top Grossing Apps
appstore.TopPaidiPad     // iPad Top Paid Apps

// Mac App Collections
appstore.TopMac          // Mac Apps Overall
appstore.TopFreeMac      // Mac Top Free Apps
appstore.TopGrossingMac  // Mac Top Grossing Apps
appstore.TopPaidMac      // Mac Top Paid Apps
```

### Category - App Categories

Defines app categories in the App Store using official Genre IDs.

```go
// Main Categories
appstore.Games               // Games
appstore.Finance             // Finance
appstore.Business            // Business
appstore.Education           // Education
appstore.Entertainment      // Entertainment
appstore.HealthAndFitness    // Health & Fitness
appstore.Lifestyle           // Lifestyle
appstore.Music              // Music
appstore.News               // News
appstore.PhotoAndVideo       // Photo & Video
appstore.Productivity        // Productivity
appstore.Shopping           // Shopping
appstore.SocialNetworking    // Social Networking
appstore.Sports             // Sports
appstore.Travel             // Travel
appstore.Utilities          // Utilities
appstore.Weather           // Weather

// Game Subcategories
appstore.GamesAction        // Games-Action
appstore.GamesAdventure     // Games-Adventure
appstore.GamesArcade        // Games-Arcade
appstore.GamesBoard         // Games-Board
appstore.GamesCard          // Games-Card
appstore.GamesCasino        // Games-Casino
appstore.GamesEducational   // Games-Educational
appstore.GamesFamily        // Games-Family
appstore.GamesKids          // Games-Kids
appstore.GamesMusic         // Games-Music
appstore.GamesPuzzle        // Games-Puzzle
appstore.GamesRacing        // Games-Racing
appstore.GamesRolePlaying   // Games-Role Playing
appstore.GamesSimulation    // Games-Simulation
appstore.GamesSports        // Games-Sports
appstore.GamesStrategy      // Games-Strategy
appstore.GamesTrivia        // Games-Trivia
appstore.GamesWord          // Games-Word
```

### Country - Country/Region Codes

Defines supported App Store country/region codes (partial list):

```go
// Asia
appstore.CountryCN  // China
appstore.CountryTW  // Taiwan
appstore.CountryHK  // Hong Kong
appstore.CountryMO  // Macau
appstore.CountryJP  // Japan
appstore.CountryKR  // South Korea
appstore.CountrySG  // Singapore
appstore.CountryTH  // Thailand
appstore.CountryMY  // Malaysia
appstore.CountryID  // Indonesia
appstore.CountryPH  // Philippines
appstore.CountryVN  // Vietnam
appstore.CountryIN  // India
appstore.CountryPK  // Pakistan

// Europe
appstore.CountryUK  // United Kingdom
appstore.CountryDE  // Germany
appstore.CountryFR  // France
appstore.CountryIT  // Italy
appstore.CountryES  // Spain
appstore.CountryRU  // Russia
appstore.CountryNL  // Netherlands
appstore.CountryBE  // Belgium
appstore.CountryCH  // Switzerland
appstore.CountrySE  // Sweden
appstore.CountryNO  // Norway
appstore.CountryDK  // Denmark
appstore.CountryFI  // Finland
appstore.CountryPL  // Poland

// North America
appstore.CountryUS  // United States
appstore.CountryCA  // Canada
appstore.CountryMX  // Mexico

// South America
appstore.CountryBR  // Brazil
appstore.CountryAR  // Argentina
appstore.CountryCL  // Chile
appstore.CountryCO  // Colombia
appstore.CountryPE  // Peru

// Oceania
appstore.CountryAU  // Australia
appstore.CountryNZ  // New Zealand

// Middle East/Africa
appstore.CountryAE  // UAE
appstore.CountrySA  // Saudi Arabia
appstore.CountryIL  // Israel
appstore.CountryZA  // South Africa
```

### Sort - Sort Order

Defines review sort orders:

```go
appstore.Recent  // Sort by most recent
appstore.Helpful // Sort by most helpful
```

## Code Examples

### Complete Usage Example

The following is a complete example demonstrating how to get app info, search for apps, and retrieve reviews:

```go
package main

import (
    "fmt"
    appstore "github.com/yililith/app-store-scraper"
)

func main() {
    scraper := appstore.NewScraper()

    // 1. Get app details
    app, err := scraper.App(appstore.AppOptions{
        ID:      553834731,
        Country: appstore.CountryUS,
    })
    if err != nil {
        fmt.Printf("Failed to get app details: %v\n", err)
        return
    }
    fmt.Printf("App name: %s\n", app.Title)
    fmt.Printf("Developer: %s\n", app.Developer)
    fmt.Printf("Version: %s\n", app.Version)
    fmt.Printf("Rating: %.1f (%d reviews)\n", app.Score, app.Reviews)

    // 2. Search for apps
    results, err := scraper.Search(appstore.SearchOptions{
        Term:    "productivity",
        Country: appstore.CountryUS,
        Num:     5,
    })
    if err != nil {
        fmt.Printf("Search failed: %v\n", err)
        return
    }
    fmt.Printf("Found %d related apps\n", len(results))

    // 3. Get reviews
    reviews, err := scraper.Reviews(appstore.ReviewsOptions{
        ID:      553834731,
        Country: appstore.CountryUS,
        Sort:    appstore.Recent,
    })
    if err != nil {
        fmt.Printf("Failed to get reviews: %v\n", err)
        return
    }
    fmt.Printf("Latest reviews (%d):\n", len(reviews))
    for i, review := range reviews {
        if i >= 3 {
            break
        }
        fmt.Printf("  [%d stars] %s\n", review.Score, review.Title)
    }

    // 4. Get rating distribution
    ratings, err := scraper.Ratings(appstore.RatingsOptions{
        ID:      553834731,
        Country: appstore.CountryUS,
    })
    if err != nil {
        fmt.Printf("Failed to get ratings: %v\n", err)
        return
    }
    fmt.Printf("Rating distribution: %d ratings\n", ratings.Ratings)
}
```

### Test Data

Test app data used in the project:

| App Name | App ID | Country | Description |
|---------|--------|---------|-------------|
| Candy Crush | 553834731 | US | Popular puzzle game |
| MariBank | 1592249158 | PH | Digital banking app |
| ChatGPT | 6448311069 | US | OpenAI's AI assistant |
| WeChat | 414478124 | CN | Tencent's social app |

## Error Handling

This library defines the following error types. It is recommended to check for errors after calling APIs:

```go
// Parameter errors
appstore.ErrInvalidParameter // Missing required parameters or invalid parameters

// Data errors
appstore.ErrNotFound          // Requested data not found
appstore.ErrInvalidResponse   // Server response format error

// Network errors
// Standard net/http errors
```

**Error Handling Example**:

```go
app, err := scraper.App(appstore.AppOptions{ID: 553834731})
if err != nil {
    switch err {
    case appstore.ErrInvalidParameter:
        fmt.Println("Parameter error: Please check App ID or Bundle ID")
    case appstore.ErrNotFound:
        fmt.Println("App not found")
    case appstore.ErrInvalidResponse:
        fmt.Println("Server response error")
    default:
        fmt.Printf("Network or other error: %v\n", err)
    }
    return
}
```

## Performance Optimization

### FullDetail Option

The List method supports two data retrieval modes:

**Fast Mode (FullDetail=false)**:
- Only calls the RSS Feed API
- Returns basic app information (ID, name, icon, rank, etc.)
- Fast response, suitable for scenarios requiring large amounts of data

**Complete Mode (FullDetail=true)**:
- First calls the RSS Feed API to get the app ID list
- Then concurrently calls the Lookup API to get detailed information for each app
- Returns complete app data, including descriptions, screenshots, etc.
- Suitable for scenarios requiring detailed information display

### Concurrent Optimization

When using FullDetail=true mode, the library internally uses goroutines to concurrently request detailed information for each app, significantly improving data retrieval efficiency. The number of concurrent requests is dynamically adjusted based on the number of target apps while ensuring the order of returned results matches the ranking order.

### Request Limits

Please note that Apple's iTunes API may have implicit rate limits. Recommendations:

- Avoid sending a large number of requests in a short period
- Use fast mode for scenarios requiring large amounts of data
- Add appropriate request intervals in production environments

## FAQ

### Q: How do I get an app's Bundle ID?

A: The Bundle ID is the unique identifier for an app in the App Store. It can be found in the app page URL. For example, WeChat's App Store URL is `https://apps.apple.com/cn/app/id414478124`, where the number after `id` (414478124) is the app ID.

### Q: How do I get a developer's ID?

A: The developer ID can be obtained by visiting the developer's App Store page. For example, OpenAI's developer page URL is `https://apps.apple.com/us/developer/openai/id284882218`, where the number after `id` (284882218) is the developer ID.

### Q: Why can't I find some apps in specific countries?

A: Some apps may only be listed on specific countries'/regions' App Stores. Please ensure you are using the correct country code, or try querying using the country where the app is listed.

### Q: How to handle request timeouts?

A: The Scraper uses default HTTP client settings. To customize the timeout, create a custom Client and set the Timeout:

```go
client := &appstore.Client{
    HTTPClient: &http.Client{
        Timeout: 30 * time.Second,
    },
}
scraper := appstore.NewScraperWithClient(client)
```

### Development Environment

- Go 1.25.6 or higher
- Code formatting tool: gofmt
- Recommended IDEs: VS Code or GoLand

### Code Standards

- Follow Go's official code standards
- All public functions must have documentation comments
- Add necessary unit tests

## License

This project is open source under the MIT License. See the [LICENSE](LICENSE) file for details.
