# Tasks

## 1. Fix WBI Signing for Rank Command
- [x] 1.1 Debug WBI signature algorithm against rank API endpoint
- [x] 1.2 Fix `wbi.go` signing logic (mixin key, encoding, filter)
- [x] 1.3 Ensure `GetRankVideos` uses `wbiGet` (already done) and verify rank works
- [x] 1.4 Verify `go build ./...` and `go run --json` returns ok:true

## 2. Fix Feed Command --max Flag
- [x] 2.1 Add `--max` / `-n` flag to feed command definition
- [x] 2.2 Wire flag value to client call limit
- [x] 2.3 Verify `aiview bilibili feed --max 3 --json` works

## 3. Fix Logout JSON Output
- [x] 3.1 Add `--json` and `--yaml` output format support to logout command
- [x] 3.2 Use `output.EmitSuccess` for JSON/YAML output
- [x] 3.3 Verify `aiview bilibili logout --json` returns `{"ok": true}`

## 4. Complete QR Login Integration
- [x] 4.1 Implement `GenerateQRCode()` — call `/x/passport-login/web/qrcode/generate`
- [x] 4.2 Implement `PollQRCode()` — call `/x/passport-login/web/qrcode/poll`
- [x] 4.3 Implement cookie extraction from response
- [x] 4.4 Wire QR functions into `commands/account.go` via `SetQRLoginFuncs`
- [x] 4.5 Wire credential save into `auth.go` store
- [x] 4.6 Wire cookie into client via `BuildCookieString()`
- [x] 4.7 Verify `aiview bilibili login --sessdata <value>` works

## 5. Translate All Source Code to English
- [x] 5.1 Translate `internal/platform/bilibili/commands/*.go` — user-facing messages, help strings, errors
- [x] 5.2 Translate `internal/platform/bilibili/client.go` — error messages
- [x] 5.3 Translate `internal/platform/bilibili/auth.go` — error messages, comments
- [x] 5.4 Translate `internal/platform/bilibili/login.go` — comments
- [x] 5.5 Translate `internal/platform/bilibili/wbi.go` — comments
- [x] 5.6 Translate `internal/platform/bilibili/bilibili.go` — comments
- [x] 5.7 Translate `internal/platform/bilibili/commands/types.go` — comments

## 6. Final Verification
- [x] 6.1 Run `go build ./...` and `go vet ./...`
- [x] 6.2 Test all commands with `--json` output
- [x] 6.3 Update checklist.md

# Task Dependencies
- Task 1, 2, 3, 5 are independent — can run in parallel
- Task 4 has sub-dependencies (4.1-4.3 implement, 4.4-4.6 wire up)
- Task 6 depends on all previous tasks