# Aiview CLI — Bilibili CLI Tool

A feature-rich command-line tool for interacting with Bilibili, built in Go with Cobra. Supports video browsing, user queries, social interactions, audio downloading, and more.

## Quick Start

```bash
# Login with cookie
aiview bilibili login --sessdata "<SESSDATA>" --bili-jct "<BILI_JCT>"

# Check status
aiview bilibili status

# Browse a video
aiview bilibili video BV1GJ411x7Rq --json

# View user info
aiview bilibili user 37737161
```

## Commands

### ✅ Account
| Command | Description | Status |
|---------|-------------|--------|
| `login --sessdata` | Login with cookie | ✅ Done |
| `login --qrcode` | QR code login | ✅ Done |
| `logout` | Clear saved credentials | ✅ Done |
| `status` | Check login status | ✅ Done |
| `whoami` | Show current user info | ✅ Done |

### ✅ User
| Command | Description | Status |
|---------|-------------|--------|
| `user <UID>` | View user info | ✅ Done |
| `user-videos <UID>` | View user's video list | ✅ Done |
| `fans <UID>` | View fans list | ✅ Done |
| `following <UID>` | View following list | ✅ Done |

### ✅ Video
| Command | Flags | Status |
|---------|-------|--------|
| `video <BV>` | `--subtitle, --ai, --comments, --related` | ✅ Done |
| `tags <BV>` | — | ✅ Done |

### ✅ Discovery
| Command | Description | Status |
|---------|-------------|--------|
| `hot` | View trending videos | ✅ Done |
| `rank` | View rankings | ✅ Done |
| `feed` | View dynamic feed (login required) | ✅ Done |
| `search <keyword>` | Search videos/users | ✅ Done |
| `suggest <keyword>` | Get search suggestions | ✅ Done |
| `recommend` | Get homepage recommendations | ✅ Done |

### ✅ Collections & Storage
| Command | Description | Status |
|---------|-------------|--------|
| `favorites` | View favorite folders | ✅ Done |
| `favorite add/del <BV>` | Add/remove video from favorites | ✅ Done |
| `collection <UID>` | View user's video collections | ✅ Done |
| `history` | View watch history | ✅ Done |
| `watch-later` | View watch later list | ✅ Done |

### ✅ Interactions
| Command | Description | Status |
|---------|-------------|--------|
| `like <BV>` | Like/unlike a video | ✅ Done |
| `coin <BV>` | Give coins to a video | ✅ Done |
| `triple <BV>` | Like + Coin + Favorite | ✅ Done |
| `unfollow <UID>` | Unfollow a user | ✅ Done |

### ✅ Comments & Danmaku
| Command | Description | Status |
|---------|-------------|--------|
| `comment <BV>` | View comments | ✅ Done |
| `comment-delete <ID>` | Delete a comment | ✅ Done |
| `danmaku <BV>` | View danmaku | ✅ Done |
| `danmaku-send <BV>` | Send a danmaku | ✅ Done |

### ✅ Dynamics
| Command | Description | Status |
|---------|-------------|--------|
| `dynamic <UID>` | View user's dynamics | ❌ B站风控限流 |
| `dynamic-post <text>` | Post a text dynamic | ✅ Done |
| `dynamic-delete <id>` | Delete a dynamic | ✅ Done |

### ✅ Audio
| Command | Description | Status |
|---------|-------------|--------|
| `audio <BV>` | Download audio and split into WAV (requires ffmpeg) | ✅ Done |

## Global Flags

```bash
--json    Output in JSON format
--yaml    Output in YAML format
-v, --verbose  Enable verbose logging
```

## Limitations & Known Issues

| Issue | Description | Workaround |
|-------|-------------|------------|
| `dynamic (user space)` | B站 rates limits user space dynamics API | Use `feed` for followed users' dynamics, or retry with delays |
| `like` | Returns "already liked" (65006) if previously liked | Use `--undo` flag to unlike first |
| `username encoding` | Chinese usernames may show garbled on Windows terminal | The JSON output contains correct Unicode data |

## TODO

- [ ] Add comprehensive test suite with mock API responses
- [ ] Implement CI/CD pipeline
- [ ] Add more platforms (e.g. Douyin, YouTube)
- [ ] Add video download capability
- [ ] Add live streaming monitoring commands
- [ ] Implement concurrent batch operations
- [ ] Write API usage documentation per command

## Build

```bash
cd aiview-cli
go build -o aiview-cli.exe .
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [ffmpeg](https://ffmpeg.org/) — Required for audio WAV splitting (`audio --segment`)
- Go 1.21+