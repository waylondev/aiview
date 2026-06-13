package douyin

import (
	"testing"
)

type mockClient struct{}

func (m *mockClient) PlatformName() string                          { return "douyin" }
func (m *mockClient) GetHotSearch() (map[string]interface{}, error) { return nil, nil }
func (m *mockClient) GetTrending() (map[string]interface{}, error)  { return nil, nil }
func (m *mockClient) Search(keyword string, page int, count int) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) GetVideoDetail(videoID string) (map[string]interface{}, error) { return nil, nil }
func (m *mockClient) GetVideoComments(videoID string, cursor int) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) GetUserPosts(uid string, cursor int) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockClient) GetUserInfo(uid string) (map[string]interface{}, error) { return nil, nil }

func TestNewHotCmd_CreatesCommand(t *testing.T) {
	cmd := NewHotCmd(&mockClient{}, func() bool { return true })
	if cmd == nil {
		t.Fatal("NewHotCmd returned nil")
	}
	if cmd.Use != "hot" {
		t.Errorf("expected 'hot', got '%s'", cmd.Use)
	}
}

func TestNewTrendingCmd_CreatesCommand(t *testing.T) {
	cmd := NewTrendingCmd(&mockClient{}, func() bool { return true })
	if cmd == nil {
		t.Fatal("NewTrendingCmd returned nil")
	}
	if cmd.Use != "trending" {
		t.Errorf("expected 'trending', got '%s'", cmd.Use)
	}
}