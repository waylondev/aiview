package bilibili

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/helper"
)

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

	d := helper.GetMap(data, "data")
	items := helper.GetSlice(d, "result")

	results := make([]SearchVideoResult, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		results = append(results, SearchVideoResult{
			BVID:     helper.GetString(m, "bvid"),
			Title:    helper.StripHTML(helper.GetString(m, "title")),
			Author:   helper.GetString(m, "author"),
			Play:     helper.GetInt(m, "play"),
			Duration: helper.GetString(m, "duration"),
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

	d := helper.GetMap(data, "data")
	items := helper.GetSlice(d, "result")

	results := make([]SearchUserResult, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		results = append(results, SearchUserResult{
			MID:    helper.GetInt(m, "mid"),
			Name:   helper.GetString(m, "uname"),
			Sign:   helper.GetString(m, "usign"),
			Fans:   helper.GetInt(m, "fans"),
			Videos: helper.GetInt(m, "videos"),
		})
	}
	return results, nil
}

// SearchSuggest fetches search suggestions.
func (c *Client) SearchSuggest(keyword string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("term", keyword)

	return c.get("/x/web-interface/search/default", params)
}

// Search implements platform.Searchable by delegating to SearchVideo.
func (c *Client) Search(query string, page int, count ...int) (map[string]interface{}, error) {
	videos, err := c.SearchVideo(query, page, "", 0, 0)
	if err != nil {
		return nil, err
	}
	// Convert to map for interface compliance
	items := make([]interface{}, len(videos))
	for i, v := range videos {
		items[i] = map[string]interface{}{
			"bvid":   v.BVID,
			"title":  v.Title,
			"author": v.Author,
			"play":   v.Play,
		}
	}
	return map[string]interface{}{"data": items}, nil
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

	d := helper.GetMap(data, "data")
	list := helper.GetSlice(d, "list")

	videos := make([]VideoInfo, 0, len(list))
	for _, v := range list {
		m := v.(map[string]interface{})
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

// GetRankVideos fetches ranking videos.
func (c *Client) GetRankVideos(rid int, day int, typeStr string) ([]VideoInfo, error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("type", "all")

	data, err := c.wbiGet("/x/web-interface/ranking/v2", params)
	if err != nil {
		return nil, err
	}

	d := helper.GetMap(data, "data")
	list := helper.GetSlice(d, "list")

	videos := make([]VideoInfo, 0, len(list))
	for _, v := range list {
		m := v.(map[string]interface{})
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

// GetRecommendVideos fetches recommended videos from homepage.
func (c *Client) GetRecommendVideos(fresh bool, page int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("fresh_type", "4")
	params.Set("fresh_idx", strconv.Itoa(page))
	params.Set("ps", "20")
	params.Set("y_num", "3")

	var lastErr error
	for i := 0; i < 3; i++ {
		data, err := c.wbiGet("/x/web-interface/wbi/index/top/feed/rcmd", params)
		if err == nil {
			return data, nil
		}
		lastErr = err
		if pe, ok := aiverr.IsPlatformError(err); ok && pe.Code == aiverr.CodeRateLimited {
			time.Sleep(time.Duration(2*(i+1)) * time.Second)
			continue
		}
		return nil, err
	}
	return nil, lastErr
}

// GetHotSearch fetches trending/hot search keywords.
func (c *Client) GetHotSearch(count ...int) (map[string]interface{}, error) {
	limit := 50
	if len(count) > 0 && count[0] > 0 {
		limit = min(count[0], 50)
	}
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	return c.wbiGet("/x/web-interface/wbi/search/square", params)
}

// GetWeeklyHotVideos fetches the weekly hot video list by series number.
func (c *Client) GetWeeklyHotVideos(number int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("number", strconv.Itoa(number))
	return c.wbiGet("/x/web-interface/popular/series/one", params)
}

// GetPreciousVideos fetches the "入站必刷" (must-watch) curated video list.
func (c *Client) GetPreciousVideos() (map[string]interface{}, error) {
	return c.get("/x/web-interface/popular/precious", nil)
}

// GetRegionVideos fetches videos by region/category.
func (c *Client) GetRegionVideos(rid int, page int, count int, sort string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", strconv.Itoa(min(count, 50)))
	if sort != "" {
		params.Set("sort", sort)
	}

	return c.get("/x/web-interface/dynamic/region", params)
}
