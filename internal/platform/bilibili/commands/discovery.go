package commands

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewHotCmd creates the hot command.
func NewHotCmd(getClient func() Client) *cobra.Command {
	var maxResults int

	cmd := &cobra.Command{
		Use:   "hot",
		Short: "查看热门视频",
		Long:  `查看 Bilibili 热门视频。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			videos, err := client.GetHotVideos(1, maxResults)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取热门视频失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": videos}, format)
			}

			fmt.Println("🔥 热门视频:")
			for i, v := range videos {
				fmt.Printf("  %d. [%s] %s\n", i+1, v.BVID, v.Title)
				fmt.Printf("     UP主: %s  播放: %s\n\n", v.Owner.Name, output.FormatCount(v.Stats.View))
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "n", 20, "最大结果数")
	return cmd
}

// NewRankCmd creates the rank command.
func NewRankCmd(getClient func() Client) *cobra.Command {
	var maxResults int

	cmd := &cobra.Command{
		Use:   "rank",
		Short: "查看排行榜",
		Long:  `查看 Bilibili 排行榜。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			videos, err := client.GetRankVideos(3)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取排行榜失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": videos[:min(maxResults, len(videos))]}, format)
			}

			fmt.Println("🏆 排行榜:")
			for i, v := range videos[:min(maxResults, len(videos))] {
				fmt.Printf("  %d. [%s] %s\n", i+1, v.BVID, v.Title)
				fmt.Printf("     UP主: %s  播放: %s\n\n", v.Owner.Name, output.FormatCount(v.Stats.View))
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "n", 20, "最大结果数")
	return cmd
}

// NewFeedCmd creates the feed command.
func NewFeedCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "feed",
		Short: "查看动态时间线",
		Long:  `查看动态时间线（需要登录）。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录，请使用 aiview bilibili login 登录", format)
				return err
			}
			_ = cred

			items, err := client.GetDynamicFeed("")
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取动态失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": items}, format)
			}

			fmt.Println("📡 动态时间线:")
			for _, d := range items {
				fmt.Printf("  %s\n", d.Author)
				if d.Text != "" {
					text := d.Text
					if len(text) > 100 {
						text = text[:100]
					}
					fmt.Printf("  %s\n\n", text)
				}
			}
			return nil
		},
	}
}