package bilibili

import (
	"testing"

	biliapi "github.com/jackwener/aiview/internal/platform/bilibili/bilibilitypes"
)

// mockClient implements Client interface for testing.
type mockClient struct{}

func (m *mockClient) GetVideoInfo(bvid string) (*biliapi.VideoInfo, error) { return nil, nil }
func (m *mockClient) GetVideoSubtitle(bvid string) (*biliapi.SubtitleInfo, error) {
	return nil, nil
}
func (m *mockClient) GetVideoAIConclusion(bvid string) (string, error) { return "", nil }
func (m *mockClient) GetVideoComments(bvid string, page int) ([]biliapi.CommentInfo, error) {
	return nil, nil
}
func (m *mockClient) GetRelatedVideos(bvid string) ([]biliapi.VideoInfo, error) { return nil, nil }
func (m *mockClient) SearchVideo(keyword string, page int, order string, duration int, tid int) ([]biliapi.SearchVideoResult, error) {
	return nil, nil
}
func (m *mockClient) SearchUser(keyword string, page int) ([]biliapi.SearchUserResult, error) {
	return nil, nil
}
func (m *mockClient) GetUserInfo(uid int) (*biliapi.UserInfo, error) { return nil, nil }
func (m *mockClient) GetUserVideos(uid int, count int, order string, tid int, keyword string) ([]biliapi.VideoInfo, error) {
	return nil, nil
}
func (m *mockClient) GetHotVideos(page int, count int) ([]biliapi.VideoInfo, error) {
	return nil, nil
}
func (m *mockClient) GetRankVideos(rid int, day int, typeStr string) ([]biliapi.VideoInfo, error) {
	return nil, nil
}
func (m *mockClient) GetFavoriteList(uid int, page int) ([]biliapi.FavoriteFolder, error) {
	return nil, nil
}
func (m *mockClient) GetFavoriteVideos(favID int, page int) ([]biliapi.FavoriteMedia, error) {
	return nil, nil
}
func (m *mockClient) GetFollowingList(uid int, page int) ([]biliapi.FollowingUser, error) {
	return nil, nil
}
func (m *mockClient) GetWatchHistory(page int, count int) ([]biliapi.HistoryItem, error) {
	return nil, nil
}
func (m *mockClient) GetWatchLater() ([]biliapi.WatchLaterItem, error) { return nil, nil }
func (m *mockClient) GetDynamicFeed(offset string) ([]biliapi.DynamicItem, error) {
	return nil, nil
}
func (m *mockClient) LikeVideo(bvid string, undo bool) error    { return nil }
func (m *mockClient) CoinVideo(bvid string, num int) error      { return nil }
func (m *mockClient) TripleVideo(bvid string) error             { return nil }
func (m *mockClient) UnfollowUser(uid int) error                { return nil }
func (m *mockClient) GetSelfInfo() (*biliapi.UserInfo, error)   { return nil, nil }
func (m *mockClient) GetAudioURL(bvid string) (string, error)   { return "", nil }
func (m *mockClient) PostComment(oid int, message string, root int, parent int) error {
	return nil
}
func (m *mockClient) DeleteComment(oid int, rpid int) error { return nil }
func (m *mockClient) GetVideoCommentsRaw(oid int, page int, sort int) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) GetVideoDanmaku(cid int) ([]byte, error) { return nil, nil }
func (m *mockClient) PostDanmaku(oid int, cid int, message string, progress int) error {
	return nil
}
func (m *mockClient) AddFavorite(bvid string, fid int) error  { return nil }
func (m *mockClient) DelFavorite(bvid string, fid int) error  { return nil }
func (m *mockClient) GetRecommendVideos(fresh bool, page int) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) GetVideoTags(bvid string) (map[string]interface{}, error)     { return nil, nil }
func (m *mockClient) SearchSuggest(keyword string) (map[string]interface{}, error) { return nil, nil }
func (m *mockClient) GetFansList(uid int, page int) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) GetUserDynamics(uid int, page int) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) PostDynamic(text string) (map[string]interface{}, error)     { return nil, nil }
func (m *mockClient) DeleteDynamic(dynamicID int) (map[string]interface{}, error) { return nil, nil }
func (m *mockClient) GetUserCollections(uid int) (map[string]interface{}, error)  { return nil, nil }
func (m *mockClient) GetRelationStat(uid int) (map[string]interface{}, error)     { return nil, nil }
func (m *mockClient) GetRegionVideos(rid int, page int, count int, sort string) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) GetLiveRoomInfo(roomID int) (map[string]interface{}, error) { return nil, nil }
func (m *mockClient) GetPreciousVideos() (map[string]interface{}, error)         { return nil, nil }
func (m *mockClient) GetHotSearch(limit int) (map[string]interface{}, error)     { return nil, nil }
func (m *mockClient) GetVideoOnlineCount(bvid string) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) GetWeeklyHotVideos(number int) (map[string]interface{}, error) {
	return nil, nil
}

func TestNewHotCmd_CreatesCommand(t *testing.T) {
	cmd := NewHotCmd(func() Client { return &mockClient{} })
	if cmd == nil {
		t.Fatal("NewHotCmd returned nil")
	}
	if cmd.Use != "hot" {
		t.Errorf("expected 'hot', got '%s'", cmd.Use)
	}
}

func TestNewSearchCmd_CreatesCommand(t *testing.T) {
	cmd := NewSearchCmd(func() Client { return &mockClient{} })
	if cmd == nil {
		t.Fatal("NewSearchCmd returned nil")
	}
	if cmd.Use != "search <keyword>" {
		t.Errorf("expected 'search <keyword>', got '%s'", cmd.Use)
	}
}