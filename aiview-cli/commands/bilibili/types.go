// Package bilibili provides CLI commands for the Bilibili platform.
package bilibili

import (
	"regexp"

	biliapi "github.com/jackwener/aiview/internal/platform/bilibili/bilibilitypes"
	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// Client is the interface that the bilibili API client must satisfy for commands.
type Client interface {
	GetVideoInfo(bvid string) (*biliapi.VideoInfo, error)
	GetVideoSubtitle(bvid string) (*biliapi.SubtitleInfo, error)
	GetVideoAIConclusion(bvid string) (string, error)
	GetVideoComments(bvid string, page int) ([]biliapi.CommentInfo, error)
	GetRelatedVideos(bvid string) ([]biliapi.VideoInfo, error)
	SearchVideo(keyword string, page int, order string, duration int, tid int) ([]biliapi.SearchVideoResult, error)
	SearchUser(keyword string, page int) ([]biliapi.SearchUserResult, error)
	GetUserInfo(uid int) (*biliapi.UserInfo, error)
	GetUserVideos(uid int, count int, order string, tid int, keyword string) ([]biliapi.VideoInfo, error)
	GetHotVideos(page int, count int) ([]biliapi.VideoInfo, error)
	GetRankVideos(rid int, day int, typeStr string) ([]biliapi.VideoInfo, error)
	GetFavoriteList(uid int, page int) ([]biliapi.FavoriteFolder, error)
	GetFavoriteVideos(favID int, page int) ([]biliapi.FavoriteMedia, error)
	GetFollowingList(uid int, page int) ([]biliapi.FollowingUser, error)
	GetWatchHistory(page int, count int) ([]biliapi.HistoryItem, error)
	GetWatchLater() ([]biliapi.WatchLaterItem, error)
	GetDynamicFeed(offset string) ([]biliapi.DynamicItem, error)
	LikeVideo(bvid string, undo bool) error
	CoinVideo(bvid string, num int) error
	TripleVideo(bvid string) error
	UnfollowUser(uid int) error
	GetSelfInfo() (*biliapi.UserInfo, error)
	GetAudioURL(bvid string) (string, error)

	// Comment methods
	PostComment(oid int, message string, root int, parent int) error
	DeleteComment(oid int, rpid int) error
	GetVideoCommentsRaw(oid int, page int, sort int) (map[string]interface{}, error)

	// Danmaku methods
	GetVideoDanmaku(cid int) ([]byte, error)
	PostDanmaku(oid int, cid int, message string, progress int) error

	// Favorite methods
	AddFavorite(bvid string, fid int) error
	DelFavorite(bvid string, fid int) error

	// Recommend
	GetRecommendVideos(fresh bool, page int) (map[string]interface{}, error)

	// Tags
	GetVideoTags(bvid string) (map[string]interface{}, error)

	// Search suggestion
	SearchSuggest(keyword string) (map[string]interface{}, error)

	// Fans
	GetFansList(uid int, page int) (map[string]interface{}, error)

	// Dynamics
	GetUserDynamics(uid int, page int) (map[string]interface{}, error)
	PostDynamic(text string) (map[string]interface{}, error)
	DeleteDynamic(dynamicID int) (map[string]interface{}, error)

	// Collections
	GetUserCollections(uid int) (map[string]interface{}, error)

	// Relation stat
	GetRelationStat(uid int) (map[string]interface{}, error)

	// Region videos
	GetRegionVideos(rid int, page int, count int, sort string) (map[string]interface{}, error)

	// Live room info
	GetLiveRoomInfo(roomID int) (map[string]interface{}, error)

	// Precious videos
	GetPreciousVideos() (map[string]interface{}, error)

	// Hot search
	GetHotSearch(limit int) (map[string]interface{}, error)

	// Video online count
	GetVideoOnlineCount(bvid string) (map[string]interface{}, error)

	// Weekly hot videos
	GetWeeklyHotVideos(number int) (map[string]interface{}, error)
}

// AuthProvider is the interface for credential management.
type AuthProvider interface {
	GetCredential() (*Credential, error)
	GetCredentialOrNil() *Credential
	RequireCredential(requireWrite bool) (*Credential, error)
	Save(c *Credential) error
	Clear() error
}

// Credential holds authentication data.
type Credential = biliapi.Credential

// ActionResult holds the result of an action command.
type ActionResult struct {
	Success bool   `json:"success"`
	Action  string `json:"action"`
}

// QRLoginSession holds QR login session data.
type QRLoginSession struct {
	QRCodeKey string
	QRCodeURL string
}

// VideoPayload is the structured output for video commands.
type VideoPayload struct {
	Video     biliapi.VideoInfo     `json:"video"`
	Subtitle  biliapi.SubtitleInfo  `json:"subtitle"`
	AISummary string                `json:"ai_summary"`
	Comments  []biliapi.CommentInfo `json:"comments"`
	Related   []biliapi.VideoInfo   `json:"related"`
	Warnings  []Warning             `json:"warnings"`
}

// Warning holds a non-fatal warning message.
type Warning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ExtractBVID extracts a BV ID from a URL or string.
var BVIDRegex = regexp.MustCompile(`\bBV[0-9A-Za-z]{10}\b`)

func ExtractBVID(urlOrBvid string) (string, error) {
	match := BVIDRegex.FindString(urlOrBvid)
	if match != "" {
		return match, nil
	}
	return "", &BVIDError{Input: urlOrBvid}
}

// BVIDError represents an invalid BV ID error.
type BVIDError struct {
	Input string
}

func (e *BVIDError) Error() string {
	return "Failed to extract BV ID: " + e.Input
}

// GetOutputFormat extracts the output format from cobra command flags.
func GetOutputFormat(cmd *cobra.Command) output.Format {
	parent := cmd
	for parent.HasParent() {
		parent = parent.Parent()
	}
	asJSON, _ := parent.Flags().GetBool("json")
	asYAML, _ := parent.Flags().GetBool("yaml")
	return output.ResolveFormat(asJSON, asYAML)
}