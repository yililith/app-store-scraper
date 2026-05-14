package appstore

type LookupResponse struct {
	ResultCount int         `json:"resultCount"`
	Results     []AppResult `json:"results"`
}

type AppResult struct {
	TrackID                            int64    `json:"trackId"`
	TrackName                          string   `json:"trackName"`
	BundleID                           string   `json:"bundleId"`
	TrackViewURL                       string   `json:"trackViewUrl"`
	Description                        string   `json:"description"`
	ArtworkURL512                      string   `json:"artworkUrl512"`
	ArtworkURL100                      string   `json:"artworkUrl100"`
	Genres                             []string `json:"genres"`
	GenreIDs                           []string `json:"genreIds"`
	PrimaryGenreName                   string   `json:"primaryGenreName"`
	PrimaryGenreID                     int64    `json:"primaryGenreId"`
	ContentAdvisoryRating              string   `json:"contentAdvisoryRating"`
	LanguageCodesISO2A                 []string `json:"languageCodesISO2A"`
	FileSizeBytes                      string   `json:"fileSizeBytes"`
	MinimumOsVersion                   string   `json:"minimumOsVersion"`
	ReleaseDate                        string   `json:"releaseDate"`
	CurrentVersionReleaseDate          string   `json:"currentVersionReleaseDate"`
	ReleaseNotes                       string   `json:"releaseNotes"`
	Version                            string   `json:"version"`
	Price                              float64  `json:"price"`
	Currency                           string   `json:"currency"`
	ArtistID                           int64    `json:"artistId"`
	ArtistName                         string   `json:"artistName"`
	ArtistViewURL                      string   `json:"artistViewUrl"`
	SellerURL                          string   `json:"sellerUrl"`
	AverageUserRating                  float64  `json:"averageUserRating"`
	UserRatingCount                    int      `json:"userRatingCount"`
	AverageUserRatingForCurrentVersion float64  `json:"averageUserRatingForCurrentVersion"`
	UserRatingCountForCurrentVersion   int      `json:"userRatingCountForCurrentVersion"`
	ScreenshotURLs                     []string `json:"screenshotUrls"`
	IpadScreenshotURLs                 []string `json:"ipadScreenshotUrls"`
	AppletvScreenshotURLs              []string `json:"appletvScreenshotUrls"`
	SupportedDevices                   []string `json:"supportedDevices"`
	Kind                               string   `json:"kind"`
	WrapperType                        string   `json:"wrapperType"`
}

type RSSFeedResponse struct {
	Feed RSSFeed `json:"feed"`
}

type RSSFeed struct {
	Entry []RSSEntry `json:"entry"`
}

type RSSEntry struct {
	ID     RSSID     `json:"id"`
	Name   RSSName   `json:"name"`
	Artist RSSArtist `json:"artist"`
}

type RSSID struct {
	Attributes RSSIDAttributes `json:"attributes"`
}

type RSSIDAttributes struct {
	ImID string `json:"im:id"`
}

type RSSName struct {
	Label string `json:"label"`
}

type RSSArtist struct {
	Label string `json:"label"`
}

type ReviewFeedResponse struct {
	Feed ReviewFeed `json:"feed"`
}

type ReviewFeed struct {
	Entry []ReviewEntry `json:"entry"`
}

type ReviewEntry struct {
	ID        ReviewID        `json:"id"`
	Author    ReviewAuthor    `json:"author"`
	ImRating  ReviewImRating  `json:"im:rating"`
	Title     ReviewTitle     `json:"title"`
	Content   ReviewContent   `json:"content"`
	Updated   ReviewUpdated   `json:"updated"`
	ImVersion ReviewImVersion `json:"im:version"`
}

type ReviewID struct {
	Label string `json:"label"`
}

type ReviewAuthor struct {
	Name ReviewAuthorName `json:"name"`
	URI  ReviewAuthorURI  `json:"uri"`
}

type ReviewAuthorName struct {
	Label string `json:"label"`
}

type ReviewAuthorURI struct {
	Label string `json:"label"`
}

type ReviewImRating struct {
	Label string `json:"label"`
}

type ReviewTitle struct {
	Label string `json:"label"`
}

type ReviewContent struct {
	Label string `json:"label"`
}

type ReviewUpdated struct {
	Label string `json:"label"`
}

type ReviewImVersion struct {
	Label string `json:"label"`
}
