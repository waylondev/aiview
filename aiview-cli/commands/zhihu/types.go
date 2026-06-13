// Package zhihu provides CLI commands for the Zhihu (知乎) platform.
package zhihu

// Client defines the interface that the Zhihu API client must satisfy for commands.
type Client interface {
	PlatformName() string
	GetHotSearch() (map[string]interface{}, error)
	Search(keyword string, page int) (map[string]interface{}, error)
	GetUserInfo(uid string) (map[string]interface{}, error)
}

// Credential holds Zhihu authentication data.
type Credential struct {
	Cookie string `json:"cookie"`
}
