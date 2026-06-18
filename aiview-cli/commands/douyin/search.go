package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/auth"
	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewSearchCmd creates the search command.
func NewSearchCmd(client Client, isLoggedIn func() bool) *cobra.Command {
	var (
		page  int
		count int
	)

	cmd := &cobra.Command{
		Use:   "search <keyword>",
		Short: "Search content on Douyin",
		Long: `Search for videos and users on Douyin by keyword.

Requires login cookie for full access.

Examples:
  aiview douyin search 美食
  aiview douyin search music --page 2 --count 20`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.RequireAuth("douyin", isLoggedIn); err != nil {
				return err
			}

			format := output.MustGetFormat(cmd)
			keyword := args[0]

			result, err := client.Search(keyword, page, count)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Search failed: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			data := helper.GetMap(result, "data")
			if data == nil {
				fmt.Println("No results found")
				return nil
			}

			fmt.Printf("🔍 Searching \"%s\" on Douyin:\n\n", keyword)
			fmt.Println("(Use --json for full results)")
			return nil
		},
	}

	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	cmd.Flags().IntVarP(&count, "count", "n", 10, "Results per page")
	return cmd
}