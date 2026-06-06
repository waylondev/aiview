# Aiview CLI Bilibili API 测试报告

**测试日期**: 2026-06-06  
**测试环境**: Windows x64, Go 1.26.4  
**测试 BV 号**: BV1GJ411x7h7 (Never Gonna Give You Up)  
**测试 UID**: 946974 (影视飓风)

---

## 1. 命令测试汇总

| # | 命令 | API 端点 | 结果 | 说明 |
|---|------|----------|------|------|
| 1 | `video BV1GJ411x7h7` | `/x/web-interface/view` | ✅ | 视频详情完美返回（标题、UP主、播放量、弹幕等） |
| 2 | `video --subtitle` | `/x/player/v2` | ⚠️ | 该视频无字幕，返回空 |
| 3 | `video --ai` | `/x/player/v2` | ⚠️ | 该视频无 AI 总结 |
| 4 | `video --comments` | `/x/v2/reply` | ✅ | 返回 20 条热门评论 |
| 5 | `video --related` | `/x/web-interface/archive/related` | ✅ | 返回 40 条相关推荐 |
| 6 | `search "Python" --type video` | `/x/web-interface/wbi/search/type` | ✅ | WBI 签名生效，搜索结果正常 |
| 7 | `search "test" --type user` | `/x/web-interface/wbi/search/type` | ✅ | 用户搜索正常 |
| 8 | `user 946974` | `/x/web-interface/card` | ✅ | 用户信息正常（名称、等级、签名、粉丝、关注数） |
| 9 | `user-videos 946974` | `/x/space/wbi/arc/search` | ❌ | 返回 HTML 页面（反爬/IP风控） |
| 10 | `hot --max 3` | `/x/web-interface/popular` | ✅ | 热门视频列表正常 |
| 11 | `rank --max 3` | `/x/web-interface/ranking/v2` | ✅ | 排行榜正常 |
| 12 | `status` | 本地检查 | ✅ | 认证状态检查正常 |
| 13 | `login` | 本地提示 | ✅ | 显示登录引导信息 |
| 14 | `logout` | 本地清除 | ✅ | 退出登录正常 |
| 15 | `whoami` | `/x/web-interface/nav` | ⚠️ | 需要登录凭证 |
| 16 | `favorites` | 需登录 | ⚠️ | 需登录，提示用户 login |
| 17 | `following` | 需登录 | ⚠️ | 需登录，提示用户 login |
| 18 | `history` | 需登录 | ⚠️ | 需登录，提示用户 login |
| 19 | `watch-later` | 需登录 | ⚠️ | 需登录，提示用户 login |
| 20 | `feed` | 需登录 | ⚠️ | 需登录，提示用户 login |
| 21 | `like` | 需登录+写权限 | ⚠️ | 需登录，提示用户 login |
| 22 | `coin` | 需登录+写权限 | ⚠️ | 需登录，提示用户 login |
| 23 | `triple` | 需登录+写权限 | ⚠️ | 需登录，提示用户 login |
| 24 | `unfollow` | 需登录+写权限 | ⚠️ | 需登录，提示用户 login |
| 25 | `audio BVxxx` | `/x/player/playurl` | ⚠️ | 需有 DASH 格式音源 |

---

## 2. API 修复记录

### 2.1 WBI 签名系统 (新增)
- **问题**: 搜索、用户信息等 API 返回 -352 (风控校验失败) 或 HTML 反爬页面
- **修复**: 新增 [wbi.go](file:///d:/Users/Administrator/workspace/aiview/internal/platform/bilibili/wbi.go)，实现完整的 WBI 签名
  - 从 `/x/web-interface/nav` 获取 `img_key` 和 `sub_key`
  - 使用 `MIXIN_KEY_ENC_TAB` (64 位) 生成 mixed key (32 位)
  - MD5 签名计算 `w_rid` + `wts`
  - 大写 hex 编码，过滤 `!'()*` 特殊字符
- **影响**: 修复了 `search` (视频/用户) 命令

### 2.2 用户信息 API 切换
- **问题**: `/x/space/wbi/acc/info` 返回 -352 风控错误
- **修复**: 切换到 `/x/web-interface/card` 公开接口
- **影响**: `user` 命令正常工作，返回昵称、等级、签名、粉丝、关注数

### 2.3 用户视频 API 待修复
- **问题**: `/x/space/wbi/arc/search` 返回 HTML 反爬页面
- **原因**: B站空间 API 反爬策略严格，需要 Cookie 或完整浏览器环境
- **状态**: 需要登录凭证才能正常调用

---

## 3. 功能覆盖率

### ✅ 完全可用（无需登录）
| 命令 | 功能 |
|------|------|
| `video` | 视频详情、评论、相关推荐 |
| `search` | 视频搜索、用户搜索 |
| `user` | 用户信息 |
| `hot` | 热门视频 |
| `rank` | 排行榜 |
| `status` | 登录状态检查 |

### ⚠️ 需要登录凭证
| 命令 | 功能 |
|------|------|
| `whoami` | 当前用户信息 |
| `favorites` | 收藏夹 |
| `following` | 关注列表 |
| `history` | 观看历史 |
| `watch-later` | 稍后再看 |
| `feed` | 动态时间线 |
| `like` | 点赞 |
| `coin` | 投币 |
| `triple` | 一键三连 |
| `unfollow` | 取消关注 |

### ❌ 待修复
| 命令 | 问题 |
|------|------|
| `user-videos` | 空间 API 反爬，返回 HTML |
| `audio` | 需要 DASH 格式音源 |

### ⚠️ 部分功能受限
| 命令/Flag | 限制 |
|-----------|------|
| `video --subtitle` | 仅返回有字幕的视频，否则为空 |
| `video --ai` | 仅返回有 AI 总结的视频 |

---

## 4. 技术债务

1. **user-videos 命令**: 需要实现 Cookie 传递机制来绕过反爬
2. **QR 码登录**: 当前仅显示文字引导，需要实现终端 QR 码显示
3. **音频下载**: 完整的分割和进度显示功能
4. **WBI 密钥缓存**: 已实现 1 小时缓存，可优化
5. **错误码映射**: 新增 API 的错误码需要完善映射

---

## 5. 已验证的 WBI 签名系统

新增文件 [wbi.go](file:///d:/Users/Administrator/workspace/aiview/internal/platform/bilibili/wbi.go)：

| 组件 | 说明 |
|------|------|
| `getWBIKey()` | 从 nav API 获取每日轮换密钥（1小时缓存） |
| `mixWBIKeys()` | 使用 MIXIN_KEY_ENC_TAB 生成 32 位 mixed key |
| `signWBI()` | 参数排序 → URL 编码(大写hex+%20) → +mixed key → MD5 |
| `wbiGet()` | 带 WBI 签名的 GET 请求 |
| `wbiEncode()` | 自定义 URL 编码（大写hex，%20代替+） |
| `filterSpecialChars()` | 过滤 `!'()*` 字符 |

已验证能正常工作的 WBI 接口：
- `/x/web-interface/wbi/search/type` — 搜索
- `/x/web-interface/card` — 用户卡片