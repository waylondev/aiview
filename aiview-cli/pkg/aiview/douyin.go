package aiview

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"
	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/douyin"
)

// DouyinClient wraps the underlying douyin client.
type DouyinClient struct {
	client platform.Client
}

// GetHotSearch fetches hot search items from Douyin.
func (d *DouyinClient) GetHotSearch() ([]HotItem, error) {
	dyClient := d.client.(*douyin.Client)
	data, err := dyClient.GetHotSearch()
	if err != nil {
		return nil, err
	}

	result := make([]HotItem, 0)
	dataMap := helper.GetMap(data, "data")
	list := helper.GetSlice(dataMap, "word_list")
	for i, item := range list {
		m := item.(map[string]interface{})
		result = append(result, HotItem{
			Keyword:  helper.GetString(m, "word"),
			HotValue: helper.GetInt(m, "hot_value"),
			Position: i + 1,
			URL:      fmt.Sprintf("https://www.douyin.com/search/%s", helper.GetString(m, "word")),
		})
	}
	return result, nil
}

// Search searches videos on Douyin.
func (d *DouyinClient) Search(keyword string, page, count int) ([]SearchItem, error) {
	dyClient := d.client.(*douyin.Client)
	data, err := dyClient.Search(keyword, page, count)
	if err != nil {
		return nil, err
	}

	result := make([]SearchItem, 0)
	list := helper.GetSlice(data, "data")
	for _, item := range list {
		m := item.(map[string]interface{})
		aweme := helper.GetMap(m, "aweme_info")
		if aweme == nil {
			continue
		}
		author := helper.GetMap(aweme, "author")
		result = append(result, SearchItem{
			ID:     helper.GetString(aweme, "aweme_id"),
			Title:  helper.GetString(aweme, "desc"),
			Author: helper.GetString(author, "nickname"),
			URL:    fmt.Sprintf("https://www.douyin.com/video/%s", helper.GetString(aweme, "aweme_id")),
		})
	}
	return result, nil
}

// GetVideoDetail fetches video details from Douyin.
func (d *DouyinClient) GetVideoDetail(videoID string) (*VideoInfo, error) {
	dyClient := d.client.(*douyin.Client)
	data, err := dyClient.GetVideoDetail(videoID)
	if err != nil {
		return nil, err
	}

	aweme := helper.GetMap(data, "aweme_detail")
	if aweme == nil {
		return nil, aiverr.NotFound("douyin", "video not found")
	}

	author := helper.GetMap(aweme, "author")
	stats := helper.GetMap(aweme, "statistics")

	return &VideoInfo{
		ID:       helper.GetString(aweme, "aweme_id"),
		Title:    helper.GetString(aweme, "desc"),
		Author:   helper.GetString(author, "nickname"),
		AuthorID: helper.GetString(author, "uid"),
		Play:     helper.GetInt(stats, "play_count"),
		Like:     helper.GetInt(stats, "digg_count"),
		Danmaku:  0,
		Coin:     0,
		Favorite: helper.GetInt(stats, "collect_count"),
		Share:    helper.GetInt(stats, "share_count"),
		Duration: "",
		URL:      fmt.Sprintf("https://www.douyin.com/video/%s", videoID),
	}, nil
}

// GetUserInfo fetches user profile from Douyin.
func (d *DouyinClient) GetUserInfo(uid string) (*UserInfo, error) {
	dyClient := d.client.(*douyin.Client)
	data, err := dyClient.GetUserInfo(uid)
	if err != nil {
		return nil, err
	}

	user := helper.GetMap(data, "user")
	return &UserInfo{
		ID:     helper.GetString(user, "uid"),
		Name:   helper.GetString(user, "nickname"),
		Avatar: helper.GetString(user, "avatar_larger", "url_list", "0"),
		Sign:   helper.GetString(user, "signature"),
		Fans:   helper.GetInt(user, "follower_count"),
		Follow: helper.GetInt(user, "following_count"),
		Videos: helper.GetInt(user, "aweme_count"),
	}, nil
}
