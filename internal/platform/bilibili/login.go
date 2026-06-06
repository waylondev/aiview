package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jackwener/aiview/internal/platform/bilibili/commands"
)

// QRLoginState represents the QR code login state.
type QRLoginState int

const (
	QRLoginPending QRLoginState = iota
	QRLoginScanned
	QRLoginExpired
	QRLoginSuccess
)

// QRLoginSession holds the QR login session data.
type QRLoginSession struct {
	QRCodeKey string `json:"qrcode_key"`
	QRCodeURL string `json:"url"`
}

// GenerateQRCode generates a new QR code login session.
func GenerateQRCode() (*QRLoginSession, error) {
	req, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/generate", nil)
	if err != nil {
		return nil, fmt.Errorf("生成二维码请求失败: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("生成二维码网络错误: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取二维码响应失败: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析二维码响应失败: %w", err)
	}

	code := getInt(result, "code")
	if code != 0 {
		return nil, fmt.Errorf("生成二维码失败 [%d]: %s", code, getString(result, "message"))
	}

	data := getMap(result, "data")
	return &QRLoginSession{
		QRCodeKey: getString(data, "qrcode_key"),
		QRCodeURL: getString(data, "url"),
	}, nil
}

// PollQRCode polls the QR code login status.
func PollQRCode(qrcodeKey string) (QRLoginState, *commands.Credential, error) {
	req, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key="+qrcodeKey, nil)
	if err != nil {
		return QRLoginPending, nil, fmt.Errorf("轮询二维码状态失败: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return QRLoginPending, nil, fmt.Errorf("轮询二维码网络错误: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return QRLoginPending, nil, fmt.Errorf("读取轮询响应失败: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return QRLoginPending, nil, fmt.Errorf("解析轮询响应失败: %w", err)
	}

	data := getMap(result, "data")
	code := getInt(data, "code")

	switch code {
	case 0:
		// Login success - extract cookies
		cred := extractCookiesFromResponse(resp)
		return QRLoginSuccess, cred, nil
	case 86038:
		return QRLoginExpired, nil, nil
	case 86090:
		return QRLoginScanned, nil, nil
	case 86101:
		return QRLoginPending, nil, nil
	default:
		msg := getString(data, "message")
		return QRLoginPending, nil, fmt.Errorf("二维码登录异常 [%d]: %s", code, msg)
	}
}

// extractCookiesFromResponse extracts cookies from HTTP response.
func extractCookiesFromResponse(resp *http.Response) *commands.Credential {
	cred := &commands.Credential{}
	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "SESSDATA":
			cred.Sessdata = cookie.Value
		case "bili_jct":
			cred.BiliJct = cookie.Value
		case "ac_time_value":
			cred.AcTimeValue = cookie.Value
		case "buvid3":
			cred.Buvid3 = cookie.Value
		case "buvid4":
			cred.Buvid4 = cookie.Value
		case "DedeUserID":
			cred.DedeUserID = cookie.Value
		}
	}
	cred.SavedAt = time.Now().Unix()
	return cred
}

// BuildCookieString builds a cookie string from credential.
func BuildCookieString(cred *commands.Credential) string {
	if cred == nil {
		return ""
	}
	var parts []string
	if cred.Sessdata != "" {
		parts = append(parts, "SESSDATA="+cred.Sessdata)
	}
	if cred.BiliJct != "" {
		parts = append(parts, "bili_jct="+cred.BiliJct)
	}
	if cred.AcTimeValue != "" {
		parts = append(parts, "ac_time_value="+cred.AcTimeValue)
	}
	if cred.Buvid3 != "" {
		parts = append(parts, "buvid3="+cred.Buvid3)
	}
	if cred.Buvid4 != "" {
		parts = append(parts, "buvid4="+cred.Buvid4)
	}
	if cred.DedeUserID != "" {
		parts = append(parts, "DedeUserID="+cred.DedeUserID)
	}
	return strings.Join(parts, "; ")
}