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
)

const (
	baseURL    = "https://api.bilibili.com"
	userAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// Client is the Bilibili API client.
type Client struct {
	httpClient  *http.Client
	credential  *Credential
	cookies     string
}

// NewClient creates a new Bilibili API client.
func NewClient(timeoutSec int, cookies string, credential *Credential) *Client {
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
	return "", fmt.Errorf("failed to extract BV ID: %s", urlOrBvid)
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

func (c *Client) buildRefererHeaders(referer string) http.Header {
	h := c.buildHeaders()
	h.Set("Referer", referer)
	return h
}

func (c *Client) get(path string, params url.Values) (map[string]interface{}, error) {
	reqURL := baseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = c.buildHeaders()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) getWithReferer(path string, params url.Values, referer string) (map[string]interface{}, error) {
	reqURL := baseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = c.buildRefererHeaders(referer)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) parseResponse(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Failed to parse response: %w", err)
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
		return nil, fmt.Errorf("API error [%d]: %s", code, msg)
	}

	return result, nil
}

func (c *Client) post(path string, params url.Values) (map[string]interface{}, error) {
	if c.credential != nil && c.credential.BiliJct != "" {
		params.Set("csrf", c.credential.BiliJct)
	}

	body := params.Encode()
	req, err := http.NewRequest("POST", baseURL+path, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = c.buildHeaders()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("Failed to parse response: %w", err)
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
		return nil, fmt.Errorf("API error [%d]: %s", code, msg)
	}

	return result, nil
}

// GetVideoInfo fetches video metadata.
func (c *Client) GetVideoInfo(bvid string) (*VideoInfo, error) {
	params := url.Values{}
	params.Set("bvid", bvid)

	data, err := c.get("/x/web-interface/view", params)
	if err != nil {
		return nil, err
	}

	info := getMap(data, "data")
	owner := getMap(info, "owner")
	stat := getMap(info, "stat")

	return &VideoInfo{
		BVID:        bvid,
		AID:         getInt(info, "aid"),
		CID:         getInt(info, "cid"),
		Title:       getString(info, "title"),
		Description: strings.TrimSpace(getString(info, "desc")),
		Duration:    getInt(info, "duration"),
		DurationStr: formatDuration(getInt(info, "duration")),
		URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", bvid),
		Owner: OwnerInfo{
			MID:  getInt(owner, "mid"),
			Name: getString(owner, "name"),
		},
		Stats: VideoStats{
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
func (c *Client) GetVideoSubtitle(bvid string) (*SubtitleInfo, error) {
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
		return &SubtitleInfo{Available: false}, nil
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
		return &SubtitleInfo{Available: false}, nil
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
	items := make([]SubtitleItem, 0, len(rawItems))
	var texts []string
	for _, item := range rawItems {
		m := item.(map[string]interface{})
		content := getString(m, "content")
		texts = append(texts, content)
		items = append(items, SubtitleItem{
			From:    getFloat(m, "from"),
			To:      getFloat(m, "to"),
			Content: content,
		})
	}

	return &SubtitleInfo{
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
	params.Set("cid", strconv.Itoa(info.CID))

	data, err := c.get("/x/player/v2", params)
	if err != nil {
		return "", err
	}

	d := getMap(data, "data")
	ai := getMap(d, "ai_summary")
	return getString(ai, "summary"), nil
}

// GetVideoComments fetches video comments.
func (c *Client) GetVideoComments(bvid string, page int) ([]CommentInfo, error) {
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

	comments := make([]CommentInfo, 0, len(replies))
	for _, r := range replies {
		m := r.(map[string]interface{})
		member := getMap(m, "member")
		content := getMap(m, "content")
		comments = append(comments, CommentInfo{
			ID: getString(m, "rpid_str"),
			Author: AuthorInfo{
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
func (c *Client) GetRelatedVideos(bvid string) ([]VideoInfo, error) {
	params := url.Values{}
	params.Set("bvid", bvid)

	data, err := c.get("/x/web-interface/archive/related", params)
	if err != nil {
		return nil, err
	}

	items := getSlice(data, "data")
	videos := make([]VideoInfo, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		owner := getMap(m, "owner")
		stat := getMap(m, "stat")
		videos = append(videos, VideoInfo{
			BVID:        getString(m, "bvid"),
			Title:       getString(m, "title"),
			Duration:    getInt(m, "duration"),
			DurationStr: formatDuration(getInt(m, "duration")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", getString(m, "bvid")),
			Owner: OwnerInfo{
				MID:  getInt(owner, "mid"),
				Name: getString(owner, "name"),
			},
			Stats: VideoStats{
				View: getInt(stat, "view"),
			},
		})
	}
	return videos, nil
}

// SearchVideo searches for videos by keyword.
func (c *Client) SearchVideo(keyword string, page int, order string, duration int, tid int) ([]SearchVideoResult, error) {
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("search_type", "video")
	params.Set("page", strconv.Itoa(page))
	if order != "" {
		params.Set("order", order)
	}
	if duration > 0 {
		params.Set("duration", strconv.Itoa(duration))
	}
	if tid > 0 {
		params.Set("tid", strconv.Itoa(tid))
	}

	data, err := c.wbiGet("/x/web-interface/wbi/search/type", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	items := getSlice(d, "result")

	results := make([]SearchVideoResult, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		results = append(results, SearchVideoResult{
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
func (c *Client) SearchUser(keyword string, page int) ([]SearchUserResult, error) {
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

	results := make([]SearchUserResult, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		results = append(results, SearchUserResult{
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
func (c *Client) GetUserInfo(uid int) (*UserInfo, error) {
	params := url.Values{}
	params.Set("mid", strconv.Itoa(uid))
	params.Set("photo", "false")

	data, err := c.wbiGet("/x/web-interface/card", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	card := getMap(d, "card")
	return &UserInfo{
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
func (c *Client) GetUserVideos(uid int, count int, order string, tid int, keyword string) ([]VideoInfo, error) {
	params := url.Values{}
	params.Set("mid", strconv.Itoa(uid))
	params.Set("ps", strconv.Itoa(min(count, 50)))
	params.Set("pn", "1")
	params.Set("order", order)
	params.Set("tid", strconv.Itoa(tid))
	params.Set("keyword", keyword)

	data, err := c.wbiGet("/x/space/wbi/arc/search", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getMap(d, "list")
	vlist := getSlice(list, "vlist")

	videos := make([]VideoInfo, 0, len(vlist))
	for _, v := range vlist {
		m := v.(map[string]interface{})
		videos = append(videos, VideoInfo{
			BVID:        getString(m, "bvid"),
			Title:       getString(m, "title"),
			Duration:    getInt(m, "length"),
			DurationStr: formatDuration(getInt(m, "length")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", getString(m, "bvid")),
			Stats: VideoStats{
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
func (c *Client) GetHotVideos(page int, count int) ([]VideoInfo, error) {
	params := url.Values{}
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", strconv.Itoa(min(count, 50)))

	data, err := c.get("/x/web-interface/popular", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	videos := make([]VideoInfo, 0, len(list))
	for _, v := range list {
		m := v.(map[string]interface{})
		owner := getMap(m, "owner")
		stat := getMap(m, "stat")
		videos = append(videos, VideoInfo{
			BVID:        getString(m, "bvid"),
			Title:       getString(m, "title"),
			Duration:    getInt(m, "duration"),
			DurationStr: formatDuration(getInt(m, "duration")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", getString(m, "bvid")),
			Owner: OwnerInfo{
				MID:  getInt(owner, "mid"),
				Name: getString(owner, "name"),
			},
			Stats: VideoStats{
				View: getInt(stat, "view"),
			},
		})
	}
	return videos, nil
}

// GetRankVideos fetches ranking videos.
func (c *Client) GetRankVideos(rid int, day int, typeStr string) ([]VideoInfo, error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("type", "all")

	data, err := c.wbiGet("/x/web-interface/ranking/v2", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	videos := make([]VideoInfo, 0, len(list))
	for _, v := range list {
		m := v.(map[string]interface{})
		owner := getMap(m, "owner")
		stat := getMap(m, "stat")
		videos = append(videos, VideoInfo{
			BVID:        getString(m, "bvid"),
			Title:       getString(m, "title"),
			Duration:    getInt(m, "duration"),
			DurationStr: formatDuration(getInt(m, "duration")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", getString(m, "bvid")),
			Owner: OwnerInfo{
				MID:  getInt(owner, "mid"),
				Name: getString(owner, "name"),
			},
			Stats: VideoStats{
				View: getInt(stat, "view"),
			},
		})
	}
	return videos, nil
}

// GetFavoriteList fetches favorite folders.
func (c *Client) GetFavoriteList(uid int, page int) ([]FavoriteFolder, error) {
	params := url.Values{}
	params.Set("up_mid", strconv.Itoa(uid))

	data, err := c.get("/x/v3/fav/folder/created/list-all", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	folders := make([]FavoriteFolder, 0, len(list))
	for _, f := range list {
		m := f.(map[string]interface{})
		folders = append(folders, FavoriteFolder{
			ID:         getInt(m, "id"),
			Title:      getString(m, "title"),
			MediaCount: getInt(m, "media_count"),
		})
	}
	return folders, nil
}

// GetFavoriteVideos fetches videos in a favorite folder.
func (c *Client) GetFavoriteVideos(favID int, page int) ([]FavoriteMedia, error) {
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

	items := make([]FavoriteMedia, 0, len(medias))
	for _, m := range medias {
		item := m.(map[string]interface{})
		upper := getMap(item, "upper")
		items = append(items, FavoriteMedia{
			BVID:     getString(item, "bvid"),
			Title:    getString(item, "title"),
			Duration: formatDuration(getInt(item, "duration")),
			Upper:    getString(upper, "name"),
		})
	}
	return items, nil
}

// GetFollowingList fetches the user's following list.
func (c *Client) GetFollowingList(uid int, page int) ([]FollowingUser, error) {
	params := url.Values{}
	params.Set("vmid", strconv.Itoa(uid))
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", "50")

	data, err := c.get("/x/relation/followings", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	users := make([]FollowingUser, 0, len(list))
	for _, u := range list {
		m := u.(map[string]interface{})
		users = append(users, FollowingUser{
			MID:  getInt(m, "mid"),
			Name: getString(m, "uname"),
			Sign: getString(m, "sign"),
		})
	}
	return users, nil
}

// GetWatchHistory fetches watch history.
func (c *Client) GetWatchHistory(page int, count int) ([]HistoryItem, error) {
	params := url.Values{}
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", strconv.Itoa(min(count, 100)))

	data, err := c.get("/x/web-interface/history/cursor", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	items := make([]HistoryItem, 0, len(list))
	for _, h := range list {
		m := h.(map[string]interface{})
		items = append(items, HistoryItem{
			BVID:   getString(m, "bvid"),
			Title:  getString(m, "title"),
			Author: getString(m, "author_name"),
		})
	}
	return items, nil
}

// GetWatchLater fetches watch-later list.
func (c *Client) GetWatchLater() ([]WatchLaterItem, error) {
	params := url.Values{}
	params.Set("ps", "20")

	data, err := c.get("/x/v2/history/toview/web", params)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	list := getSlice(d, "list")

	items := make([]WatchLaterItem, 0, len(list))
	for _, w := range list {
		m := w.(map[string]interface{})
		owner := getMap(m, "owner")
		items = append(items, WatchLaterItem{
			BVID:     getString(m, "bvid"),
			Title:    getString(m, "title"),
			Author:   getString(owner, "name"),
			Duration: formatDuration(getInt(m, "duration")),
		})
	}
	return items, nil
}

// GetDynamicFeed fetches the dynamic timeline.
func (c *Client) GetDynamicFeed(offset string) ([]DynamicItem, error) {
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

	dynamics := make([]DynamicItem, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		modules := getMap(m, "modules")
		author := getMap(modules, "module_author")
		desc := getMap(getMap(modules, "module_dynamic"), "desc")

		dynamics = append(dynamics, DynamicItem{
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
	if undo {
		params.Set("like", "2")
	} else {
		params.Set("like", "1")
	}

	_, err := c.post("/x/web-interface/archive/like", params)
	return err
}

// CoinVideo gives coins to a video.
func (c *Client) CoinVideo(bvid string, num int) error {
	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("multiply", strconv.Itoa(num))
	params.Set("select_like", "1")

	_, err := c.post("/x/web-interface/coin/add", params)
	return err
}

// TripleVideo does like + coin + favorite on a video.
func (c *Client) TripleVideo(bvid string) error {
	params := url.Values{}
	params.Set("bvid", bvid)

	_, err := c.post("/x/web-interface/archive/like/triple", params)
	return err
}

// UnfollowUser unfollows a user.
func (c *Client) UnfollowUser(uid int) error {
	params := url.Values{}
	params.Set("fid", strconv.Itoa(uid))
	params.Set("act", "2")
	params.Set("re_src", "11")

	_, err := c.post("/x/relation/modify", params)
	return err
}

// GetSelfInfo fetches the logged-in user's own info.
func (c *Client) GetSelfInfo() (*UserInfo, error) {
	data, err := c.get("/x/web-interface/nav", nil)
	if err != nil {
		return nil, err
	}

	d := getMap(data, "data")
	if !getBool(d, "isLogin") {
		return nil, fmt.Errorf("not_authenticated: not logged in")
	}

	return &UserInfo{
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
	params.Set("cid", strconv.Itoa(info.CID))
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
		return "", fmt.Errorf("not_found: unable to get audio stream")
	}

	// Get the best quality audio
	best := audio[0].(map[string]interface{})
	return getString(best, "baseUrl"), nil
}

// PostComment posts a comment on a video.
// root=0, parent=0 means top-level comment.
func (c *Client) PostComment(oid int, message string, root int, parent int) error {
	params := url.Values{}
	params.Set("oid", strconv.Itoa(oid))
	params.Set("type", "1")
	params.Set("message", message)
	params.Set("root", strconv.Itoa(root))
	params.Set("parent", strconv.Itoa(parent))

	_, err := c.post("/x/v2/comment/add", params)
	return err
}

// DeleteComment deletes a comment.
func (c *Client) DeleteComment(oid int, rpid int) error {
	params := url.Values{}
	params.Set("oid", strconv.Itoa(oid))
	params.Set("type", "1")
	params.Set("rpid", strconv.Itoa(rpid))

	_, err := c.post("/x/v2/comment/del", params)
	return err
}

// GetVideoCommentsRaw fetches raw comment data with sort option.
// sort: 0=time, 2=hot
func (c *Client) GetVideoCommentsRaw(oid int, page int, sort int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("oid", strconv.Itoa(oid))
	params.Set("type", "1")
	params.Set("next", strconv.Itoa(page-1))
	params.Set("ps", "20")
	if sort == 2 {
		params.Set("mode", "2")
	} else {
		params.Set("mode", "1")
	}

	return c.get("/x/v2/reply/main", params)
}

// GetVideoDanmaku fetches danmaku (bullet comments) for a video segment.
func (c *Client) GetVideoDanmaku(cid int) ([]byte, error) {
	reqURL := "https://api.bilibili.com/x/v1/dm/list.so?oid=" + strconv.Itoa(cid)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = c.buildHeaders()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// PostDanmaku sends a danmaku (bullet comment) on a video.
func (c *Client) PostDanmaku(oid int, cid int, message string, progress int) error {
	params := url.Values{}
	params.Set("oid", strconv.Itoa(oid))
	params.Set("type", "1")
	params.Set("message", message)
	params.Set("progress", strconv.Itoa(progress))
	params.Set("color", "16777215") // white
	params.Set("fontsize", "25")

	_, err := c.post("/x/v2/dm/post", params)
	return err
}

// AddFavorite adds a video to a favorite folder.
func (c *Client) AddFavorite(bvid string, fid int) error {
	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("add_media_ids", strconv.Itoa(fid))

	_, err := c.post("/x/v3/fav/resource/deal", params)
	return err
}

// DelFavorite removes a video from a favorite folder.
func (c *Client) DelFavorite(bvid string, fid int) error {
	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("del_media_ids", strconv.Itoa(fid))

	_, err := c.post("/x/v3/fav/resource/deal", params)
	return err
}

// GetRecommendVideos fetches recommended videos from homepage.
func (c *Client) GetRecommendVideos(fresh bool, page int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("fresh_type", "4")
	params.Set("fresh_idx", strconv.Itoa(page))
	params.Set("ps", "20")
	params.Set("y_num", "3")

	var lastErr error
	for i := 0; i < 3; i++ {
		data, err := c.wbiGet("/x/web-interface/index/top/rcmd", params)
		if err == nil {
			return data, nil
		}
		lastErr = err
		// Check if rate limited
		if strings.Contains(err.Error(), "rate_limited") {
			time.Sleep(time.Duration(2*(i+1)) * time.Second)
			continue
		}
		return nil, err
	}
	return nil, lastErr
}

// GetVideoTags fetches tags of a video.
func (c *Client) GetVideoTags(bvid string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("bvid", bvid)

	return c.get("/x/web-interface/view/detail/tag", params)
}

// SearchSuggest fetches search suggestions.
func (c *Client) SearchSuggest(keyword string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("term", keyword)

	return c.get("/x/web-interface/search/default", params)
}

// GetFansList fetches a user's fans list.
func (c *Client) GetFansList(uid int, page int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("vmid", strconv.Itoa(uid))
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", "50")

	return c.get("/x/relation/fans", params)
}

// GetUserDynamics fetches a user's dynamics list from their space.
func (c *Client) GetUserDynamics(uid int, page int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("host_uid", strconv.Itoa(uid))
	params.Set("offset_dynamic_id", "0")
	if page > 1 {
		params.Set("offset_dynamic_id", strconv.Itoa(page))
	}
	params.Set("features", "itemOpusStyle,listOnlyfans,opusBigCover,onlyfansVote,decorationCard,onlyfansAssetsV2,forwardListHidden,ugcDelete")

	data, err := c.wbiGet("/x/polymer/web-dynamic/v1/feed/space", params)
	if err == nil {
		return data, nil
	}

	// If rate limited or wbi error, try old API as fallback
	if err != nil {
		// Try the old API endpoint (no WBI signing required)
		oldParams := url.Values{}
		oldParams.Set("host_uid", strconv.Itoa(uid))
		oldParams.Set("offset_dynamic_id", "0")
		if page > 1 {
			oldParams.Set("offset_dynamic_id", strconv.Itoa(page))
		}
		oldParams.Set("need_top", "0")

		oldData, oldErr := c.getWithReferer("/dynamic_svr/v1/dynamic_svr/space_history", oldParams,
			fmt.Sprintf("https://space.bilibili.com/%d/dynamic", uid))
		if oldErr == nil {
			return oldData, nil
		}
	}

	return data, err
}

// GetUserCollections fetches a user's series/collection list.
func (c *Client) GetUserCollections(uid int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("mid", strconv.Itoa(uid))
	params.Set("page_num", "1")
	params.Set("page_size", "30")
	params.Set("web_location", "333.999")
	params.Set("sort_reverse", "false")

	referer := fmt.Sprintf("https://space.bilibili.com/%d", uid)

	// Try primary endpoint first
	data, err := c.wbiGetWithReferer("/x/polymer/web-space/home/seasons_series", params, referer)
	if err == nil {
		return data, nil
	}

	// Fallback 1: try the old endpoint (may not need referer)
	oldData, oldErr := c.wbiGet("/x/polymer/space/seasons_series_list", params)
	if oldErr == nil {
		return oldData, nil
	}

	// If both fail, return the primary error
	return data, err
}

// PostDynamic publishes a plain text dynamic.
// Requires logged-in credential with write permission.
func (c *Client) PostDynamic(text string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("content", text)
	return c.post("/x/polymer/web-dynamic/v1/opus/schedule", params)
}

// DeleteDynamic deletes a dynamic by its ID.
// Requires logged-in credential with write permission.
func (c *Client) DeleteDynamic(dynamicID int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("dynamic_id", strconv.Itoa(dynamicID))
	return c.post("/x/polymer/web-dynamic/v1/retcode", params)
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