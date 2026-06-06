package bilibili

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/spf13/cobra"
)

// BilibiliPlatform implements the platform.Platform interface for Bilibili.
type BilibiliPlatform struct {
	authStore *AuthStore
}

// NewPlatform creates a new Bilibili platform instance.
func NewPlatform() *BilibiliPlatform {
	return &BilibiliPlatform{
		authStore: NewAuthStore(),
	}
}

func init() {
	platform.Register(NewPlatform())
}

// Name returns the platform identifier.
func (p *BilibiliPlatform) Name() string {
	return "bilibili"
}

// NewClient creates a new Bilibili API client.
func (p *BilibiliPlatform) NewClient(cfg *config.Config) (platform.Client, error) {
	return p.buildClient(), nil
}

// buildClient creates a client using the current credential.
func (p *BilibiliPlatform) buildClient() *Client {
	cred := p.authStore.GetCredentialOrNil()
	return NewClient(30, BuildCookieString(cred), cred)
}

// Commands returns all Bilibili commands.
func (p *BilibiliPlatform) Commands() []*cobra.Command {
	// Commands are registered by the CLI layer to avoid import cycles.
	return nil
}

// GetAuthStore returns the platform's auth store.
func (p *BilibiliPlatform) GetAuthStore() *AuthStore {
	return p.authStore
}

// BuildClient creates a client using the current credential.
func (p *BilibiliPlatform) BuildClient() *Client {
	return p.buildClient()
}