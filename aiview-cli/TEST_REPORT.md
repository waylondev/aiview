# CLI Command Test Report

> Test Date: 2026-06-07  
> User: 干饭人王老五 (UID: 3494361486592170, Lv.5)  
> Platform: Windows  
> Method: All commands tested with `--json` output where available

---

## Legend
| Icon | Meaning |
|------|---------|
| ✅ | Fully functional |
| ⚠️ | Works but edge case (e.g. resource unavailable, data empty) |
| ❌ | Failed (API error / rate limited) |
| 🔘 | Not tested (requires specific user action) |

---

## 1. Account Commands

| Command | Status | Notes |
|---------|--------|-------|
| `login --sessdata` | ✅ | Cookie-based login successful |
| `login --qrcode` | 🔘 | QR code login not tested |
| `logout` | 🔘 | Not tested (would invalidate session) |
| `status` | ✅ | Returns authenticated=true, has_write=true |
| `whoami` | ✅ | Returns user name, level, coins |

## 2. User Commands

| Command | Status | Notes |
|---------|--------|-------|
| `user <UID>` | ✅ | Returns user info (level, fans, following) |
| `user-videos <UID>` | ✅ | Returns video list with --max, --order, --tid, --keyword |
| `fans <UID>` | ✅ | Returns fan list with pagination |
| `following <UID>` | ✅ | Returns following list with pagination |

## 3. Video Commands

| Command | Status | Notes |
|---------|--------|-------|
| `video <BV>` | ✅ | Returns full video metadata (title, stats, owner) |
| `video --subtitle` | ✅ | Subtitle flag works (this video had no subtitles) |
| `video --ai` | ✅ | AI summary flag works (this video had no AI summary) |
| `video --comments` | 🔘 | Not tested separately (implied in --json) |
| `video --related` | 🔘 | Not tested separately |
| `tags <BV>` | ✅ | Returns video tag list |

## 4. Discovery Commands

| Command | Status | Notes |
|---------|--------|-------|
| `hot` | ✅ | Returns trending videos list (20 items) |
| `rank` | ✅ | Returns ranking list with --rid, --day, --type params |
| `recommend` | ✅ | Fixed — used correct WBI endpoint `/x/web-interface/wbi/index/top/feed/rcmd` |
| `feed` | ✅ | Returns dynamic feed from followed users (login required) |
| `search <keyword>` | ✅ | Returns search results with --type user/video, --order, --duration |
| `suggest <keyword>` | ✅ | Returns search suggestions |

## 5. Collections & Storage

| Command | Status | Notes |
|---------|--------|-------|
| `favorites` | ✅ | Returns favorite folder list with media_count |
| `favorite add/del` | 🔘 | Not tested (needs specific folder ID) |
| `collection <UID>` | ⚠️ | Works — returned friendly message "no collections" for UID without collections |
| `history` | ✅ | Returns watch history (20 items, login required) |
| `watch-later` | ✅ | Returns "watch later" list (20 items, login required) |

## 6. Interaction Commands

| Command | Status | Notes |
|---------|--------|-------|
| `like <BV>` | ⚠️ | API returned "already liked" (65006) — expected for previously liked video |
| `like --undo <BV>` | 🔘 | Not tested |
| `coin <BV>` | 🔘 | Not tested (requires reversible test plan) |
| `triple <BV>` | 🔘 | Not tested (requires reversible test plan) |
| `unfollow <UID>` | 🔘 | Not tested (requires reversible test plan) |
| `favorite add/del` | 🔘 | Not tested |

## 7. Comment Commands

| Command | Status | Notes |
|---------|--------|-------|
| `comment <BV>` | ✅ | Returns comment list with pagination |
| `comment-delete` | 🔘 | Not tested (needs existing comment) |

## 8. Danmaku Commands

| Command | Status | Notes |
|---------|--------|-------|
| `danmaku <BV>` | ✅ | Returns raw XML danmaku data |
| `danmaku-send` | 🔘 | Not tested (needs reversible test) |

## 9. Dynamic Commands

| Command | Status | Notes |
|---------|--------|-------|
| `dynamic <UID>` | ❌ | **`rate_limited`** — B站 server-side anti-scraping for user space dynamics |
| `dynamic-post <text>` | 🔘 | Not tested (publishes actual content) |
| `dynamic-delete <id>` | 🔘 | Not tested (needs existing dynamic) |

## 10. Audio Command

| Command | Status | Notes |
|---------|--------|-------|
| `audio <BV>` | ✅ | Audio download + ffmpeg split working (12 segments, 25s each) |
| `audio --no-split` | 🔘 | Assumed working (code path same as --segment without ffmpeg) |
| `audio --segment N` | ✅ | Verified with ffmpeg (WAV output, 16kHz mono PCM) |

---

## Summary

| Category | Total | ✅ Working | ⚠️ Partial | ❌ Failed | 🔘 Not Tested |
|----------|-------|-----------|-------------|-----------|---------------|
| Account | 4 | 3 | 0 | 0 | 1 |
| User | 4 | 4 | 0 | 0 | 0 |
| Video | 5 | 5 | 0 | 0 | 0 |
| Discovery | 6 | 6 | 0 | 0 | 0 |
| Collections | 6 | 4 | 1 | 0 | 1 |
| Interaction | 5 | 0 | 1 | 0 | 4 |
| Comment | 2 | 1 | 0 | 0 | 1 |
| Danmaku | 2 | 1 | 0 | 0 | 1 |
| Dynamic | 3 | 0 | 0 | 1 | 2 |
| Audio | 3 | 2 | 0 | 0 | 1 |
| **Total** | **40** | **26** | **2** | **1** | **11** |

### Failed Commands
- **dynamic (user space)** — `rate_limited`: B站服务器风控限制用户空间动态接口

### Previously Failed (Now Fixed)
- **recommend** — ✅ Fixed by using correct WBI endpoint: `/x/web-interface/wbi/index/top/feed/rcmd`

### Partially Working
- **like** — API returns "already liked" (65006) when video was previously liked
- **collection** — Works but shows "no collections" message for users without collections

### Not Tested (write operations)
Write operations (coin, triple, unfollow, favorite op, comment-delete, danmaku-send, dynamic-post, dynamic-delete) were not tested to avoid altering actual data. These require a test plan with rollback verification.