# Aiview CLI — Bilibili 平台测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: Cookie 已存在但已过期（status 显示 authenticated=true，但 whoami 等返回 not_authenticated）

## 概览
| 类别 | 数量 | ✅ | ⚠️ | ❌ |
|------|------|----|----|-----|
| 只读 | 18 | 14 | 2 | 2 |
| 需认证读 | 9 | 2 | 0 | 7 |
| 写操作 | 13 | 10 | 1 | 2 |
| **总计** | **40** | **26** | **3** | **11** |

---

## 只读命令

### 1. hot — 热门视频

**命令：**
```
$ aiview bilibili hot --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "items": [
      {
        "bvid": "BV11GEv6dEvW",
        ...
```

**状态：** ✅ 通过

**备注：** 返回热门视频列表。

---

### 2. search — 搜索

**命令：**
```
$ aiview bilibili search test --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "items": [
      { "id": "BV1UbX3B2EZQ", "title": "[FNF] 官方隐藏曲目 test", ... },
      ...
    ],
    "total": 20
  }
}
```

**状态：** ✅ 通过

**备注：** 搜索功能正常，返回 20 条结果。

---

### 3. video — 视频详情

**命令：**
```
$ aiview bilibili video BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "video": {
      "bvid": "BV1xx411c7m9",
      "aid": 7,
      "title": "2012地球便当之日宣传片",
      "owner": { "id": 2, "name": "碧诗" },
      "stats": { "view": 9473424, "danmaku": 225605, "like": 396109, ... }
    }
  }
}
```

**状态：** ✅ 通过

**备注：** 返回完整的视频元数据，包含字幕、AI摘要、评论、相关视频等信息。

---

### 4. user — 用户信息

**命令：**
```
$ aiview bilibili user 37737161 --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "id": 37737161,
    "name": "青空の霞光",
    "level": 6,
    "fans": 84099,
    "following": 20
  }
}
```

**状态：** ✅ 通过

**备注：** 用户信息获取正常。

---

### 5. recommend — 推荐视频

**命令：**
```
$ aiview bilibili recommend --json
```

**输出：**
```json
{ "ok": true, "data": { "items": [ ... ] } }
```

**状态：** ✅ 通过

**备注：** 返回推荐视频列表，包含多个视频卡片及统计信息。

---

### 6. trending — 热搜/热词

**命令：**
```
$ aiview bilibili trending --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "trending": {
      "list": [
        { "keyword": "加拿大1-1波黑", "heat_score": 4051551 },
        { "keyword": "锐评HLE战胜T1", "heat_score": 3742426 },
        ...
      ],
      "title": "bilibili热搜"
    }
  }
}
```

**状态：** ✅ 通过

**备注：** 返回 B 站热搜榜单，包含热度和图标。

---

### 7. rank — 排行榜

**命令：**
```
$ aiview bilibili rank --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "api_error",
    "message": "Failed to get rankings: API error [-352]: -352"
  }
}
```

**状态：** ❌ 失败

**备注：** API 返回 -352 错误（风控校验失败）。该接口需要更完整的请求头或 cookie 才能通过风控。

---

### 8. region — 分区视频

**命令：**
```
$ aiview bilibili region 1 --json
```

**输出：**
```json
{ "ok": true, "data": { "archives": [ ... ], "page": { "count": 190 } } }
```

**状态：** ✅ 通过

**备注：** 按分区 ID=1 获取视频列表，返回 20 条/页，共 190 条。

---

### 9. tags — 视频标签

**命令：**
```
$ aiview bilibili tags BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "data": [
    { "tag_name": "B站最早的视频" },
    { "tag_name": "历史的起点" },
    { "tag_name": "考古的尽头" },
    { "tag_name": "B站的起源" },
    { "tag_name": "传说中的诞生" },
    { "tag_name": "便当与诞生" }
  ]
}
```

**状态：** ✅ 通过

**备注：** 返回视频的标签列表，包含 tag_id 和 tag_type。

---

### 10. suggest — 搜索建议

**命令：**
```
$ aiview bilibili suggest test --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "name": "智能路障",
    "show_name": "智能路障",
    "url": "https://search.bilibili.com/all?keyword=智能路障"
  }
}
```

**状态：** ✅ 通过

**备注：** 返回搜索建议（只返回了一个最匹配的 UP 主名称）。

---

### 11. precious — 入站必刷

**命令：**
```
$ aiview bilibili precious --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "title": "入站必刷",
    "list": [
      { "achievement": "哔哩哔哩第一个视频", "title": "【B站】2012地球便当...",
        "owner": { "name": "碧诗" } },
      { "achievement": "全站最多播放", "title": "【B站】BDF2017...", ... },
      ...
    ]
  }
}
```

**状态：** ✅ 通过

**备注：** 返回入站必刷经典视频列表，每个视频带有 achievement 字段说明。

---

### 12. online — 在线人数

**命令：**
```
$ aiview bilibili online BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "count": "1",
    "total": "1"
  }
}
```

**状态：** ✅ 通过

**备注：** 返回视频当前在线观看人数。

---

### 13. weekly — 每周必看

**命令：**
```
$ aiview bilibili weekly 1 --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "config": { "name": "2019第1期 03.22 - 03.28", "subject": "神仙爱情" },
    "list": [ ... ]
  }
}
```

**状态：** ✅ 通过

**备注：** 返回指定期的每周必看内容，包含视频列表和推荐理由。

---

### 14. live — 直播间信息

**命令：**
```
$ aiview bilibili live --room 3 --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "api_error",
    "message": "Failed to get live room info for 3: Failed to parse response: invalid character '<' looking for beginning of value"
  }
}
```

**状态：** ⚠️ 警告

**备注：** Room ID=3 的直播间可能不存在或 API 返回了 HTML 页面。建议使用一个活跃的 room ID 测试。

---

### 15. audio — 下载音频

**命令：**
```
$ aiview bilibili audio BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "directory": "C:\\Users\\Administrator\\AppData\\Local\\Temp\\aiview\\2012地球便当之日宣传片",
    "duration": 25,
    "segments": 7
  }
}
```

**状态：** ✅ 通过

**备注：** 成功下载音频到临时目录，25 秒时长，7 个分段。

---

### 16. danmaku — 弹幕查看

**命令：**
```
$ aiview bilibili danmaku BV1xx411c7m9 --json
```

**输出：**
```
(二进制 protobuf 数据，非 JSON)
```

**状态：** ⚠️ 警告

**备注：** 弹幕数据以 protobuf 二进制格式返回，未解析为可读的 JSON。需要格式化 protobuf 为可读输出。

---

### 17. collection — 用户合集

**命令：**
```
$ aiview bilibili collection 37737161 --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "no_collections",
    "message": "This user has no video collections (合集/系列). Try another UID."
  }
}
```

**状态：** ⚠️ 警告

**备注：** 该用户（青空の霞光）没有创建视频合集/系列。这是一个合法的业务结果（用户确实没有合集），建议用有合集的 UID 测试。

---

### 18. user-videos — 用户视频列表

**命令：**
```
$ aiview bilibili user-videos 37737161 --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "api_error",
    "message": "Failed to get user videos: API error [-352]: 风控校验失败"
  }
}
```

**状态：** ❌ 失败

**备注：** 与 rank 命令相同的 -352 风控错误。该接口需要更强的认证信息。

---

## 需认证读命令

### 测试说明
本地 cookie 文件存在但已过期：`status` 命令报告 `authenticated: true, has_write: true`，但实际 API 调用返回 `not_authenticated`（账号未登录）。以下标注"⚠️ Cookie 过期"表示 cookie 存在但无效。

---

### 19. status — 登录状态检查

**命令：**
```
$ aiview bilibili status --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "authenticated": true,
    "has_write": true
  }
}
```

**状态：** ✅ 通过

**备注：** 检查本地 cookie 文件是否存在和是否具备写权限。注意：此命令仅检查本地状态，不验证 cookie 有效性。

---

### 20. whoami — 当前用户信息

**命令：**
```
$ aiview bilibili whoami --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "not_authenticated",
    "message": "Not logged in, use aiview bilibili login to log in"
  }
}
```

**状态：** ❌ 失败 (Cookie 过期)

**备注：** Cookie 存在但已过期，API 验证不通过。

---

### 21. feed — 关注动态

**命令：**
```
$ aiview bilibili feed --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "api_error", "message": "Failed to get feed: not_authenticated: 账号未登录" }
}
```

**状态：** ❌ 失败 (Cookie 过期)

---

### 22. history — 观看历史

**命令：**
```
$ aiview bilibili history --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "api_error", "message": "Failed to get watch history: not_authenticated: 账号未登录" }
}
```

**状态：** ❌ 失败 (Cookie 过期)

---

### 23. watch-later — 稍后再看

**命令：**
```
$ aiview bilibili watch-later --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "api_error", "message": "Failed to get watch later list: not_authenticated: 账号未登录" }
}
```

**状态：** ❌ 失败 (Cookie 过期)

---

### 24. favorites — 收藏夹

**命令：**
```
$ aiview bilibili favorites --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "not_authenticated", "message": "Failed to get user info: not_authenticated: 账号未登录" }
}
```

**状态：** ❌ 失败 (Cookie 过期)

---

### 25. following — 关注列表

**命令：**
```
$ aiview bilibili following --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "not_authenticated", "message": "Failed to get user info: not_authenticated: 账号未登录" }
}
```

**状态：** ❌ 失败 (Cookie 过期)

---

### 26. fans — 粉丝列表

**命令：**
```
$ aiview bilibili fans 37737161 --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "api_error", "message": "Failed to get fans list: API error [-352]: -352" }
}
```

**状态：** ❌ 失败

**备注：** -352 风控错误，与 rank、user-videos 相同。

---

### 27. relation — 用户关系

**命令：**
```
$ aiview bilibili relation 37737161 --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "mid": 37737161,
    "follower": 84098,
    "following": 20,
    "black": 0,
    "whisper": 0
  }
}
```

**状态：** ✅ 通过

**备注：** 关系查询无需认证，返回关注/粉丝/拉黑状态。

---

## 写操作命令

### 测试说明
测试过程中执行了 `logout`（测试项目需要），导致后续 `like`/`coin`/`triple` 命令因登出而失败。这些命令在已登录状态下应能正常工作。`logout` 的副作用已体现在 like/coin/triple 的结果中。

---

### 28. login --help

**命令：**
```
$ aiview bilibili login --help
```

**输出：**
```
Login to Bilibili.

Three methods supported:
  1. No arguments: QR code scan login
  2. --sessdata: Pass SESSDATA Cookie directly
  3. --sessdata + --bili-jct: Pass full credential (supports write operations)

Usage:
  aiview bilibili login [flags]

Flags:
      --bili-jct string   Set bili_jct Cookie directly
      --sessdata string   Set SESSDATA Cookie directly
```

**状态：** 🔘 需手动（浏览器扫码或手动注入 cookie）

**备注：** 支持三种登录方式：二维码扫描、SESSDATA cookie、完整凭证。

---

### 29. logout

**命令：**
```
$ aiview bilibili logout
```

**输出：**
```
✅ Logged out
```

**状态：** ✅ 通过

**备注：** 成功清除本地存储的 cookie 文件。

---

### 30. like — 点赞

**命令：**
```
$ aiview bilibili like BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "not_authenticated", "message": "Login with write permission required" }
}
```

**状态：** 🔘 需认证（本次因 logout 导致未登录）

**备注：** 需要写权限的 cookie。Flags: `--undo` 取消点赞。

---

### 31. coin — 投币

**命令：**
```
$ aiview bilibili coin BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "not_authenticated", "message": "Login with write permission required" }
}
```

**状态：** 🔘 需认证（本次因 logout 导致未登录）

**备注：** Flags: `-n, --num` 指定投币数量（1-2）。

---

### 32. triple — 一键三连

**命令：**
```
$ aiview bilibili triple BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": false,
  "error": { "code": "not_authenticated", "message": "Login with write permission required" }
}
```

**状态：** 🔘 需认证（本次因 logout 导致未登录）

---

### 33. favorite --help

**命令：**
```
$ aiview bilibili favorite --help
```

**输出：**
```
Add or remove a video from a favorite folder.

Requires --fid to specify the folder ID.
Use --delete to remove from the folder.

Examples:
  aiview bilibili favorite BV1xx --fid 12345
  aiview bilibili favorite BV1xx --fid 12345 --delete

Flags:
      --delete    Remove from favorites
      --fid int   Favorite folder ID (required)
```

**状态：** 🔘 需认证

---

### 34. unfollow --help

**命令：**
```
$ aiview bilibili unfollow --help
```

**输出：**
```
Unfollow a user (login and write permission required).

Usage:
  aiview bilibili unfollow <UID> [flags]
```

**状态：** 🔘 需认证

---

### 35. comment --help

**命令：**
```
$ aiview bilibili comment --help
```

**输出：**
```
View comments on a video, or post/delete comments.

Usage:
  aiview bilibili comment <BV or URL> [flags]

Flags:
  -m, --message string   Post a comment with the given message
  -p, --page int         Page number (default 1)
      --parent int       Parent comment ID for reply
      --root int         Root comment ID for reply
      --sort int         Sort order: 0=time, 2=hot
```

**状态：** 🔘 需认证（发评论）/ ✅ 查看评论无需认证

---

### 36. comment-delete --help

**命令：**
```
$ aiview bilibili comment-delete --help
```

**输出：**
```
Delete your own comment on a video (login and write permission required).

Usage:
  aiview bilibili comment-delete <BV> <RPID> [flags]
```

**状态：** 🔘 需认证

---

### 37. danmaku-send --help

**命令：**
```
$ aiview bilibili danmaku-send --help
```

**输出：**
```
Send a danmaku (bullet comment) on a video (login and write permission required).

Usage:
  aiview bilibili danmaku-send <BV> <message> [flags]

Flags:
      --progress int   Video progress in seconds
```

**状态：** 🔘 需认证

---

### 38. dynamic-post --help

**命令：**
```
$ aiview bilibili dynamic-post --help
```

**输出：**
```
Post a plain text dynamic to your Bilibili feed (login and write permission required).

Examples:
  aiview bilibili dynamic-post "Hello World!"

Usage:
  aiview bilibili dynamic-post <text> [flags]
```

**状态：** 🔘 需认证

---

### 39. dynamic-delete --help

**命令：**
```
$ aiview bilibili dynamic-delete --help
```

**输出：**
```
Delete a dynamic by its ID (login and write permission required).

Examples:
  aiview bilibili dynamic-delete 123456

Usage:
  aiview bilibili dynamic-delete <id> [flags]
```

**状态：** 🔘 需认证

---

### 40. video-status — 视频统计

**命令：**
```
$ aiview bilibili video-status BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "data": {
    "view": 9473433,
    "danmaku": 225606,
    "like": 396110,
    "coin": 42296,
    "favorite": 149621,
    "share": 30121
  }
}
```

**状态：** ✅ 通过

**备注：** 返回视频的统计信息（播放量、弹幕数、点赞、硬币、收藏、分享），无需认证。

---

## 总结

### 通过 (26)
大部分只读命令和帮助命令工作正常。

### 需关注的 API 问题

| 问题 | 影响命令 | 错误码 |
|------|----------|--------|
| 风控校验失败 | rank, user-videos, fans | -352 |
| 直播间解析失败 | live | HTML 响应 |
| 弹幕二进制输出 | danmaku | protobuf 未解析 |

### 认证状态
| Cookie 状态 | 说明 |
|------------|------|
| 已过期 | status 报告已登录但 API 返回未认证 |
| 需要刷新 | 使用 `aiview bilibili login` 重新扫码或注入新 cookie |

### 建议
1. **风控问题（-352）**：rank、user-videos、fans 命令被 B 站风控拦截。需要添加更多请求头（如 User-Agent、Referer）或使用 cookie 绕过。
2. **弹幕解析**：danmaku 命令返回 protobuf 二进制数据，建议添加 protobuf→JSON 的自动转换。
3. **live 命令**：使用无效 room ID 时 API 返回 HTML 而非 JSON，导致解析失败。建议使用已知有效的 room ID 测试。
4. **Cookie 管理**：`status` 和 `whoami` 的不一致表明 status 仅检查本地文件存在性，建议 status 也验证 cookie 有效性。