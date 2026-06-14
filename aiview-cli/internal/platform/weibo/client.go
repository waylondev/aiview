package weibo

import (
	"fmt"
	"net/url"

	"github.com/jackwener/aiview/internal/platform/base"
	aiverr "github.com/jackwener/aiview/internal/errors"
)

const (
	baseURL   = "https://weibo.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Weibo API client.
type Client struct {
	*base.Client
}

// NewClient creates a new Weibo API client.
func NewClient(timeoutSec int, cookies string) *Client {
	return &Client{
		Client: base.NewClient(timeoutSec, cookies, userAgent, baseURL, "https://weibo.com/", "weibo"),
	}
}

// PlatformName implements the platform.Client interface.
func (c *Client) PlatformName() string {
	return "weibo"
}

// checkAuth verifies that the client has authentication credentials.
func (c *Client) checkAuth() error {
	if c.Cookies == "" {
		return aiverr.NotAuthenticated("weibo", "Authentication required").
			WithSuggestion("Please login with: aiview weibo login --cookie <cookie>")
	}
	return nil
}

// GetHotSearch fetches weibo hot search terms.
// API: /ajax/side/hotSearch
func (c *Client) GetHotSearch() (map[string]interface{}, error) {
	return c.Get("/ajax/side/hotSearch", nil)
}

// Search performs a search on Weibo.
// API: /ajax/statuses/search
func (c *Client) Search(keyword string, page int) (map[string]interface{}, error) {
	return c.Get("/ajax/statuses/search", url.Values{
		"q":    {keyword},
		"page": {fmt.Sprintf("%d", page)},
	})
}

// GetUserInfo fetches user info by user ID.
// API: /ajax/profile/info
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	return c.Get("/ajax/profile/info", url.Values{
		"uid": {uid},
	})
}
