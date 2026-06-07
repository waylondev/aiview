package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewVideoCmd creates the video command.
func NewVideoCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "video <share_url>",
		Short: "View Douyin video details",
		Long: `View details of a Douyin video by its share URL.

Requires login cookie for full access.

Examples:
  aiview douyin video https://www.douyin.com/video/123456789`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := GetOutputFormat(cmd)
			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{
					"note":      "Video details require login. Use --cookie to authenticate.",
					"share_url": args[0],
				}, format)
			}
			fmt.Printf("🎬 Douyin Video\n\n")
			fmt.Printf("  URL: %s\n", args[0])
			fmt.Printf("  Note: Video details require login\n")
			fmt.Printf("  Use: aiview douyin login --cookie <cookie>\n")
			return nil
		},
	}
}