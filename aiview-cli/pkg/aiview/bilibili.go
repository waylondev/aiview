package aiview

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/bilibili"
)

// BilibiliClient wraps the underlying bilibili client.
type BilibiliClient struct {
	client platform.Client
}

// GetHotVideos fetches popular/hot videos from Bilibili.
func (b *BilibiliClient) GetHotVideos(page, count int) ([]VideoInfo, error) {
	biliClient := b.client.(*bilibili.Client)
	videos, err := biliClient.GetHotVideos(page, count)
	if err != nil {
		return nil, err
	}

	result := make([]VideoInfo, 0, len(videos))
	for _, v := range videos {
		result = append(result, VideoInfo{
			ID:       v.BVID,
			Title:    v.Title,
			Author:   v.Owner.Name,
			AuthorID: strconv.Itoa(v.Owner.MID),
			Play:     v.Stats.View,
			Danmaku:  v.Stats.Danmaku,
			Like:     v.Stats.Like,
			Coin:     v.Stats.Coin,
			Favorite: v.Stats.Favorite,
			Share:    v.Stats.Share,
			Duration: v.DurationStr,
			URL:      v.URL,
		})
	}
	return result, nil
}

// GetVideoInfo fetches video metadata from Bilibili.
func (b *BilibiliClient) GetVideoInfo(bvid string) (*VideoInfo, error) {
	biliClient := b.client.(*bilibili.Client)
	v, err := biliClient.GetVideoInfo(bvid)
	if err != nil {
		return nil, err
	}

	return &VideoInfo{
		ID:       v.BVID,
		Title:    v.Title,
		Author:   v.Owner.Name,
		AuthorID: strconv.Itoa(v.Owner.MID),
		Play:     v.Stats.View,
		Danmaku:  v.Stats.Danmaku,
		Like:     v.Stats.Like,
		Coin:     v.Stats.Coin,
		Favorite: v.Stats.Favorite,
		Share:    v.Stats.Share,
		Duration: v.DurationStr,
		URL:      v.URL,
	}, nil
}

// GetUserInfo fetches user profile from Bilibili.
func (b *BilibiliClient) GetUserInfo(uid int) (*UserInfo, error) {
	biliClient := b.client.(*bilibili.Client)
	u, err := biliClient.GetUserInfo(uid)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:     strconv.Itoa(u.MID),
		Name:   u.Name,
		Sign:   u.Sign,
		Fans:   u.Fans,
		Follow: u.Following,
	}, nil
}

// SearchVideo searches videos by keyword.
func (b *BilibiliClient) SearchVideo(keyword string, page int) ([]SearchItem, error) {
	biliClient := b.client.(*bilibili.Client)
	videos, err := biliClient.SearchVideo(keyword, page, "", 0, 0)
	if err != nil {
		return nil, err
	}

	result := make([]SearchItem, 0, len(videos))
	for _, v := range videos {
		result = append(result, SearchItem{
			ID:     v.BVID,
			Title:  v.Title,
			Author: v.Author,
			URL:    fmt.Sprintf("https://www.bilibili.com/video/%s", v.BVID),
		})
	}
	return result, nil
}
