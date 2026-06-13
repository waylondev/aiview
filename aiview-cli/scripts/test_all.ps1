# scripts/test_all.ps1
# End-to-end automated test script for aiview CLI (Windows PowerShell)
# Usage: .\scripts\test_all.ps1 [-SkipBuild]

param(
    [switch]$SkipBuild
)

$ErrorActionPreference = "Continue"

# ─── Config ──────────────────────────────────────────────────────────────────
$ScriptDir   = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
Set-Location $ProjectRoot

$Binary    = ".\aiview.exe"
$PassCount = 0
$FailCount = 0
$SkipCount = 0
$Results   = @()

# ─── Helpers ─────────────────────────────────────────────────────────────────
function Write-Info  { param($msg) Write-Host "[INFO]  $msg" -ForegroundColor Cyan }
function Write-Ok    { param($msg) Write-Host "[PASS]  $msg" -ForegroundColor Green }
function Write-Fail  { param($msg) Write-Host "[FAIL]  $msg" -ForegroundColor Red }
function Write-Skip  { param($msg) Write-Host "[SKIP]  $msg" -ForegroundColor Yellow }

function Record-Result {
    param($Name, $Status, $Detail = "")
    $script:Results += [PSCustomObject]@{ Name=$Name; Status=$Status; Detail=$Detail }
    switch ($Status) {
        "PASS" { $script:PassCount++ }
        "FAIL" { $script:FailCount++ }
        default { $script:SkipCount++ }
    }
}

function Invoke-Test {
    param($Name, $Command)
    Write-Info "Running: $Command"
    try {
        $output = Invoke-Expression "$Command 2>&1" | Out-String
        $exitCode = $LASTEXITCODE
        if ($null -eq $exitCode) { $exitCode = 0 }
    } catch {
        $output = $_.Exception.Message
        $exitCode = 1
    }

    if ($exitCode -ne 0) {
        Write-Fail "$Name (exit code: $exitCode)"
        $lines = ($output -split "`n") | Select-Object -First 5
        Write-Host "  Output: $($lines -join ' ')"
        Record-Result $Name "FAIL" "exit_code=$exitCode"
        return $false
    }

    Write-Ok $Name
    Record-Result $Name "PASS" ""
    return $true
}

function Invoke-TestExpect {
    param($Name, $Pattern, $Command)
    Write-Info "Running: $Command"
    try {
        $output = Invoke-Expression "$Command 2>&1" | Out-String
        $exitCode = $LASTEXITCODE
        if ($null -eq $exitCode) { $exitCode = 0 }
    } catch {
        $output = $_.Exception.Message
        $exitCode = 1
    }

    if ($exitCode -ne 0) {
        Write-Fail "$Name (exit code: $exitCode)"
        $lines = ($output -split "`n") | Select-Object -First 5
        Write-Host "  Output: $($lines -join ' ')"
        Record-Result $Name "FAIL" "exit_code=$exitCode"
        return $false
    }

    if ($output -match $Pattern) {
        Write-Ok $Name
        Record-Result $Name "PASS" ""
        return $true
    } else {
        Write-Fail "$Name (pattern '$Pattern' not found)"
        $lines = ($output -split "`n") | Select-Object -First 10
        Write-Host "  Output: $($lines -join ' ')"
        Record-Result $Name "FAIL" "pattern_not_found=$Pattern"
        return $false
    }
}

# ─── Build ───────────────────────────────────────────────────────────────────
function Build-Binary {
    if ($SkipBuild) {
        Write-Info "Skipping build (-SkipBuild)"
        if (-not (Test-Path $Binary)) {
            Write-Fail "Binary not found: $Binary"
            exit 1
        }
        return
    }

    Write-Info "Building aiview..."
    $buildOutput = go build -o aiview.exe . 2>&1 | Out-String
    if ($LASTEXITCODE -ne 0) {
        Write-Fail "Build failed"
        Write-Host $buildOutput
        exit 1
    }
    Write-Ok "Build successful"
}

# ─── Test: Hot commands ─────────────────────────────────────────────────────
function Test-HotCommands {
    Write-Host ""
    Write-Host "============================================"
    Write-Host "  Platform Hot Commands"
    Write-Host "============================================"

    Invoke-Test "bilibili hot"          "$Binary bilibili hot"          | Out-Null
    Invoke-Test "bilibili hot --json"   "$Binary bilibili hot --json"   | Out-Null
    Invoke-Test "douyin hot"            "$Binary douyin hot"            | Out-Null
    Invoke-Test "douyin hot --json"     "$Binary douyin hot --json"     | Out-Null
    Invoke-Test "xiaohongshu hot"       "$Binary xiaohongshu hot"       | Out-Null
    Invoke-Test "xiaohongshu hot --json" "$Binary xiaohongshu hot --json" | Out-Null
}

# ─── Test: Global command help ───────────────────────────────────────────────
function Test-GlobalHelp {
    Write-Host ""
    Write-Host "============================================"
    Write-Host "  Global Command Help"
    Write-Host "============================================"

    Invoke-TestExpect "analyze trend --help" "Usage|usage|trend"    "$Binary analyze trend --help" | Out-Null
    Invoke-TestExpect "compare --help"       "Usage|usage|compare"  "$Binary compare --help"       | Out-Null
    Invoke-TestExpect "schedule --help"      "Usage|usage|schedule" "$Binary schedule --help"      | Out-Null
    Invoke-TestExpect "dashboard --help"     "Usage|usage|dashboard" "$Binary dashboard --help"    | Out-Null
    Invoke-TestExpect "tui --help"           "Usage|usage|tui"      "$Binary tui --help"           | Out-Null
    Invoke-TestExpect "export --help"        "Usage|usage|export"   "$Binary export --help"        | Out-Null
    Invoke-TestExpect "root --help"          "Usage|usage|aiview"   "$Binary --help"               | Out-Null
}

# ─── Test: Subcommand help ──────────────────────────────────────────────────
function Test-SubcommandHelp {
    Write-Host ""
    Write-Host "============================================"
    Write-Host "  Subcommand Help"
    Write-Host "============================================"

    Invoke-TestExpect "schedule add --help"    "Usage|usage|add"    "$Binary schedule add --help"    | Out-Null
    Invoke-TestExpect "schedule list --help"   "Usage|usage|list"   "$Binary schedule list --help"   | Out-Null
    Invoke-TestExpect "schedule remove --help" "Usage|usage|remove" "$Binary schedule remove --help" | Out-Null
}

# ─── Test: Miscellaneous ────────────────────────────────────────────────────
function Test-Misc {
    Write-Host ""
    Write-Host "============================================"
    Write-Host "  Miscellaneous"
    Write-Host "============================================"

    Write-Info "Running: $Binary nonexistent_command_xyz"
    try {
        $output = Invoke-Expression "$Binary nonexistent_command_xyz 2>&1" | Out-String
        $exitCode = $LASTEXITCODE
        if ($null -eq $exitCode) { $exitCode = 0 }
    } catch {
        $exitCode = 1
    }

    if ($exitCode -ne 0) {
        Write-Ok "unknown command returns non-zero exit"
        Record-Result "unknown command" "PASS" ""
    } else {
        Write-Fail "unknown command should return non-zero exit"
        Record-Result "unknown command" "FAIL" "exit_code=$exitCode"
    }
}

# ─── Summary ─────────────────────────────────────────────────────────────────
function Show-Summary {
    $total = $PassCount + $FailCount + $SkipCount
    Write-Host ""
    Write-Host "============================================"
    Write-Host "  Test Summary"
    Write-Host "============================================"
    Write-Host "  Total : $total"
    Write-Host "  Passed: $PassCount" -ForegroundColor Green
    Write-Host "  Failed: $FailCount" -ForegroundColor Red
    Write-Host "  Skipped: $SkipCount" -ForegroundColor Yellow
    Write-Host "============================================"

    if ($FailCount -gt 0) {
        Write-Host ""
        Write-Host "Failed tests:"
        foreach ($r in $Results) {
            if ($r.Status -eq "FAIL") {
                Write-Host "  X $($r.Name)  ($($r.Detail))" -ForegroundColor Red
            }
        }
    }

    Write-Host ""
    if ($FailCount -gt 0) {
        Write-Host "TEST SUITE FAILED" -ForegroundColor Red
        exit 1
    } else {
        Write-Host "ALL TESTS PASSED" -ForegroundColor Green
        exit 0
    }
}

# ─── Main ────────────────────────────────────────────────────────────────────
function Main {
    Write-Host "============================================"
    Write-Host "  aiview E2E Test Suite (Windows)"
    Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
    Write-Host "============================================"

    Build-Binary
    Test-HotCommands
    Test-GlobalHelp
    Test-SubcommandHelp
    Test-Misc
    Show-Summary
}

Main
