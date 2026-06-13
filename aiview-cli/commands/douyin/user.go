package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewUserCmd creates the user command.
func NewUserCmd(client Client) *cobra.Command {
	return &cobra.Command{
		Use:   "user <uid>",
		Short: "View Douyin user info",
		Long: `View Douyin user profile information by user ID.

Requires login cookie for full access.

Examples:
  aiview douyin user 123456789`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := GetOutputFormat(cmd)

			result, err := client.GetUserInfo(args[0])
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get user info: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			fmt.Printf("👤 Douyin User\n\n")
			fmt.Printf("  UID: %s\n", args[0])
			if note, ok := result["note"].(string); ok && note != "" {
				fmt.Printf("  Note: %s\n", note)
			}
			return nil
		},
	}
}