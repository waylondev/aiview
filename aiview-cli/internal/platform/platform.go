// Package platform defines the core Platform interface and registry for multi-platform CLI applications.
package platform

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/spf13/cobra"
)

// Client is the interface that all platform API clients must implement.
type Client interface {
	// PlatformName returns the platform name.
	PlatformName() string
}

// Platform defines the interface for a platform implementation.
type Platform interface {
	// Name returns the platform identifier (e.g., "bilibili", "twitter").
	Name() string
	// Commands returns all cobra commands for this platform.
	Commands() []*cobra.Command
	// NewClient creates a new API client for this platform.
	NewClient(cfg *config.Config) (Client, error)
}

// HotSearchable is the interface for platforms that support hot search/trending.
type HotSearchable interface {
	// GetHotSearch returns hot/trending search items. The optional count parameter
	// specifies the maximum number of results to return.
	GetHotSearch(count ...int) (map[string]interface{}, error)
}

// Searchable is the interface for platforms that support content search.
type Searchable interface {
	// Search performs a content search with the given query and page number.
	// The optional count parameter specifies the number of results per page.
	Search(query string, page int, count ...int) (map[string]interface{}, error)
}

// UserQueryable is the interface for platforms that support user info queries.
type UserQueryable interface {
	// GetUserInfo fetches user profile information by user ID.
	GetUserInfo(uid string) (map[string]interface{}, error)
}