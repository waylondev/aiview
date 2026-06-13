package auth

import (
	"fmt"
)

// loginHelp holds login instructions for a platform.
type loginHelp struct {
	loginCmd    string
	platformURL string
}

// platformLoginHelp maps platform names to their login instructions.
var platformLoginHelp = map[string]loginHelp{
	"douyin": {
		loginCmd:    `aiview douyin login --cookie "<your_cookie>"`,
		platformURL: "douyin.com",
	},
	"xiaohongshu": {
		loginCmd:    `aiview xiaohongshu login --cookie "<your_cookie>"`,
		platformURL: "xiaohongshu.com",
	},
	"bilibili": {
		loginCmd:    `aiview bilibili login --sessdata "<your_sessdata>"`,
		platformURL: "bilibili.com",
	},
	"weibo": {
		loginCmd:    `aiview weibo login --cookie "<your_cookie>"`,
		platformURL: "weibo.com",
	},
	"kuaishou": {
		loginCmd:    `aiview kuaishou login --cookie "<your_cookie>"`,
		platformURL: "kuaishou.com",
	},
	"zhihu": {
		loginCmd:    `aiview zhihu login --cookie "<your_cookie>"`,
		platformURL: "zhihu.com",
	},
}

func getHelp(platform string) loginHelp {
	if h, ok := platformLoginHelp[platform]; ok {
		return h
	}
	return loginHelp{
		loginCmd:    fmt.Sprintf(`aiview %s login`, platform),
		platformURL: platform,
	}
}

// NotAuthenticatedError is returned when a command requires authentication but the user is not logged in.
type NotAuthenticatedError struct {
	Platform string
}

func (e *NotAuthenticatedError) Error() string {
	h := getHelp(e.Platform)
	return fmt.Sprintf(`Error: 此命令需要登录才能使用

请登录：
$ %s

如何获取 Cookie：
1. 在浏览器中登录 %s
2. 打开开发者工具 (F12)
3. 在 Network 标签中找到任意请求
4. 复制 Cookie 请求头的值`, h.loginCmd, h.platformURL)
}

// RequireAuth checks if the user is logged in. If not, returns a NotAuthenticatedError
// with friendly login instructions for the given platform.
func RequireAuth(platform string, isLoggedIn func() bool) error {
	if isLoggedIn == nil || !isLoggedIn() {
		return &NotAuthenticatedError{Platform: platform}
	}
	return nil
}

// WritePermissionError is returned when a command requires write permission but the credential doesn't have it.
type WritePermissionError struct {
	Platform string
}

func (e *WritePermissionError) Error() string {
	h := getHelp(e.Platform)
	return fmt.Sprintf(`Error: 当前登录凭证不支持写操作

请使用完整凭证重新登录：
$ %s

写操作需要额外的凭证字段（如 bilibili 的 bili_jct）。
获取方式同 Cookie 获取步骤，在开发者工具中找到对应字段的值即可。`, h.loginCmd)
}

// RequireWriteAuth checks if the user is logged in and has write permission.
func RequireWriteAuth(platform string, isLoggedIn func() bool, hasWrite func() bool) error {
	if isLoggedIn == nil || !isLoggedIn() {
		return &NotAuthenticatedError{Platform: platform}
	}
	if hasWrite != nil && !hasWrite() {
		return &WritePermissionError{Platform: platform}
	}
	return nil
}

// IsNotAuthenticatedError checks if an error is a NotAuthenticatedError.
func IsNotAuthenticatedError(err error) bool {
	_, ok := err.(*NotAuthenticatedError)
	return ok
}

// IsWritePermissionError checks if an error is a WritePermissionError.
func IsWritePermissionError(err error) bool {
	_, ok := err.(*WritePermissionError)
	return ok
}
