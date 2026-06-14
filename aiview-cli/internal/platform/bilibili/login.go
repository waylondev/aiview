package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/helper"
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
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("Failed to create QR code request: %v", err))
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Network error while generating QR code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read QR code response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Failed to parse QR code response: %w", err)
	}

	code := helper.GetInt(result, "code")
	if code != 0 {
		return nil, aiverr.APIError("bilibili", fmt.Sprintf("Failed to generate QR code [%d]: %s", code, helper.GetString(result, "message")))
	}

	data := helper.GetMap(result, "data")
	return &QRLoginSession{
		QRCodeKey: helper.GetString(data, "qrcode_key"),
		QRCodeURL: helper.GetString(data, "url"),
	}, nil
}

// PollQRCode polls the QR code login status.
func PollQRCode(qrcodeKey string) (QRLoginState, *Credential, error) {
	req, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key="+qrcodeKey, nil)
	if err != nil {
		return QRLoginPending, nil, fmt.Errorf("Failed to poll QR code status: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return QRLoginPending, nil, fmt.Errorf("Network error while polling QR code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return QRLoginPending, nil, aiverr.NetworkError("bilibili", fmt.Sprintf("Failed to read poll response: %v", err))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return QRLoginPending, nil, aiverr.ParseError("bilibili", fmt.Sprintf("Failed to parse poll response: %v", err))
	}

	data := helper.GetMap(result, "data")
	code := helper.GetInt(data, "code")

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
		msg := helper.GetString(data, "message")
		return QRLoginPending, nil, aiverr.APIError("bilibili", fmt.Sprintf("QR code login error [%d]: %s", code, msg))
	}
}

// extractCookiesFromResponse extracts cookies from HTTP response.
func extractCookiesFromResponse(resp *http.Response) *Credential {
	cred := &Credential{}
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
func BuildCookieString(cred *Credential) string {
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