package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackwener/aiview/internal/platform/bilibili/commands"
)

const (
	baseURL    = "https://api.bilibili.com"
	userAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Bilibili API client.
type Client struct {
	httpClient  *http.Client
	credential  *commands.Credential
	cookies     string
}

// NewClient creates a new Bilibili API client.
func NewClient(timeoutSec int, cookies string, credential *commands.Credential) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
		cookies:    cookies,
		credential: credential,
	}
}

// PlatformName implements the platform.Client interface.
func (c *Client) PlatformName() string {
	return "bilibili"
}

// BVIDRegex matches BV IDs.
var BVIDRegex = regexp.MustCompile(`\bBV[0-9A-Za-z]{10}\b`)

// ExtractBVID extracts a BV ID from a URL or returns the input if it's already a BV ID.
func ExtractBVID(urlOrBvid string) (string, error) {
	match := BVIDRegex.FindString(urlOrBvid)
	if match != "" {
		return match, nil
	}
	return "", fmt.Errorf("无法提取 BV 号: %s", urlOrBvid)
}

func (c *Client) buildHeaders() http.Header {
	h := http.Header{}
	h.Set("User-Agent", userAgent)
	h.Set("Origin", "https://www.bilibili.com")
	h.Set("Referer", "https://www.bilibili.com")
	h.Set("Accept", "application/json, text/plain, */*")
	h.Set("Accept-Language", "zh-CN,zh;q=0.9")
	h.Set("sec-ch-ua", "\"Chromium\";v=\"133\", \"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\"")
	h.Set("sec-ch-ua-mobile", "?0")
	h.Set("sec-ch-ua-platform", "\"Windows\"")
	if c.cookies != "" {
		h.Set("Cookie", c.cookies)
	}
	return h
}

func (c *Client) get(path string, params url.Values) (map[string]interface{}, error) {
	reqURL := baseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header = c.buildHeaders()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("网络请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	code := getInt(result, "code")
	if code != 0 {
		msg := getString(result, "message")
		if code == -101 || code == -111 {
			return nil, fmt.Errorf("not_authenticated: %s", msg)
		}
		if code == -404 || code == 62002 || code == 62004 {
			return nil, fmt.Errorf("not_found: %s", msg)
		}
		if code == -412 || code == 412 {
			return nil, fmt.Errorf("rate_limited: %s", msg)
		}
		return nil, fmt.Errorf("API 错误 [%d]: %s", code, msg)
	}

	return result, nil
}

// GetVideoInfo fetches video metadata.
func (c *Client) GetVideoInfo(bvid string) (*commands.VideoInfo, error) {
	params := url.Values{}
	params.Set("bvid", bvid)

	data, err := c.get("/x/web-interface/view", params)
	if err != nil {
		return nil, err
	}

	info := getMap(data, "data")
	owner := getMap(info, "owner")
	stat := getMap(info, "stat")

	return &commands.VideoInfo{
		BVID:        bvid,
		AID:         getInt(info, "aid"),
		Title:       getString(info, "title"),
		Description: strings.TrimSpace(getString(info, "desc")),
		Duration:    getInt(info, "duration"),
		DurationStr: formatDuration(getInt(info, "duration")),
		URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", bvid),
		Owner: commands.OwnerInfo{
			MID:  getInt(owner, "mid"),
			Name: getString(owner, "name"),
		},
		Stats: commands.VideoStats{
			View:     getInt(stat, "view"),
			Danmaku:  getInt(stat, "danmaku"),
			Like:     getInt(stat, "like"),
			Coin:     getInt(stat, "coin"),
			Favorite: getInt(stat, "favorite"),
			Share:    getInt(stat, "share"),
		},
	}, nil
}

// GetVideoSubtitle fetches video subtitle content.
func (c *Client) GetVideoSubtitle(bvid string) (*commands.SubtitleInfo, error) {
	// First get video info to get cid
	info, err := c.GetVideoInfo(bvid)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("cid", strconv.Itoa(info.AID))

	playerData, err := c.get("/x/player/v2", params)
	if err != nil {
		return nil, err
	}
	data := getMap(playerData, "data")
	subtitle := getMap(data, "subtitle")
	subtitles := getSlice(subtitle, "subtitles")

	if len(subtitles) == 0 {
		return &commands.SubtitleInfo{Available: false}, nil
	}

	// Prefer Chinese subtitles
	var subtitleURL string
	for _, sub := range subtitles {
		s := sub.(map[string]interface{})
		lan := getString(s, "lan")
		if strings.Contains(strings.ToLower(lan), "zh") {
			subtitleURL = getString(s, "subtitle_url")
			break
		}
	}
	if subtitleURL == "" {
		s := subtitles[0].(map[string]interface{})
		subtitleURL = getString(s, "subtitle_url")
	}

	if subtitleURL == "" {
		return &commands.SubtitleInfo{Available: false}, nil
	}

	if strings.HasPrefix(subtitleURL, "//") {
		subtitleURL = "https:" + subtitleURL
	}

	// Download subtitle JSON
	req, err := http.NewRequest("GET", subtitleURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var subData map[string]interface{}
	if err := json.Unmarshal(body, &subData); err != nil {
		return nil, err
	}

	rawItems := getSlice(subData, "body")
	items := make([]commands.SubtitleItem, 0, len(rawItems))
	var texts []string
	for _, item := range rawItems {
		m := item.(map[string]interface{})
		content := getString(m, "content")
		texts = append(texts, content)
		items = append(items, commands.SubtitleItem{
			From:    getFloat(m, "from"),
			To:      getFloat(m, "to"),
			Content: content,
		})
	}

	return &commands.SubtitleInfo{
		Available: true,
		Format:    "plain",
		Text:      strings.Join(texts, "\n"),
		Items:     items,
	}, nil
}

// GetVideoAIConclusion fetches AI-generated video summary.
func (c *Client) GetVideoAIConclusion(bvid string) (string, error) {
	info, err := c.GetVideoInfo(bvid)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("cid", strconv.Itoa(info.AID))

	data, err := c.get("/x/player/v2", params)
	if err != nil {
		return "", err
	}

	d := getMap(data, "data")
	ai := getMap(d, "ai_summary")
	return getString(ai, "summary"), nil
}

// GetVideoComments fetches video comments.
func (c *Client) GetVideoComments(bvid string, page int) ([]commands.CommentInfo, error) {
	info, err := c.GetVideoInfo(bvid)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("oid", strconv.Itoa(info.AID))
	params.Set("type", "1")
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", "20")
	params.Set("sort", "2")

	data, err := c.get("/x/v2/reply", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	replies := getSlice(d, "replies")

	comments := make([]commands.CommentInfo, 0, len(replies))
	for _, r := range replies {
		m := r.(map[string]interface{})
		member := getMap(m, "member")
		content := getMap(m, "content")
		comments = append(comments, commands.CommentInfo{
			ID: getString(m, "rpid_str"),
			Author: commands.AuthorInfo{
				MID:  getInt(member, "mid"),
				Name: getString(member, "uname"),
			},
			Message: getString(content, "message"),
			Like:    getInt(m, "like"),
		})
	}
	return comments, nil
}

// GetRelatedVideos fetches related/recommended videos.
func (c *Client) GetRelatedVideos(bvid string) ([]commands.VideoInfo, error) {
	params := url.Values{}
	params.Set("bvid", bvid)

	data, err := c.get("/x/web-interface/archive/related", params)
	if err != nil {
		return nil, err
	}

	items := getSlice(data, "data")
	videos := make([]commands.VideoInfo, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		owner := getMap(m, "owner")
		stat := getMap(m, "stat")
		videos = append(videos, commands.VideoInfo{
			BVID:        getString(m, "bvid"),
			Title:       getString(m, "title"),
			Duration:    getInt(m, "duration"),
			DurationStr: formatDuration(getInt(m, "duration")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", getString(m, "bvid")),
			Owner: commands.OwnerInfo{
				MID:  getInt(owner, "mid"),
				Name: getString(owner, "name"),
			},
			Stats: commands.VideoStats{
				View: getInt(stat, "view"),
			},
		})
	}
	return videos, nil
}

// SearchVideo searches for videos by keyword.
func (c *Client) SearchVideo(keyword string, page int) ([]commands.SearchVideoResult, error) {
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("search_type", "video")
	params.Set("page", strconv.Itoa(page))

	data, err := c.wbiGet("/x/web-interface/wbi/search/type", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	items := getSlice(d, "result")

	results := make([]commands.SearchVideoResult, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		results = append(results, commands.SearchVideoResult{
			BVID:     getString(m, "bvid"),
			Title:    stripHTML(getString(m, "title")),
			Author:   getString(m, "author"),
			Play:     getInt(m, "play"),
			Duration: getString(m, "duration"),
		})
	}
	return results, nil
}

// SearchUser searches for users by keyword.
func (c *Client) SearchUser(keyword string, page int) ([]commands.SearchUserResult, error) {
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("search_type", "bili_user")
	params.Set("page", strconv.Itoa(page))

	data, err := c.wbiGet("/x/web-interface/wbi/search/type", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	items := getSlice(d, "result")

	results := make([]commands.SearchUserResult, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		results = append(results, commands.SearchUserResult{
			MID:    getInt(m, "mid"),
			Name:   getString(m, "uname"),
			Sign:   getString(m, "usign"),
			Fans:   getInt(m, "fans"),
			Videos: getInt(m, "videos"),
		})
	}
	return results, nil
}

// GetUserInfo fetches user profile information.
func (c *Client) GetUserInfo(uid int) (*commands.UserInfo, error) {
	params := url.Values{}
	params.Set("mid", strconv.Itoa(uid))
	params.Set("photo", "false")

	data, err := c.wbiGet("/x/web-interface/card", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	card := getMap(d, "card")
	return &commands.UserInfo{
		MID:       getInt(card, "mid"),
		Name:      getString(card, "name"),
		Level:     getInt(card, "level_info", "current_level"),
		Coins:     getInt(card, "coins"),
		Sign:      getString(card, "sign"),
		Fans:      getInt(card, "fans"),
		Following: getInt(card, "attention"),
	}, nil
}

// GetUserVideos fetches a user's latest videos.
func (c *Client) GetUserVideos(uid int, count int) ([]commands.VideoInfo, error) {
	params := url.Values{}
	params.Set("mid", strconv.Itoa(uid))
	params.Set("ps", strconv.Itoa(min(count, 50)))
	params.Set("pn", "1")
	params.Set("order", "pubdate")
	params.Set("tid", "0")
	params.Set("keyword", "")

	data, err := c.wbiGet("/x/space/wbi/arc/search", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getMap(d, "list")
	vlist := getSlice(list, "vlist")

	videos := make([]commands.VideoInfo, 0, len(vlist))
	for _, v := range vlist {
		m := v.(map[string]interface{})
		videos = append(videos, commands.VideoInfo{
			BVID:        getString(m, "bvid"),
			Title:       getString(m, "title"),
			Duration:    getInt(m, "length"),
			DurationStr: formatDuration(getInt(m, "length")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", getString(m, "bvid")),
			Stats: commands.VideoStats{
				View: getInt(m, "play"),
			},
		})
		if len(videos) >= count {
			break
		}
	}
	return videos, nil
}

// GetHotVideos fetches popular/hot videos.
func (c *Client) GetHotVideos(page int, count int) ([]commands.VideoInfo, error) {
	params := url.Values{}
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", strconv.Itoa(min(count, 50)))

	data, err := c.get("/x/web-interface/popular", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	videos := make([]commands.VideoInfo, 0, len(list))
	for _, v := range list {
		m := v.(map[string]interface{})
		owner := getMap(m, "owner")
		stat := getMap(m, "stat")
		videos = append(videos, commands.VideoInfo{
			BVID:        getString(m, "bvid"),
			Title:       getString(m, "title"),
			Duration:    getInt(m, "duration"),
			DurationStr: formatDuration(getInt(m, "duration")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", getString(m, "bvid")),
			Owner: commands.OwnerInfo{
				MID:  getInt(owner, "mid"),
				Name: getString(owner, "name"),
			},
			Stats: commands.VideoStats{
				View: getInt(stat, "view"),
			},
		})
	}
	return videos, nil
}

// GetRankVideos fetches ranking videos.
func (c *Client) GetRankVideos(day int) ([]commands.VideoInfo, error) {
	rid := 0 // all
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("type", "all")

	data, err := c.get("/x/web-interface/ranking/v2", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	videos := make([]commands.VideoInfo, 0, len(list))
	for _, v := range list {
		m := v.(map[string]interface{})
		owner := getMap(m, "owner")
		stat := getMap(m, "stat")
		videos = append(videos, commands.VideoInfo{
			BVID:        getString(m, "bvid"),
			Title:       getString(m, "title"),
			Duration:    getInt(m, "duration"),
			DurationStr: formatDuration(getInt(m, "duration")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", getString(m, "bvid")),
			Owner: commands.OwnerInfo{
				MID:  getInt(owner, "mid"),
				Name: getString(owner, "name"),
			},
			Stats: commands.VideoStats{
				View: getInt(stat, "view"),
			},
		})
	}
	return videos, nil
}

// GetFavoriteList fetches favorite folders.
func (c *Client) GetFavoriteList(uid int) ([]commands.FavoriteFolder, error) {
	params := url.Values{}
	params.Set("up_mid", strconv.Itoa(uid))

	data, err := c.get("/x/v3/fav/folder/created/list-all", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	folders := make([]commands.FavoriteFolder, 0, len(list))
	for _, f := range list {
		m := f.(map[string]interface{})
		folders = append(folders, commands.FavoriteFolder{
			ID:         getInt(m, "id"),
			Title:      getString(m, "title"),
			MediaCount: getInt(m, "media_count"),
		})
	}
	return folders, nil
}

// GetFavoriteVideos fetches videos in a favorite folder.
func (c *Client) GetFavoriteVideos(favID int, page int) ([]commands.FavoriteMedia, error) {
	params := url.Values{}
	params.Set("media_id", strconv.Itoa(favID))
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", "20")

	data, err := c.get("/x/v3/fav/resource/list", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	medias := getSlice(d, "medias")

	items := make([]commands.FavoriteMedia, 0, len(medias))
	for _, m := range medias {
		item := m.(map[string]interface{})
		upper := getMap(item, "upper")
		items = append(items, commands.FavoriteMedia{
			BVID:     getString(item, "bvid"),
			Title:    getString(item, "title"),
			Duration: formatDuration(getInt(item, "duration")),
			Upper:    getString(upper, "name"),
		})
	}
	return items, nil
}

// GetFollowingList fetches the user's following list.
func (c *Client) GetFollowingList(uid int, page int) ([]commands.FollowingUser, error) {
	params := url.Values{}
	params.Set("vmid", strconv.Itoa(uid))
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", "20")

	data, err := c.get("/x/relation/followings", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	users := make([]commands.FollowingUser, 0, len(list))
	for _, u := range list {
		m := u.(map[string]interface{})
		users = append(users, commands.FollowingUser{
			MID:  getInt(m, "mid"),
			Name: getString(m, "uname"),
			Sign: getString(m, "sign"),
		})
	}
	return users, nil
}

// GetWatchHistory fetches watch history.
func (c *Client) GetWatchHistory(page int, count int) ([]commands.HistoryItem, error) {
	params := url.Values{}
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", strconv.Itoa(min(count, 100)))

	data, err := c.get("/x/web-interface/history/cursor", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	items := make([]commands.HistoryItem, 0, len(list))
	for _, h := range list {
		m := h.(map[string]interface{})
		items = append(items, commands.HistoryItem{
			BVID:   getString(m, "bvid"),
			Title:  getString(m, "title"),
			Author: getString(m, "author_name"),
		})
	}
	return items, nil
}

// GetWatchLater fetches watch-later list.
func (c *Client) GetWatchLater() ([]commands.WatchLaterItem, error) {
	params := url.Values{}
	params.Set("ps", "20")

	data, err := c.get("/x/v2/history/toview/web", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	items := make([]commands.WatchLaterItem, 0, len(list))
	for _, w := range list {
		m := w.(map[string]interface{})
		owner := getMap(m, "owner")
		items = append(items, commands.WatchLaterItem{
			BVID:     getString(m, "bvid"),
			Title:    getString(m, "title"),
			Author:   getString(owner, "name"),
			Duration: formatDuration(getInt(m, "duration")),
		})
	}
	return items, nil
}

// GetDynamicFeed fetches the dynamic timeline.
func (c *Client) GetDynamicFeed(offset string) ([]commands.DynamicItem, error) {
	params := url.Values{}
	params.Set("type", "all")
	if offset != "" {
		params.Set("offset", offset)
	}

	data, err := c.get("/x/polymer/web-dynamic/v1/feed/all", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	items := getSlice(d, "items")

	dynamics := make([]commands.DynamicItem, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		modules := getMap(m, "modules")
		author := getMap(modules, "module_author")
		desc := getMap(getMap(modules, "module_dynamic"), "desc")

		dynamics = append(dynamics, commands.DynamicItem{
			ID:     getString(m, "id_str"),
			Author: getString(author, "name"),
			Text:   getString(desc, "text"),
		})
	}
	return dynamics, nil
}

// LikeVideo likes a video.
func (c *Client) LikeVideo(bvid string, undo bool) error {
	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("type", "1")
	if undo {
		params.Set("act", "2")
	} else {
		params.Set("act", "1")
	}

	_, err := c.get("/x/web-interface/archive/like", params)
	return err
}

// CoinVideo gives coins to a video.
func (c *Client) CoinVideo(bvid string, num int) error {
	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("multiply", strconv.Itoa(num))
	params.Set("select_like", "1")

	_, err := c.get("/x/web-interface/coin/add", params)
	return err
}

// TripleVideo does like + coin + favorite on a video.
func (c *Client) TripleVideo(bvid string) error {
	params := url.Values{}
	params.Set("bvid", bvid)

	_, err := c.get("/x/web-interface/archive/like/triple", params)
	return err
}

// UnfollowUser unfollows a user.
func (c *Client) UnfollowUser(uid int) error {
	params := url.Values{}
	params.Set("fid", strconv.Itoa(uid))
	params.Set("act", "2")
	params.Set("re_src", "11")

	_, err := c.get("/x/relation/modify", params)
	return err
}

// GetSelfInfo fetches the logged-in user's own info.
func (c *Client) GetSelfInfo() (*commands.UserInfo, error) {
	data, err := c.get("/x/web-interface/nav", nil)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	if !getBool(d, "isLogin") {
		return nil, fmt.Errorf("not_authenticated: 未登录")
	}

	return &commands.UserInfo{
		MID:   getInt(d, "mid"),
		Name:  getString(d, "uname"),
		Level: getInt(d, "level_info", "current_level"),
		Coins: getInt(d, "money"),
	}, nil
}

// GetAudioURL gets the audio stream URL for a video.
func (c *Client) GetAudioURL(bvid string) (string, error) {
	info, err := c.GetVideoInfo(bvid)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("cid", strconv.Itoa(info.AID))
	params.Set("qn", "0")
	params.Set("fnval", "4048")
	params.Set("fourk", "1")

	data, err := c.get("/x/player/playurl", params)
	if err != nil {
		return "", err
	}

	d := getMap(data, "data")
	dash := getMap(d, "dash")
	audio := getSlice(dash, "audio")

	if len(audio) == 0 {
		return "", fmt.Errorf("not_found: 无法获取音频流")
	}

	// Get the best quality audio
	best := audio[0].(map[string]interface{})
	return getString(best, "baseUrl"), nil
}

// Helper functions for parsing JSON

func getString(m map[string]interface{}, keys ...string) string {
	if m == nil {
		return ""
	}
	val, ok := m[keys[0]]
	if !ok {
		return ""
	}
	if len(keys) == 1 {
		s, _ := val.(string)
		return s
	}
	sub, ok := val.(map[string]interface{})
	if !ok {
		return ""
	}
	return getString(sub, keys[1:]...)
}

func getInt(m map[string]interface{}, keys ...string) int {
	if m == nil {
		return 0
	}
	val, ok := m[keys[0]]
	if !ok {
		return 0
	}
	if len(keys) == 1 {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case string:
			n, _ := strconv.Atoi(v)
			return n
		}
		return 0
	}
	sub, ok := val.(map[string]interface{})
	if !ok {
		return 0
	}
	return getInt(sub, keys[1:]...)
}

func getFloat(m map[string]interface{}, key string) float64 {
	if m == nil {
		return 0
	}
	val, ok := m[key]
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}
	return 0
}

func getBool(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	val, ok := m[key]
	if !ok {
		return false
	}
	b, _ := val.(bool)
	return b
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok {
		return nil
	}
	sub, _ := val.(map[string]interface{})
	return sub
}

func getSlice(m map[string]interface{}, key string) []interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok {
		return nil
	}
	sub, _ := val.([]interface{})
	return sub
}

func formatDuration(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}
	if seconds >= 3600 {
		h := seconds / 3600
		m := (seconds % 3600) / 60
		s := seconds % 60
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	m := seconds / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func stripHTML(text string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return strings.TrimSpace(re.ReplaceAllString(text, ""))
}