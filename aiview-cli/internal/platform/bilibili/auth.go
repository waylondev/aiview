package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/auth"
)

// AuthStore wraps the generic auth store for bilibili.
type AuthStore struct {
	store *auth.Store
}

// NewAuthStore creates a new bilibili auth store.
func NewAuthStore() *AuthStore {
	return &AuthStore{
		store: auth.NewStore("bilibili"),
	}
}

// Save saves the credential.
func (a *AuthStore) Save(c *Credential) error {
	ac := &auth.Credential{
		Sessdata:    c.Sessdata,
		BiliJct:     c.BiliJct,
		AcTimeValue: c.AcTimeValue,
		Buvid3:      c.Buvid3,
		Buvid4:      c.Buvid4,
		DedeUserID:  c.DedeUserID,
		SavedAt:     c.SavedAt,
	}
	return a.store.Save(ac)
}

// Load loads the credential.
func (a *AuthStore) Load() (*Credential, error) {
	c, err := a.store.Load()
	if err != nil || c == nil {
		return nil, err
	}
	return &Credential{
		Sessdata:    c.Sessdata,
		BiliJct:     c.BiliJct,
		AcTimeValue: c.AcTimeValue,
		Buvid3:      c.Buvid3,
		Buvid4:      c.Buvid4,
		DedeUserID:  c.DedeUserID,
		SavedAt:     c.SavedAt,
	}, nil
}

// Clear removes the credential.
func (a *AuthStore) Clear() error {
	return a.store.Clear()
}

// GetCredential tries to load the saved credential.
func (a *AuthStore) GetCredential() (*Credential, error) {
	cred, err := a.Load()
	if err != nil {
		return nil, err
	}
	if cred == nil {
		return nil, nil
	}
	if cred.IsStale(7) {
		return cred, fmt.Errorf("Credential expired, please log in again")
	}
	return cred, nil
}

// GetCredentialOrNil returns the credential or nil if not available.
func (a *AuthStore) GetCredentialOrNil() *Credential {
	cred, _ := a.GetCredential()
	return cred
}

// RequireCredential returns credential or an error if not logged in.
func (a *AuthStore) RequireCredential(requireWrite bool) (*Credential, error) {
	cred, err := a.GetCredential()
	if err != nil {
		return nil, &auth.NotAuthenticatedError{Platform: "bilibili"}
	}
	if cred == nil {
		return nil, &auth.NotAuthenticatedError{Platform: "bilibili"}
	}
	if requireWrite && !cred.HasWriteCapability() {
		return nil, &auth.WritePermissionError{Platform: "bilibili"}
	}
	return cred, nil
}