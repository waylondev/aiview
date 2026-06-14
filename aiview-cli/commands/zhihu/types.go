// Package zhihu provides CLI commands for the Zhihu (知乎) platform.
package zhihu

// Client defines the interface that the Zhihu API client must satisfy for commands.
type Client interface {
	PlatformName() string
	GetHotSearch(count ...int) (map[string]interface{}, error)
	Search(query string, page int, count ...int) (map[string]interface{}, error)
	GetUserInfo(uid string) (map[string]interface{}, error)
}

// Credential holds Zhihu authentication data.
type Credential struct {
	Cookie string `json:"cookie"`
}
