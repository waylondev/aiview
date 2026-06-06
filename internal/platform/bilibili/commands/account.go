package commands

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLoginCmd creates the login command.
func NewLoginCmd(authStore AuthProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "登录 Bilibili 账号",
		Long:  `通过 QR 码登录 Bilibili 账号。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("📱 请从浏览器复制 Cookie 并保存，或使用以下方式:")
			fmt.Println()
			fmt.Println("  1. 打开浏览器登录 Bilibili")
			fmt.Println("  2. 按 F12 打开开发者工具")
			fmt.Println("  3. 在 Application > Cookies > bilibili.com 中找到 SESSDATA")
			fmt.Println("  4. 运行: aiview bilibili login --sessdata <你的SESSDATA>")
			fmt.Println()
			fmt.Println("  QR 码登录功能开发中...")
			return nil
		},
	}
}

// NewLogoutCmd creates the logout command.
func NewLogoutCmd(authStore AuthProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "退出登录",
		Long:  `清除本地保存的登录凭证。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := authStore.Clear(); err != nil {
				return err
			}
			fmt.Println("✅ 已退出登录")
			return nil
		},
	}
}

// NewStatusCmd creates the status command.
func NewStatusCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "检查登录状态",
		Long:  `检查当前 Bilibili 登录状态。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := GetOutputFormat(cmd)
			cred, err := authStore.GetCredential()

			if format == output.FormatJSON || format == output.FormatYAML {
				status := map[string]interface{}{
					"authenticated": err == nil && cred != nil,
				}
				if cred != nil {
					status["has_write"] = cred.HasWriteCapability()
				}
				return output.EmitSuccess(status, format)
			}

			if err != nil || cred == nil {
				fmt.Println("❌ 未登录")
				fmt.Println("   使用 aiview bilibili login 登录")
				return nil
			}
			fmt.Println("✅ 已登录")
			return nil
		},
	}
}

// NewWhoamiCmd creates the whoami command.
func NewWhoamiCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "查看当前用户信息",
		Long:  `查看当前登录的 Bilibili 用户信息。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := GetOutputFormat(cmd)
			client := getClient()

			info, err := client.GetSelfInfo()
			if err != nil {
				output.EmitError("not_authenticated", "未登录，请使用 aiview bilibili login 登录", format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(info, format)
			}

			fmt.Printf("👤 %s (UID: %d)\n", info.Name, info.MID)
			fmt.Printf("  等级: %d\n", info.Level)
			fmt.Printf("  硬币: %d\n", info.Coins)
			return nil
		},
	}
}