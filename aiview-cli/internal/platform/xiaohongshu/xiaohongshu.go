// Package xiaohongshu provides the Xiaohongshu (小红书/RED) platform implementation.
package xiaohongshu

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	xhsCmd "github.com/jackwener/aiview/commands/xiaohongshu"
	"github.com/spf13/cobra"
)

// XiaohongshuPlatform implements the platform.Platform interface for Xiaohongshu.
type XiaohongshuPlatform struct {
	config       *config.Config
	authStore    *AuthStore
	cachedClient *Client
}

// NewPlatform creates a new Xiaohongshu platform instance.
func NewPlatform() *XiaohongshuPlatform {
	return &XiaohongshuPlatform{
		authStore: NewAuthStore(),
	}
}

func init() {
	platform.Register(NewPlatform())
}

// Name returns the platform identifier.
func (p *XiaohongshuPlatform) Name() string {
	return "xiaohongshu"
}

// NewClient creates a new Xiaohongshu API client.
func (p *XiaohongshuPlatform) NewClient(cfg *config.Config) (platform.Client, error) {
	p.config = cfg
	return p.getClient(), nil
}

// Commands returns all Xiaohongshu commands.
func (p *XiaohongshuPlatform) Commands() []*cobra.Command {
	xhsCommand := &cobra.Command{
		Use:   "xiaohongshu",
		Short: "Xiaohongshu (小红书) platform commands",
		Long:  `Commands for interacting with Xiaohongshu (小红书/RED) content.`,
	}

	c := p.getClient()
	xhsCommand.AddCommand(xhsCmd.NewLoginCmd(p.SaveCookie, func() xhsCmd.Client {
		return NewClient(30, p.authStore.GetCookie())
	}))
	xhsCommand.AddCommand(xhsCmd.NewHotCmd(c))
	xhsCommand.AddCommand(xhsCmd.NewSearchCmd(c))
	xhsCommand.AddCommand(xhsCmd.NewNoteCmd(c))
	xhsCommand.AddCommand(xhsCmd.NewUserCmd(c))

	return []*cobra.Command{xhsCommand}
}

// getClient lazily creates and caches a client using the current credential and config timeout.
func (p *XiaohongshuPlatform) getClient() *Client {
	if p.cachedClient != nil {
		return p.cachedClient
	}
	timeout := 30
	if p.config != nil && p.config.Platforms.Douyin.Timeout > 0 {
		// Xiaohongshu doesn't have its own config yet, use default
		timeout = 30
	}
	p.cachedClient = NewClient(timeout, p.authStore.GetCookie())
	return p.cachedClient
}

// SaveCookie saves the Xiaohongshu cookie to local credential storage.
func (p *XiaohongshuPlatform) SaveCookie(cookie string) error {
	return p.authStore.SaveCookie(cookie)
}
