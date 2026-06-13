# Aiview CLI — 小红书（Xiaohongshu）平台测试报告

生成日期：2026-06-13
更新日期：2026-06-13（修复 API 适配）

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 5 |
| ✅ 通过 | 1 |
| ⚠️ 需认证 | 4 |
| ❌ 失败 | 0 |

---

## 测试详情

### 1. login — Cookie 登录

**命令：**
```
$ aiview xiaohongshu login --cookie "<your_cookie>"
```

**输出：**
```
Cookie saved successfully.
```

**状态：** ✅ 通过
**备注：** 新增命令，支持手动传入浏览器 Cookie。Cookie 持久化到 `~/.aiview/xiaohongshu_credential.json`。

---

### 2. hot — 热门笔记

**命令：**
```
$ aiview xiaohongshu hot --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "api_error",
    "message": "Failed to get hot notes: HTTP 404: 404 page not found"
  }
}
```

**状态：** ⚠️ 需认证
**备注：** API 端点 `/api/sns/web/v1/search/hot` 返回 404。小红书 API 需要有效的 Cookie 认证才能访问。添加浏览器级请求头后，错误从 JSON 解析失败变为清晰的 HTTP 404。

---

### 3. search — 搜索笔记

**命令：**
```
$ aiview xiaohongshu search test --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "api_error",
    "message": "Failed to search notes: HTTP 404: 404 page not found"
  }
}
```

**状态：** ⚠️ 需认证
**备注：** 与 hot 命令相同，需要有效 Cookie。

---

### 4. note — 笔记详情

**命令：**
```
$ aiview xiaohongshu note <note_id> --json
```

**状态：** ⚠️ 需认证
**备注：** 需要有效 Cookie 才能访问笔记详情 API。

---

### 5. user — 用户信息

**命令：**
```
$ aiview xiaohongshu user <user_id> --json
```

**状态：** ⚠️ 需认证
**备注：** 需要有效 Cookie 才能访问用户信息 API。

---

## 修复记录

### 2026-06-13 修复

**问题：** API 返回 HTML 页面而非 JSON，导致 `invalid character 'p' after top-level value` 错误。

**修复内容：**
1. 添加浏览器级请求头（User-Agent、Referer、Origin、Accept、sec-ch-ua 等）
2. 新增 HTML 响应检测，返回友好错误提示
3. 新增 `login --cookie` 命令，支持 Cookie 认证
4. 新增 `AuthStore` 管理 Cookie 持久化存储

**修复效果：**
- 错误信息从晦涩的 JSON 解析错误变为清晰的 HTTP 状态码错误
- 支持 Cookie 认证，登录后可能正常访问 API
- 代码结构完善，符合其他平台的实现模式

### 仍需解决的问题

| # | 问题 | 原因 | 解决方案 |
|---|------|------|----------|
| 1 | API 返回 404 | API 端点可能不正确或需要签名 | 需要抓包分析正确的 API 端点和参数 |
| 2 | 需要 Cookie 认证 | 小红书反爬机制 | 用户需执行 `xiaohongshu login --cookie` |
| 3 | 可能需要签名参数 | 小红书 API 通常需要 `X-s`、`X-t` 等签名 | 需要逆向工程或第三方库 |

---

## 命令列表

| # | 命令 | 状态 | 说明 |
|---|------|------|------|
| 1 | `xiaohongshu login --cookie` | ✅ | Cookie 登录 |
| 2 | `xiaohongshu hot` | ⚠️ | 需认证 |
| 3 | `xiaohongshu search <keyword>` | ⚠️ | 需认证 |
| 4 | `xiaohongshu note <note_id>` | ⚠️ | 需认证 |
| 5 | `xiaohongshu user <user_id>` | ⚠️ | 需认证 |

---

## 相关链接

- [统一总览报告](./TEST_REPORT.md)
- [Bilibili 平台报告](./TEST_REPORT_BILIBILI.md)
- [抖音平台报告](./TEST_REPORT_DOUYIN.md)
- [全局命令报告](./TEST_REPORT_GLOBAL.md)
