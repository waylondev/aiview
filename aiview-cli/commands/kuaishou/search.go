package kuaishou

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewSearchCmd creates the search command.
func NewSearchCmd(client Client, isLoggedIn func() bool) *cobra.Command {
	var page int

	cmd := &cobra.Command{
		Use:   "search <keyword>",
		Short: "Search content on Kuaishou",
		Long: `Search for videos by keyword on Kuaishou.

Examples:
  aiview kuaishou search 美食
  aiview kuaishou search travel --page 2`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.MustGetFormat(cmd)
			keyword := args[0]

			result, err := client.Search(keyword, page)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Search failed: %v", err), format)
				return err
			}

			return output.EmitSuccess(result, format)
		},
	}

	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	return cmd
}
