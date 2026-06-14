package douyin

import (
	"fmt"
	"net/url"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
)

const (
	baseURL   = "https://www.douyin.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Douyin API client.
type Client struct {
	*base.Client
}

// Compile-time interface assertions
var _ platform.HotSearchable = (*Client)(nil)
var _ platform.Searchable = (*Client)(nil)
var _ platform.UserQueryable = (*Client)(nil)

// NewClient creates a new Douyin API client.
func NewClient(timeoutSec int, cookies string) *Client {
	return &Client{
		Client: base.NewClient(timeoutSec, cookies, userAgent, baseURL, "https://www.douyin.com/", "douyin"),
	}
}

// PlatformName implements the platform.Client interface.
func (c *Client) PlatformName() string {
	return "douyin"
}

// getResponse is a generic response wrapper for Douyin API responses.
type getResponse struct {
	StatusCode    int                    `json:"status_code"`
	StatusMsg     string                 `json:"status_msg"`
	Data          map[string]interface{} `json:"data"`
}

// checkAuth verifies that the client has authentication credentials.
func (c *Client) checkAuth() error {
	if c.Cookies == "" {
		return aiverr.NotAuthenticated("douyin", "Authentication required").
			WithSuggestion("Please login with: aiview douyin login --cookie <cookie>")
	}
	return nil
}

// GetHotSearch fetches hot/douyin trending search terms.
// Uses the hot search API endpoint.
func (c *Client) GetHotSearch(count ...int) (map[string]interface{}, error) {
	// Douyin hot search API
	return c.Get("/aweme/v1/web/hot/search/list/", url.Values{
		"detail_list": {"1"},
		"source":      {"6"},
	})
}

// GetTrending fetches the trending/challenge list.
func (c *Client) GetTrending() (map[string]interface{}, error) {
	return c.Get("/aweme/v1/web/hot/search/list/", url.Values{
		"detail_list": {"1"},
		"source":      {"2"},
	})
}

// GetVideoDetail fetches video details by video ID.
// API: /aweme/v1/web/aweme/detail/?aweme_id={videoID}
func (c *Client) GetVideoDetail(videoID string) (map[string]interface{}, error) {
	if err := c.checkAuth(); err != nil {
		return nil, err
	}
	return c.Get("/aweme/v1/web/aweme/detail/", url.Values{
		"aweme_id": {videoID},
	})
}

// GetUserInfo fetches user info by uid.
// API: /aweme/v1/web/user/profile/other/?sec_user_id={uid}
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	if err := c.checkAuth(); err != nil {
		return nil, err
	}
	return c.Get("/aweme/v1/web/user/profile/other/", url.Values{
		"sec_user_id": {uid},
	})
}

// GetVideoComments fetches comments for a video.
// API: /aweme/v1/web/comment/list/?aweme_id={videoID}&cursor={cursor}&count=20
func (c *Client) GetVideoComments(videoID string, cursor int) (map[string]interface{}, error) {
	return c.Get("/aweme/v1/web/comment/list/", url.Values{
		"aweme_id": {videoID},
		"cursor":   {fmt.Sprintf("%d", cursor)},
		"count":    {"20"},
	})
}

// GetUserPosts fetches posts by a user.
// API: /aweme/v1/web/aweme/post/?sec_user_id={uid}&cursor={cursor}&count=20
func (c *Client) GetUserPosts(uid string, cursor int) (map[string]interface{}, error) {
	if err := c.checkAuth(); err != nil {
		return nil, err
	}
	return c.Get("/aweme/v1/web/aweme/post/", url.Values{
		"sec_user_id": {uid},
		"cursor":      {fmt.Sprintf("%d", cursor)},
		"count":       {"20"},
	})
}

// Search performs a search on Douyin for videos/users.
func (c *Client) Search(query string, page int, count ...int) (map[string]interface{}, error) {
	cnt := 20
	if len(count) > 0 && count[0] > 0 {
		cnt = count[0]
	}
	params := url.Values{}
	params.Set("keyword", query)
	params.Set("search_source", "normal_search")
	params.Set("sort_type", "0")
	params.Set("publish_time", "0")
	params.Set("offset", fmt.Sprintf("%d", (page-1)*cnt))
	params.Set("count", fmt.Sprintf("%d", cnt))
	return c.Get("/aweme/v1/web/general/search/single/", params)
}