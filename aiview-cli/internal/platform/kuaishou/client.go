package kuaishou

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jackwener/aiview/internal/cache"
	"github.com/jackwener/aiview/internal/ratelimit"
)

const (
	baseURL   = "https://www.kuaishou.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Kuaishou API client.
type Client struct {
	httpClient *http.Client
	cookies    string
	limiter    *ratelimit.Limiter
	cache      *cache.Cache
}

// NewClient creates a new Kuaishou API client.
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
	return "kuaishou"
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
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.kuaishou.com/")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	if c.cookies != "" {
		req.Header.Set("Cookie", c.cookies)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
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
		return fmt.Errorf("请先执行 'aiview kuaishou login --cookie <cookie>' 登录")
	}
	return nil
}

// GetHotSearch fetches kuaishou hot search terms.
// API: /graphql (with query for hot search)
func (c *Client) GetHotSearch() (map[string]interface{}, error) {
	return c.get("/graphql", url.Values{
		"operationName": []string{"visionHotRank"},
		"variables":     []string{"{}"},
		"query":         []string{"query visionHotRank { visionHotRank { items { name hotValue } } }"},
	})
}

// Search performs a search on Kuaishou.
// API: /graphql (with search query)
func (c *Client) Search(keyword string, page int) (map[string]interface{}, error) {
	return c.get("/graphql", url.Values{
		"operationName": []string{"visionSearchPhoto"},
		"variables":     []string{fmt.Sprintf(`{"keyword":"%s","page":%d}`, keyword, page)},
		"query":         []string{"query visionSearchPhoto($keyword: String, $page: Int) { visionSearchPhoto(keyword: $keyword, page: $page) { feeds { id caption } } }"},
	})
}

// GetUserInfo fetches user info by user ID.
// API: /graphql (with user profile query)
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	return c.get("/graphql", url.Values{
		"operationName": []string{"visionProfile"},
		"variables":     []string{fmt.Sprintf(`{"userId":"%s"}`, uid)},
		"query":         []string{"query visionProfile($userId: String) { visionProfile(userId: $userId) { user { id name followerCount } } }"},
	})
}
