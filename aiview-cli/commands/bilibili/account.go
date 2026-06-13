package bilibili

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackwener/aiview/internal/browser"
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
		auto     bool
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Bilibili",
		Long: `Login to Bilibili.

Three methods supported:
  1. No arguments: QR code scan login
  2. --sessdata: Pass SESSDATA Cookie directly
  3. --sessdata + --bili-jct: Pass full credential (supports write operations)
  4. --auto: Automatically open browser to get cookie

How to get bili_jct:
  1. Open browser, login to bilibili.com
  2. Press F12 → Application → Cookies → bilibili.com
  3. Copy the value of bili_jct`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)

			// Auto browser login
			if auto {
				fmt.Println("🌐 Opening browser for login...")
				cookies, err := browser.GetCookies("https://www.bilibili.com", 2*time.Minute)
				if err != nil {
					output.EmitError("internal_error", fmt.Sprintf("Failed to get cookies from browser: %v", err), format)
					return err
				}
				fmt.Printf("🍪 Got cookies from browser (%d bytes)\n", len(cookies))
				cred := &Credential{
					Sessdata: extractCookieValue(cookies, "SESSDATA"),
					BiliJct:  extractCookieValue(cookies, "bili_jct"),
					SavedAt:  time.Now().Unix(),
				}
				if err := authStore.Save(cred); err != nil {
					output.EmitError("internal_error", fmt.Sprintf("Failed to save credential: %v", err), format)
					return err
				}
				fmt.Println("✅ Logged in via browser")
				return nil
			}

			// Cookie-based login
			if sessdata != "" {
				cred := &Credential{
					Sessdata: sessdata,
					BiliJct:  biliJct,
					SavedAt:  time.Now().Unix(),
				}
				if err := authStore.Save(cred); err != nil {
					output.EmitError("internal_error", fmt.Sprintf("Failed to save credential: %v", err), format)
					return err
				}
				fmt.Println("✅ Logged in via Cookie")
				return nil
			}

			// QR code login
			fmt.Println("📱 Generating QR code...")
			session, err := qrGenerate()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to generate QR code: %v", err), format)
				return err
			}

			fmt.Println("\nScan the QR code with the Bilibili App:")
			fmt.Printf("  %s\n\n", session.QRCodeURL)
			fmt.Println("⭐ Confirm login on your phone after scanning...")

			// Poll QR code state
			for i := 0; i < 60; i++ {
				state, cred, err := qrPoll(session.QRCodeKey)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  ⚠️ Poll error: %v\n", err)
					time.Sleep(2 * time.Second)
					continue
				}

				switch state {
				case 1: // Scanned
					fmt.Println("  📲 Scanned, confirm on your phone...")
				case 2: // Expired
					fmt.Println("\n❌ QR code expired, please retry")
					return nil
				case 3: // Success
					if err := authStore.Save(cred); err != nil {
						output.EmitError("internal_error", fmt.Sprintf("Failed to save credential: %v", err), format)
						return err
					}
					fmt.Println("\n✅ Login successful! Credential saved")
					return nil
				}

				time.Sleep(2 * time.Second)
			}

			fmt.Println("\n❌ Login timeout, please retry")
			return nil
		},
	}

	cmd.Flags().StringVar(&sessdata, "sessdata", "", "Set SESSDATA Cookie directly")
	cmd.Flags().StringVar(&biliJct, "bili-jct", "", "Set bili_jct Cookie directly")
	cmd.Flags().BoolVar(&auto, "auto", false, "Automatically open browser to get cookie")

	return cmd
}

// extractCookieValue extracts a cookie value from a cookie string.
func extractCookieValue(cookies, name string) string {
	for _, part := range strings.Split(cookies, ";") {
		part = strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 && kv[0] == name {
			return kv[1]
		}
	}
	return ""
}

// NewLogoutCmd creates the logout command.
func NewLogoutCmd(authStore AuthProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout",
		Long:  `Clear locally saved login credentials.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)
			if err := authStore.Clear(); err != nil {
				output.EmitError("internal_error", fmt.Sprintf("Failed to clear credential: %v", err), format)
				return err
			}
			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]string{"message": "logged out"}, format)
			}
			fmt.Println("✅ Logged out")
			return nil
		},
	}
	return cmd
}

// NewStatusCmd creates the status command.
func NewStatusCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check login status",
		Long:  `Check current Bilibili login status.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)
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
				fmt.Println("❌ Not logged in")
				fmt.Println("   Use aiview bilibili login to log in")
				return nil
			}
			fmt.Println("✅ Logged in")
			return nil
		},
	}
}

// NewWhoamiCmd creates the whoami command.
func NewWhoamiCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "View current user info",
		Long:  `View currently logged-in Bilibili user info.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)
			client := getClient()

			info, err := client.GetSelfInfo()
			if err != nil {
				output.EmitError("not_authenticated", err.Error(), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(info, format)
			}

			fmt.Printf("👤 %s (UID: %d)\n", info.Name, info.MID)
			fmt.Printf("  Level: %d\n", info.Level)
			fmt.Printf("  Coins: %d\n", info.Coins)
			return nil
		},
	}
}