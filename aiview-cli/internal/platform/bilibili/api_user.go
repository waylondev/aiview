package bilibili

import (
	"fmt"
	"net/url"
	"strconv"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/helper"
)

// GetUserInfoCard fetches user profile information and returns a typed struct.
func (c *Client) GetUserInfoCard(uid int) (*UserInfo, error) {
	params := url.Values{}
	params.Set("mid", strconv.Itoa(uid))
	params.Set("photo", "false")

	data, err := c.wbiGet("/x/web-interface/card", params)
	if err != nil {
		return nil, err
	}

	d := helper.GetMap(data, "data")
	card := helper.GetMap(d, "card")
	return &UserInfo{
		MID:       helper.GetInt(card, "mid"),
		Name:      helper.GetString(card, "name"),
		Level:     helper.GetInt(card, "level_info", "current_level"),
		Coins:     helper.GetInt(card, "coins"),
		Sign:      helper.GetString(card, "sign"),
		Fans:      helper.GetInt(card, "fans"),
		Following: helper.GetInt(card, "attention"),
	}, nil
}

// GetUserInfo implements platform.UserQueryable.
func (c *Client) GetUserInfo(uid string) (map[string]interface{}, error) {
	id, err := strconv.Atoi(uid)
	if err != nil {
		return nil, aiverr.InvalidInput("bilibili", fmt.Sprintf("invalid uid %q: %v", uid, err))
	}
	info, err := c.GetUserInfoCard(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"mid":       info.MID,
		"name":      info.Name,
		"level":     info.Level,
		"coins":     info.Coins,
		"sign":      info.Sign,
		"fans":      info.Fans,
		"following": info.Following,
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

	d := helper.GetMap(data, "data")
	list := helper.GetMap(d, "list")
	vlist := helper.GetSlice(list, "vlist")

	videos := make([]VideoInfo, 0, len(vlist))
	for _, v := range vlist {
		m := v.(map[string]interface{})
		videos = append(videos, VideoInfo{
			BVID:        helper.GetString(m, "bvid"),
			Title:       helper.GetString(m, "title"),
			Duration:    helper.GetInt(m, "length"),
			DurationStr: helper.FormatDuration(helper.GetInt(m, "length")),
			URL:         fmt.Sprintf("https://www.bilibili.com/video/%s", helper.GetString(m, "bvid")),
			Stats: VideoStats{
				View: helper.GetInt(m, "play"),
			},
		})
		if len(videos) >= count {
			break
		}
	}
	return videos, nil
}

// GetSelfInfo fetches the logged-in user's own info.
func (c *Client) GetSelfInfo() (*UserInfo, error) {
	data, err := c.get("/x/web-interface/nav", nil)
	if err != nil {
		return nil, err
	}

	d := helper.GetMap(data, "data")
	if !helper.GetBool(d, "isLogin") {
		return nil, aiverr.NotAuthenticated("bilibili", "not logged in")
	}

	return &UserInfo{
		MID:   helper.GetInt(d, "mid"),
		Name:  helper.GetString(d, "uname"),
		Level: helper.GetInt(d, "level_info", "current_level"),
		Coins: helper.GetInt(d, "money"),
	}, nil
}

// GetFansList fetches a user's fans list.
func (c *Client) GetFansList(uid int, page int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("vmid", strconv.Itoa(uid))
	params.Set("pn", strconv.Itoa(page))
	params.Set("ps", "50")

	return c.get("/x/relation/fans", params)
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

	d := helper.GetMap(data, "data")
	list := helper.GetSlice(d, "list")

	users := make([]FollowingUser, 0, len(list))
	for _, u := range list {
		m := u.(map[string]interface{})
		users = append(users, FollowingUser{
			MID:  helper.GetInt(m, "mid"),
			Name: helper.GetString(m, "uname"),
			Sign: helper.GetString(m, "sign"),
		})
	}
	return users, nil
}

// GetRelationStat fetches the relation status between the current user and another user.
func (c *Client) GetRelationStat(uid int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("vmid", strconv.Itoa(uid))

	return c.get("/x/relation/stat", params)
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

	return c.wbiGet("/x/polymer/web-dynamic/v1/feed/space", params)
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
	return c.wbiGetWithReferer("/x/polymer/web-space/home/seasons_series", params, referer)
}
