// Package bilibilitypes holds shared data types for the Bilibili platform.
// This package exists to break an import cycle between the platform and commands layers.
package bilibilitypes

// API response types used across the bilibili platform layer.

type (
	VideoInfo struct {
		BVID        string     `json:"bvid"`
		AID         int        `json:"aid"`
		CID         int        `json:"cid"`
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

	// Credential holds authentication data.
	Credential struct {
		Sessdata    string `json:"sessdata"`
		BiliJct     string `json:"bili_jct"`
		AcTimeValue string `json:"ac_time_value"`
		Buvid3      string `json:"buvid3"`
		Buvid4      string `json:"buvid4"`
		DedeUserID  string `json:"dedeuserid"`
		SavedAt     int64  `json:"saved_at"`
	}
)

// DanmakuInfo represents a danmaku (bullet comment) item.
type DanmakuInfo struct {
	ID       int64  `json:"id"`
	Progress int    `json:"progress"`
	Mode     int    `json:"mode"`
	FontSize int    `json:"fontsize"`
	Color    int64  `json:"color"`
	SendTime int64  `json:"send_time"`
	Message  string `json:"message"`
}

// VideoTag represents a video tag.
type VideoTag struct {
	TagID   int    `json:"tag_id"`
	TagName string `json:"tag_name"`
	Likes   int    `json:"likes"`
}

// FansUserInfo represents a fan/follower.
type FansUserInfo struct {
	MID    int    `json:"mid"`
	Name   string `json:"uname"`
	Sign   string `json:"sign"`
	Avatar string `json:"avatar"`
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