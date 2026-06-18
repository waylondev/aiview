package xiaohongshu

import (
	"fmt"
	"net/http"
	"net/url"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
)

const (
	baseURL = "https://edith.xiaohongshu.com"

	// HTTP client constants
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
	defaultReferer   = "https://www.xiaohongshu.com"
)

// Client is the Xiaohongshu API client.
type Client struct {
	*base.Client
}

// Compile-time interface assertions
var _ platform.HotSearchable = (*Client)(nil)
var _ platform.Searchable = (*Client)(nil)
var _ platform.UserQueryable = (*Client)(nil)

// NewClient creates a new Xiaohongshu API client.
func NewClient(timeoutSec int, cookies string) *Client {
	return &Client{
		Client: base.NewClient(timeoutSec, cookies, defaultUserAgent, baseURL, defaultReferer, "xiaohongshu"),
	}
}

// PlatformName implements the platform.Client interface.
func (c *Client) PlatformName() string {
	return "xiaohongshu"
}

// BuildHeaders overrides the default to add xiaohongshu-specific headers.
func (c *Client) BuildHeaders() http.Header {
	h := c.Client.BuildHeaders()
	h.Set("Origin", defaultReferer)
	h.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	h.Set("sec-ch-ua", `"Chromium";v="133", "Not(A:Brand";v="99", "Google Chrome";v="133"`)
	h.Set("sec-ch-ua-mobile", "?0")
	h.Set("sec-ch-ua-platform", `"Windows"`)
	return h
}

// checkAuth verifies that the client has authentication credentials.
func (c *Client) checkAuth() error {
	if c.Cookies == "" {
		return aiverr.NotAuthenticated("xiaohongshu", "Authentication required").
			WithSuggestion("Please login with: aiview xiaohongshu login --cookie <cookie>")
	}
	return nil
}

// toValues converts a map[string]string to url.Values.
func toValues(params map[string]string) url.Values {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return v
}

// GetHotNotes returns hot/trending notes.
func (c *Client) GetHotNotes() (map[string]interface{}, error) {
	return c.Get("/api/sns/web/v1/search/hot", nil)
}

// GetHotSearch implements platform.HotSearchable.
func (c *Client) GetHotSearch(count ...int) (map[string]interface{}, error) {
	return c.Get("/api/sns/web/v1/search/hot", nil)
}

// SearchNotes searches notes by keyword.
func (c *Client) SearchNotes(keyword string, page int) (map[string]interface{}, error) {
	return c.Get("/api/sns/web/v1/search/notes", toValues(map[string]string{
		"keyword":   keyword,
		"page":      fmt.Sprintf("%d", page),
		"sort":      "general",
		"note_type": "0",
	}))
}

// Search implements platform.Searchable.
func (c *Client) Search(query string, page int, count ...int) (map[string]interface{}, error) {
	return c.SearchNotes(query, page)
}

// GetNoteDetail returns note detail by note ID.
func (c *Client) GetNoteDetail(noteID string) (map[string]interface{}, error) {
	return c.Get("/api/sns/web/v1/feed", toValues(map[string]string{
		"source_note_id": noteID,
	}))
}

// GetUserInfo returns user info by user ID.
func (c *Client) GetUserInfo(userID string) (map[string]interface{}, error) {
	return c.Get("/api/sns/web/v1/user/otherinfo", toValues(map[string]string{
		"target_user_id": userID,
	}))
}
