package xiaohongshu

import (
	"fmt"

	"github.com/jackwener/aiview/internal/auth"
	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewUserCmd creates the user info command.
func NewUserCmd(client Client, isLoggedIn func() bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user <user_id>",
		Short: "Show user info",
		Long: `Show Xiaohongshu user profile information by user ID.

Requires login cookie for full access.

Examples:
  aiview xiaohongshu user 123456789
  aiview xiaohongshu user abcdef123456`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.RequireAuth("xiaohongshu", isLoggedIn); err != nil {
				return err
			}

			format := output.MustGetFormat(cmd)

			result, err := client.GetUserInfo(args[0])
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get user info: %v", err), format)
				return err
			}

			return output.EmitSuccess(result, format)
		},
	}
	return cmd
}
