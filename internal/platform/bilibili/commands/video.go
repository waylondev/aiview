package commands

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewVideoCmd creates the video command.
func NewVideoCmd(getClient func() Client) *cobra.Command {
	var (
		subtitle     bool
		ai           bool
		comments     bool
		related      bool
	)

	cmd := &cobra.Command{
		Use:   "video <BV号或URL>",
		Short: "查看视频详情",
		Long:  `查看 Bilibili 视频详情，包括标题、UP主、播放量等。`,
		Args:  cobra.ExactArgs(1),
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
				output.EmitError("api_error", fmt.Sprintf("获取视频信息失败: %v", err), format)
				return err
			}

			payload := &VideoPayload{
				Video:    *info,
				Warnings: []Warning{},
			}

			// Subtitle
			if subtitle {
				sub, err := client.GetVideoSubtitle(bvid)
				if err == nil && sub != nil {
					payload.Subtitle = *sub
				} else {
					payload.Warnings = append(payload.Warnings, Warning{Code: "subtitle_unavailable", Message: "获取字幕失败"})
				}
			}

			// AI summary
			if ai {
				summary, err := client.GetVideoAIConclusion(bvid)
				if err == nil {
					payload.AISummary = summary
				} else {
					payload.Warnings = append(payload.Warnings, Warning{Code: "ai_summary_unavailable", Message: "获取 AI 总结失败"})
				}
			}

			// Comments
			if comments {
				cm, err := client.GetVideoComments(bvid, 1)
				if err == nil {
					payload.Comments = cm
				} else {
					payload.Warnings = append(payload.Warnings, Warning{Code: "comments_unavailable", Message: "获取评论失败"})
				}
			}

			// Related
			if related {
				rel, err := client.GetRelatedVideos(bvid)
				if err == nil {
					payload.Related = rel
				} else {
					payload.Warnings = append(payload.Warnings, Warning{Code: "related_unavailable", Message: "获取相关推荐失败"})
				}
			}

			// Structured output
			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(payload, format)
			}

			// Table output
			fmt.Printf("📺 %s\n\n", info.Title)
			fmt.Printf("  BV号:   %s\n", bvid)
			fmt.Printf("  UP主:   %s (UID: %d)\n", info.Owner.Name, info.Owner.MID)
			fmt.Printf("  时长:   %s\n", info.DurationStr)
			fmt.Printf("  播放:   %s\n", output.FormatCount(info.Stats.View))
			fmt.Printf("  弹幕:   %s\n", output.FormatCount(info.Stats.Danmaku))
			fmt.Printf("  点赞:   %s\n", output.FormatCount(info.Stats.Like))
			fmt.Printf("  投币:   %s\n", output.FormatCount(info.Stats.Coin))
			fmt.Printf("  收藏:   %s\n", output.FormatCount(info.Stats.Favorite))
			fmt.Printf("  分享:   %s\n", output.FormatCount(info.Stats.Share))
			fmt.Printf("  链接:   %s\n", info.URL)
			if info.Description != "" {
				desc := info.Description
				if len(desc) > 200 {
					desc = desc[:200]
				}
				fmt.Printf("  简介:   %s\n", desc)
			}

			if subtitle && payload.Subtitle.Available {
				fmt.Printf("\n📝 字幕内容:\n\n")
				fmt.Println(payload.Subtitle.Text)
			}

			if ai && payload.AISummary != "" {
				fmt.Printf("\n🤖 AI 总结:\n\n")
				fmt.Println(payload.AISummary)
			}

			if comments && len(payload.Comments) > 0 {
				fmt.Printf("\n💬 热门评论:\n\n")
				for _, c := range payload.Comments[:min(10, len(payload.Comments))] {
					fmt.Printf("  %s (👍 %d)\n", c.Author.Name, c.Like)
					msg := c.Message
					if len(msg) > 120 {
						msg = msg[:120]
					}
					fmt.Printf("  %s\n\n", msg)
				}
			}

			if related && len(payload.Related) > 0 {
				fmt.Printf("\n📎 相关推荐:\n\n")
				for i, r := range payload.Related[:min(10, len(payload.Related))] {
					fmt.Printf("  %d. [%s] %s (UP: %s, 播放: %s)\n",
						i+1, r.BVID, truncate(r.Title, 40), r.Owner.Name, output.FormatCount(r.Stats.View))
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&subtitle, "subtitle", "s", false, "显示字幕内容")
	cmd.Flags().BoolVar(&ai, "ai", false, "显示 AI 总结")
	cmd.Flags().BoolVarP(&comments, "comments", "c", false, "显示评论")
	cmd.Flags().BoolVarP(&related, "related", "r", false, "显示相关推荐视频")

	return cmd
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// Note: these types are imported from the parent package
// User code should import them via the parent package