package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewUserCmd creates the user command.
func NewUserCmd(getClient func() Client) *cobra.Command {
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
			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{
					"note": "User details require login. Use --cookie to authenticate.",
					"uid":  args[0],
				}, format)
			}
			fmt.Printf("👤 Douyin User\n\n")
			fmt.Printf("  UID: %s\n", args[0])
			fmt.Printf("  Note: User details require login\n")
			return nil
		},
	}
}