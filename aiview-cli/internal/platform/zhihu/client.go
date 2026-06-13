package zhihu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/cache"
	"github.com/jackwener/aiview/internal/ratelimit"
)

const (
	baseURL   = "https://www.zhihu.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Zhihu API client.
type Client struct {
	httpClient *http.Client
	cookies    string
	limiter    *ratelimit.Limiter
	cache      *cache.Cache
}

// NewClient creates a new Zhihu API client.
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
	return "zhihu"
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
		return nil, aiverr.NetworkError("zhihu", fmt.Sprintf("failed to create request: %v", err))
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.zhihu.com/")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	if c.cookies != "" {
		req.Header.Set("Cookie", c.cookies)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, aiverr.NetworkError("zhihu", fmt.Sprintf("request failed: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, aiverr.NetworkError("zhihu", fmt.Sprintf("failed to read response body: %v", err))
	}

	// Check if response is HTML instead of JSON (indicates authentication issue)
	if strings.Contains(string(body), "<!DOCTYPE") || strings.Contains(string(body), "<html") {
		return nil, aiverr.NotAuthenticated("zhihu", "API returned HTML instead of JSON; authentication may be required").WithHTTPStatus(resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 401, 403:
			return nil, aiverr.NotAuthenticated("zhihu", "Authentication required").WithHTTPStatus(resp.StatusCode)
		case 429:
			return nil, aiverr.RateLimited("zhihu", "Request too frequent").WithHTTPStatus(resp.StatusCode)
		default:
			return nil, aiverr.APIError("zhihu", fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncateBody(string(body), 200))).WithHTTPStatus(resp.StatusCode)
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
		return aiverr.NotAuthenticated("zhihu", "Authentication required").WithSuggestion("Please login with: aiview zhihu login --cookie <cookie>")
	}
	return nil
}

// truncateBody truncates a string to maxLen characters, appending "..." if truncated.
func truncateBody(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// GetHotSearch fetches zhihu hot search terms.
// API: /api/v4/search/top_search
func (c *Client) GetHotSearch() (map[string]interface{}, error) {
	return c.get("/api/v4/search/top_search", nil)
}

// Search performs a search on Zhihu.
// API: /api/v4/search_v3
func (c *Client) Search(keyword string, page int) (map[string]interface{}, error) {
	return c.get("/api/v4/search_v3", url.Values{
		"q":         {keyword},
		"t":         {"general"},
		"correction": {"1"},
		"offset":    {fmt.Sprintf("%d", (page-1)*10)},
		"limit":     {"10"},
	})
}

// GetUserInfo fetches user info by user ID.
// API: /api/v4/members/{uid}
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	return c.get("/api/v4/members/"+uid, nil)
}
