package aiview

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/kuaishou"
)

// KuaishouClient wraps the underlying kuaishou client.
type KuaishouClient struct {
	client platform.Client
}

// GetHotSearch fetches hot search items from Kuaishou.
func (k *KuaishouClient) GetHotSearch() ([]HotItem, error) {
	ksClient := k.client.(*kuaishou.Client)
	data, err := ksClient.GetHotSearch()
	if err != nil {
		return nil, err
	}

	result := make([]HotItem, 0)
	dataMap := helper.GetMap(data, "data")
	visionHotRank := helper.GetMap(dataMap, "visionHotRank")
	items := helper.GetSlice(visionHotRank, "items")
	for i, item := range items {
		m := item.(map[string]interface{})
		result = append(result, HotItem{
			Keyword:  helper.GetString(m, "name"),
			HotValue: helper.GetInt(m, "hotValue"),
			Position: i + 1,
			URL:      fmt.Sprintf("https://www.kuaishou.com/search/video?searchKey=%s", helper.GetString(m, "name")),
		})
	}
	return result, nil
}

// Search searches videos on Kuaishou.
func (k *KuaishouClient) Search(keyword string, page int) ([]SearchItem, error) {
	ksClient := k.client.(*kuaishou.Client)
	data, err := ksClient.Search(keyword, page)
	if err != nil {
		return nil, err
	}

	result := make([]SearchItem, 0)
	dataMap := helper.GetMap(data, "data")
	visionSearchPhoto := helper.GetMap(dataMap, "visionSearchPhoto")
	feeds := helper.GetSlice(visionSearchPhoto, "feeds")
	for _, item := range feeds {
		m := item.(map[string]interface{})
		result = append(result, SearchItem{
			ID:     helper.GetString(m, "id"),
			Title:  helper.GetString(m, "caption"),
			Author: "",
			URL:    fmt.Sprintf("https://www.kuaishou.com/short-video/%s", helper.GetString(m, "id")),
		})
	}
	return result, nil
}

// GetUserInfo fetches user profile from Kuaishou.
func (k *KuaishouClient) GetUserInfo(uid string) (*UserInfo, error) {
	ksClient := k.client.(*kuaishou.Client)
	data, err := ksClient.GetUserInfo(uid)
	if err != nil {
		return nil, err
	}

	dataMap := helper.GetMap(data, "data")
	visionProfile := helper.GetMap(dataMap, "visionProfile")
	user := helper.GetMap(visionProfile, "user")
	return &UserInfo{
		ID:     helper.GetString(user, "id"),
		Name:   helper.GetString(user, "name"),
		Avatar: "",
		Sign:   "",
		Fans:   helper.GetInt(user, "followerCount"),
		Follow: 0,
		Videos: 0,
	}, nil
}
