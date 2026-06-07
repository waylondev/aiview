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
			client := getClient()
			format := GetOutputFormat(cmd)

			result, err := client.GetVideoInfo(args[0])
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get video info: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			fmt.Printf("🎬 Douyin Video\n\n")
			fmt.Printf("  URL: %s\n", args[0])
			if note, ok := result["note"].(string); ok && note != "" {
				fmt.Printf("  Note: %s\n", note)
			}
			fmt.Printf("  Use: aiview douyin login --cookie <cookie>\n")
			return nil
		},
	}
}