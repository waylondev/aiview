package weibo

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
		Short: "View trending/hot search on Weibo",
		Long: `View the current hot search/trending list on Weibo.

Examples:
  aiview weibo hot
  aiview weibo hot --json`,
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

			realtime := helper.GetSlice(data, "realtime")
			if len(realtime) == 0 {
				fmt.Println("No hot search results")
				return nil
			}

			fmt.Printf("🔥 Weibo Hot Search:\n\n")
			for i, item := range realtime {
				m, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				word := helper.GetString(m, "word")
				hotValue := helper.GetInt(m, "num")

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
