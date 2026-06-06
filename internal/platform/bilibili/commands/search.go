package commands

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewSearchCmd creates the search command.
func NewSearchCmd(getClient func() Client) *cobra.Command {
	var (
		searchType string
		maxResults int
		page       int
	)

	cmd := &cobra.Command{
		Use:   "search <关键词>",
		Short: "搜索内容",
		Long:  `搜索 Bilibili 视频或用户。`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)
			keyword := args[0]

			switch searchType {
			case "video":
				results, err := client.SearchVideo(keyword, page)
				if err != nil {
					output.EmitError("api_error", fmt.Sprintf("搜索失败: %v", err), format)
					return err
				}

				if format == output.FormatJSON || format == output.FormatYAML {
					type searchResult struct {
						Items []SearchVideoResult `json:"items"`
						Total int                  `json:"total"`
					}
					return output.EmitSuccess(searchResult{
						Items: results[:min(maxResults, len(results))],
						Total: len(results),
					}, format)
				}

				fmt.Printf("🔍 搜索 \"%s\" (视频):\n\n", keyword)
				for i, r := range results[:min(maxResults, len(results))] {
					fmt.Printf("  %d. [%s] %s\n", i+1, r.BVID, r.Title)
					fmt.Printf("     UP主: %s  播放: %s  时长: %s\n\n", r.Author, output.FormatCount(r.Play), r.Duration)
				}

			case "user":
				results, err := client.SearchUser(keyword, page)
				if err != nil {
					output.EmitError("api_error", fmt.Sprintf("搜索失败: %v", err), format)
					return err
				}

				if format == output.FormatJSON || format == output.FormatYAML {
					type searchResult struct {
						Items []SearchUserResult `json:"items"`
						Total int                `json:"total"`
					}
					return output.EmitSuccess(searchResult{
						Items: results[:min(maxResults, len(results))],
						Total: len(results),
					}, format)
				}

				fmt.Printf("🔍 搜索 \"%s\" (用户):\n\n", keyword)
				for i, r := range results[:min(maxResults, len(results))] {
					fmt.Printf("  %d. %s (UID: %d)\n", i+1, r.Name, r.MID)
					fmt.Printf("     粉丝: %s  视频: %d\n", output.FormatCount(r.Fans), r.Videos)
					if r.Sign != "" {
						fmt.Printf("     签名: %s\n", r.Sign)
					}
					fmt.Println()
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&searchType, "type", "t", "video", "搜索类型: video 或 user")
	cmd.Flags().IntVarP(&maxResults, "max", "n", 10, "最大结果数")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "页码")

	return cmd
}