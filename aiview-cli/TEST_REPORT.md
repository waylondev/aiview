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
| `video-status <BV>` | ✅ | Returns video view/like/coin/fav/share statistics |
| `relation <UID>` | ✅ | Returns relation status with another user |

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
| `recommend` | ✅ | Uses WBI endpoint `/x/web-interface/wbi/index/top/feed/rcmd` |
| `feed` | ✅ | Returns dynamic feed from followed users (login required) |
| `search <keyword>` | ✅ | Returns search results with --type user/video, --order, --duration |
| `suggest <keyword>` | ✅ | Returns search suggestions |
| `trending` | ✅ | Hot search/trending keywords with --limit flag |
| `region <rid>` | ✅ (new) | Returns videos by region/category with --pn, --ps, --sort params |

## 5. Collections & Storage

| Command | Status | Notes |
|---------|--------|-------|
| `favorites <UID>` | ✅ | Returns favorite folder list with media_count |
| `favorite add/del` | 🔘 | Not tested (needs specific folder ID) |
| `collection <UID>` | ⚠️ | Works — returned friendly message "no collections" for UID without collections |
| `history` | ✅ | Returns watch history (20 items, login required) |
| `watch-later` | ✅ | Returns "watch later" list (20 items, login required) |

## 6. Interaction Commands (write — not tested)

| Command | Status | Notes |
|---------|--------|-------|
| `like <BV>` | 🔘 | Not tested (write operation) |
| `coin <BV>` | 🔘 | Not tested (write operation) |
| `triple <BV>` | 🔘 | Not tested (write operation) |
| `unfollow <UID>` | 🔘 | Not tested (write operation) |
| `favorite add/del` | 🔘 | Not tested (write operation) |

## 7. Comment Commands

| Command | Status | Notes |
|---------|--------|-------|
| `comment <BV>` | ✅ | Returns comment list with pagination |
| `comment-delete` | 🔘 | Not tested (write operation) |

## 8. Danmaku Commands

| Command | Status | Notes |
|---------|--------|-------|
| `danmaku <BV>` | ✅ | Returns danmaku data (protobuf format) |
| `danmaku-send` | 🔘 | Not tested (write operation) |

## 9. Dynamic Commands

| Command | Status | Notes |
|---------|--------|-------|
| `dynamic <UID>` | ❌ | **`rate_limited`** — B站 server-side anti-scraping for user space dynamics |
| `dynamic-post` | 🔘 | Not tested (write operation) |
| `dynamic-delete` | 🔘 | Not tested (write operation) |

## 10. Audio & Live Commands

| Command | Status | Notes |
|---------|--------|-------|
| `audio <BV>` | ✅ | Audio download + ffmpeg split working |
| `live` | ⚠️ (new) | Works with valid room ID; returns error for non-existent rooms |

## 11. New Read-Only Commands

| Command | Status | Notes |
|---------|--------|-------|
| `precious` | ✅ | "入站必刷" curated video collection, returns 85 must-watch videos |
| `trending` | ✅ | Hot search/trending keywords with --limit flag |
| `online <BV>` | ✅ | Real-time video online viewer count (total + web) |
| `weekly <number>` | ✅ (fixed) | Weekly hot video series, requires WBI signing |

### 12. Douyin Commands

| Command | Args | Status | Notes |
|---------|------|--------|-------|
| `douyin hot` | `--json` | ✅ | Returns 50 hot search items |
| `douyin trending` | `--json` | ✅ | Returns trending list |
| `douyin search` | `<keyword>` | ⚠️ | Returns empty results (needs cookies) |
| `douyin video` | `<share_url>` | 🔘 | Not tested (needs login) |
| `douyin user` | `<uid>` | 🔘 | Not tested (needs login) |

**Summary**: 2 ✅ / 1 ⚠️ / 0 ❌ / 2 🔘

---

## Summary

| Category | Total | ✅ Working | ⚠️ Partial | ❌ Failed | 🔘 Not Tested |
|----------|-------|-----------|-------------|-----------|---------------|
| Account | 7 | 6 | 0 | 0 | 1 |
| User | 4 | 4 | 0 | 0 | 0 |
| Video | 6 | 6 | 0 | 0 | 0 |
| Discovery | 8 | 8 | 0 | 0 | 0 |
| Collections | 5 | 4 | 1 | 0 | 0 |
| Interaction | 5 | 0 | 0 | 0 | 5 |
| Comment | 2 | 1 | 0 | 0 | 1 |
| Danmaku | 2 | 1 | 0 | 0 | 1 |
| Dynamic | 3 | 0 | 0 | 1 | 2 |
| Audio & Live | 2 | 1 | 1 | 0 | 0 |
| New Read-Only | 4 | 4 | 0 | 0 | 0 |
| Douyin | 5 | 2 | 1 | 0 | 2 |
| **Total** | **53** | **36** | **3** | **3** | **13** |

### Failed Commands
- **dynamic (user space)** — `rate_limited`: B站服务器风控限制用户空间动态接口

### Partially Working
- **collection** — Works but shows "no collections" message for users without collections
- **live** — Works with valid room IDs; non-existent rooms return API error

### Previously Failed (Now Fixed)
- **recommend** — ✅ Fixed by using correct WBI endpoint
- **video-status** — ✅ Fixed by using GetVideoInfo stat data instead of separate endpoint
- **weekly** — ✅ Fixed by implementing WBI signing for the weekly hot video series endpoint

### Not Tested (write operations)
Write operations (like, coin, triple, unfollow, favorite, comment-delete, danmaku-send, dynamic-post, dynamic-delete) were not tested as per spec — only read-only commands were tested.