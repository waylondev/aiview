package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewTrendingCmd creates the trending command.
func NewTrendingCmd(client Client) *cobra.Command {
	return &cobra.Command{
		Use:   "trending",
		Short: "View trending topics/challenges on Douyin",
		Long: `View trending topics, challenges and hot content on Douyin.

Examples:
  aiview douyin trending
  aiview douyin trending --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)

			result, err := client.GetTrending()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get trending: %v", err), format)
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

			wordList := helper.GetSlice(data, "word_list")
			if len(wordList) == 0 {
				fmt.Println("No trending data")
				return nil
			}

			fmt.Printf("📈 Douyin Trending:\n\n")
			for i, item := range wordList {
				m := item.(map[string]interface{})
				word := helper.GetString(m, "word")
				hotValue := helper.GetInt(m, "hot_value")

				fmt.Printf("  %2d. %s", i+1, word)
				if hotValue > 0 {
					fmt.Printf(" (热度: %d)", hotValue)
				}
				fmt.Println()
			}
			return nil
		},
	}
}