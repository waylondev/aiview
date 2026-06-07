package douyin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL   = "https://www.douyin.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Douyin API client.
type Client struct {
	httpClient *http.Client
	cookies    string
}

// NewClient creates a new Douyin API client.
func NewClient(timeoutSec int, cookies string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
		cookies: cookies,
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

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.douyin.com/")
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
		// Try parsing as raw response
		return map[string]interface{}{"raw": string(body)}, nil
	}
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