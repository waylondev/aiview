package douyin

import (
	"path/filepath"
	"testing"
)

func TestAuthStore_SaveAndGetCookie(t *testing.T) {
	dir := t.TempDir()
	store := &AuthStore{
		dir:  dir,
		file: filepath.Join(dir, "douyin_credential.json"),
	}

	// Initially empty
	if v := store.GetCookie(); v != "" {
		t.Errorf("expected empty, got '%s'", v)
	}

	// Save
	if err := store.SaveCookie("test_cookie=abc"); err != nil {
		t.Fatalf("SaveCookie failed: %v", err)
	}

	// Get
	if v := store.GetCookie(); v != "test_cookie=abc" {
		t.Errorf("expected 'test_cookie=abc', got '%s'", v)
	}
}

func TestAuthStore_ClearCookie(t *testing.T) {
	dir := t.TempDir()
	store := &AuthStore{
		dir:  dir,
		file: filepath.Join(dir, "douyin_credential.json"),
	}

	store.SaveCookie("test_cookie=abc")
	if err := store.ClearCookie(); err != nil {
		t.Fatalf("ClearCookie failed: %v", err)
	}

	if v := store.GetCookie(); v != "" {
		t.Errorf("expected empty after clear, got '%s'", v)
	}

	// Clear when file doesn't exist should not error
	if err := store.ClearCookie(); err != nil {
		t.Errorf("ClearCookie on non-existent file should not error: %v", err)
	}
}

func TestAuthStore_ClearNonExistent(t *testing.T) {
	dir := t.TempDir()
	store := &AuthStore{
		dir:  dir,
		file: filepath.Join(dir, "nonexistent.json"),
	}

	if err := store.ClearCookie(); err != nil {
		t.Errorf("ClearCookie on non-existent file should not error: %v", err)
	}
}