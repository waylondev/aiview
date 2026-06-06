package bilibili

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/bilibili/commands"
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
	// Wire up QR login functions
	commands.SetQRLoginFuncs(
		func() (*commands.QRLoginSession, error) {
			sess, err := GenerateQRCode()
			if err != nil {
				return nil, err
			}
			return &commands.QRLoginSession{
				QRCodeKey: sess.QRCodeKey,
				QRCodeURL: sess.QRCodeURL,
			}, nil
		},
		func(key string) (int, *commands.Credential, error) {
			state, cred, err := PollQRCode(key)
			return int(state), cred, err
		},
	)

	getClient := func() commands.Client {
		return p.buildClient()
	}

	bilibiliCmd := &cobra.Command{
		Use:   "bilibili",
		Short: "Bilibili platform commands",
		Long:  `Bilibili platform commands, including video, search, user, favorites, interaction and more.`,
	}

	bilibiliCmd.AddCommand(commands.NewLoginCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewLogoutCmd(p.authStore))
	bilibiliCmd.AddCommand(commands.NewStatusCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewWhoamiCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewVideoCmd(getClient))
	bilibiliCmd.AddCommand(commands.NewSearchCmd(getClient))
	bilibiliCmd.AddCommand(commands.NewUserCmd(getClient))
	bilibiliCmd.AddCommand(commands.NewUserVideosCmd(getClient))
	bilibiliCmd.AddCommand(commands.NewFavoritesCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewFollowingCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewHistoryCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewWatchLaterCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewHotCmd(getClient))
	bilibiliCmd.AddCommand(commands.NewRankCmd(getClient))
	bilibiliCmd.AddCommand(commands.NewFeedCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewLikeCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewCoinCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewTripleCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewUnfollowCmd(p.authStore, getClient))
	bilibiliCmd.AddCommand(commands.NewAudioCmd(getClient))

	return []*cobra.Command{bilibiliCmd}
}

func init() {
	platform.Register(NewPlatform())
}