package xiaohongshu

import (
	"fmt"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/browser"
	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLoginCmd creates the login command for Xiaohongshu.
func NewLoginCmd(saveCookie func(cookie string) error, getClient func() Client) *cobra.Command {
	var (
		cookie string
		auto   bool
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Xiaohongshu with cookie",
		Long: `Login to Xiaohongshu using a browser cookie.

Xiaohongshu API requires authentication for most endpoints.
You can get the cookie from your browser's developer tools after logging in to xiaohongshu.com.

Examples:
  aiview xiaohongshu login --cookie "your_cookie_here"
  aiview xiaohongshu login --auto`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.MustGetFormat(cmd)

			// Auto browser login
			if auto {
				fmt.Println("🌐 Opening browser for login...")
				cookies, err := browser.GetCookies("https://www.xiaohongshu.com", 2*time.Minute)
				if err != nil {
					output.EmitError("internal_error", fmt.Sprintf("Failed to get cookies from browser: %v", err), format)
					return err
				}
				fmt.Printf("🍪 Got cookies from browser (%d bytes)\n", len(cookies))
				if err := saveCookie(cookies); err != nil {
					output.EmitError("internal_error", fmt.Sprintf("failed to save credential: %v", err), format)
					return err
				}
				fmt.Println("登录成功")
				return nil
			}

			if cookie == "" {
				return aiverr.NotAuthenticated("xiaohongshu", "cookie is required, use --cookie flag or --auto")
			}
			if err := saveCookie(cookie); err != nil {
				output.EmitError("internal_error", fmt.Sprintf("failed to save credential: %v", err), format)
				return err
			}

			// Verify cookie by calling GetHotNotes
			client := getClient()
			_, err := client.GetHotNotes()
			if err != nil {
				fmt.Printf("Cookie 已保存，但验证失败: %v\n", err)
				return nil
			}
			fmt.Println("登录成功")
			return nil
		},
	}

	cmd.Flags().StringVar(&cookie, "cookie", "", "Xiaohongshu browser cookie")
	cmd.Flags().BoolVar(&auto, "auto", false, "Automatically open browser to get cookie")
	return cmd
}
