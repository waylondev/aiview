package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// QR login helpers — injected by the platform package via SetQRLoginFuncs.
var (
	qrGenerate = func() (*QRLoginSession, error) {
		return nil, fmt.Errorf("qr login not implemented")
	}
	qrPoll = func(key string) (int, *Credential, error) {
		return 0, nil, fmt.Errorf("qr login not implemented")
	}
)

// SetQRLoginFuncs sets the QR login implementation functions.
func SetQRLoginFuncs(
	gen func() (*QRLoginSession, error),
	poll func(key string) (int, *Credential, error),
) {
	qrGenerate = gen
	qrPoll = poll
}

// NewLoginCmd creates the login command.
func NewLoginCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var (
		sessdata string
		biliJct  string
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "登录 Bilibili 账号",
		Long: `登录 Bilibili 账号。

支持三种方式:
  1. 无参数: QR 码扫码登录
  2. --sessdata: 直接传入 SESSDATA Cookie
  3. --sessdata + --bili-jct: 传入完整凭证（支持写操作）`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := GetOutputFormat(cmd)

			// Cookie-based login
			if sessdata != "" {
				cred := &Credential{
					Sessdata: sessdata,
					BiliJct:  biliJct,
					SavedAt:  time.Now().Unix(),
				}
				if err := authStore.Save(cred); err != nil {
					output.EmitError("internal_error", fmt.Sprintf("保存凭证失败: %v", err), format)
					return err
				}
				fmt.Println("✅ 已通过 Cookie 登录")
				return nil
			}

			// QR code login
			fmt.Println("📱 正在生成二维码...")
			session, err := qrGenerate()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("生成二维码失败: %v", err), format)
				return err
			}

			fmt.Println("\n请使用 Bilibili App 扫描以下二维码登录:")
			fmt.Printf("  %s\n\n", session.QRCodeURL)
			fmt.Println("⭐ 扫码后请在手机上确认登录...")

			// Poll QR code state
			for i := 0; i < 60; i++ {
				state, cred, err := qrPoll(session.QRCodeKey)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  ⚠️ 轮询错误: %v\n", err)
					time.Sleep(2 * time.Second)
					continue
				}

				switch state {
				case 1: // Scanned
					fmt.Println("  📲 已扫码，请在手机上确认...")
				case 2: // Expired
					fmt.Println("\n❌ 二维码已过期，请重试")
					return nil
				case 3: // Success
					if err := authStore.Save(cred); err != nil {
						output.EmitError("internal_error", fmt.Sprintf("保存凭证失败: %v", err), format)
						return err
					}
					fmt.Println("\n✅ 登录成功！凭证已保存")
					return nil
				}

				time.Sleep(2 * time.Second)
			}

			fmt.Println("\n❌ 登录超时，请重试")
			return nil
		},
	}

	cmd.Flags().StringVar(&sessdata, "sessdata", "", "直接传入 SESSDATA Cookie")
	cmd.Flags().StringVar(&biliJct, "bili-jct", "", "直接传入 bili_jct Cookie")

	return cmd
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