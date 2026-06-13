// Package xiaohongshu provides CLI commands for the Xiaohongshu (小红书/RED) platform.
package xiaohongshu

// Client defines the Xiaohongshu API client interface.
type Client interface {
	GetHotNotes() (map[string]interface{}, error)
	SearchNotes(keyword string, page int) (map[string]interface{}, error)
	GetNoteDetail(noteID string) (map[string]interface{}, error)
	GetUserInfo(userID string) (map[string]interface{}, error)
}
