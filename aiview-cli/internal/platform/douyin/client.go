package douyin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jackwener/aiview/internal/cache"
	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/ratelimit"
)

const (
	baseURL   = "https://www.douyin.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Douyin API client.
type Client struct {
	httpClient *http.Client
	cookies    string
	limiter    *ratelimit.Limiter
	cache      *cache.Cache
}

// NewClient creates a new Douyin API client.
func NewClient(timeoutSec int, cookies string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
		cookies: cookies,
		limiter: ratelimit.New(2, 5),
		cache:   cache.New(5 * time.Minute),
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

// get sends a GET request to the specified path.
func (c *Client) get(path string, params url.Values) (map[string]interface{}, error) {
	u := baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}

	// Check cache first
	if cached, ok := c.cache.Get(u); ok {
		return cached.(map[string]interface{}), nil
	}

	// Rate limit
	c.limiter.Wait()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, aiverr.NetworkError("douyin", fmt.Sprintf("failed to create request: %v", err))
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.douyin.com/")
	if c.cookies != "" {
		req.Header.Set("Cookie", c.cookies)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, aiverr.NetworkError("douyin", fmt.Sprintf("request failed: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, aiverr.NetworkError("douyin", fmt.Sprintf("failed to read response body: %v", err))
	}

	if resp.StatusCode != http.StatusOK {
		bodyStr := string(body)
		// Check for authentication errors
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
			return nil, aiverr.NotAuthenticated("douyin", "Authentication required").
				WithHTTPStatus(resp.StatusCode).
				WithCommand("video")
		}
		// Check for rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, aiverr.RateLimited("douyin", "Request too frequent").
				WithHTTPStatus(resp.StatusCode)
		}
		// Check for HTML response (authentication page)
		if strings.Contains(bodyStr, "<!DOCTYPE") || strings.Contains(bodyStr, "<html") {
			return nil, aiverr.NotAuthenticated("douyin", "API returned HTML instead of JSON; authentication may be required").
				WithHTTPStatus(resp.StatusCode)
		}
		return nil, aiverr.APIError("douyin", fmt.Sprintf("HTTP %d: %s", resp.StatusCode, bodyStr)).
			WithHTTPStatus(resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// Try parsing as raw response
		bodyStr := string(body)
		if strings.Contains(bodyStr, "<!DOCTYPE") || strings.Contains(bodyStr, "<html") {
			return nil, aiverr.NotAuthenticated("douyin", "API returned HTML instead of JSON; authentication may be required")
		}
		result = map[string]interface{}{"raw": bodyStr}
	}

	// Store in cache
	c.cache.Set(u, result)
	return result, nil
}

// GetHotSearch fetches hot/douyin trending search terms.
// Uses the hot search API endpoint.
func (c *Client) GetHotSearch() (map[string]interface{}, error) {
	// Douyin hot search API
	return c.get("/aweme/v1/web/hot/search/list/", url.Values{
		"detail_list": {"1"},
		"source":      {"6"},
	})
}

// GetTrending fetches the trending/challenge list.
func (c *Client) GetTrending() (map[string]interface{}, error) {
	return c.get("/aweme/v1/web/hot/search/list/", url.Values{
		"detail_list": {"1"},
		"source":      {"2"},
	})
}

// checkAuth verifies that the client has authentication credentials.
func (c *Client) checkAuth() error {
	if c.cookies == "" {
		return aiverr.NotAuthenticated("douyin", "Authentication required").
			WithSuggestion("Please login with: aiview douyin login --cookie <cookie>")
	}
	return nil
}

// GetVideoDetail fetches video details by video ID.
// API: /aweme/v1/web/aweme/detail/?aweme_id={videoID}
func (c *Client) GetVideoDetail(videoID string) (map[string]interface{}, error) {
	if err := c.checkAuth(); err != nil {
		return nil, err
	}
	return c.get("/aweme/v1/web/aweme/detail/", url.Values{
		"aweme_id": {videoID},
	})
}

// GetUserInfo fetches user info by uid.
// API: /aweme/v1/web/user/profile/other/?sec_user_id={uid}
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	if err := c.checkAuth(); err != nil {
		return nil, err
	}
	return c.get("/aweme/v1/web/user/profile/other/", url.Values{
		"sec_user_id": {uid},
	})
}

// GetVideoComments fetches comments for a video.
// API: /aweme/v1/web/comment/list/?aweme_id={videoID}&cursor={cursor}&count=20
func (c *Client) GetVideoComments(videoID string, cursor int) (map[string]interface{}, error) {
	return c.get("/aweme/v1/web/comment/list/", url.Values{
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
	return c.get("/aweme/v1/web/aweme/post/", url.Values{
		"sec_user_id": {uid},
		"cursor":      {fmt.Sprintf("%d", cursor)},
		"count":       {"20"},
	})
}

// Search performs a search on Douyin for videos/users.
func (c *Client) Search(keyword string, page int, count int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("search_source", "normal_search")
	params.Set("sort_type", "0")
	params.Set("publish_time", "0")
	params.Set("offset", fmt.Sprintf("%d", (page-1)*count))
	params.Set("count", fmt.Sprintf("%d", count))
	return c.get("/aweme/v1/web/general/search/single/", params)
}