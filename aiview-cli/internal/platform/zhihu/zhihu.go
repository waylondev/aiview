// Package zhihu provides the Zhihu (知乎) platform implementation.
package zhihu

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	zhihuCommands "github.com/jackwener/aiview/commands/zhihu"
	"github.com/spf13/cobra"
)

// ZhihuPlatform implements the platform.Platform interface for Zhihu.
type ZhihuPlatform struct {
	authStore    *AuthStore
	config       *config.Config
	cachedClient *Client
}

// NewPlatform creates a new Zhihu platform instance.
func NewPlatform() *ZhihuPlatform {
	return &ZhihuPlatform{
		authStore: NewAuthStore(),
	}
}

func init() {
	platform.Register(NewPlatform())
}

// Name returns the platform identifier.
func (p *ZhihuPlatform) Name() string {
	return "zhihu"
}

// NewClient creates a new Zhihu API client.
func (p *ZhihuPlatform) NewClient(cfg *config.Config) (platform.Client, error) {
	p.config = cfg
	return p.getClient(), nil
}

// Commands returns all Zhihu commands.
func (p *ZhihuPlatform) Commands() []*cobra.Command {
	zhihuCmd := &cobra.Command{
		Use:   "zhihu",
		Short: "Zhihu (知乎) platform commands",
		Long:  `Commands for interacting with Zhihu (知乎) content.`,
	}

	c := p.getClient()
	isLoggedIn := func() bool {
		return p.authStore.GetCookie() != ""
	}
	zhihuCmd.AddCommand(zhihuCommands.NewLoginCmd(p.SaveCookie, func() zhihuCommands.Client {
		return NewClient(30, p.authStore.GetCookie())
	}))
	zhihuCmd.AddCommand(zhihuCommands.NewHotCmd(c, isLoggedIn))
	zhihuCmd.AddCommand(zhihuCommands.NewSearchCmd(c, isLoggedIn))
	zhihuCmd.AddCommand(zhihuCommands.NewUserCmd(c, isLoggedIn))

	return []*cobra.Command{zhihuCmd}
}

// getClient lazily creates and caches a client using the current credential and config timeout.
func (p *ZhihuPlatform) getClient() *Client {
	if p.cachedClient != nil {
		return p.cachedClient
	}
	timeout := 30
	if p.config != nil && p.config.Platforms.Douyin.Timeout > 0 {
		timeout = 30
	}
	p.cachedClient = NewClient(timeout, p.authStore.GetCookie())
	return p.cachedClient
}

// SaveCookie saves the Zhihu cookie to local credential storage.
func (p *ZhihuPlatform) SaveCookie(cookie string) error {
	return p.authStore.SaveCookie(cookie)
}
