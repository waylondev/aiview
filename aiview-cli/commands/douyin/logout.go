package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLogoutCmd creates the logout command.
func NewLogoutCmd(logoutFn func() error) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from Douyin",
		Long: `Logout from Douyin and clear saved credentials.

Examples:
  aiview douyin logout`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.MustGetFormat(cmd)

			if err := logoutFn(); err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to logout: %v", err), format)
				return err
			}

			fmt.Println("Logged out")
			return nil
		},
	}
}