package zhihu

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewUserCmd creates the user command.
func NewUserCmd(client Client, isLoggedIn func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "user <uid>",
		Short: "View Zhihu user info",
		Long: `View Zhihu user profile information by user ID.

Examples:
  aiview zhihu user 1234567890`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.MustGetFormat(cmd)

			result, err := client.GetUserInfo(args[0])
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get user info: %v", err), format)
				return err
			}

			return output.EmitSuccess(result, format)
		},
	}
}
