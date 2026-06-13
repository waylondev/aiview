package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLoginCmd creates the login command for Douyin.
func NewLoginCmd(saveCookie func(cookie string) error, getClient func() Client) *cobra.Command {
	var cookie string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Douyin with cookie",
		Long: `Login to Douyin using a browser cookie.

Examples:
  aiview douyin login --cookie "your_cookie_here"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)
			if cookie == "" {
				return fmt.Errorf("cookie is required, use --cookie flag")
			}
			if err := saveCookie(cookie); err != nil {
				output.EmitError("internal_error", fmt.Sprintf("failed to save credential: %v", err), format)
				return err
			}

			// Verify cookie by calling GetHotSearch
			client := getClient()
			_, err := client.GetHotSearch()
			if err != nil {
				fmt.Printf("Cookie 已保存，但验证失败: %v\n", err)
				return nil
			}
			fmt.Println("登录成功")
			return nil
		},
	}

	cmd.Flags().StringVar(&cookie, "cookie", "", "Douyin browser cookie")
	return cmd
}