package douyin

import (
	"testing"

	"github.com/jackwener/aiview/internal/auth"
)

func TestAuthStore_SaveAndGetCookie(t *testing.T) {
	store := auth.NewCookieStore("douyin")

	// Save
	if err := store.SaveCookie("test_cookie=abc"); err != nil {
		t.Fatalf("SaveCookie failed: %v", err)
	}

	// Get
	v, err := store.GetCookie()
	if err != nil {
		t.Fatalf("GetCookie failed: %v", err)
	}
	if v != "test_cookie=abc" {
		t.Errorf("expected 'test_cookie=abc', got '%s'", v)
	}

	// Cleanup
	store.ClearCookie()
}

func TestAuthStore_ClearCookie(t *testing.T) {
	store := auth.NewCookieStore("douyin")

	store.SaveCookie("test_cookie=abc")
	if err := store.ClearCookie(); err != nil {
		t.Fatalf("ClearCookie failed: %v", err)
	}

	v, _ := store.GetCookie()
	if v != "" {
		t.Errorf("expected empty after clear, got '%s'", v)
	}

	// Clear when file doesn't exist should not error
	if err := store.ClearCookie(); err != nil {
		t.Errorf("ClearCookie on non-existent file should not error: %v", err)
	}
}