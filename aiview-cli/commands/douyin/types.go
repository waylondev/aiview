// Package douyin provides CLI commands for the Douyin (抖音) platform.
package douyin

import (
	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// Client is the interface that the Douyin API client must satisfy for commands.
type Client interface {
	PlatformName() string
	GetHotSearch() (map[string]interface{}, error)
	GetTrending() (map[string]interface{}, error)
	Search(keyword string, page int, count int) (map[string]interface{}, error)
	GetVideoDetail(videoID string) (map[string]interface{}, error)
	GetVideoComments(videoID string, cursor int) (map[string]interface{}, error)
	GetUserPosts(uid string, cursor int) (map[string]interface{}, error)
	GetUserInfo(uid string) (map[string]interface{}, error)
}

// Credential holds Douyin authentication data.
type Credential struct {
	Cookie string `json:"cookie"`
}

// GetOutputFormat extracts the output format from cobra command flags.
func GetOutputFormat(cmd *cobra.Command) output.Format {
	parent := cmd
	for parent.HasParent() {
		parent = parent.Parent()
	}
	asJSON, _ := parent.Flags().GetBool("json")
	asYAML, _ := parent.Flags().GetBool("yaml")
	return output.ResolveFormat(asJSON, asYAML)
}