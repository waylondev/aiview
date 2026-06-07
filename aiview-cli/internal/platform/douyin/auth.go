package douyin

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AuthStore manages Douyin cookie persistence.
type AuthStore struct {
	dir  string
	file string
}

// credential holds Douyin authentication data on disk.
type credential struct {
	Cookie string `json:"cookie"`
}

// NewAuthStore creates a new auth store for Douyin.
func NewAuthStore() *AuthStore {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".aiview")
	return &AuthStore{
		dir:  dir,
		file: filepath.Join(dir, "douyin_credential.json"),
	}
}

// SaveCookie saves the cookie to disk.
func (a *AuthStore) SaveCookie(cookie string) error {
	if err := os.MkdirAll(a.dir, 0700); err != nil {
		return err
	}
	cred := credential{Cookie: cookie}
	data, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.file, data, 0600)
}

// GetCookie loads the saved cookie from disk.
func (a *AuthStore) GetCookie() string {
	data, err := os.ReadFile(a.file)
	if err != nil {
		return ""
	}
	var cred credential
	if err := json.Unmarshal(data, &cred); err != nil {
		return ""
	}
	return cred.Cookie
}

// ClearCookie removes the stored credential from disk.
func (a *AuthStore) ClearCookie() error {
	if err := os.Remove(a.file); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}