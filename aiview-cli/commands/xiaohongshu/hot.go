package xiaohongshu

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewHotCmd creates the hot notes command.
func NewHotCmd(client Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hot",
		Short: "Show hot/trending notes",
		Long: `Show the current hot/trending notes on Xiaohongshu.

Examples:
  aiview xiaohongshu hot
  aiview xiaohongshu hot --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)

			result, err := client.GetHotNotes()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get hot notes: %v", err), format)
				return err
			}

			return output.EmitSuccess(result, format)
		},
	}
	return cmd
}
