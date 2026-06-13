package xiaohongshu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
)

const (
	baseURL   = "https://edith.xiaohongshu.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Xiaohongshu API client.
type Client struct {
	httpClient *http.Client
	cookies    string
}

// NewClient creates a new Xiaohongshu API client.
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
	return "xiaohongshu"
}

// get sends a GET request to the specified path with query parameters.
func (c *Client) get(path string, params map[string]string) (map[string]interface{}, error) {
	url := baseURL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, aiverr.NetworkError("xiaohongshu", fmt.Sprintf("create request: %v", err))
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Set browser-like headers to avoid HTML responses
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.xiaohongshu.com")
	req.Header.Set("Origin", "https://www.xiaohongshu.com")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("sec-ch-ua", `"Chromium";v="133", "Not(A:Brand";v="99", "Google Chrome";v="133"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	if c.cookies != "" {
		req.Header.Set("Cookie", c.cookies)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, aiverr.NetworkError("xiaohongshu", fmt.Sprintf("request failed: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, aiverr.NetworkError("xiaohongshu", fmt.Sprintf("read response: %v", err))
	}

	if resp.StatusCode != http.StatusOK {
		bodyStr := truncate(string(body), 200)
		// Check for authentication errors
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
			return nil, aiverr.NotAuthenticated("xiaohongshu", "Authentication required").
				WithHTTPStatus(resp.StatusCode)
		}
		// Check for rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, aiverr.RateLimited("xiaohongshu", "Request too frequent").
				WithHTTPStatus(resp.StatusCode)
		}
		return nil, aiverr.APIError("xiaohongshu", fmt.Sprintf("HTTP %d: %s", resp.StatusCode, bodyStr)).
			WithHTTPStatus(resp.StatusCode)
	}

	return c.parseResponse(resp.Header.Get("Content-Type"), body)
}

// parseResponse detects HTML responses and returns a friendly error, otherwise parses JSON.
func (c *Client) parseResponse(contentType string, body []byte) (map[string]interface{}, error) {
	bodyStr := string(body)
	trimmed := strings.TrimSpace(bodyStr)

	// Detect HTML response (API returns login page or error page instead of JSON)
	if strings.Contains(contentType, "text/html") ||
		strings.HasPrefix(trimmed, "<!DOCTYPE") ||
		strings.HasPrefix(trimmed, "<html") ||
		strings.HasPrefix(trimmed, "<!doctype") {
		return nil, aiverr.NotAuthenticated("xiaohongshu", "API returned HTML instead of JSON; authentication may be required").
			WithSuggestion("Please login with: aiview xiaohongshu login --cookie <your_cookie>")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, aiverr.ParseError("xiaohongshu", fmt.Sprintf("Failed to parse JSON response: %v", err))
	}

	return result, nil
}

// truncate shortens a string to maxLen characters.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// GetHotNotes returns hot/trending notes.
func (c *Client) GetHotNotes() (map[string]interface{}, error) {
	return c.get("/api/sns/web/v1/search/hot", nil)
}

// SearchNotes searches notes by keyword.
func (c *Client) SearchNotes(keyword string, page int) (map[string]interface{}, error) {
	return c.get("/api/sns/web/v1/search/notes", map[string]string{
		"keyword":   keyword,
		"page":      fmt.Sprintf("%d", page),
		"sort":      "general",
		"note_type": "0",
	})
}

// GetNoteDetail returns note detail by note ID.
func (c *Client) GetNoteDetail(noteID string) (map[string]interface{}, error) {
	return c.get("/api/sns/web/v1/feed", map[string]string{
		"source_note_id": noteID,
	})
}

// GetUserInfo returns user info by user ID.
func (c *Client) GetUserInfo(userID string) (map[string]interface{}, error) {
	return c.get("/api/sns/web/v1/user/otherinfo", map[string]string{
		"target_user_id": userID,
	})
}
