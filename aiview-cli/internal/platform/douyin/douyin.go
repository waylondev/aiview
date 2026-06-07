package douyin

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
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

// Commands returns all Douyin commands (populated by the CLI layer).
func (p *DouyinPlatform) Commands() []*cobra.Command {
	// Will be populated by the CLI layer after import initialization
	return nil
}