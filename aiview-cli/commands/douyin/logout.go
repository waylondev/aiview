package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLogoutCmd creates the logout command.
func NewLogoutCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from Douyin",
		Long: `Logout from Douyin and clear saved credentials.

Examples:
  aiview douyin logout`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			if err := client.Logout(); err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to logout: %v", err), format)
				return err
			}

			fmt.Println("Logged out")
			return nil
		},
	}
}