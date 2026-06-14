package aiview

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/weibo"
)

// WeiboClient wraps the underlying weibo client.
type WeiboClient struct {
	client platform.Client
}

// GetHotSearch fetches hot search items from Weibo.
func (w *WeiboClient) GetHotSearch() ([]HotItem, error) {
	wbClient := w.client.(*weibo.Client)
	data, err := wbClient.GetHotSearch()
	if err != nil {
		return nil, err
	}

	result := make([]HotItem, 0)
	dataMap := helper.GetMap(data, "data")
	list := helper.GetSlice(dataMap, "realtime")
	for i, item := range list {
		m := item.(map[string]interface{})
		result = append(result, HotItem{
			Keyword:  helper.GetString(m, "word"),
			HotValue: helper.GetInt(m, "num"),
			Position: i + 1,
			URL:      fmt.Sprintf("https://s.weibo.com/weibo?q=%s", helper.GetString(m, "word")),
		})
	}
	return result, nil
}

// Search searches content on Weibo.
func (w *WeiboClient) Search(keyword string, page int) ([]SearchItem, error) {
	wbClient := w.client.(*weibo.Client)
	data, err := wbClient.Search(keyword, page)
	if err != nil {
		return nil, err
	}

	result := make([]SearchItem, 0)
	list := helper.GetSlice(data, "data")
	for _, item := range list {
		m := item.(map[string]interface{})
		result = append(result, SearchItem{
			ID:     helper.GetString(m, "id"),
			Title:  helper.GetString(m, "text"),
			Author: helper.GetString(helper.GetMap(m, "user"), "name"),
			URL:    fmt.Sprintf("https://weibo.com/%s", helper.GetString(m, "id")),
		})
	}
	return result, nil
}

// GetUserInfo fetches user profile from Weibo.
func (w *WeiboClient) GetUserInfo(uid string) (*UserInfo, error) {
	wbClient := w.client.(*weibo.Client)
	data, err := wbClient.GetUserInfo(uid)
	if err != nil {
		return nil, err
	}

	dataMap := helper.GetMap(data, "data")
	user := helper.GetMap(dataMap, "user")
	return &UserInfo{
		ID:     helper.GetString(user, "idstr"),
		Name:   helper.GetString(user, "screen_name"),
		Avatar: helper.GetString(user, "avatar_large"),
		Sign:   helper.GetString(user, "description"),
		Fans:   helper.GetInt(user, "followers_count"),
		Follow: helper.GetInt(user, "friends_count"),
		Videos: helper.GetInt(user, "statuses_count"),
	}, nil
}
