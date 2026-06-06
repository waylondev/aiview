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
		order      string
		duration   int
		tid        int
	)

	cmd := &cobra.Command{
		Use:   "search <keyword>",
		Short: "Search content",
		Long:  `Search Bilibili videos or users.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)
			keyword := args[0]

			switch searchType {
			case "video":
				results, err := client.SearchVideo(keyword, page, order, duration, tid)
				if err != nil {
					output.EmitError("api_error", fmt.Sprintf("Search failed: %v", err), format)
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

				fmt.Printf("🔍 Searching \"%s\" (videos):\n\n", keyword)
				for i, r := range results[:min(maxResults, len(results))] {
					fmt.Printf("  %d. [%s] %s\n", i+1, r.BVID, r.Title)
					fmt.Printf("     Uploader: %s  Views: %s  Duration: %s\n\n", r.Author, output.FormatCount(r.Play), r.Duration)
				}

			case "user":
				results, err := client.SearchUser(keyword, page)
				if err != nil {
					output.EmitError("api_error", fmt.Sprintf("Search failed: %v", err), format)
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

				fmt.Printf("🔍 Searching \"%s\" (users):\n\n", keyword)
				for i, r := range results[:min(maxResults, len(results))] {
					fmt.Printf("  %d. %s (UID: %d)\n", i+1, r.Name, r.MID)
					fmt.Printf("     Fans: %s  Videos: %d\n", output.FormatCount(r.Fans), r.Videos)
					if r.Sign != "" {
						fmt.Printf("     Sign: %s\n", r.Sign)
					}
					fmt.Println()
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&searchType, "type", "t", "video", "Search type: video or user")
	cmd.Flags().IntVarP(&maxResults, "max", "n", 10, "Max results")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	cmd.Flags().StringVarP(&order, "order", "o", "", "Sort order: click/pubdate/dm/score")
	cmd.Flags().IntVarP(&duration, "duration", "d", 0, "Duration filter: 0=all, 1=<5min, 2=5-30min, 3=>30min")
	cmd.Flags().IntVar(&tid, "tid", 0, "Category ID filter")

	return cmd
}