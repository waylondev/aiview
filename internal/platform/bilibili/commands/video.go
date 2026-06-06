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
		Use:   "video <BV or URL>",
		Short: "View video details",
		Long:  `View Bilibili video details, including title, uploader, view count, etc.`,
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
				output.EmitError("api_error", fmt.Sprintf("Failed to get video info: %v", err), format)
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
					payload.Warnings = append(payload.Warnings, Warning{Code: "subtitle_unavailable", Message: "Failed to get subtitles"})
				}
			}

			// AI summary
			if ai {
				summary, err := client.GetVideoAIConclusion(bvid)
				if err == nil {
					payload.AISummary = summary
				} else {
					payload.Warnings = append(payload.Warnings, Warning{Code: "ai_summary_unavailable", Message: "Failed to get AI summary"})
				}
			}

			// Comments
			if comments {
				cm, err := client.GetVideoComments(bvid, 1)
				if err == nil {
					payload.Comments = cm
				} else {
					payload.Warnings = append(payload.Warnings, Warning{Code: "comments_unavailable", Message: "Failed to get comments"})
				}
			}

			// Related
			if related {
				rel, err := client.GetRelatedVideos(bvid)
				if err == nil {
					payload.Related = rel
				} else {
					payload.Warnings = append(payload.Warnings, Warning{Code: "related_unavailable", Message: "Failed to get related videos"})
				}
			}

			// Structured output
			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(payload, format)
			}

			// Table output
			fmt.Printf("📺 %s\n\n", info.Title)
			fmt.Printf("  BV:      %s\n", bvid)
			fmt.Printf("  Uploader: %s (UID: %d)\n", info.Owner.Name, info.Owner.MID)
			fmt.Printf("  Duration: %s\n", info.DurationStr)
			fmt.Printf("  Views:   %s\n", output.FormatCount(info.Stats.View))
			fmt.Printf("  Danmaku: %s\n", output.FormatCount(info.Stats.Danmaku))
			fmt.Printf("  Likes:   %s\n", output.FormatCount(info.Stats.Like))
			fmt.Printf("  Coins:   %s\n", output.FormatCount(info.Stats.Coin))
			fmt.Printf("  Favs:    %s\n", output.FormatCount(info.Stats.Favorite))
			fmt.Printf("  Shares:  %s\n", output.FormatCount(info.Stats.Share))
			fmt.Printf("  URL:     %s\n", info.URL)
			if info.Description != "" {
				desc := info.Description
				if len(desc) > 200 {
					desc = desc[:200]
				}
				fmt.Printf("  Desc:    %s\n", desc)
			}

			if subtitle && payload.Subtitle.Available {
				fmt.Printf("\n📝 Subtitle:\n\n")
				fmt.Println(payload.Subtitle.Text)
			}

			if ai && payload.AISummary != "" {
				fmt.Printf("\n🤖 AI Summary:\n\n")
				fmt.Println(payload.AISummary)
			}

			if comments && len(payload.Comments) > 0 {
				fmt.Printf("\n💬 Top Comments:\n\n")
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
				fmt.Printf("\n📎 Related Videos:\n\n")
				for i, r := range payload.Related[:min(10, len(payload.Related))] {
					fmt.Printf("  %d. [%s] %s (Uploader: %s, Views: %s)\n",
						i+1, r.BVID, truncate(r.Title, 40), r.Owner.Name, output.FormatCount(r.Stats.View))
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&subtitle, "subtitle", "s", false, "Show subtitle content")
	cmd.Flags().BoolVar(&ai, "ai", false, "Show AI summary")
	cmd.Flags().BoolVarP(&comments, "comments", "c", false, "Show comments")
	cmd.Flags().BoolVarP(&related, "related", "r", false, "Show related video recommendations")

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