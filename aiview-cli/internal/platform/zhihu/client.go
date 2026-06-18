package zhihu

import (
	"fmt"
	"net/url"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
)

const (
	baseURL = "https://www.zhihu.com"

	// HTTP client constants
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
	defaultReferer   = "https://www.zhihu.com/"
)

// Client is the Zhihu API client.
type Client struct {
	*base.Client
}

// Compile-time interface assertions
var _ platform.HotSearchable = (*Client)(nil)
var _ platform.Searchable = (*Client)(nil)
var _ platform.UserQueryable = (*Client)(nil)

// NewClient creates a new Zhihu API client.
func NewClient(timeoutSec int, cookies string) *Client {
	return &Client{
		Client: base.NewClient(timeoutSec, cookies, defaultUserAgent, baseURL, defaultReferer, "zhihu"),
	}
}

// PlatformName implements the platform.Client interface.
func (c *Client) PlatformName() string {
	return "zhihu"
}

// checkAuth verifies that the client has authentication credentials.
func (c *Client) checkAuth() error {
	if c.Cookies == "" {
		return aiverr.NotAuthenticated("zhihu", "Authentication required").WithSuggestion("Please login with: aiview zhihu login --cookie <cookie>")
	}
	return nil
}

// GetHotSearch fetches zhihu hot search terms.
// API: /api/v4/search/top_search
func (c *Client) GetHotSearch(count ...int) (map[string]interface{}, error) {
	return c.Get("/api/v4/search/top_search", nil)
}

// Search performs a search on Zhihu.
// API: /api/v4/search_v3
func (c *Client) Search(query string, page int, count ...int) (map[string]interface{}, error) {
	return c.Get("/api/v4/search_v3", url.Values{
		"q":          {query},
		"t":          {"general"},
		"correction": {"1"},
		"offset":     {fmt.Sprintf("%d", (page-1)*10)},
		"limit":      {"10"},
	})
}

// GetUserInfo fetches user info by user ID.
// API: /api/v4/members/{uid}
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	return c.Get("/api/v4/members/"+uid, nil)
}
