package bilibili

import (
	"fmt"
	"strconv"

	biliapi "github.com/jackwener/aiview/internal/platform/bilibili/bilibilitypes"
	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewUserCmd creates the user command.
func NewUserCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "user <UID>",
		Short: "View user info",
		Long:  `View Bilibili user information.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID must be a number", format)
				return err
			}

			info, err := client.GetUserInfo(uid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get user info: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(info, format)
			}

			fmt.Printf("👤 %s (UID: %d)\n", info.Name, info.MID)
			fmt.Printf("  Level: %d\n", info.Level)
			fmt.Printf("  Coins: %d\n", info.Coins)
			if info.Sign != "" {
				fmt.Printf("  Sign: %s\n", info.Sign)
			}
			return nil
		},
	}
}

// NewUserVideosCmd creates the user-videos command.
func NewUserVideosCmd(getClient func() Client) *cobra.Command {
	var maxResults int
	var order string
	var tid int
	var keyword string

	cmd := &cobra.Command{
		Use:   "user-videos <UID>",
		Short: "View user video list",
		Long:  `View a Bilibili user's video list.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID must be a number", format)
				return err
			}

			videos, err := client.GetUserVideos(uid, maxResults, order, tid, keyword)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get user videos: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				type result struct {
					Items []biliapi.VideoInfo `json:"items"`
				}
				return output.EmitSuccess(result{Items: videos}, format)
			}

			fmt.Printf("📹 User Videos (UID: %d):\n\n", uid)
			for i, v := range videos {
				fmt.Printf("  %d. [%s] %s\n", i+1, v.BVID, v.Title)
				fmt.Printf("     Views: %s  Duration: %s\n\n", output.FormatCount(v.Stats.View), v.DurationStr)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "n", 10, "Max results")
	cmd.Flags().StringVar(&order, "order", "pubdate", "Sort order: pubdate/click/stow")
	cmd.Flags().IntVar(&tid, "tid", 0, "Category ID filter")
	cmd.Flags().StringVar(&keyword, "keyword", "", "Search keyword within user videos")
	return cmd
}