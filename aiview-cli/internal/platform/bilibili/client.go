package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/helper"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/jackwener/aiview/internal/platform/base"
)

const (
	baseURL = "https://api.bilibili.com"

	// HTTP client constants
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
	defaultReferer   = "https://www.bilibili.com"

	// Bilibili API error codes
	CodeNotAuthenticated = -101
	CodeAccountBanned    = -111
	CodeNotFound         = -404
	CodeVideoUnavailable = 62002
	CodeLoginRequired    = 62004
	CodeRiskControl      = -352
	CodeForbidden        = -403
	CodeRateLimited      = -412
	CodeRateLimitedAlt   = 412
	CodeForbiddenAlt     = 403
)

// Client is the Bilibili API client.
type Client struct {
	*base.Client
	credential *Credential
}

// Compile-time interface assertions
var _ platform.HotSearchable = (*Client)(nil)
var _ platform.Searchable = (*Client)(nil)
var _ platform.UserQueryable = (*Client)(nil)

// NewClient creates a new Bilibili API client.
func NewClient(timeoutSec int, cookies string, credential *Credential) *Client {
	return &Client{
		Client:     base.NewClient(timeoutSec, cookies, defaultUserAgent, baseURL, defaultReferer, "bilibili"),
		credential: credential,
	}
}

// PlatformName implements the platform.Client interface.
func (c *Client) PlatformName() string {
	return "bilibili"
}

// bvidRegex matches BV IDs.
var bvidRegex = regexp.MustCompile(`\bBV[0-9A-Za-z]{10}\b`)

// BuildHeaders overrides the default to add bilibili-specific headers.
func (c *Client) BuildHeaders() http.Header {
	h := c.Client.BuildHeaders()
	h.Set("Origin", defaultReferer)
	h.Set("sec-ch-ua", "\"Chromium\";v=\"133\", \"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\"")
	h.Set("sec-ch-ua-mobile", "?0")
	h.Set("sec-ch-ua-platform", "\"Windows\"")
	return h
}

func (c *Client) buildRefererHeaders(referer string) http.Header {
	h := c.BuildHeaders()
	h.Set("Referer", referer)
	return h
}

// classifyErrorCode maps bilibili API error codes to typed errors.
func classifyErrorCode(code int, msg string) error {
	if code == CodeNotAuthenticated || code == CodeAccountBanned {
		return aiverr.NotAuthenticated("bilibili", msg)
	}
	if code == CodeNotFound || code == CodeVideoUnavailable || code == CodeLoginRequired {
		return aiverr.NotFound("bilibili", msg)
	}
	if code == CodeRateLimited || code == CodeRateLimitedAlt {
		return aiverr.RateLimited("bilibili", msg)
	}
	if code == CodeRiskControl {
		return aiverr.RateLimited("bilibili", fmt.Sprintf("Risk control triggered (code %d): %s", code, msg))
	}
	if code == CodeForbidden || code == CodeForbiddenAlt {
		return aiverr.Forbidden("bilibili", msg)
	}
	return aiverr.APIError("bilibili", fmt.Sprintf("API error [%d]: %s", code, msg))
}

func (c *Client) checkCache(key string) (map[string]interface{}, bool) {
	if c.Cache != nil {
		if cached, ok := c.Cache.Get(key); ok {
			return cached.(map[string]interface{}), true
		}
	}
	return nil, false
}

func (c *Client) applyRateLimit() {
	if c.Limiter != nil {
		c.Limiter.Wait()
	}
}

func (c *Client) executeRequest(reqURL string, headers http.Header) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("failed to create request: %v", err))
	}
	req.Header = headers

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("network request failed: %v", err))
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) handleRiskControl(reqURL string, headers http.Header, originalErr error) (map[string]interface{}, error) {
	if !isRiskControlError(originalErr) {
		return nil, originalErr
	}
	// -352 风控绕过：添加 buvid3 cookie 重试一次
	retryReq, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("failed to create retry request: %v", err))
	}
	retryReq.Header = headers
	existingCookie := retryReq.Header.Get("Cookie")
	if existingCookie != "" {
		retryReq.Header.Set("Cookie", existingCookie+"; buvid3=placeholder")
	} else {
		retryReq.Header.Set("Cookie", "buvid3=placeholder")
	}
	retryResp, err := c.HTTPClient.Do(retryReq)
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("network request failed on retry: %v", err))
	}
	defer retryResp.Body.Close()
	return c.parseResponse(retryResp)
}

func (c *Client) doGet(path string, params url.Values, headers http.Header, useCache bool) (map[string]interface{}, error) {
	reqURL := c.BaseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	// Check cache first
	if useCache {
		if cached, ok := c.checkCache(reqURL); ok {
			return cached, nil
		}
		c.applyRateLimit()
	}

	result, err := c.executeRequest(reqURL, headers)
	if err != nil {
		result, err = c.handleRiskControl(reqURL, headers, err)
		if err != nil {
			return nil, err
		}
	}

	// Store in cache
	if useCache && c.Cache != nil {
		c.Cache.Set(reqURL, result)
	}
	return result, nil
}

func (c *Client) get(path string, params url.Values) (map[string]interface{}, error) {
	return c.doGet(path, params, c.BuildHeaders(), true)
}

func (c *Client) getWithReferer(path string, params url.Values, referer string) (map[string]interface{}, error) {
	return c.doGet(path, params, c.buildRefererHeaders(referer), false)
}

func (c *Client) parseResponse(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("Failed to read response: %v", err))
	}

	// 检测 HTML 响应（直播间不存在或已关闭等情况）
	if c.DetectHTMLResponse(body) {
		return nil, aiverr.NotFound("bilibili", "直播间不存在或已关闭")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, aiverr.ParseError("bilibili", fmt.Sprintf("Failed to parse response: %v", err))
	}

	code := helper.GetInt(result, "code")
	if code != 0 {
		msg := helper.GetString(result, "message")
		return nil, classifyErrorCode(code, msg)
	}

	return result, nil
}

func (c *Client) post(path string, params url.Values) (map[string]interface{}, error) {
	if c.credential != nil && c.credential.BiliJct != "" {
		params.Set("csrf", c.credential.BiliJct)
	}

	body := params.Encode()
	req, err := http.NewRequest("POST", c.BaseURL+path, strings.NewReader(body))
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("failed to create request: %v", err))
	}
	req.Header = c.BuildHeaders()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("network request failed: %v", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, aiverr.NetworkError("bilibili", fmt.Sprintf("Failed to read response: %v", err))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, aiverr.ParseError("bilibili", fmt.Sprintf("Failed to parse response: %v", err))
	}

	code := helper.GetInt(result, "code")
	if code != 0 {
		msg := helper.GetString(result, "message")
		return nil, classifyErrorCode(code, msg)
	}

	return result, nil
}

// isRiskControlError checks if the error is a -352 risk control error.
func isRiskControlError(err error) bool {
	if err == nil {
		return false
	}
	if pe, ok := aiverr.IsPlatformError(err); ok {
		return pe.Code == aiverr.CodeRateLimited && strings.Contains(pe.Message, "-352")
	}
	return strings.Contains(err.Error(), "-352")
}
