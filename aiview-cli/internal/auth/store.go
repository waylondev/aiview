package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Credential holds authentication data for a platform.
type Credential struct {
	Sessdata      string `json:"sessdata"`
	BiliJct       string `json:"bili_jct"`
	AcTimeValue   string `json:"ac_time_value"`
	Buvid3        string `json:"buvid3"`
	Buvid4        string `json:"buvid4"`
	DedeUserID    string `json:"dedeuserid"`
	SavedAt       int64  `json:"saved_at"`
}

// Store manages credential persistence.
type Store struct {
	Dir  string
	File string
}

// NewStore creates a new credential store.
func NewStore(platform string) *Store {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".aiview")
	return &Store{
		Dir:  dir,
		File: filepath.Join(dir, platform+"_credential.json"),
	}
}

// IsStale checks if the credential is older than TTL days.
func (c *Credential) IsStale(ttlDays int) bool {
	if c.SavedAt == 0 {
		return true
	}
	return time.Now().Unix()-c.SavedAt > int64(ttlDays*86400)
}

// IsValid checks if the credential has the minimum required fields.
func (c *Credential) IsValid() bool {
	return c.Sessdata != ""
}

// HasWriteCapability checks if the credential supports write operations.
func (c *Credential) HasWriteCapability() bool {
	return c.BiliJct != ""
}

// Save saves a credential to disk.
func (s *Store) Save(c *Credential) error {
	if err := os.MkdirAll(s.Dir, 0700); err != nil {
		return err
	}
	c.SavedAt = time.Now().Unix()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.File, data, 0600)
}

// Load loads a credential from disk.
func (s *Store) Load() (*Credential, error) {
	data, err := os.ReadFile(s.File)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var c Credential
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if !c.IsValid() {
		return nil, nil
	}
	return &c, nil
}

// Clear removes the saved credential file.
func (s *Store) Clear() error {
	if err := os.Remove(s.File); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// CookieStore provides cookie-based authentication storage for any platform.
type CookieStore struct {
	platform string
}

// credential holds simple cookie-based authentication data on disk.
type cookieCredential struct {
	Cookie string `json:"cookie"`
}

// NewCookieStore creates a new cookie store for the specified platform.
func NewCookieStore(platform string) *CookieStore {
	return &CookieStore{
		platform: platform,
	}
}

// SaveCookie saves the cookie to disk.
func (s *CookieStore) SaveCookie(cookie string) error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".aiview")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	file := filepath.Join(dir, s.platform+"_credential.json")
	cred := cookieCredential{Cookie: cookie}
	data, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0600)
}

// GetCookie loads the saved cookie from disk.
func (s *CookieStore) GetCookie() (string, error) {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".aiview")
	file := filepath.Join(dir, s.platform+"_credential.json")
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	var cred cookieCredential
	if err := json.Unmarshal(data, &cred); err != nil {
		return "", err
	}
	return cred.Cookie, nil
}

// ClearCookie removes the stored credential from disk.
func (s *CookieStore) ClearCookie() error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".aiview")
	file := filepath.Join(dir, s.platform+"_credential.json")
	if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}