package base

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	aiverr "github.com/jackwener/aiview/internal/errors"
)

func TestNewBaseClient(t *testing.T) {
	c := NewClient(30, "cookies", "ua", "https://example.com", "https://referer.com", "test-platform")

	if c == nil {
		t.Fatal("expected non-nil Client")
	}
	if c.UserAgent != "ua" {
		t.Errorf("expected UserAgent %q, got %q", "ua", c.UserAgent)
	}
	if c.BaseURL != "https://example.com" {
		t.Errorf("expected BaseURL %q, got %q", "https://example.com", c.BaseURL)
	}
	if c.Referer != "https://referer.com" {
		t.Errorf("expected Referer %q, got %q", "https://referer.com", c.Referer)
	}
	if c.Platform != "test-platform" {
		t.Errorf("expected Platform %q, got %q", "test-platform", c.Platform)
	}
	if c.Cookies != "cookies" {
		t.Errorf("expected Cookies %q, got %q", "cookies", c.Cookies)
	}
	if c.HTTPClient == nil {
		t.Error("expected non-nil HTTPClient")
	}
	if c.Limiter == nil {
		t.Error("expected non-nil Limiter")
	}
	if c.Cache == nil {
		t.Error("expected non-nil Cache")
	}
}

func TestBuildHeaders(t *testing.T) {
	c := NewClient(30, "token=abc123", "MyApp/1.0", "https://example.com", "https://referer.com", "test")

	headers := c.BuildHeaders()

	if headers.Get("User-Agent") != "MyApp/1.0" {
		t.Errorf("expected User-Agent %q, got %q", "MyApp/1.0", headers.Get("User-Agent"))
	}
	if headers.Get("Accept") == "" {
		t.Error("expected non-empty Accept header")
	}
	if headers.Get("Referer") != "https://referer.com" {
		t.Errorf("expected Referer %q, got %q", "https://referer.com", headers.Get("Referer"))
	}
	if headers.Get("Cookie") != "token=abc123" {
		t.Errorf("expected Cookie %q, got %q", "token=abc123", headers.Get("Cookie"))
	}
}

func TestBuildHeaders_NoRefererNoCookie(t *testing.T) {
	c := NewClient(30, "", "App/1.0", "https://example.com", "", "test")

	headers := c.BuildHeaders()

	if headers.Get("Referer") != "" {
		t.Error("expected no Referer header")
	}
	if headers.Get("Cookie") != "" {
		t.Error("expected no Cookie header")
	}
}

func TestParseResponse(t *testing.T) {
	t.Run("valid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "count": 42})
		}))
		defer server.Close()

		c := NewClient(5, "", "test", server.URL, "", "test-platform")
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		result, err := c.ParseResponse(resp)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["status"] != "ok" {
			t.Errorf("expected status=ok, got %v", result["status"])
		}
		if result["count"] != float64(42) {
			t.Errorf("expected count=42, got %v", result["count"])
		}
	})

	t.Run("401 response returns NotAuthenticated", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}))
		defer server.Close()

		c := NewClient(5, "", "test", server.URL, "", "test-platform")
		resp, _ := http.Get(server.URL)
		defer resp.Body.Close()

		_, err := c.ParseResponse(resp)
		if err == nil {
			t.Fatal("expected error for 401")
		}
		// Should be a PlatformError with not_authenticated code
		pe, ok := err.(*aiverr.PlatformError)
		if !ok {
			t.Fatalf("expected *aiverr.PlatformError, got %T", err)
		}
		if pe.Code != aiverr.CodeNotAuthenticated {
			t.Errorf("expected code %q, got %q", aiverr.CodeNotAuthenticated, pe.Code)
		}
	})

	t.Run("429 response returns RateLimited", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too Many Requests"))
		}))
		defer server.Close()

		c := NewClient(5, "", "test", server.URL, "", "test-platform")
		resp, _ := http.Get(server.URL)
		defer resp.Body.Close()

		_, err := c.ParseResponse(resp)
		if err == nil {
			t.Fatal("expected error for 429")
		}
		pe, ok := err.(*aiverr.PlatformError)
		if !ok {
			t.Fatalf("expected *aiverr.PlatformError, got %T", err)
		}
		if pe.Code != aiverr.CodeRateLimited {
			t.Errorf("expected code %q, got %q", aiverr.CodeRateLimited, pe.Code)
		}
	})

	t.Run("500 response returns APIError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		c := NewClient(5, "", "test", server.URL, "", "test-platform")
		resp, _ := http.Get(server.URL)
		defer resp.Body.Close()

		_, err := c.ParseResponse(resp)
		if err == nil {
			t.Fatal("expected error for 500")
		}
		pe, ok := err.(*aiverr.PlatformError)
		if !ok {
			t.Fatalf("expected *aiverr.PlatformError, got %T", err)
		}
		if pe.Code != aiverr.CodeAPIError {
			t.Errorf("expected code %q, got %q", aiverr.CodeAPIError, pe.Code)
		}
	})

	t.Run("invalid JSON returns ParseError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not-json{{{bad{"))
		}))
		defer server.Close()

		c := NewClient(5, "", "test", server.URL, "", "test-platform")
		resp, _ := http.Get(server.URL)
		defer resp.Body.Close()

		_, err := c.ParseResponse(resp)
		if err == nil {
			t.Fatal("expected parse error for invalid JSON")
		}
		pe, ok := err.(*aiverr.PlatformError)
		if !ok {
			t.Fatalf("expected *aiverr.PlatformError, got %T", err)
		}
		if pe.Code != aiverr.CodeParseError {
			t.Errorf("expected code %q, got %q", aiverr.CodeParseError, pe.Code)
		}
	})
}

func TestDetectHTML(t *testing.T) {
	c := NewClient(5, "", "test", "https://example.com", "", "test-platform")

	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{"DOCTYPE uppercase", "<!DOCTYPE html><html></html>", true},
		{"DOCTYPE lowercase", "<!doctype html><html></html>", true},
		{"html tag", "<html><body>test</body></html>", true},
		{"html with spaces", "  <html lang=\"en\">", true},
		{"JSON body", `{"key":"value"}`, false},
		{"plain text", "just some text", false},
		{"empty", "", false},
		{"XML", "<?xml version=\"1.0\"?>", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.DetectHTMLResponse([]byte(tt.body))
			if result != tt.expected {
				t.Errorf("expected %v for body %q, got %v", tt.expected, tt.body, result)
			}
		})
	}
}