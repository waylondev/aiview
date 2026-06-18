// Package douyin provides the Douyin (抖音) platform implementation.
package douyin

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
	douyinCommands "github.com/jackwener/aiview/commands/douyin"
	"github.com/spf13/cobra"
)

// DouyinPlatform implements the platform.Platform interface for Douyin.
type DouyinPlatform struct {
	*base.CookiePlatformBase
	config       *config.Config
	cachedClient *Client
}

// NewPlatform creates a new Douyin platform instance.
func NewPlatform() *DouyinPlatform {
	return &DouyinPlatform{
		CookiePlatformBase: base.NewCookiePlatformBase("douyin"),
	}
}

func init() {
	platform.Register(NewPlatform())
}

// Name returns the platform identifier.
func (p *DouyinPlatform) Name() string {
	return "douyin"
}

// NewClient creates a new Douyin API client.
func (p *DouyinPlatform) NewClient(cfg *config.Config) (platform.Client, error) {
	p.config = cfg
	return p.getClient(), nil
}

// Commands returns all Douyin commands.
func (p *DouyinPlatform) Commands() []*cobra.Command {
	douyinCmd := &cobra.Command{
		Use:   "douyin",
		Short: "Douyin platform commands",
		Long:  `Commands for interacting with Douyin (抖音) content.`,
	}

	c := p.getClient()
	isLoggedIn := func() bool {
		return p.GetCookie() != ""
	}
	douyinCmd.AddCommand(douyinCommands.NewLoginCmd(p.SaveCookie, func() douyinCommands.Client {
		return NewClient(30, p.GetCookie())
	}))
	douyinCmd.AddCommand(douyinCommands.NewHotCmd(c, isLoggedIn))
	douyinCmd.AddCommand(douyinCommands.NewTrendingCmd(c, isLoggedIn))
	douyinCmd.AddCommand(douyinCommands.NewSearchCmd(c, isLoggedIn))
	douyinCmd.AddCommand(douyinCommands.NewVideoCmd(c, isLoggedIn))
	douyinCmd.AddCommand(douyinCommands.NewUserCmd(c, isLoggedIn))
	douyinCmd.AddCommand(douyinCommands.NewCommentCmd(c, isLoggedIn))
	douyinCmd.AddCommand(douyinCommands.NewUserPostsCmd(c, isLoggedIn))
	douyinCmd.AddCommand(douyinCommands.NewStatusCmd(func() (map[string]interface{}, error) { return p.Status() }))
	douyinCmd.AddCommand(douyinCommands.NewLogoutCmd(func() error { return p.Logout() }))
	douyinCmd.AddCommand(douyinCommands.NewCollectCmd(func() douyinCommands.Client {
		return NewClient(30, p.GetCookie())
	}, isLoggedIn))

	return []*cobra.Command{douyinCmd}
}

// getClient lazily creates and caches a client using the current credential and config timeout.
func (p *DouyinPlatform) getClient() *Client {
	if p.cachedClient != nil {
		return p.cachedClient
	}
	timeout := config.DefaultTimeout
	if p.config != nil && p.config.Platforms["douyin"].Timeout > 0 {
		timeout = p.config.Platforms["douyin"].Timeout
	}
	p.cachedClient = NewClient(timeout, p.GetCookie())
	return p.cachedClient
}
