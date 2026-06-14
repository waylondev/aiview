// Package base provides a shared HTTP client abstraction for platform API clients.
package base

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

// Client contains the common fields shared by all platform API clients.
type Client struct {
	HTTPClient *http.Client
	Cookies    string
	Limiter    *ratelimit.Limiter
	Cache      *cache.Cache
	UserAgent  string
	BaseURL    string
	Referer    string
	Platform   string
}

// NewClient creates a new base Client with the given parameters.
func NewClient(timeoutSec int, cookies, userAgent, baseURL, referer, platform string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
		Cookies:   cookies,
		Limiter:   ratelimit.New(2, 5),
		Cache:     cache.New(5 * time.Minute),
		UserAgent: userAgent,
		BaseURL:   baseURL,
		Referer:   referer,
		Platform:  platform,
	}
}

// BuildHeaders returns a default set of browser-like HTTP headers.
func (c *Client) BuildHeaders() http.Header {
	h := http.Header{}
	h.Set("User-Agent", c.UserAgent)
	h.Set("Accept", "application/json, text/plain, */*")
	h.Set("Accept-Language", "zh-CN,zh;q=0.9")
	if c.Referer != "" {
		h.Set("Referer", c.Referer)
	}
	if c.Cookies != "" {
		h.Set("Cookie", c.Cookies)
	}
	return h
}

// Get sends a GET request to the specified path with query parameters.
// It handles caching, rate limiting, and standard error parsing.
func (c *Client) Get(path string, params url.Values) (map[string]interface{}, error) {
	reqURL := c.BaseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	// Check cache first
	if cached, ok := c.Cache.Get(reqURL); ok {
		return cached.(map[string]interface{}), nil
	}

	// Rate limit
	c.Limiter.Wait()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, aiverr.NetworkError(c.Platform, fmt.Sprintf("failed to create request: %v", err))
	}
	req.Header = c.BuildHeaders()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, aiverr.NetworkError(c.Platform, fmt.Sprintf("network request failed: %v", err))
	}
	defer resp.Body.Close()

	result, err := c.ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	c.Cache.Set(reqURL, result)
	return result, nil
}

// ParseResponse reads the HTTP response body and parses it as JSON.
// It detects HTML responses and maps common HTTP status codes to platform errors.
func (c *Client) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, aiverr.NetworkError(c.Platform, fmt.Sprintf("failed to read response: %v", err))
	}

	if c.DetectHTMLResponse(body) {
		return nil, aiverr.NotAuthenticated(c.Platform, "API returned HTML instead of JSON; authentication may be required").
			WithHTTPStatus(resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		bodyStr := truncate(string(body), 200)
		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return nil, aiverr.NotAuthenticated(c.Platform, "Authentication required").
				WithHTTPStatus(resp.StatusCode)
		case http.StatusTooManyRequests:
			return nil, aiverr.RateLimited(c.Platform, "Request too frequent").
				WithHTTPStatus(resp.StatusCode)
		default:
			return nil, aiverr.APIError(c.Platform, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, bodyStr)).
				WithHTTPStatus(resp.StatusCode)
		}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, aiverr.ParseError(c.Platform, fmt.Sprintf("failed to parse response: %v", err))
	}

	return result, nil
}

// DetectHTMLResponse returns true if the response body appears to be HTML.
func (c *Client) DetectHTMLResponse(body []byte) bool {
	bodyStr := string(body)
	trimmed := strings.TrimSpace(bodyStr)
	return strings.HasPrefix(trimmed, "<!DOCTYPE") ||
		strings.HasPrefix(trimmed, "<!doctype") ||
		strings.HasPrefix(trimmed, "<html")
}

// truncate shortens a string to maxLen characters.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
