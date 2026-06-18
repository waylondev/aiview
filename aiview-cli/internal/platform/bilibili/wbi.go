package bilibili

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/helper"
)

// WBIKey holds the WBI signing keys.
type WBIKey struct {
	ImgKey string
	SubKey string
	MixKey string
}

var (
	wbiKeyCache     *WBIKey
	wbiKeyCacheTime time.Time
	wbiKeyMu        sync.Mutex
)

// getWBIKey fetches the WBI signing keys from Bilibili.
func getWBIKey() (*WBIKey, error) {
	wbiKeyMu.Lock()
	defer wbiKeyMu.Unlock()

	// Cache for 1 hour
	if wbiKeyCache != nil && time.Since(wbiKeyCacheTime) < time.Hour {
		return wbiKeyCache, nil
	}

	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/web-interface/nav", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	data := helper.GetMap(result, "data")
	wbi := helper.GetMap(data, "wbi_img")
	imgURL := helper.GetString(wbi, "img_url")
	subURL := helper.GetString(wbi, "sub_url")

	imgKey := extractKeyFromURL(imgURL)
	subKey := extractKeyFromURL(subURL)

	if imgKey == "" || subKey == "" {
		return nil, aiverr.APIError("bilibili", "failed to extract WBI keys")
	}

	mixKey := mixWBIKeys(imgKey, subKey)

	wbiKeyCache = &WBIKey{
		ImgKey: imgKey,
		SubKey: subKey,
		MixKey: mixKey,
	}
	wbiKeyCacheTime = time.Now()

	return wbiKeyCache, nil
}

// extractKeyFromURL extracts the key from a URL like https://i0.hdslb.com/bfs/wbi/xxx.png
func extractKeyFromURL(rawURL string) string {
	idx := strings.LastIndex(rawURL, "/")
	if idx < 0 {
		return ""
	}
	name := rawURL[idx+1:]
	dotIdx := strings.LastIndex(name, ".")
	if dotIdx < 0 {
		return name
	}
	return name[:dotIdx]
}

// mixinKeyEncTab is the fixed 64-element WBI mixing table.
var mixinKeyEncTab = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35,
	27, 43, 5, 49, 33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13,
	37, 48, 7, 16, 24, 55, 40, 61, 26, 17, 0, 1, 60, 51, 30, 4,
	22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11, 36, 20, 34, 44, 52,
}

// mixWBIKeys generates the 32-character mixed key.
func mixWBIKeys(imgKey, subKey string) string {
	rawKey := imgKey + subKey
	var sb strings.Builder
	for _, idx := range mixinKeyEncTab {
		if idx < len(rawKey) {
			sb.WriteByte(rawKey[idx])
		}
	}
	result := sb.String()
	if len(result) > 32 {
		result = result[:32]
	}
	return result
}

// signWBI signs the request parameters with WBI signing.
func (c *Client) signWBI(params url.Values) url.Values {
	key, err := getWBIKey()
	if err != nil {
		params.Set("wts", fmt.Sprintf("%d", time.Now().Unix()))
		return params
	}

	params.Set("wts", fmt.Sprintf("%d", time.Now().Unix()))

	// Sort params by key
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string with proper URL encoding
	// Bilibili WBI requires: uppercase hex, %20 for space, filter !'()*
	var queryParts []string
	for _, k := range keys {
		v := params.Get(k)
		// Filter characters: !'()*
		v = filterSpecialChars(v)
		queryParts = append(queryParts, wbiEncode(k)+"="+wbiEncode(v))
	}
	query := strings.Join(queryParts, "&")

	signStr := query + key.MixKey
	hash := md5.Sum([]byte(signStr))
	wrid := hex.EncodeToString(hash[:])

	params.Set("w_rid", wrid)
	return params
}

// filterSpecialChars removes !'()* characters from values.
func filterSpecialChars(s string) string {
	s = strings.ReplaceAll(s, "!", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "*", "")
	return s
}

// wbiEncode URL-encodes with uppercase hex and %20 for spaces.
func wbiEncode(s string) string {
	encoded := url.QueryEscape(s)
	// Replace + with %20 for spaces
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	// Uppercase hex
	var result strings.Builder
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == '%' && i+2 < len(encoded) {
			result.WriteByte('%')
			result.WriteString(strings.ToUpper(encoded[i+1 : i+3]))
			i += 2
		} else {
			result.WriteByte(encoded[i])
		}
	}
	return result.String()
}

// wbiGet performs a GET request with WBI signing.
func (c *Client) wbiGet(path string, params url.Values) (map[string]interface{}, error) {
	signedParams := c.signWBI(params)
	return c.get(path, signedParams)
}

// wbiGetWithReferer performs a GET request with WBI signing and custom Referer.
func (c *Client) wbiGetWithReferer(path string, params url.Values, referer string) (map[string]interface{}, error) {
	signedParams := c.signWBI(params)
	return c.getWithReferer(path, signedParams, referer)
}