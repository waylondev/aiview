package xiaohongshu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://edith.xiaohongshu.com"
)

// Client is the Xiaohongshu API client.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Xiaohongshu API client.
func NewClient(timeoutSec int) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
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
		return nil, fmt.Errorf("create request: %w", err)
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.xiaohongshu.com")
	req.Header.Set("Origin", "https://www.xiaohongshu.com")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return result, nil
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
