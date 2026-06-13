package xiaohongshu

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewSearchCmd creates the search command.
func NewSearchCmd(client Client) *cobra.Command {
	var page int
	cmd := &cobra.Command{
		Use:   "search <keyword>",
		Short: "Search notes by keyword",
		Long: `Search for notes on Xiaohongshu by keyword.

Examples:
  aiview xiaohongshu search 美食
  aiview xiaohongshu search travel --page 2`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)

			result, err := client.SearchNotes(args[0], page)
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
