// Package aiview provides a Go library API for accessing multiple social media platforms.
package aiview

// VideoInfo represents a video item.
type VideoInfo struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	AuthorID string `json:"author_id"`
	Play     int    `json:"play"`
	Danmaku  int    `json:"danmaku"`
	Like     int    `json:"like"`
	Coin     int    `json:"coin"`
	Favorite int    `json:"favorite"`
	Share    int    `json:"share"`
	Duration string `json:"duration"`
	URL      string `json:"url"`
}

// UserInfo represents a user profile.
type UserInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Sign   string `json:"sign"`
	Fans   int    `json:"fans"`
	Follow int    `json:"follow"`
	Videos int    `json:"videos"`
}

// HotItem represents a hot/trending item.
type HotItem struct {
	Keyword  string `json:"keyword"`
	HotValue int    `json:"hot_value"`
	Position int    `json:"position"`
	URL      string `json:"url"`
}

// SearchItem represents a search result.
type SearchItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`
}

// Platform represents a supported platform.
type Platform string

const (
	PlatformBilibili    Platform = "bilibili"
	PlatformDouyin      Platform = "douyin"
	PlatformXiaohongshu Platform = "xiaohongshu"
	PlatformWeibo       Platform = "weibo"
	PlatformKuaishou    Platform = "kuaishou"
	PlatformZhihu       Platform = "zhihu"
)

// SupportedPlatforms returns all supported platforms.
func SupportedPlatforms() []Platform {
	return []Platform{
		PlatformBilibili,
		PlatformDouyin,
		PlatformXiaohongshu,
		PlatformWeibo,
		PlatformKuaishou,
		PlatformZhihu,
	}
}

// HotSearchable defines the interface for platforms that support hot search.
type HotSearchable interface {
	GetHotSearch() ([]HotItem, error)
}

// Searchable defines the interface for platforms that support search.
type Searchable interface {
	Search(keyword string, page int) ([]SearchItem, error)
}

// UserQueryable defines the interface for platforms that support user queries.
type UserQueryable interface {
	GetUserInfo(uid string) (*UserInfo, error)
}

// VideoQueryable defines the interface for platforms that support video queries (Bilibili, Douyin).
type VideoQueryable interface {
	GetVideoInfo(videoID string) (*VideoInfo, error)
}

// NoteQueryable defines the interface for platforms that support note queries (Xiaohongshu).
type NoteQueryable interface {
	GetNoteDetail(noteID string) (*VideoInfo, error)
}
