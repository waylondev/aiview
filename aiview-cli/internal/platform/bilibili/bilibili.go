// Package bilibili provides the Bilibili platform implementation.
package bilibili

import (
	"encoding/json"
	"fmt"

	"github.com/jackwener/aiview/internal/auth"
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
	biliCommands "github.com/jackwener/aiview/commands/bilibili"
	"github.com/spf13/cobra"
)

// BilibiliPlatform implements the platform.Platform interface for Bilibili.
type BilibiliPlatform struct {
	*base.CookiePlatformBase
	credStore    *auth.Store
	cachedClient *Client
}

// NewPlatform creates a new Bilibili platform instance.
func NewPlatform() *BilibiliPlatform {
	credStore, _ := auth.NewStore("bilibili")
	return &BilibiliPlatform{
		CookiePlatformBase: base.NewCookiePlatformBase("bilibili"),
		credStore:          credStore,
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
	if p.cachedClient != nil {
		return p.cachedClient
	}
	cred := p.getCredentialOrNil()
	p.cachedClient = NewClient(config.DefaultTimeout, BuildCookieString(cred), cred)
	return p.cachedClient
}

// getCredentialOrNil returns the credential or nil if not available.
func (p *BilibiliPlatform) getCredentialOrNil() *Credential {
	cred, _ := p.credStore.Load()
	if cred == nil {
		return nil
	}
	return authCredToBili(cred)
}

// authCredToBili converts auth.Credential to bilibili Credential.
func authCredToBili(c *auth.Credential) *Credential {
	if c == nil {
		return nil
	}
	var bc Credential
	if err := json.Unmarshal([]byte(c.Cookie), &bc); err != nil {
		return nil
	}
	bc.SavedAt = c.SavedAt
	return &bc
}

// biliCredToAuth converts bilibili Credential to auth.Credential.
func biliCredToAuth(c *Credential) *auth.Credential {
	if c == nil {
		return nil
	}
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return &auth.Credential{
		Platform: "bilibili",
		Cookie:   string(data),
		SavedAt:  c.SavedAt,
	}
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

	// Collect
	biliCmd.AddCommand(biliCommands.NewCollectCmd(func() biliCommands.Client { return p.BuildClient() }))

	return []*cobra.Command{biliCmd}
}

// GetAuthStore returns an AuthProvider adapter for the platform.
func (p *BilibiliPlatform) GetAuthStore() biliCommands.AuthProvider {
	return &authAdapter{store: p.credStore}
}

// authAdapter adapts auth.Store to the AuthProvider interface.
type authAdapter struct {
	store *auth.Store
}

func (a *authAdapter) GetCredential() (*Credential, error) {
	cred, err := a.store.Load()
	if err != nil {
		return nil, err
	}
	return authCredToBili(cred), nil
}

func (a *authAdapter) GetCredentialOrNil() *Credential {
	cred, _ := a.store.Load()
	return authCredToBili(cred)
}

func (a *authAdapter) RequireCredential(requireWrite bool) (*Credential, error) {
	cred, err := a.store.Load()
	if err != nil {
		return nil, err
	}
	if cred == nil {
		return nil, fmt.Errorf("not logged in")
	}
	bc := authCredToBili(cred)
	if requireWrite && !bc.HasWriteCapability() {
		return nil, fmt.Errorf("write permission required")
	}
	return bc, nil
}

func (a *authAdapter) Save(c *Credential) error {
	return a.store.Save(biliCredToAuth(c))
}

func (a *authAdapter) Clear() error {
	return a.store.Clear()
}

// BuildClient creates a client using the current credential.
func (p *BilibiliPlatform) BuildClient() *Client {
	return p.buildClient()
}