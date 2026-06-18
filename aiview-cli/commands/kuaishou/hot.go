package kuaishou

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
		Short: "View trending/hot search on Kuaishou",
		Long: `View the current hot search/trending list on Kuaishou.

Examples:
  aiview kuaishou hot
  aiview kuaishou hot --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.MustGetFormat(cmd)

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

			visionHotRank := helper.GetMap(data, "visionHotRank")
			if visionHotRank == nil {
				fmt.Println("No hot search results")
				return nil
			}

			items := helper.GetSlice(visionHotRank, "items")
			if len(items) == 0 {
				fmt.Println("No hot search results")
				return nil
			}

			fmt.Printf("🔥 Kuaishou Hot Search:\n\n")
			for i, item := range items {
				m, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				name := helper.GetString(m, "name")
				hotValue := helper.GetInt(m, "hotValue")

				fmt.Printf("  %2d. %s", i+1, name)
				if hotValue > 0 {
					fmt.Printf(" (热度: %d)", hotValue)
				}
				fmt.Println()
			}
			return nil
		},
	}
}
