package douyin

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	douyinCommands "github.com/jackwener/aiview/commands/douyin"
	"github.com/spf13/cobra"
)

// DouyinPlatform implements the platform.Platform interface for Douyin.
type DouyinPlatform struct{}

// NewPlatform creates a new Douyin platform instance.
func NewPlatform() *DouyinPlatform {
	return &DouyinPlatform{}
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
	return NewClient(30, ""), nil
}

// Commands returns all Douyin commands.
func (p *DouyinPlatform) Commands() []*cobra.Command {
	douyinCmd := &cobra.Command{
		Use:   "douyin",
		Short: "Douyin platform commands",
		Long:  `Commands for interacting with Douyin (抖音) content.`,
	}

	douyinCmd.AddCommand(douyinCommands.NewHotCmd(func() douyinCommands.Client { return p.BuildClient() }))
	douyinCmd.AddCommand(douyinCommands.NewTrendingCmd(func() douyinCommands.Client { return p.BuildClient() }))
	douyinCmd.AddCommand(douyinCommands.NewSearchCmd(func() douyinCommands.Client { return p.BuildClient() }))
	douyinCmd.AddCommand(douyinCommands.NewVideoCmd(func() douyinCommands.Client { return p.BuildClient() }))
	douyinCmd.AddCommand(douyinCommands.NewUserCmd(func() douyinCommands.Client { return p.BuildClient() }))

	return []*cobra.Command{douyinCmd}
}

// BuildClient creates a client using the current credential.
func (p *DouyinPlatform) BuildClient() *Client {
	return NewClient(30, "")
}