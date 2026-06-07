package douyin

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	douyinCommands "github.com/jackwener/aiview/commands/douyin"
	"github.com/spf13/cobra"
)

// DouyinPlatform implements the platform.Platform interface for Douyin.
type DouyinPlatform struct {
	authStore *AuthStore
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
	return p.BuildClient(), nil
}

// Commands returns all Douyin commands.
func (p *DouyinPlatform) Commands() []*cobra.Command {
	douyinCmd := &cobra.Command{
		Use:   "douyin",
		Short: "Douyin platform commands",
		Long:  `Commands for interacting with Douyin (抖音) content.`,
	}

	douyinCmd.AddCommand(douyinCommands.NewLoginCmd(p.SaveCookie))
	douyinCmd.AddCommand(douyinCommands.NewHotCmd(func() douyinCommands.Client { return p }))
	douyinCmd.AddCommand(douyinCommands.NewTrendingCmd(func() douyinCommands.Client { return p }))
	douyinCmd.AddCommand(douyinCommands.NewSearchCmd(func() douyinCommands.Client { return p }))
	douyinCmd.AddCommand(douyinCommands.NewVideoCmd(func() douyinCommands.Client { return p }))
	douyinCmd.AddCommand(douyinCommands.NewUserCmd(func() douyinCommands.Client { return p }))
	douyinCmd.AddCommand(douyinCommands.NewCommentCmd(func() douyinCommands.Client { return p }))
	douyinCmd.AddCommand(douyinCommands.NewUserPostsCmd(func() douyinCommands.Client { return p }))
	douyinCmd.AddCommand(douyinCommands.NewStatusCmd(func() douyinCommands.Client { return p }))
	douyinCmd.AddCommand(douyinCommands.NewLogoutCmd(func() douyinCommands.Client { return p }))

	return []*cobra.Command{douyinCmd}
}

// BuildClient creates a client using the current credential.
func (p *DouyinPlatform) BuildClient() *Client {
	return NewClient(30, p.authStore.GetCookie())
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

// --- commands/douyin.Client interface delegation ---

// PlatformName returns the platform name.
func (p *DouyinPlatform) PlatformName() string {
	return p.BuildClient().PlatformName()
}

// GetHotSearch fetches hot/douyin trending search terms.
func (p *DouyinPlatform) GetHotSearch() (map[string]interface{}, error) {
	return p.BuildClient().GetHotSearch()
}

// GetTrending fetches the trending/challenge list.
func (p *DouyinPlatform) GetTrending() (map[string]interface{}, error) {
	return p.BuildClient().GetTrending()
}

// Search performs a search on Douyin for videos/users.
func (p *DouyinPlatform) Search(keyword string, page int, count int) (map[string]interface{}, error) {
	return p.BuildClient().Search(keyword, page, count)
}

// GetVideoDetail fetches video details by video ID.
func (p *DouyinPlatform) GetVideoDetail(videoID string) (map[string]interface{}, error) {
	return p.BuildClient().GetVideoDetail(videoID)
}

// GetVideoComments fetches comments for a video.
func (p *DouyinPlatform) GetVideoComments(videoID string, cursor int) (map[string]interface{}, error) {
	return p.BuildClient().GetVideoComments(videoID, cursor)
}

// GetUserPosts fetches posts by a user.
func (p *DouyinPlatform) GetUserPosts(uid string, cursor int) (map[string]interface{}, error) {
	return p.BuildClient().GetUserPosts(uid, cursor)
}

// GetUserInfo fetches user info by uid.
func (p *DouyinPlatform) GetUserInfo(uid string) (map[string]interface{}, error) {
	return p.BuildClient().GetUserInfo(uid)
}