package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/helper"
)

// GetVideoInfo fetches video metadata.
func (c *Client) GetVideoInfo(bvid string) (*VideoInfo, error) {
	params := url.Values{}
	params.Set("bvid", bvid)

	data, err := c.get("/x/web-interface/view", params)
	if err != nil {
		return nil, err
	}

	info := helper.GetMap(data, "data")
	owner := helper.GetMap(info, "owner")
	stat := helper.GetMap(info, "stat")

	return &VideoInfo{
		BVID:        bvid,
		AID:         helper.GetInt(info, "aid"),
		CID:         helper.GetInt(info, "cid"),
		Title:       helper.GetString(info, "title"),
		Description: strings.TrimSpace(helper.GetString(info, "desc")),
		Duration:    helper.GetInt(info, "duration"),
		DurationStr: helper.FormatDuration(helper.GetInt(info, "duration")),
		URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", bvid),
		Owner: OwnerInfo{
			MID:  helper.GetInt(owner, "mid"),
			Name: helper.GetString(owner, "name"),
		},
		Stats: VideoStats{
			View:     helper.GetInt(stat, "view"),
			Danmaku:  helper.GetInt(stat, "danmaku"),
			Like:     helper.GetInt(stat, "like"),
			Coin:     helper.GetInt(stat, "coin"),
			Favorite: helper.GetInt(stat, "favorite"),
			Share:    helper.GetInt(stat, "share"),
		},
	}, nil
}

// GetVideoSubtitle fetches video subtitle content.
func (c *Client) GetVideoSubtitle(bvid string) (*SubtitleInfo, error) {
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
	data := helper.GetMap(playerData, "data")
	subtitle := helper.GetMap(data, "subtitle")
	subtitles := helper.GetSlice(subtitle, "subtitles")

	if len(subtitles) == 0 {
		return &SubtitleInfo{Available: false}, nil
	}

	var subtitleURL string
	for _, sub := range subtitles {
		s := sub.(map[string]interface{})
		lan := helper.GetString(s, "lan")
		if strings.Contains(strings.ToLower(lan), "zh") {
			subtitleURL = helper.GetString(s, "subtitle_url")
			break
		}
	}
	if subtitleURL == "" {
		s := subtitles[0].(map[string]interface{})
		subtitleURL = helper.GetString(s, "subtitle_url")
	}

	if subtitleURL == "" {
		return &SubtitleInfo{Available: false}, nil
	}

	if strings.HasPrefix(subtitleURL, "//") {
		subtitleURL = "https:" + subtitleURL
	}

	req, err := http.NewRequest("GET", subtitleURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.HTTPClient.Do(req)
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

	rawItems := helper.GetSlice(subData, "body")
	items := make([]SubtitleItem, 0, len(rawItems))
	var texts []string
	for _, item := range rawItems {
		m := item.(map[string]interface{})
		content := helper.GetString(m, "content")
		texts = append(texts, content)
		items = append(items, SubtitleItem{
			From:    helper.GetFloat(m, "from"),
			To:      helper.GetFloat(m, "to"),
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

	d := helper.GetMap(data, "data")
	ai := helper.GetMap(d, "ai_summary")
	return helper.GetString(ai, "summary"), nil
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

	d := helper.GetMap(data, "data")
	replies := helper.GetSlice(d, "replies")

	comments := make([]CommentInfo, 0, len(replies))
	for _, r := range replies {
		m := r.(map[string]interface{})
		member := helper.GetMap(m, "member")
		content := helper.GetMap(m, "content")
		comments = append(comments, CommentInfo{
			ID: helper.GetString(m, "rpid_str"),
			Author: AuthorInfo{
				MID:  helper.GetInt(member, "mid"),
				Name: helper.GetString(member, "uname"),
			},
			Message: helper.GetString(content, "message"),
			Like:    helper.GetInt(m, "like"),
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

	items := helper.GetSlice(data, "data")
	videos := make([]VideoInfo, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		owner := helper.GetMap(m, "owner")
		stat := helper.GetMap(m, "stat")
		videos = append(videos, VideoInfo{
			BVID:        helper.GetString(m, "bvid"),
			Title:       helper.GetString(m, "title"),
			Duration:    helper.GetInt(m, "duration"),
			DurationStr: helper.FormatDuration(helper.GetInt(m, "duration")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", helper.GetString(m, "bvid")),
			Owner: OwnerInfo{
				MID:  helper.GetInt(owner, "mid"),
				Name: helper.GetString(owner, "name"),
			},
			Stats: VideoStats{
				View: helper.GetInt(stat, "view"),
			},
		})
	}
	return videos, nil
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
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("failed to create request: %v", err))
	}
	req.Header = c.BuildHeaders()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("network request failed: %v", err))
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

	d := helper.GetMap(data, "data")
	dash := helper.GetMap(d, "dash")
	audio := helper.GetSlice(dash, "audio")

	if len(audio) == 0 {
		return "", aiverr.NotFound("bilibili", "unable to get audio stream")
	}

	best := audio[0].(map[string]interface{})
	return helper.GetString(best, "baseUrl"), nil
}

// GetVideoTags fetches tags of a video.
func (c *Client) GetVideoTags(bvid string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("bvid", bvid)

	return c.get("/x/web-interface/view/detail/tag", params)
}

// GetVideoOnlineCount fetches the real-time online viewer count for a video.
func (c *Client) GetVideoOnlineCount(bvid string) (map[string]interface{}, error) {
	info, err := c.GetVideoInfo(bvid)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("cid", strconv.Itoa(info.CID))
	return c.get("/x/player/online/total", params)
}
