package bilibili

import (
	"net/url"
	"strconv"

	"github.com/jackwener/aiview/internal/helper"
)

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

// GetFavoriteList fetches favorite folders.
func (c *Client) GetFavoriteList(uid int, page int) ([]FavoriteFolder, error) {
	params := url.Values{}
	params.Set("up_mid", strconv.Itoa(uid))

	data, err := c.get("/x/v3/fav/folder/created/list-all", params)
	if err != nil {
		return nil, err
	}

	d := helper.GetMap(data, "data")
	list := helper.GetSlice(d, "list")

	folders := make([]FavoriteFolder, 0, len(list))
	for _, f := range list {
		m := f.(map[string]interface{})
		folders = append(folders, FavoriteFolder{
			ID:         helper.GetInt(m, "id"),
			Title:      helper.GetString(m, "title"),
			MediaCount: helper.GetInt(m, "media_count"),
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

	d := helper.GetMap(data, "data")
	medias := helper.GetSlice(d, "medias")

	items := make([]FavoriteMedia, 0, len(medias))
	for _, m := range medias {
		item := m.(map[string]interface{})
		upper := helper.GetMap(item, "upper")
		items = append(items, FavoriteMedia{
			BVID:     helper.GetString(item, "bvid"),
			Title:    helper.GetString(item, "title"),
			Duration: helper.FormatDuration(helper.GetInt(item, "duration")),
			Upper:    helper.GetString(upper, "name"),
		})
	}
	return items, nil
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

	d := helper.GetMap(data, "data")
	list := helper.GetSlice(d, "list")

	items := make([]HistoryItem, 0, len(list))
	for _, h := range list {
		m := h.(map[string]interface{})
		items = append(items, HistoryItem{
			BVID:   helper.GetString(m, "bvid"),
			Title:  helper.GetString(m, "title"),
			Author: helper.GetString(m, "author_name"),
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

	d := helper.GetMap(data, "data")
	list := helper.GetSlice(d, "list")

	items := make([]WatchLaterItem, 0, len(list))
	for _, w := range list {
		m := w.(map[string]interface{})
		owner := helper.GetMap(m, "owner")
		items = append(items, WatchLaterItem{
			BVID:     helper.GetString(m, "bvid"),
			Title:    helper.GetString(m, "title"),
			Author:   helper.GetString(owner, "name"),
			Duration: helper.FormatDuration(helper.GetInt(m, "duration")),
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

	d := helper.GetMap(data, "data")
	items := helper.GetSlice(d, "items")

	dynamics := make([]DynamicItem, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		modules := helper.GetMap(m, "modules")
		author := helper.GetMap(modules, "module_author")
		desc := helper.GetMap(helper.GetMap(modules, "module_dynamic"), "desc")

		dynamics = append(dynamics, DynamicItem{
			ID:     helper.GetString(m, "id_str"),
			Author: helper.GetString(author, "name"),
			Text:   helper.GetString(desc, "text"),
		})
	}
	return dynamics, nil
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

// GetLiveRoomInfo fetches live room information by room ID.
func (c *Client) GetLiveRoomInfo(roomID int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("room_id", strconv.Itoa(roomID))

	return c.get("/room/v1/Room/get_info", params)
}
