package xiaohongshu

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewNoteCmd creates the note detail command.
func NewNoteCmd(client Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note <note_id>",
		Short: "Show note detail",
		Long: `Show details of a Xiaohongshu note by note ID.

Examples:
  aiview xiaohongshu note 123456789
  aiview xiaohongshu note abcdef123456`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)

			result, err := client.GetNoteDetail(args[0])
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get note detail: %v", err), format)
				return err
			}

			return output.EmitSuccess(result, format)
		},
	}
	return cmd
}
