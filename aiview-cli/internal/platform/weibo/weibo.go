// Package weibo provides the Weibo (微博) platform implementation.
package weibo

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
	weiboCommands "github.com/jackwener/aiview/commands/weibo"
	"github.com/spf13/cobra"
)

// WeiboPlatform implements the platform.Platform interface for Weibo.
type WeiboPlatform struct {
	*base.CookiePlatformBase
	config       *config.Config
	cachedClient *Client
}

// NewPlatform creates a new Weibo platform instance.
func NewPlatform() *WeiboPlatform {
	return &WeiboPlatform{
		CookiePlatformBase: base.NewCookiePlatformBase("weibo"),
	}
}

func init() {
	platform.Register(NewPlatform())
}

// Name returns the platform identifier.
func (p *WeiboPlatform) Name() string {
	return "weibo"
}

// NewClient creates a new Weibo API client.
func (p *WeiboPlatform) NewClient(cfg *config.Config) (platform.Client, error) {
	p.config = cfg
	return p.getClient(), nil
}

// Commands returns all Weibo commands.
func (p *WeiboPlatform) Commands() []*cobra.Command {
	weiboCmd := &cobra.Command{
		Use:   "weibo",
		Short: "Weibo (微博) platform commands",
		Long:  `Commands for interacting with Weibo (微博) content.`,
	}

	c := p.getClient()
	isLoggedIn := func() bool {
		return p.GetCookie() != ""
	}
	weiboCmd.AddCommand(weiboCommands.NewLoginCmd(p.SaveCookie, func() weiboCommands.Client {
		return NewClient(30, p.GetCookie())
	}))
	weiboCmd.AddCommand(weiboCommands.NewHotCmd(c, isLoggedIn))
	weiboCmd.AddCommand(weiboCommands.NewSearchCmd(c, isLoggedIn))
	weiboCmd.AddCommand(weiboCommands.NewUserCmd(c, isLoggedIn))

	return []*cobra.Command{weiboCmd}
}

// getClient lazily creates and caches a client using the current credential and config timeout.
func (p *WeiboPlatform) getClient() *Client {
	if p.cachedClient != nil {
		return p.cachedClient
	}
	timeout := 30
	if p.config != nil && p.config.Platforms.Weibo.Timeout > 0 {
		timeout = p.config.Platforms.Weibo.Timeout
	}
	p.cachedClient = NewClient(timeout, p.GetCookie())
	return p.cachedClient
}
