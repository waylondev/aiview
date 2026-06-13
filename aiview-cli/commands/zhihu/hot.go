package zhihu

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"
	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewHotCmd creates the hot search command.
func NewHotCmd(client Client, isLoggedIn func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "hot",
		Short: "View trending/hot search on Zhihu",
		Long: `View the current hot search/trending list on Zhihu.

Examples:
  aiview zhihu hot
  aiview zhihu hot --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)

			result, err := client.GetHotSearch()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get hot search: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			data := helper.GetMap(result, "data")
			if data == nil {
				fmt.Println("No data returned")
				return nil
			}

			topSearch := helper.GetSlice(data, "top_search")
			if len(topSearch) == 0 {
				fmt.Println("No hot search results")
				return nil
			}

			fmt.Printf("🔥 Zhihu Hot Search:\n\n")
			for i, item := range topSearch {
				m, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				query := helper.GetString(m, "query")
				hotScore := helper.GetInt(m, "hot_score")

				fmt.Printf("  %2d. %s", i+1, query)
				if hotScore > 0 {
					fmt.Printf(" (热度: %d)", hotScore)
				}
				fmt.Println()
			}
			return nil
		},
	}
}
