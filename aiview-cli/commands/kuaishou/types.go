// Package kuaishou provides CLI commands for the Kuaishou (快手) platform.
package kuaishou

// Client defines the interface that the Kuaishou API client must satisfy for commands.
type Client interface {
	PlatformName() string
	GetHotSearch(count ...int) (map[string]interface{}, error)
	Search(query string, page int, count ...int) (map[string]interface{}, error)
	GetUserInfo(uid string) (map[string]interface{}, error)
}

// Credential holds Kuaishou authentication data.
type Credential struct {
	Cookie string `json:"cookie"`
}
