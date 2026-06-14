// Package weibo provides CLI commands for the Weibo (微博) platform.
package weibo

// Client defines the interface that the Weibo API client must satisfy for commands.
type Client interface {
	PlatformName() string
	GetHotSearch(count ...int) (map[string]interface{}, error)
	Search(query string, page int, count ...int) (map[string]interface{}, error)
	GetUserInfo(uid string) (map[string]interface{}, error)
}

// Credential holds Weibo authentication data.
type Credential struct {
	Cookie string `json:"cookie"`
}
