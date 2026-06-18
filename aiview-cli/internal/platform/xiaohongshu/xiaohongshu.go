// Package xiaohongshu provides the Xiaohongshu (小红书/RED) platform implementation.
package xiaohongshu

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
	xhsCmd "github.com/jackwener/aiview/commands/xiaohongshu"
	"github.com/spf13/cobra"
)

// XiaohongshuPlatform implements the platform.Platform interface for Xiaohongshu.
type XiaohongshuPlatform struct {
	*base.CookiePlatformBase
	config       *config.Config
	cachedClient *Client
}

// NewPlatform creates a new Xiaohongshu platform instance.
func NewPlatform() *XiaohongshuPlatform {
	return &XiaohongshuPlatform{
		CookiePlatformBase: base.NewCookiePlatformBase("xiaohongshu"),
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
	isLoggedIn := func() bool {
		return p.GetCookie() != ""
	}
	xhsCommand.AddCommand(xhsCmd.NewLoginCmd(p.SaveCookie, func() xhsCmd.Client {
		return NewClient(30, p.GetCookie())
	}))
	xhsCommand.AddCommand(xhsCmd.NewHotCmd(c, isLoggedIn))
	xhsCommand.AddCommand(xhsCmd.NewSearchCmd(c, isLoggedIn))
	xhsCommand.AddCommand(xhsCmd.NewNoteCmd(c, isLoggedIn))
	xhsCommand.AddCommand(xhsCmd.NewUserCmd(c, isLoggedIn))

	return []*cobra.Command{xhsCommand}
}

// getClient lazily creates and caches a client using the current credential and config timeout.
func (p *XiaohongshuPlatform) getClient() *Client {
	if p.cachedClient != nil {
		return p.cachedClient
	}
	timeout := config.DefaultTimeout
	if p.config != nil && p.config.Platforms["xiaohongshu"].Timeout > 0 {
		timeout = p.config.Platforms["xiaohongshu"].Timeout
	}
	p.cachedClient = NewClient(timeout, p.GetCookie())
	return p.cachedClient
}
