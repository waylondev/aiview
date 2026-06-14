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