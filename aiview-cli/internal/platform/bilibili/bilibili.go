package bilibili

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	biliCommands "github.com/jackwener/aiview/commands/bilibili"
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
	biliCmd := &cobra.Command{
		Use:   "bilibili",
		Short: "Bilibili platform commands",
		Long:  `Commands for interacting with Bilibili content.`,
	}

	// Auth commands
	biliCmd.AddCommand(biliCommands.NewLoginCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewLogoutCmd(p.GetAuthStore()))
	biliCmd.AddCommand(biliCommands.NewStatusCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewWhoamiCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))

	// Video commands
	biliCmd.AddCommand(biliCommands.NewVideoCmd(func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewAudioCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Discovery commands
	biliCmd.AddCommand(biliCommands.NewHotCmd(func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewRankCmd(func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewFeedCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))

	// Search
	biliCmd.AddCommand(biliCommands.NewSearchCmd(func() biliCommands.Client { return p.BuildClient() }))

	// User commands
	biliCmd.AddCommand(biliCommands.NewUserCmd(func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewUserVideosCmd(func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewDynamicCmd(func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewDynamicPostCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewDynamicDeleteCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewCollectionCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Collections
	biliCmd.AddCommand(biliCommands.NewFavoritesCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewFollowingCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewHistoryCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewWatchLaterCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))

	// Interactions
	biliCmd.AddCommand(biliCommands.NewLikeCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewCoinCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewTripleCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewUnfollowCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewFavoriteCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))

	// Comment
	biliCmd.AddCommand(biliCommands.NewCommentCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewCommentDeleteCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))

	// Danmaku
	biliCmd.AddCommand(biliCommands.NewDanmakuCmd(func() biliCommands.Client { return p.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewDanmakuSendCmd(p.GetAuthStore(), func() biliCommands.Client { return p.BuildClient() }))

	// Recommend
	biliCmd.AddCommand(biliCommands.NewRecommendCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Tags
	biliCmd.AddCommand(biliCommands.NewTagsCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Suggest
	biliCmd.AddCommand(biliCommands.NewSuggestCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Fans
	biliCmd.AddCommand(biliCommands.NewFansCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Video stat
	biliCmd.AddCommand(biliCommands.NewVideoStatCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Relation
	biliCmd.AddCommand(biliCommands.NewRelationCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Region
	biliCmd.AddCommand(biliCommands.NewRegionCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Live
	biliCmd.AddCommand(biliCommands.NewLiveCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Precious
	biliCmd.AddCommand(biliCommands.NewPreciousCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Trending
	biliCmd.AddCommand(biliCommands.NewTrendingCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Online
	biliCmd.AddCommand(biliCommands.NewOnlineCmd(func() biliCommands.Client { return p.BuildClient() }))

	// Weekly
	biliCmd.AddCommand(biliCommands.NewWeeklyCmd(func() biliCommands.Client { return p.BuildClient() }))

	return []*cobra.Command{biliCmd}
}

// GetAuthStore returns the platform's auth store.
func (p *BilibiliPlatform) GetAuthStore() *AuthStore {
	return p.authStore
}

// BuildClient creates a client using the current credential.
func (p *BilibiliPlatform) BuildClient() *Client {
	return p.buildClient()
}