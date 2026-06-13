package weibo

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
	baseURL   = "https://weibo.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Weibo API client.
type Client struct {
	httpClient *http.Client
	cookies    string
	limiter    *ratelimit.Limiter
	cache      *cache.Cache
}

// NewClient creates a new Weibo API client.
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
	return "weibo"
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
		return nil, aiverr.NetworkError("weibo", fmt.Sprintf("failed to create request: %v", err))
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://weibo.com/")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	if c.cookies != "" {
		req.Header.Set("Cookie", c.cookies)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, aiverr.NetworkError("weibo", fmt.Sprintf("request failed: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, aiverr.NetworkError("weibo", fmt.Sprintf("failed to read response body: %v", err))
	}

	// 检测 HTML 响应（未登录或被风控）
	bodyStr := string(body)
	if strings.Contains(bodyStr, "<!DOCTYPE") || strings.Contains(bodyStr, "<html") {
		return nil, aiverr.NotAuthenticated("weibo", "API returned HTML instead of JSON; authentication may be required").WithHTTPStatus(resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return nil, aiverr.NotAuthenticated("weibo", "Authentication required").WithHTTPStatus(resp.StatusCode)
		case http.StatusTooManyRequests:
			return nil, aiverr.RateLimited("weibo", "Request too frequent").WithHTTPStatus(resp.StatusCode)
		default:
			return nil, aiverr.APIError("weibo", fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncateBody(bodyStr, 200))).WithHTTPStatus(resp.StatusCode)
		}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		result = map[string]interface{}{"raw": string(body)}
	}

	// Store in cache
	c.cache.Set(u, result)
	return result, nil
}

// checkAuth verifies that the client has authentication credentials.
func (c *Client) checkAuth() error {
	if c.cookies == "" {
		return aiverr.NotAuthenticated("weibo", "Authentication required").
			WithSuggestion("Please login with: aiview weibo login --cookie <cookie>")
	}
	return nil
}

// GetHotSearch fetches weibo hot search terms.
// API: /ajax/side/hotSearch
func (c *Client) GetHotSearch() (map[string]interface{}, error) {
	return c.get("/ajax/side/hotSearch", nil)
}

// Search performs a search on Weibo.
// API: /ajax/statuses/search
func (c *Client) Search(keyword string, page int) (map[string]interface{}, error) {
	return c.get("/ajax/statuses/search", url.Values{
		"q":    {keyword},
		"page": {fmt.Sprintf("%d", page)},
	})
}

// GetUserInfo fetches user info by user ID.
// API: /ajax/profile/info
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	return c.get("/ajax/profile/info", url.Values{
		"uid": {uid},
	})
}

// truncateBody truncates a string to maxLen characters, appending "..." if truncated.
func truncateBody(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
