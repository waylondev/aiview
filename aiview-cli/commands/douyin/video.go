package douyin

import (
	"fmt"
	"regexp"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// extractVideoID extracts the video ID from a Douyin share URL or plain ID.
func extractVideoID(input string) (string, error) {
	// If it's already a plain number, use it directly
	if matched, _ := regexp.MatchString(`^\d+$`, input); matched {
		return input, nil
	}

	// Try to extract from URL pattern /video/<digits>
	re := regexp.MustCompile(`/video/(\d+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("unable to extract video ID from '%s'. Expected a plain video ID or a URL like https://www.douyin.com/video/123456789", input)
}

// NewVideoCmd creates the video command.
func NewVideoCmd(client Client) *cobra.Command {
	return &cobra.Command{
		Use:   "video <share_url_or_id>",
		Short: "View Douyin video details",
		Long: `View details of a Douyin video by its share URL or video ID.

Requires login cookie for full access.

Examples:
  aiview douyin video https://www.douyin.com/video/123456789
  aiview douyin video 123456789`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := GetOutputFormat(cmd)

			videoID, err := extractVideoID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			result, err := client.GetVideoDetail(videoID)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get video info: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			fmt.Printf("🎬 Douyin Video\n\n")
			fmt.Printf("  Video ID: %s\n", videoID)
			if note, ok := result["note"].(string); ok && note != "" {
				fmt.Printf("  Note: %s\n", note)
			}
			fmt.Printf("  Use: aiview douyin login --cookie <cookie>\n")
			return nil
		},
	}
}