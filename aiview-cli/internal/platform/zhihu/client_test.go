package zhihu

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackwener/aiview/internal/cache"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
	"github.com/jackwener/aiview/internal/ratelimit"
)

func TestNewClient(t *testing.T) {
	client := NewClient(30, "")
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Client == nil {
		t.Fatal("expected non-nil base Client")
	}
	if client.BaseURL != "https://www.zhihu.com" {
		t.Errorf("expected baseURL 'https://www.zhihu.com', got '%s'", client.BaseURL)
	}
	if client.Platform != "zhihu" {
		t.Errorf("expected platform 'zhihu', got '%s'", client.Platform)
	}
}

func TestPlatformName(t *testing.T) {
	client := NewClient(30, "")
	if name := client.PlatformName(); name != "zhihu" {
		t.Errorf("expected platform name 'zhihu', got '%s'", name)
	}
}

func TestInterfaceAssertions(t *testing.T) {
	var c *Client
	var _ platform.HotSearchable = c
	var _ platform.Searchable = c
	var _ platform.UserQueryable = c
}

func TestGetHotSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"top_search": map[string]interface{}{
				"words": []interface{}{
					map[string]interface{}{
						"query":       "知乎热搜测试",
						"display_name": "知乎热搜测试",
						"hot_value":    9999,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		Client: &base.Client{
			HTTPClient: &http.Client{
				Transport: &testTransport{serverURL: server.URL},
			},
			Limiter: ratelimit.New(100, 10),
			Cache:   cache.New(5 * time.Minute),
			BaseURL: "https://www.zhihu.com",
		},
	}

	result, err := client.GetHotSearch()
	if err != nil {
		t.Fatalf("GetHotSearch failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"data": []interface{}{
				map[string]interface{}{
					"type":  "search_result",
					"object": map[string]interface{}{
						"id":    12345,
						"title": "测试知乎内容",
					},
				},
			},
			"paging": map[string]interface{}{
				"is_end": false,
				"totals": 100,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		Client: &base.Client{
			HTTPClient: &http.Client{
				Transport: &testTransport{serverURL: server.URL},
			},
			Limiter: ratelimit.New(100, 10),
			Cache:   cache.New(5 * time.Minute),
			BaseURL: "https://www.zhihu.com",
		},
	}

	result, err := client.Search("测试", 1)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGetUserInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"id":      "zhihu_user_001",
			"name":    "测试知乎用户",
			"headline": "这是一个测试用户",
			"follower_count": 3000,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		Client: &base.Client{
			HTTPClient: &http.Client{
				Transport: &testTransport{serverURL: server.URL},
			},
			Limiter: ratelimit.New(100, 10),
			Cache:   cache.New(5 * time.Minute),
			BaseURL: "https://www.zhihu.com",
		},
	}

	result, err := client.GetUserInfo("zhihu_user_001")
	if err != nil {
		t.Fatalf("GetUserInfo failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
	}))
	defer server.Close()

	client := &Client{
		Client: &base.Client{
			HTTPClient: &http.Client{
				Transport: &testTransport{serverURL: server.URL},
			},
			Limiter: ratelimit.New(100, 10),
			Cache:   cache.New(5 * time.Minute),
			BaseURL: "https://www.zhihu.com",
		},
	}

	_, err := client.GetHotSearch()
	if err == nil {
		t.Error("expected error for HTTP 500, got nil")
	}
}

// testTransport rewrites requests to a test server.
type testTransport struct {
	serverURL string
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	newURL := t.serverURL + req.URL.Path
	if req.URL.RawQuery != "" {
		newURL += "?" + req.URL.RawQuery
	}
	newReq, err := http.NewRequest(req.Method, newURL, req.Body)
	if err != nil {
		return nil, err
	}
	newReq.Header = req.Header
	return http.DefaultTransport.RoundTrip(newReq)
}