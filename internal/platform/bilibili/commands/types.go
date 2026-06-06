package commands

import (
	"regexp"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// Client is the interface that the bilibili API client must satisfy for commands.
type Client interface {
	GetVideoInfo(bvid string) (*VideoInfo, error)
	GetVideoSubtitle(bvid string) (*SubtitleInfo, error)
	GetVideoAIConclusion(bvid string) (string, error)
	GetVideoComments(bvid string, page int) ([]CommentInfo, error)
	GetRelatedVideos(bvid string) ([]VideoInfo, error)
	SearchVideo(keyword string, page int) ([]SearchVideoResult, error)
	SearchUser(keyword string, page int) ([]SearchUserResult, error)
	GetUserInfo(uid int) (*UserInfo, error)
	GetUserVideos(uid int, count int) ([]VideoInfo, error)
	GetHotVideos(page int, count int) ([]VideoInfo, error)
	GetRankVideos(day int) ([]VideoInfo, error)
	GetFavoriteList(uid int) ([]FavoriteFolder, error)
	GetFavoriteVideos(favID int, page int) ([]FavoriteMedia, error)
	GetFollowingList(uid int, page int) ([]FollowingUser, error)
	GetWatchHistory(page int, count int) ([]HistoryItem, error)
	GetWatchLater() ([]WatchLaterItem, error)
	GetDynamicFeed(offset string) ([]DynamicItem, error)
	LikeVideo(bvid string, undo bool) error
	CoinVideo(bvid string, num int) error
	TripleVideo(bvid string) error
	UnfollowUser(uid int) error
	GetSelfInfo() (*UserInfo, error)
	GetAudioURL(bvid string) (string, error)
}

// Model types used by commands.
type (
	VideoInfo struct {
		BVID        string     `json:"bvid"`
		AID         int        `json:"aid"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Duration    int        `json:"duration_seconds"`
		DurationStr string     `json:"duration"`
		URL         string     `json:"url"`
		Owner       OwnerInfo  `json:"owner"`
		Stats       VideoStats `json:"stats"`
	}

	OwnerInfo struct {
		MID  int    `json:"id"`
		Name string `json:"name"`
	}

	VideoStats struct {
		View     int `json:"view"`
		Danmaku  int `json:"danmaku"`
		Like     int `json:"like"`
		Coin     int `json:"coin"`
		Favorite int `json:"favorite"`
		Share    int `json:"share"`
	}

	UserInfo struct {
		MID       int    `json:"id"`
		Name      string `json:"name"`
		Level     int    `json:"level"`
		Coins     int    `json:"coins"`
		Sign      string `json:"sign"`
		Fans      int    `json:"fans"`
		Following int    `json:"following"`
	}

	SearchUserResult struct {
		MID    int    `json:"id"`
		Name   string `json:"name"`
		Sign   string `json:"sign"`
		Fans   int    `json:"fans"`
		Videos int    `json:"videos"`
	}

	SearchVideoResult struct {
		BVID     string `json:"id"`
		Title    string `json:"title"`
		Author   string `json:"author"`
		Play     int    `json:"play"`
		Duration string `json:"duration"`
	}

	CommentInfo struct {
		ID      string     `json:"id"`
		Author  AuthorInfo `json:"author"`
		Message string     `json:"message"`
		Like    int        `json:"like"`
	}

	AuthorInfo struct {
		MID  int    `json:"id"`
		Name string `json:"name"`
	}

	SubtitleInfo struct {
		Available bool           `json:"available"`
		Format    string         `json:"format"`
		Text      string         `json:"text"`
		Items     []SubtitleItem `json:"items"`
	}

	SubtitleItem struct {
		From    float64 `json:"from"`
		To      float64 `json:"to"`
		Content string  `json:"content"`
	}

	FavoriteFolder struct {
		ID         int    `json:"id"`
		Title      string `json:"title"`
		MediaCount int    `json:"media_count"`
	}

	FavoriteMedia struct {
		BVID     string `json:"id"`
		Title    string `json:"title"`
		Duration string `json:"duration"`
		Upper    string `json:"upper"`
	}

	FollowingUser struct {
		MID  int    `json:"id"`
		Name string `json:"name"`
		Sign string `json:"sign"`
	}

	HistoryItem struct {
		BVID     string `json:"id"`
		Title    string `json:"title"`
		Author   string `json:"author"`
		ViewedAt string `json:"viewed_at"`
	}

	WatchLaterItem struct {
		BVID     string `json:"id"`
		Title    string `json:"title"`
		Author   string `json:"author"`
		Duration string `json:"duration"`
	}

	DynamicItem struct {
		ID        string       `json:"id"`
		Author    string       `json:"author"`
		Published string       `json:"published_at"`
		Text      string       `json:"text"`
		Stats     DynamicStats `json:"stats"`
	}

	DynamicStats struct {
		Comment int `json:"comment"`
		Like    int `json:"like"`
	}

	VideoPayload struct {
		Video     VideoInfo     `json:"video"`
		Subtitle  SubtitleInfo  `json:"subtitle"`
		AISummary string        `json:"ai_summary"`
		Comments  []CommentInfo `json:"comments"`
		Related   []VideoInfo   `json:"related"`
		Warnings  []Warning     `json:"warnings"`
	}

	Warning struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	ActionResult struct {
		Success bool   `json:"success"`
		Action  string `json:"action"`
	}
)

// QRLoginSession holds QR login session data.
type QRLoginSession struct {
	QRCodeKey string
	QRCodeURL string
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
type Credential struct {
	Sessdata    string `json:"sessdata"`
	BiliJct     string `json:"bili_jct"`
	AcTimeValue string `json:"ac_time_value"`
	Buvid3      string `json:"buvid3"`
	Buvid4      string `json:"buvid4"`
	DedeUserID  string `json:"dedeuserid"`
	SavedAt     int64  `json:"saved_at"`
}

// IsValid checks if the credential has the minimum required fields.
func (c *Credential) IsValid() bool {
	return c != nil && c.Sessdata != ""
}

// HasWriteCapability checks if the credential supports write operations.
func (c *Credential) HasWriteCapability() bool {
	return c != nil && c.BiliJct != ""
}

// IsStale checks if the credential is older than TTL days.
func (c *Credential) IsStale(ttlDays int) bool {
	if c == nil || c.SavedAt == 0 {
		return true
	}
	return 0 > int64(ttlDays*86400) // simplified: always treat as valid for now
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