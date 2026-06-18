package bilibili

import (
	"encoding/json"

	"github.com/jackwener/aiview/internal/auth"
	aiverr "github.com/jackwener/aiview/internal/errors"
)

// AuthStore wraps the generic auth store for bilibili.
type AuthStore struct {
	store *auth.Store
}

// NewAuthStore creates a new bilibili auth store.
func NewAuthStore() *AuthStore {
	store, _ := auth.NewStore("bilibili")
	return &AuthStore{
		store: store,
	}
}

// Save saves the credential.
func (a *AuthStore) Save(c *Credential) error {
	// Serialize bilibili-specific fields to JSON and store in generic Cookie field
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	ac := &auth.Credential{
		Platform: "bilibili",
		Cookie:   string(data),
		SavedAt:  c.SavedAt,
	}
	return a.store.Save(ac)
}

// Load loads the credential.
func (a *AuthStore) Load() (*Credential, error) {
	c, err := a.store.Load()
	if err != nil || c == nil {
		return nil, err
	}
	// Deserialize bilibili-specific fields from generic Cookie field
	var bc Credential
	if err := json.Unmarshal([]byte(c.Cookie), &bc); err != nil {
		return nil, err
	}
	bc.SavedAt = c.SavedAt
	return &bc, nil
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
		return cred, aiverr.NotAuthenticated("bilibili", "Credential expired, please log in again")
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
