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
)
