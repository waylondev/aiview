package douyin

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"status_code": 0,
			"status_msg":  "success",
			"data": map[string]interface{}{
				"word_list": []interface{}{
					map[string]interface{}{
						"word":      "测试热搜",
						"hot_value": 9999.0,
						"position":  1.0,
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
			BaseURL: server.URL,
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

func TestClient_HTTPError(t *testing.T) {
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
			BaseURL: server.URL,
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