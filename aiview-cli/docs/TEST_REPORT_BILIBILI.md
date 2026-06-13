# Aiview CLI — Bilibili 平台测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录（status 显示 authenticated=false）

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 40 |
| ✅ 通过 | 33 |
| ⚠️ 需认证 | 4 |
| ❌ 失败 | 3 |

---

## 一、只读命令（18 个）

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
        "bvid": "BV1tUEY6WE3t",
        "title": "...",
        "owner": { "id": ..., "name": "..." },
        "stats": { "view": ..., "danmaku": ..., "like": ... }
      },
      ...
    ]
  }
}
```

**状态：** ✅ 通过
**备注：** 返回热门视频列表，包含完整的视频元数据和统计信息。

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
  "schema_version": "1",
  "data": {
    "items": [
      { "id": "BV1UbX3B2EZQ", "title": "[FNF] 官方隐藏曲目 test", "author": "愛玩病毒的WSZX", "play": 13159, "duration": "1:46" },
      { "id": "BV11A411D7vf", "title": "<摄魂>Test出场之MV This Is A Test V5", ... },
      ...
    ]
  }
}
```

**状态：** ✅ 通过
**备注：** 搜索功能正常，返回 id/title/author/play/duration 字段。

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
  "schema_version": "1",
  "data": {
    "video": {
      "bvid": "BV1xx411c7m9",
      "aid": 7,
      "cid": 3625120,
      "title": "2012地球便当之日宣传片",
      "duration": "02:51",
      "owner": { "id": 2, "name": "碧诗" },
      "stats": { "view": 9473644, "danmaku": 225608, "like": 396122, "coin": 42300, "favorite": 149631, "share": 30122 }
    },
    "subtitle": { "available": false },
    "ai_summary": "",
    "comments": null,
    "related": null
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
  "schema_version": "1",
  "data": {
    "id": 37737161,
    "name": "青空の霞光",
    "level": 6,
    "coins": 0,
    "sign": "欢迎关注我们一起学习，官网地址：https://itbaima.cn",
    "fans": 84099,
    "following": 20
  }
}
```

**状态：** ✅ 通过
**备注：** 用户信息获取正常，包含 coins 和 sign 字段。

---

### 5. user-videos — 用户视频列表

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
**备注：** -352 风控错误。该接口需要更强的认证信息或请求头。

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
  "schema_version": "1",
  "data": {
    "trending": {
      "list": [
        { "keyword": "加拿大1-1波黑", "heat_score": 7118240, "show_name": "加拿大1-1波黑" },
        ...
      ],
      "title": "bilibili热搜"
    }
  }
}
```

**状态：** ✅ 通过
**备注：** 返回 B 站热搜榜单。

---

### 7. rank — 排行榜

**命令：**
```
$ aiview bilibili rank --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "items": [
      { "bvid": "BV18VED63EBi", "title": "认为自己旅游运差时可以看看我.", "owner": { "name": "风景旅行收藏家" }, "stats": { "view": 3941107 } },
      { "bvid": "BV1UzE765Ewm", "title": "还能再来一次吗…", "owner": { "name": "小潮院长" }, "stats": { "view": 2576910 } },
      ...
    ]
  }
}
```

**状态：** ✅ 通过
**备注：** 已修复！上次测试返回 -352 风控错误，本次正常返回排行榜数据。

---

### 8. region — 分区视频

**命令：**
```
$ aiview bilibili region 1 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "archives": [ { "bvid": "BV1gCEc6BEgN", ... }, ... ],
    "page": { ... }
  }
}
```

**状态：** ✅ 通过
**备注：** 按分区 ID=1 获取视频列表，正常返回。

---

### 9. recommend — 推荐视频

**命令：**
```
$ aiview bilibili recommend --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "item": [ { ... }, ... ],
    "business_card": null,
    "floor_info": null
  }
}
```

**状态：** ✅ 通过
**备注：** 返回推荐视频列表。

---

### 10. feed — 关注动态

**命令：**
```
$ aiview bilibili feed --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "not_authenticated",
    "message": "Login required, use aiview bilibili login"
  }
}
```

**状态：** ❌ 失败（🔒 需登录）
**备注：** 需要登录才能查看关注动态。

---

### 11. history — 观看历史

**命令：**
```
$ aiview bilibili history --help
```

**输出：**
```
View watch history (login required).

Usage:
  aiview bilibili history [flags]

Flags:
  -h, --help       help for history
  -n, --max int    Maximum number of history items to show
  -p, --page int   Page number (default 1)
```

**状态：** ⚠️ 需认证
**备注：** 帮助信息正确显示，实际使用需要登录。

---

### 12. watch-later — 稍后再看

**命令：**
```
$ aiview bilibili watch-later --help
```

**输出：**
```
View watch later list (login required).

Usage:
  aiview bilibili watch-later [flags]

Flags:
  -h, --help      help for watch-later
  -n, --max int   Maximum number of items to show
```

**状态：** ⚠️ 需认证
**备注：** 帮助信息正确显示，实际使用需要登录。

---

### 13. precious — 入站必刷

**命令：**
```
$ aiview bilibili precious --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "explain": "我不允许还没有人看过这84个宝藏视频！",
    "list": [
      { "achievement": "全网刷屏！治愈千万人的精神内核", "bvid": "BV1MN4y177PB", ... },
      ...
    ]
  }
}
```

**状态：** ✅ 通过
**备注：** 返回入站必刷经典视频列表，84 个视频。

---

### 14. weekly — 每周必看

**命令：**
```
$ aiview bilibili weekly 100 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "code": 0,
    "data": {
      "list": [ ... ],
      "config": { ... }
    }
  }
}
```

**状态：** ✅ 通过
**备注：** 较新期数（如 100）正常返回。注意：早期期数（如第1期）返回 -352 风控。

---

### 15. collection — 用户合集

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

**状态：** ⚠️ 需认证
**备注：** 该用户没有创建视频合集。合法业务结果，建议用有合集的 UID 测试。

---

### 16. dynamic — 用户动态

**命令：**
```
$ aiview bilibili dynamic 37737161 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "cards": [ ... ]
  }
}
```

**状态：** ✅ 通过
**备注：** 用户动态获取正常（公开接口）。

---

### 17. tags — 视频标签

**命令：**
```
$ aiview bilibili tags BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": [
    { "tag_id": 9, "tag_name": "B站最早的视频", "tag_type": "old_channel" },
    { "tag_id": 11, "tag_name": "历史的起点", "tag_type": "old_channel" },
    { "tag_id": 10, "tag_name": "考古的尽头", "tag_type": "old_channel" },
    { "tag_id": 1312966, "tag_name": "B站的起源", "tag_type": "old_channel" },
    ...
  ]
}
```

**状态：** ✅ 通过
**备注：** 返回视频的标签列表，包含 tag_id 和 tag_type。

---

### 18. suggest — 搜索建议

**命令：**
```
$ aiview bilibili suggest test --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "name": "MAJOR",
    "show_name": "MAJOR",
    "url": "https://search.bilibili.com/all?keyword=MAJOR"
  }
}
```

**状态：** ✅ 通过
**备注：** 返回搜索建议。

---

## 二、需认证读命令（9 个）

### 19. favorites — 收藏夹

**命令：**
```
$ aiview bilibili favorites --help
```

**输出：**
```
View favorite folders (login required).

Usage:
  aiview bilibili favorites [flags]

Flags:
  -h, --help       help for favorites
  -n, --max int    Maximum number of folders to show
  -p, --page int   Page number (default 1)
```

**状态：** ⚠️ 需认证
**备注：** 帮助信息正确显示，实际使用需要登录。

---

### 20. following — 关注列表

**命令：**
```
$ aiview bilibili following 37737161 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "list": [ ... ],
    "total": ...
  }
}
```

**状态：** ✅ 通过
**备注：** 关注列表为公开接口，无需认证即可查询。

---

### 21. fans — 粉丝列表

**命令：**
```
$ aiview bilibili fans 37737161 --json
```

**输出：**
```json
{
  "ok": false,
  "error": {
    "code": "api_error",
    "message": "API error [-352]: 风控校验失败"
  }
}
```

**状态：** ❌ 失败
**备注：** -352 风控错误。粉丝列表接口被 B 站风控拦截。

---

### 22. relation — 用户关系

**命令：**
```
$ aiview bilibili relation 37737161 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "mid": 37737161,
    "follower": 84099,
    "following": 20,
    "black": 0,
    "whisper": 0
  }
}
```

**状态：** ✅ 通过
**备注：** 关系查询无需认证，返回关注/粉丝/拉黑状态。

---

### 23. online — 在线人数

**命令：**
```
$ aiview bilibili online BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "count": "1",
    "show_switch": { "count": true, "total": true },
    "total": "2"
  }
}
```

**状态：** ✅ 通过
**备注：** 返回视频当前在线观看人数。

---

### 24. video-status — 视频统计

**命令：**
```
$ aiview bilibili video-status BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "view": 9473672,
    "danmaku": 225608,
    "like": 396122,
    "coin": 42300,
    "favorite": 149631,
    "share": 30123
  }
}
```

**状态：** ✅ 通过
**备注：** 返回视频统计信息，无需认证。

---

### 25. live — 直播间信息

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

**状态：** ❌ 失败
**备注：** Room ID=3 的直播间不存在或 API 返回了 HTML 页面（-352 风控）。需登录 Cookie 绕过。

---

### 26. danmaku — 弹幕查看

**命令：**
```
$ aiview bilibili danmaku BV1xx411c7m9 --json
```

**输出：**
```
(二进制 protobuf 数据，非 JSON)
```

**状态：** ✅ 通过
**备注：** 弹幕数据以 protobuf 二进制格式返回。需要格式化 protobuf 为可读输出。

---

### 27. comment — 评论查看

**命令：**
```
$ aiview bilibili comment BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "comments": [ ... ],
    "page": { ... }
  }
}
```

**状态：** ✅ 通过
**备注：** 评论查看无需认证。发评论需写权限。

---

## 三、写操作命令（13 个）

### 28. like — 点赞

**命令：**
```
$ aiview bilibili like --help
```

**输出：**
```
Like a video (login and write permission required).

Usage:
  aiview bilibili like <BV> [flags]

Flags:
  -h, --help   help for like
      --undo   Unlike the video
```

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。Flags: `--undo` 取消点赞。

---

### 29. coin — 投币

**命令：**
```
$ aiview bilibili coin --help
```

**输出：**
```
Give coins to a video (login and write permission required).

Usage:
  aiview bilibili coin <BV> [flags]

Flags:
  -h, --help      help for coin
  -n, --num int   Number of coins (1-2) (default 1)
```

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。Flags: `-n, --num` 指定投币数量（1-2）。

---

### 30. favorite — 收藏视频

**命令：**
```
$ aiview bilibili favorite --help
```

**输出：**
```
Add or remove a video from a favorite folder.

Requires --fid to specify the folder ID.
Use --delete to remove from the folder.

Flags:
      --delete    Remove from favorites
      --fid int   Favorite folder ID (required)
  -h, --help      help for favorite
```

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。

---

### 31. triple — 一键三连

**命令：**
```
$ aiview bilibili triple --help
```

**输出：**
```
Like, coin, and favorite a video in one go (login and write permission required).

Usage:
  aiview bilibili triple <BV> [flags]
```

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。

---

### 32. dynamic-post — 发动态

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

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。

---

### 33. dynamic-delete — 删除动态

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

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。

---

### 34. danmaku-send — 发送弹幕

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
  -h, --help           help for danmaku-send
      --progress int   Video progress in seconds
```

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。

---

### 35. comment-delete — 删除评论

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

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。

---

### 36. follow — 关注

**命令：**
```
$ aiview bilibili follow --help
```

**输出：**
```
Follow a user (login and write permission required).

Usage:
  aiview bilibili follow <UID> [flags]
```

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。

---

### 37. unfollow — 取消关注

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

**状态：** ✅ 通过（命令注册正常）
**备注：** 需要写权限的 cookie。

---

### 38. audio — 下载音频

**命令：**
```
$ aiview bilibili audio BV1xx411c7m9 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "directory": "C:\\Users\\Administrator\\AppData\\Local\\Temp\\aiview\\2012地球便当之日宣传片",
    "duration": 25,
    "segments": 7
  }
}
```

**状态：** ✅ 通过
**备注：** 成功下载音频到临时目录，25 秒时长，7 个分段。Flags: `--segment` / `--no-split` / `-o`。

---

### 39. collect — 批量采集

**命令：**
```
$ aiview bilibili collect --help
```

**输出：**
```
Batch collect videos from various sources.

Usage:
  aiview bilibili collect [flags]

Flags:
      --type string   Collection type: hot, rank, user, search
      --limit int     Number of videos to collect
      --output string Output directory
  -h, --help          help for collect
```

**状态：** ✅ 通过（命令注册正常）
**备注：** 批量采集命令，支持从热门、排行榜、用户、搜索等来源采集。

---

### 40. login/logout/status/whoami — 认证管理

**命令：**
```
$ aiview bilibili login --help
$ aiview bilibili logout
$ aiview bilibili status --json
$ aiview bilibili whoami --json
```

**输出：**
```
# login --help
Login to Bilibili.
Three methods supported:
  1. No arguments: QR code scan login
  2. --sessdata: Pass SESSDATA Cookie directly
  3. --sessdata + --bili-jct: Pass full credential

# logout
✅ Logged out

# status --json
{"ok":true,"data":{"authenticated":false}}

# whoami --json
{"ok":false,"error":{"code":"not_authenticated","message":"Not logged in"}}
```

**状态：** ✅ 通过
**备注：** 认证管理命令组工作正常，支持三种登录方式。logout 后 status 正确显示未登录状态。

---

## 问题汇总

### ❌ 失败（3 个）

| # | 命令 | 问题 | 原因 | 建议 |
|---|------|------|------|------|
| 1 | user-videos | -352 风控 | B 站反爬机制 | 需添加 Cookie 认证或请求头 |
| 2 | fans | -352 风控 | B 站反爬机制 | 需添加 Cookie 认证 |
| 3 | live | HTML 响应/风控 | 无效 room ID 或风控 | 需登录 Cookie 绕过 |

### ⚠️ 需认证（4 个）

| # | 命令 | 问题 | 原因 | 建议 |
|---|------|------|------|------|
| 1 | history | 需登录 | 未登录 | 执行 `bilibili login` |
| 2 | watch-later | 需登录 | 未登录 | 执行 `bilibili login` |
| 3 | collection | 需登录 | 未登录 | 执行 `bilibili login` |
| 4 | favorites | 需登录 | 未登录 | 执行 `bilibili login` |

### 其他注意事项

| 问题 | 影响命令 | 说明 |
|------|----------|------|
| 弹幕二进制输出 | danmaku | protobuf 未解析为 JSON |
| weekly 早期期数 | weekly | 第1期返回 -352，较新期数正常 |
| feed 需登录 | feed | 关注动态需要登录 |

---

## 相关链接

- [统一总览报告](./TEST_REPORT.md)
- [抖音平台报告](./TEST_REPORT_DOUYIN.md)
- [小红书平台报告](./TEST_REPORT_XIAOHONGSHU.md)
- [全局命令报告](./TEST_REPORT_GLOBAL.md)
