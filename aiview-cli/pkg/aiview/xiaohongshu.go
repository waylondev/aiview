package aiview

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/xiaohongshu"
)

// XiaohongshuClient wraps the underlying xiaohongshu client.
type XiaohongshuClient struct {
	client platform.Client
}

// GetHotNotes fetches hot notes from Xiaohongshu.
func (x *XiaohongshuClient) GetHotNotes() ([]HotItem, error) {
	xhsClient := x.client.(*xiaohongshu.Client)
	data, err := xhsClient.GetHotNotes()
	if err != nil {
		return nil, err
	}

	result := make([]HotItem, 0)
	dataMap := helper.GetMap(data, "data")
	items := helper.GetSlice(dataMap, "items")
	for i, item := range items {
		m := item.(map[string]interface{})
		result = append(result, HotItem{
			Keyword:  helper.GetString(m, "title"),
			HotValue: helper.GetInt(m, "hot_value"),
			Position: i + 1,
			URL:      helper.GetString(m, "url"),
		})
	}
	return result, nil
}

// SearchNotes searches notes on Xiaohongshu.
func (x *XiaohongshuClient) SearchNotes(keyword string, page int) ([]SearchItem, error) {
	xhsClient := x.client.(*xiaohongshu.Client)
	data, err := xhsClient.SearchNotes(keyword, page)
	if err != nil {
		return nil, err
	}

	result := make([]SearchItem, 0)
	dataMap := helper.GetMap(data, "data")
	items := helper.GetSlice(dataMap, "items")
	for _, item := range items {
		m := item.(map[string]interface{})
		note := helper.GetMap(m, "note_card")
		if note == nil {
			continue
		}
		user := helper.GetMap(note, "user")
		result = append(result, SearchItem{
			ID:     helper.GetString(m, "id"),
			Title:  helper.GetString(note, "title"),
			Author: helper.GetString(user, "nickname"),
			URL:    fmt.Sprintf("https://www.xiaohongshu.com/explore/%s", helper.GetString(m, "id")),
		})
	}
	return result, nil
}

// GetNoteDetail fetches note details from Xiaohongshu.
func (x *XiaohongshuClient) GetNoteDetail(noteID string) (*VideoInfo, error) {
	xhsClient := x.client.(*xiaohongshu.Client)
	data, err := xhsClient.GetNoteDetail(noteID)
	if err != nil {
		return nil, err
	}

	dataMap := helper.GetMap(data, "data")
	items := helper.GetSlice(dataMap, "items")
	if len(items) == 0 {
		return nil, fmt.Errorf("note not found")
	}

	item := items[0].(map[string]interface{})
	note := helper.GetMap(item, "note_card")
	user := helper.GetMap(note, "user")
	interact := helper.GetMap(note, "interact_info")

	return &VideoInfo{
		ID:       noteID,
		Title:    helper.GetString(note, "title"),
		Author:   helper.GetString(user, "nickname"),
		AuthorID: helper.GetString(user, "user_id"),
		Play:     0,
		Like:     helper.GetInt(interact, "liked_count"),
		Coin:     0,
		Favorite: helper.GetInt(interact, "collected_count"),
		Share:    helper.GetInt(interact, "share_count"),
		Duration: "",
		URL:      fmt.Sprintf("https://www.xiaohongshu.com/explore/%s", noteID),
	}, nil
}

// GetUserInfo fetches user profile from Xiaohongshu.
func (x *XiaohongshuClient) GetUserInfo(userID string) (*UserInfo, error) {
	xhsClient := x.client.(*xiaohongshu.Client)
	data, err := xhsClient.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	dataMap := helper.GetMap(data, "data")
	user := helper.GetMap(dataMap, "user")
	return &UserInfo{
		ID:     helper.GetString(user, "user_id"),
		Name:   helper.GetString(user, "nickname"),
		Avatar: helper.GetString(user, "avatar"),
		Sign:   helper.GetString(user, "desc"),
		Fans:   helper.GetInt(user, "fans"),
		Follow: helper.GetInt(user, "follows"),
		Videos: helper.GetInt(user, "notes"),
	}, nil
}
