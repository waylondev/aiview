# Fix CLI Bugs & Improve Login Spec

## Why
Three CLI bugs were identified in testing: rank (-352 WBI error), feed (missing --max flag), logout (no JSON output). Additionally, login needs completion and all source code should use English.

## What Changes
- Fix WBI signing implementation so `rank` command works
- Add `--max` flag to `feed` command
- Add JSON format output to `logout` command
- Add `--max` flag support to `feed` command (currently missing)
- Translate all command strings, comments, and user-facing messages to English
- Complete QR login integration with proper cookie credential flow
- **BREAKING**: Command output messages change from Chinese to English

## Impact
- Affected specs: Bilibili command implementations
- Affected code:
  - `internal/platform/bilibili/wbi.go` — WBI signing algorithm
  - `internal/platform/bilibili/client.go` — WBI call for rank
  - `internal/platform/bilibili/commands/discovery.go` — feed command flags
  - `internal/platform/bilibili/commands/account.go` — logout JSON output
  - `internal/platform/bilibili/` all `commands/*.go` files — English translations
  - `internal/platform/bilibili/login.go` — QR login flow
  - `internal/platform/bilibili/auth.go` — credential save
  - `internal/platform/bilibili/bilibili.go` — platform integration

## ADDED Requirements

### Requirement: Fix WBI Signing for Rank
The system SHALL correctly sign requests to `/x/web-interface/ranking/v2` with WBI signature.

#### Scenario: Rank command succeeds
- **WHEN** user runs `aiview bilibili rank --json`
- **THEN** response has `"ok": true` with ranking data

### Requirement: Fix Feed --max Flag
The system SHALL support `--max` / `-n` flag on feed command.

#### Scenario: Feed with limit
- **WHEN** user runs `aiview bilibili feed --max 3 --json`
- **THEN** response has `"ok": true` with at most 3 items

### Requirement: English Source Code
The system SHALL use English for all Go source code including:
- Comments
- User-facing output strings (error messages, help text)
- Variable/function names (already English)

## MODIFIED Requirements

### Requirement: Logout JSON Output
The `logout` command SHALL support `--json` and `--yaml` output formats with proper agent envelope format.

#### Scenario: Logout with JSON
- **WHEN** user runs `aiview bilibili logout --json`
- **THEN** response has `"ok": true`

### Requirement: Complete QR Login
The QR login flow SHALL:
1. Generate QR code URL from Bilibili API
2. Poll for scan status
3. Extract cookies on success
4. Save credential to auth store
5. Set cookie on client for subsequent requests

#### Scenario: QR Login
- **WHEN** user runs `aiview bilibili login`
- **THEN** QR code is generated and displayed as URL
- **WHEN** user scans with Bilibili app
- **THEN** credential is saved and client cookie is set