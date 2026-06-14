// Package douyin provides CLI commands for the Douyin (抖音) platform.
package douyin

// Client is the interface that the Douyin API client must satisfy for commands.
type Client interface {
	PlatformName() string
	GetHotSearch(count ...int) (map[string]interface{}, error)
	GetTrending() (map[string]interface{}, error)
	Search(query string, page int, count ...int) (map[string]interface{}, error)
	GetVideoDetail(videoID string) (map[string]interface{}, error)
	GetVideoComments(videoID string, cursor int) (map[string]interface{}, error)
	GetUserPosts(uid string, cursor int) (map[string]interface{}, error)
	GetUserInfo(uid string) (map[string]interface{}, error)
}

// Credential holds Douyin authentication data.
type Credential struct {
	Cookie string `json:"cookie"`
}

