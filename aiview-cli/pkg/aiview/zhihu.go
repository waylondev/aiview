package aiview

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/zhihu"
)

// ZhihuClient wraps the underlying zhihu client.
type ZhihuClient struct {
	client platform.Client
}

// GetHotSearch fetches hot search items from Zhihu.
func (z *ZhihuClient) GetHotSearch() ([]HotItem, error) {
	zhClient := z.client.(*zhihu.Client)
	data, err := zhClient.GetHotSearch()
	if err != nil {
		return nil, err
	}

	result := make([]HotItem, 0)
	dataMap := helper.GetMap(data, "data")
	list := helper.GetSlice(dataMap, "top_search")
	for i, item := range list {
		m := item.(map[string]interface{})
		result = append(result, HotItem{
			Keyword:  helper.GetString(m, "query"),
			HotValue: helper.GetInt(m, "hot_value"),
			Position: i + 1,
			URL:      fmt.Sprintf("https://www.zhihu.com/search?q=%s", helper.GetString(m, "query")),
		})
	}
	return result, nil
}

// Search searches content on Zhihu.
func (z *ZhihuClient) Search(keyword string, page int) ([]SearchItem, error) {
	zhClient := z.client.(*zhihu.Client)
	data, err := zhClient.Search(keyword, page)
	if err != nil {
		return nil, err
	}

	result := make([]SearchItem, 0)
	list := helper.GetSlice(data, "data")
	for _, item := range list {
		m := item.(map[string]interface{})
		result = append(result, SearchItem{
			ID:     helper.GetString(m, "id"),
			Title:  helper.GetString(m, "title"),
			Author: helper.GetString(helper.GetMap(m, "author"), "name"),
			URL:    helper.GetString(m, "url"),
		})
	}
	return result, nil
}

// GetUserInfo fetches user profile from Zhihu.
func (z *ZhihuClient) GetUserInfo(uid string) (*UserInfo, error) {
	zhClient := z.client.(*zhihu.Client)
	data, err := zhClient.GetUserInfo(uid)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:     helper.GetString(data, "id"),
		Name:   helper.GetString(data, "name"),
		Avatar: helper.GetString(data, "avatar_url"),
		Sign:   helper.GetString(data, "headline"),
		Fans:   helper.GetInt(data, "follower_count"),
		Follow: helper.GetInt(data, "following_count"),
		Videos: 0,
	}, nil
}
