package douyin

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	douyinCommands "github.com/jackwener/aiview/commands/douyin"
	"github.com/spf13/cobra"
)

// DouyinPlatform implements the platform.Platform interface for Douyin.
type DouyinPlatform struct {
	authStore    *AuthStore
	config       *config.Config
	cachedClient *Client
}

// NewPlatform creates a new Douyin platform instance.
func NewPlatform() *DouyinPlatform {
	return &DouyinPlatform{
		authStore: NewAuthStore(),
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
	douyinCmd.AddCommand(douyinCommands.NewLoginCmd(p.SaveCookie))
	douyinCmd.AddCommand(douyinCommands.NewHotCmd(c))
	douyinCmd.AddCommand(douyinCommands.NewTrendingCmd(c))
	douyinCmd.AddCommand(douyinCommands.NewSearchCmd(c))
	douyinCmd.AddCommand(douyinCommands.NewVideoCmd(c))
	douyinCmd.AddCommand(douyinCommands.NewUserCmd(c))
	douyinCmd.AddCommand(douyinCommands.NewCommentCmd(c))
	douyinCmd.AddCommand(douyinCommands.NewUserPostsCmd(c))
	douyinCmd.AddCommand(douyinCommands.NewStatusCmd(func() (map[string]interface{}, error) { return p.Status() }))
	douyinCmd.AddCommand(douyinCommands.NewLogoutCmd(func() error { return p.Logout() }))

	return []*cobra.Command{douyinCmd}
}

// getClient lazily creates and caches a client using the current credential and config timeout.
func (p *DouyinPlatform) getClient() *Client {
	if p.cachedClient != nil {
		return p.cachedClient
	}
	timeout := 30
	if p.config != nil && p.config.Platforms.Douyin.Timeout > 0 {
		timeout = p.config.Platforms.Douyin.Timeout
	}
	p.cachedClient = NewClient(timeout, p.authStore.GetCookie())
	return p.cachedClient
}

// SaveCookie saves the Douyin cookie to local credential storage.
func (p *DouyinPlatform) SaveCookie(cookie string) error {
	return p.authStore.SaveCookie(cookie)
}

// Status returns the current login status.
func (p *DouyinPlatform) Status() (map[string]interface{}, error) {
	cookie := p.authStore.GetCookie()
	loggedIn := cookie != ""
	return map[string]interface{}{
		"platform":  "douyin",
		"logged_in": loggedIn,
	}, nil
}

// Logout clears the stored credential.
func (p *DouyinPlatform) Logout() error {
	return p.authStore.ClearCookie()
}