package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewTrendingCmd creates the trending command.
func NewTrendingCmd(getClient func() Client) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "trending",
		Short: "View trending/hot search keywords",
		Long: `View Bilibili's trending hot search keywords.

Shows the current hot/trending search terms on Bilibili.

Examples:
  aiview bilibili trending
  aiview bilibili trending --limit 20`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			result, err := client.GetHotSearch(limit)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get trending: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			data := getMap(result, "data")
			if data == nil {
				fmt.Println("No data returned")
				return nil
			}

			trending := getMap(data, "trending")
			if trending == nil {
				fmt.Println("No trending data returned")
				return nil
			}

			list := getSlice(trending, "list")
			if len(list) == 0 {
				fmt.Println("No trending keywords found")
				return nil
			}

			fmt.Printf("🔥 Bilibili Trending:\n\n")
			for i, item := range list {
				m := item.(map[string]interface{})
				keyword := getString(m, "keyword")
				showName := getString(m, "show_name")
				icon := getString(m, "icon")

				display := keyword
				if showName != "" && showName != keyword {
					display = showName
				}
				badge := ""
				if icon != "" {
					badge = " [" + icon + "]"
				}
				fmt.Printf("  %2d. %s%s\n", i+1, display, badge)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "Number of trending results (1-50)")
	return cmd
}