// Package kuaishou provides the Kuaishou (快手) platform implementation.
package kuaishou

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
	kuaishouCommands "github.com/jackwener/aiview/commands/kuaishou"
	"github.com/spf13/cobra"
)

// KuaishouPlatform implements the platform.Platform interface for Kuaishou.
type KuaishouPlatform struct {
	*base.CookiePlatformBase
	config       *config.Config
	cachedClient *Client
}

// NewPlatform creates a new Kuaishou platform instance.
func NewPlatform() *KuaishouPlatform {
	return &KuaishouPlatform{
		CookiePlatformBase: base.NewCookiePlatformBase("kuaishou"),
	}
}

func init() {
	platform.Register(NewPlatform())
}

// Name returns the platform identifier.
func (p *KuaishouPlatform) Name() string {
	return "kuaishou"
}

// NewClient creates a new Kuaishou API client.
func (p *KuaishouPlatform) NewClient(cfg *config.Config) (platform.Client, error) {
	p.config = cfg
	return p.getClient(), nil
}

// Commands returns all Kuaishou commands.
func (p *KuaishouPlatform) Commands() []*cobra.Command {
	kuaishouCmd := &cobra.Command{
		Use:   "kuaishou",
		Short: "Kuaishou (快手) platform commands",
		Long:  `Commands for interacting with Kuaishou (快手) content.`,
	}

	c := p.getClient()
	isLoggedIn := func() bool {
		return p.GetCookie() != ""
	}
	kuaishouCmd.AddCommand(kuaishouCommands.NewLoginCmd(p.SaveCookie, func() kuaishouCommands.Client {
		return NewClient(30, p.GetCookie())
	}))
	kuaishouCmd.AddCommand(kuaishouCommands.NewHotCmd(c, isLoggedIn))
	kuaishouCmd.AddCommand(kuaishouCommands.NewSearchCmd(c, isLoggedIn))
	kuaishouCmd.AddCommand(kuaishouCommands.NewUserCmd(c, isLoggedIn))

	return []*cobra.Command{kuaishouCmd}
}

// getClient lazily creates and caches a client using the current credential and config timeout.
func (p *KuaishouPlatform) getClient() *Client {
	if p.cachedClient != nil {
		return p.cachedClient
	}
	timeout := 30
	if p.config != nil && p.config.Platforms.Kuaishou.Timeout > 0 {
		timeout = p.config.Platforms.Kuaishou.Timeout
	}
	p.cachedClient = NewClient(timeout, p.GetCookie())
	return p.cachedClient
}
