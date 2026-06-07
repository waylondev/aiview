package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLoginCmd creates the login command for Douyin.
func NewLoginCmd(saveCookie func(cookie string) error) *cobra.Command {
	var cookie string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Douyin with cookie",
		Long: `Login to Douyin using a browser cookie.

Examples:
  aiview douyin login --cookie "your_cookie_here"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := GetOutputFormat(cmd)
			if cookie == "" {
				return fmt.Errorf("cookie is required, use --cookie flag")
			}
			if err := saveCookie(cookie); err != nil {
				output.EmitError("internal_error", fmt.Sprintf("failed to save credential: %v", err), format)
				return err
			}
			fmt.Println("Login credential saved")
			return nil
		},
	}

	cmd.Flags().StringVar(&cookie, "cookie", "", "Douyin browser cookie")
	return cmd
}