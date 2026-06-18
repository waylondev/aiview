package kuaishou

import (
	"fmt"
	"net/url"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
)

const (
	baseURL = "https://www.kuaishou.com"

	// HTTP client constants
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
	defaultReferer   = "https://www.kuaishou.com/"
)

// Client is the Kuaishou API client.
type Client struct {
	*base.Client
}

// Compile-time interface assertions
var _ platform.HotSearchable = (*Client)(nil)
var _ platform.Searchable = (*Client)(nil)
var _ platform.UserQueryable = (*Client)(nil)

// NewClient creates a new Kuaishou API client.
func NewClient(timeoutSec int, cookies string) *Client {
	return &Client{
		Client: base.NewClient(timeoutSec, cookies, defaultUserAgent, baseURL, defaultReferer, "kuaishou"),
	}
}

// PlatformName implements the platform.Client interface.
func (c *Client) PlatformName() string {
	return "kuaishou"
}

// checkAuth verifies that the client has authentication credentials.
func (c *Client) checkAuth() error {
	if c.Cookies == "" {
		return aiverr.NotAuthenticated("kuaishou", "Authentication required").WithSuggestion("Please login with: aiview kuaishou login --cookie <cookie>")
	}
	return nil
}

// GetHotSearch fetches kuaishou hot search terms.
// API: /graphql (with query for hot search)
func (c *Client) GetHotSearch(count ...int) (map[string]interface{}, error) {
	return c.Get("/graphql", url.Values{
		"operationName": []string{"visionHotRank"},
		"variables":     []string{"{}"},
		"query":         []string{"query visionHotRank { visionHotRank { items { name hotValue } } }"},
	})
}

// Search performs a search on Kuaishou.
// API: /graphql (with search query)
func (c *Client) Search(query string, page int, count ...int) (map[string]interface{}, error) {
	return c.Get("/graphql", url.Values{
		"operationName": []string{"visionSearchPhoto"},
		"variables":     []string{fmt.Sprintf(`{"keyword":"%s","page":%d}`, query, page)},
		"query":         []string{"query visionSearchPhoto($keyword: String, $page: Int) { visionSearchPhoto(keyword: $keyword, page: $page) { feeds { id caption } } }"},
	})
}

// GetUserInfo fetches user info by user ID.
// API: /graphql (with user profile query)
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	return c.Get("/graphql", url.Values{
		"operationName": []string{"visionProfile"},
		"variables":     []string{fmt.Sprintf(`{"userId":"%s"}`, uid)},
		"query":         []string{"query visionProfile($userId: String) { visionProfile(userId: $userId) { user { id name followerCount } } }"},
	})
}
