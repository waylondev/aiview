# Aiview CLI — 抖音（Douyin）平台测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 11 |
| ✅ 通过 | 6 |
| ⚠️ 需认证 | 5 |
| ❌ 失败 | 0 |

---

## 测试详情

### 1. hot — 热搜榜

**命令：**
```
$ aiview douyin hot --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "data": [
      {"hot_value": 7717101, "word": "世界杯精彩程度下降了吗", "position": 31, ...},
      {"hot_value": 7680387, "word": "这个榴莲馒头看起来很诱人", "position": 42, ...},
      ... (共 50 条)
      {"hot_value": 7655134, "word": "石纪元登月篇千空成功登月", "position": 49, ...}
    ],
    "extra": {
      "logid": "202606131039582004EADB45EF68BF671F",
      "now": 1781318398000,
      "time_cost": {"stream_inner": 179}
    },
    "status_code": 0
  }
}
```

**状态：** ✅ 通过
**备注：** 成功返回 50 条热搜数据，包含 hot_value、word、position、label、word_cover 等字段。耗时约 179ms。

---

### 2. trending — 热点榜

**命令：**
```
$ aiview douyin trending --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "banner_dark": {"uri": "32e5200049fac86168117", "url_list": [...]},
    "banner_light": {"uri": "321f4000f703369c24604", "url_list": [...]},
    "data": {
      "active_time": "2026-06-13 10:37:14",
      "trending_desc": "实时上升热点",
      "word_list": [
        {"word": "美加墨世界杯美国开幕式", "hot_value": 10454035, "position": 5, ...},
        {"word": "河南3岁男童走失已超40小时", "hot_value": 7699158, "position": 37, ...},
        ...
      ]
    },
    "extra": {
      "now": 1781318412000,
      "time_cost": {"stream_inner": 226}
    },
    "status_code": 0
  }
}
```

**状态：** ✅ 通过
**备注：** 成功返回热点榜数据，包含 word_list 列表、banner 图片、share_info 分享信息和 active_time 时间戳。耗时约 226ms。

---

### 3. search — 搜索

**命令：**
```
$ aiview douyin search test --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "cursor": 0,
    "data": [],
    "has_more": 0,
    "search_nil_info": {
      "is_load_more": "first_flush",
      "search_nil_item": "invalid_app",
      "search_nil_type": "params_check",
      "text_type": 60
    },
    "status_code": 0,
    "time_cost": {"stream_inner": 59}
  }
}
```

**状态：** ⚠️ 需认证
**备注：** 返回有效 JSON，但搜索结果为空数组。`search_nil_info.search_nil_type` 为 `"params_check"`，`search_nil_item` 为 `"invalid_app"`，表明 API 需要额外的客户端参数或 Cookie 认证。

---

### 4. video — 视频详情

**命令：**
```
$ aiview douyin video 123456789 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "aweme_detail": null,
    "filter_detail": {
      "aweme_id": "123456789",
      "detail_msg": "",
      "filter_reason": "core_dep",
      "icon": "",
      "notice": ""
    },
    "status_code": 0
  }
}
```

**状态：** ⚠️ 需认证
**备注：** 使用假 ID 返回 `aweme_detail: null`，`filter_reason` 为 `"core_dep"`，表明需要 Cookie 登录后才能获取视频详情。

---

### 5. user — 用户信息

**命令：**
```
$ aiview douyin user 123456789 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "raw": ""
  }
}
```

**状态：** ⚠️ 需认证
**备注：** 使用假 UID 返回 `{"raw":""}`，表明需要 Cookie 登录后才能获取用户信息。

---

### 6. user-posts — 用户作品

**命令：**
```
$ aiview douyin user-posts 123456789 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "raw": ""
  }
}
```

**状态：** ⚠️ 需认证
**备注：** 返回 `{"raw":""}`，需要 Cookie 登录。

---

### 7. comment — 评论列表

**命令：**
```
$ aiview douyin comment 123456789 --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "raw": ""
  }
}
```

**状态：** ⚠️ 需认证
**备注：** 返回 `{"raw":""}`，需要 Cookie 登录。

---

### 8. login — 登录

**命令：**
```
$ aiview douyin login --help
```

**输出：**
```
Login to Douyin using a browser cookie.

Examples:
  aiview douyin login --cookie "your_cookie_here"

Usage:
  aiview douyin login [flags]

Flags:
      --cookie string   Douyin browser cookie
  -h, --help            help for login
```

**状态：** ✅ 通过
**备注：** 帮助信息完整，支持 `--cookie` 参数进行 Cookie 登录。

---

### 9. logout — 登出

**命令：**
```
$ aiview douyin logout
```

**输出：**
```
Logged out
```

**状态：** ✅ 通过
**备注：** 输出 "Logged out"，退出码为 0。登出后再次运行 `status` 确认仍为未登录状态。

---

### 10. status — 登录状态

**命令：**
```
$ aiview douyin status --json
```

**输出：**
```json
{
  "ok": true,
  "schema_version": "1",
  "data": {
    "logged_in": false,
    "platform": "douyin"
  }
}
```

**状态：** ✅ 通过
**备注：** 正确反映当前未登录状态。

---

### 11. collect — 批量采集

**命令：**
```
$ aiview douyin collect --help
```

**输出：**
```
Batch collect videos from Douyin.

Usage:
  aiview douyin collect [flags]

Flags:
      --type string   Collection type: hot, trending, search
      --limit int     Number of videos to collect
      --output string Output directory
  -h, --help          help for collect
```

**状态：** ✅ 通过
**备注：** 批量采集命令，支持从热门、热点、搜索等来源采集。

---

## 问题汇总

### ⚠️ 需认证（5 个）

| # | 命令 | 问题 | 原因 | 建议 |
|---|------|------|------|------|
| 1 | search | 搜索无结果 | 未登录，API 返回空 data | 执行 `douyin login --cookie` |
| 2 | video | 视频详情为空 | 未登录，filter_reason=core_dep | 执行 `douyin login --cookie` |
| 3 | user | 用户信息为空 | 未登录，返回 raw="" | 执行 `douyin login --cookie` |
| 4 | user-posts | 用户作品为空 | 未登录，返回 raw="" | 执行 `douyin login --cookie` |
| 5 | comment | 评论列表为空 | 未登录，返回 raw="" | 执行 `douyin login --cookie` |

### 已知问题

1. **搜索功能无返回结果**：`search` 命令虽然返回有效 JSON，但 `data` 字段为空数组。响应的 `search_nil_info` 显示 `search_nil_type: "params_check"`、`search_nil_item: "invalid_app"`，说明 API 需要额外的客户端参数或 Cookie 认证。
2. **需认证的 API 返回格式不一致**：`video` 命令返回 `{"aweme_detail": null, "filter_detail": {...}}`，而 `user`、`comment`、`user-posts` 返回 `{"raw":""}`。两种格式都不友好，建议统一返回明确的错误提示信息。
3. **`raw` 字段语义不明确**：`user`、`comment`、`user-posts` 认证失败时返回 `{"raw":""}` 而非错误消息，不利于用户诊断问题。

### 修复计划

1. **P0**：统一认证失败时的错误返回格式，返回明确的 `not_authenticated` 错误码和提示信息
2. **P1**：确保 Cookie 登录后所有命令可用
3. **P2**：优化 `raw` 字段的解析逻辑，提供结构化的输出

---

## 相关链接

- [统一总览报告](./TEST_REPORT.md)
- [Bilibili 平台报告](./TEST_REPORT_BILIBILI.md)
- [小红书平台报告](./TEST_REPORT_XIAOHONGSHU.md)
- [全局命令报告](./TEST_REPORT_GLOBAL.md)
