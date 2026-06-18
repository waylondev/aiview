// Package base provides shared base types for platform implementations.
package base

import (
	"github.com/jackwener/aiview/internal/auth"
)

// CookiePlatformBase provides common cookie-based authentication functionality.
type CookiePlatformBase struct {
	PlatformName string
	AuthStore    *auth.CookieStore
}

// NewCookiePlatformBase creates a new cookie platform base for the specified platform.
func NewCookiePlatformBase(platform string) *CookiePlatformBase {
	return &CookiePlatformBase{
		PlatformName: platform,
		AuthStore:    auth.NewCookieStore(platform),
	}
}

// GetCookie returns the stored cookie, or empty string on error.
func (b *CookiePlatformBase) GetCookie() string {
	cookie, _ := b.AuthStore.GetCookie()
	return cookie
}

// SaveCookie saves the cookie to local credential storage.
func (b *CookiePlatformBase) SaveCookie(cookie string) error {
	return b.AuthStore.SaveCookie(cookie)
}

// Logout clears the stored credential.
func (b *CookiePlatformBase) Logout() error {
	return b.AuthStore.ClearCookie()
}

// Status returns the current login status.
func (b *CookiePlatformBase) Status() (map[string]interface{}, error) {
	cookie := b.GetCookie()
	return map[string]interface{}{
		"platform":  b.PlatformName,
		"logged_in": cookie != "",
	}, nil
}
