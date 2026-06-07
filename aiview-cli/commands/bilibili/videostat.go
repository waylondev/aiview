package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewVideoStatCmd creates the video-status command.
func NewVideoStatCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "video-status <BV>",
		Short: "View video statistics",
		Long: `View video statistics including view, like, coin, favorite and share counts.

Examples:
  aiview bilibili video-status BV1xx411c7m9`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			info, err := client.GetVideoInfo(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get video info: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(info.Stats, format)
			}

			s := info.Stats
			fmt.Printf("📊 Video Statistics: %s\n\n", bvid)
			fmt.Printf("  Views:    %s\n", output.FormatCount(s.View))
			fmt.Printf("  Danmaku:  %s\n", output.FormatCount(s.Danmaku))
			fmt.Printf("  Likes:    %s\n", output.FormatCount(s.Like))
			fmt.Printf("  Coins:    %s\n", output.FormatCount(s.Coin))
			fmt.Printf("  Favs:     %s\n", output.FormatCount(s.Favorite))
			fmt.Printf("  Shares:   %s\n", output.FormatCount(s.Share))
			return nil
		},
	}
}