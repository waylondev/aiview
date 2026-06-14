package bilibili

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackwener/aiview/internal/cache"
	"github.com/jackwener/aiview/internal/platform/base"
	"github.com/jackwener/aiview/internal/ratelimit"
)

func TestClient_GetHotSearch(t *testing.T) {
	// Mock server returning hot search JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code":    0,
			"message": "0",
			"data": map[string]interface{}{
				"trending": map[string]interface{}{
					"list": []interface{}{
						map[string]interface{}{
							"keyword":   "test hot",
							"show_name": "测试热搜",
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client with custom transport that rewrites URL to test server
	client := &Client{
		Client: &base.Client{
			HTTPClient: &http.Client{
				Transport: &testTransport{serverURL: server.URL},
			},
			Limiter: ratelimit.New(100, 10),
			Cache:   cache.New(5 * time.Minute),
		},
	}

	result, err := client.GetHotSearch(10)
	if err != nil {
		t.Fatalf("GetHotSearch failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestClient_GetVideoInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code":    0,
			"message": "0",
			"data": map[string]interface{}{
				"aid":      12345.0,
				"cid":      67890.0,
				"title":    "Test Video",
				"desc":     "Test Description",
				"duration": 120.0,
				"owner": map[string]interface{}{
					"mid":  100.0,
					"name": "testuser",
					"face": "https://example.com/face.jpg",
				},
				"stat": map[string]interface{}{
					"view":     1000.0,
					"danmaku":  50.0,
					"reply":    30.0,
					"favorite": 20.0,
					"coin":     10.0,
					"share":    5.0,
					"like":     40.0,
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
		},
	}

	info, err := client.GetVideoInfo("BV1xx411c7m9")
	if err != nil {
		t.Fatalf("GetVideoInfo failed: %v", err)
	}
	if info.Title != "Test Video" {
		t.Errorf("expected 'Test Video', got '%s'", info.Title)
	}
	if info.AID != 12345 {
		t.Errorf("expected 12345, got %d", info.AID)
	}
}

func TestClient_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := &Client{
		Client: &base.Client{
			HTTPClient: &http.Client{
				Transport: &testTransport{serverURL: server.URL},
			},
			Limiter: ratelimit.New(100, 10),
			Cache:   cache.New(5 * time.Minute),
		},
	}

	_, err := client.GetHotSearch(10)
	if err == nil {
		t.Error("expected error for HTTP 500, got nil")
	}
}

func TestGetVideoDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code":    0,
			"message": "0",
			"data": map[string]interface{}{
				"aid":      99999.0,
				"cid":      88888.0,
				"title":    "Detailed Video Test",
				"desc":     "Long description for testing",
				"duration": 300.0,
				"owner": map[string]interface{}{
					"mid":  200.0,
					"name": "owner_user",
					"face": "https://example.com/avatar.jpg",
				},
				"stat": map[string]interface{}{
					"view":     5000.0,
					"danmaku":  200.0,
					"reply":    100.0,
					"favorite": 80.0,
					"coin":     30.0,
					"share":    15.0,
					"like":     150.0,
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
		},
	}

	info, err := client.GetVideoInfo("BV1xx411c7mD")
	if err != nil {
		t.Fatalf("GetVideoInfo failed: %v", err)
	}
	if info.Title != "Detailed Video Test" {
		t.Errorf("expected 'Detailed Video Test', got '%s'", info.Title)
	}
	if info.Duration != 300 {
		t.Errorf("expected 300, got %d", info.Duration)
	}
	if info.Stats.View != 5000 {
		t.Errorf("expected 5000 views, got %d", info.Stats.View)
	}
	if info.Owner.Name != "owner_user" {
		t.Errorf("expected 'owner_user', got '%s'", info.Owner.Name)
	}
}

func TestGetHotVideos(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code":    0,
			"message": "0",
			"data": map[string]interface{}{
				"list": []interface{}{
					map[string]interface{}{
						"bvid":     "BV1xx00000001",
						"title":    "Hot Video 1",
						"duration": 180,
						"owner": map[string]interface{}{
							"mid":  100,
							"name": "creator1",
						},
						"stat": map[string]interface{}{
							"view": 100000,
						},
					},
					map[string]interface{}{
						"bvid":     "BV1xx00000002",
						"title":    "Hot Video 2",
						"duration": 240,
						"owner": map[string]interface{}{
							"mid":  200,
							"name": "creator2",
						},
						"stat": map[string]interface{}{
							"view": 80000,
						},
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
		},
	}

	videos, err := client.GetHotVideos(1, 10)
	if err != nil {
		t.Fatalf("GetHotVideos failed: %v", err)
	}
	if len(videos) != 2 {
		t.Fatalf("expected 2 videos, got %d", len(videos))
	}
	if videos[0].BVID != "BV1xx00000001" {
		t.Errorf("expected 'BV1xx00000001', got '%s'", videos[0].BVID)
	}
	if videos[0].Title != "Hot Video 1" {
		t.Errorf("expected 'Hot Video 1', got '%s'", videos[0].Title)
	}
	if videos[0].Owner.Name != "creator1" {
		t.Errorf("expected 'creator1', got '%s'", videos[0].Owner.Name)
	}
	if videos[0].Stats.View != 100000 {
		t.Errorf("expected 100000 views, got %d", videos[0].Stats.View)
	}
}

func TestSearchVideo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code":    0,
			"message": "0",
			"data": map[string]interface{}{
				"result": []interface{}{
					map[string]interface{}{
						"bvid":     "BV1search001",
						"title":    "搜索结果1",
						"author":   "up主A",
						"play":     12345,
						"duration": "05:30",
					},
					map[string]interface{}{
						"bvid":     "BV1search002",
						"title":    "搜索结果2",
						"author":   "up主B",
						"play":     6789,
						"duration": "12:15",
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
		},
	}

	results, err := client.SearchVideo("测试", 1, "", 0, 0)
	if err != nil {
		t.Fatalf("SearchVideo failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].BVID != "BV1search001" {
		t.Errorf("expected 'BV1search001', got '%s'", results[0].BVID)
	}
	if results[0].Title != "搜索结果1" {
		t.Errorf("expected '搜索结果1', got '%s'", results[0].Title)
	}
	if results[0].Author != "up主A" {
		t.Errorf("expected 'up主A', got '%s'", results[0].Author)
	}
}

func TestGetUserVideos(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code":    0,
			"message": "0",
			"data": map[string]interface{}{
				"list": map[string]interface{}{
					"vlist": []interface{}{
						map[string]interface{}{
							"bvid":   "BV1user001",
							"title":  "用户视频1",
							"length": 150,
							"play":   5000,
						},
						map[string]interface{}{
							"bvid":   "BV1user002",
							"title":  "用户视频2",
							"length": 320,
							"play":   3000,
						},
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
		},
	}

	videos, err := client.GetUserVideos(12345, 10, "pubdate", 0, "")
	if err != nil {
		t.Fatalf("GetUserVideos failed: %v", err)
	}
	if len(videos) != 2 {
		t.Fatalf("expected 2 videos, got %d", len(videos))
	}
	if videos[0].BVID != "BV1user001" {
		t.Errorf("expected 'BV1user001', got '%s'", videos[0].BVID)
	}
	if videos[0].Title != "用户视频1" {
		t.Errorf("expected '用户视频1', got '%s'", videos[0].Title)
	}
	if videos[0].Stats.View != 5000 {
		t.Errorf("expected 5000 views, got %d", videos[0].Stats.View)
	}
}

// testTransport rewrites requests to a test server.
type testTransport struct {
	serverURL string
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the URL to point to the test server, preserving the path
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