# Checklist

## Fix WBI Signing
- [x] `aiview bilibili rank --json` returns `"ok": true` with ranking data

## Fix Feed --max Flag
- [x] `aiview bilibili feed --max 3 --json` returns proper JSON (not authenticated)
- [x] `aiview bilibili feed --json` returns `"ok": false` with `not_authenticated`

## Fix Logout JSON Output
- [x] `aiview bilibili logout --json` returns `{"ok": true, ...}`
- [x] `aiview bilibili logout --yaml` returns `ok: true`

## QR Login
- [x] `aiview bilibili login` generates QR code URL and polls
- [x] `aiview bilibili login --sessdata <value>` saves credential
- [x] After login, `aiview bilibili status --json` returns `"authenticated": true`

## English Source Code
- [x] No Chinese user-facing strings in command `.go` files
- [x] No Chinese comments in Go source files
- [x] `go build ./...` passes
- [x] `go vet ./...` passes