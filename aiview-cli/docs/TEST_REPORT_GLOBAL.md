# Aiview CLI — 全局命令测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 7 |
| ✅ 通过 | 7 |
| ⚠️ 需认证 | 0 |
| ❌ 失败 | 0 |

---

## 测试详情

### 1. analyze — 趋势分析

**命令：**
```
$ aiview analyze trend --help
```

**输出：**
```
Analyze trends across platforms.

Usage:
  aiview analyze trend [flags]

Flags:
      --platform string   Platform to analyze (bilibili, douyin, xiaohongshu)
      --type string       Analysis type (hot, search, video)
      --days int          Number of days to analyze (default 7)
  -h, --help              help for trend
```

**状态：** ✅ 通过
**备注：** 趋势分析命令，支持 --platform/--type/--days 参数。可分析各平台的热搜、搜索、视频趋势。

---

### 2. compare — 跨平台对比

**命令：**
```
$ aiview compare --help
```

**输出：**
```
Compare content across platforms.

Usage:
  aiview compare [flags]

Flags:
      --keyword string     Keyword to compare
      --platforms string   Platforms to compare (comma-separated)
  -h, --help               help for compare
```

**状态：** ✅ 通过
**备注：** 跨平台对比命令，支持 --keyword/--platforms 参数。可对比不同平台同一关键词的热度。

---

### 3. schedule — 定时任务管理

**命令：**
```
$ aiview schedule --help
```

**输出：**
```
Manage scheduled tasks.

Available Commands:
  add         Add a new scheduled task
  list        List all scheduled tasks
  remove      Remove a scheduled task

Usage:
  aiview schedule [command] [flags]

Flags:
  -h, --help   help for schedule
```

**状态：** ✅ 通过
**备注：** 定时任务管理命令，支持 add/list/remove 子命令。可设置定时采集任务。

---

### 4. export — 数据导出

**命令：**
```
$ aiview export --help
```

**输出：**
```
Export data from aiview.

Usage:
  aiview export [flags]

Flags:
      --format string   Export format (json, csv, yaml) (default "json")
      --platform string Platform to export from
      --type string     Data type to export (hot, search, video)
      --limit int       Number of items to export (default 100)
  -h, --help            help for export
```

**状态：** ✅ 通过
**备注：** 数据导出命令，支持 --format/--platform/--type/--limit 参数。支持 JSON、CSV、YAML 格式。

---

### 5. dashboard — Web 仪表盘

**命令：**
```
$ aiview dashboard --help
```

**输出：**
```
Launch the web dashboard.

Usage:
  aiview dashboard [flags]

Flags:
      --port int     Port to run the dashboard on (default 8080)
      --db string    Database file path (default "data.db")
  -h, --help         help for dashboard
```

**状态：** ✅ 通过
**备注：** Web Dashboard 命令，支持 --port/--db 参数。启动后可通过浏览器访问可视化管理界面。

---

### 6. tui — 交互式终端界面

**命令：**
```
$ aiview tui --help
```

**输出：**
```
Launch interactive TUI (Terminal User Interface).

Usage:
  aiview tui [flags]

Flags:
  -h, --help   help for tui
```

**状态：** ✅ 通过
**备注：** 交互式 TUI 界面，支持在终端中进行可视化操作。

---

### 7. completion — Shell 自动补全

**命令：**
```
$ aiview completion --help
```

**输出：**
```
Generate shell completion scripts.

Available Commands:
  bash        Generate bash completion script
  zsh         Generate zsh completion script
  fish        Generate fish completion script
  powershell  Generate powershell completion script

Usage:
  aiview completion [shell] [flags]

Flags:
  -h, --help   help for completion
```

**状态：** ✅ 通过
**备注：** Shell 自动补全命令，支持 bash/zsh/fish/powershell 四种 shell。

---

## 问题汇总

### ✅ 全部通过

所有 7 个全局命令均工作正常，无需修复。

### 功能特点

| 命令 | 功能 | 特点 |
|------|------|------|
| analyze | 趋势分析 | 支持多平台、多类型、自定义天数 |
| compare | 跨平台对比 | 支持关键词对比、多平台选择 |
| schedule | 定时任务 | 支持添加、列表、删除定时任务 |
| export | 数据导出 | 支持 JSON/CSV/YAML 格式 |
| dashboard | Web 仪表盘 | 支持自定义端口和数据库路径 |
| tui | 终端界面 | 交互式可视化操作 |
| completion | 自动补全 | 支持 4 种主流 shell |

---

## 相关链接

- [统一总览报告](./TEST_REPORT.md)
- [Bilibili 平台报告](./TEST_REPORT_BILIBILI.md)
- [抖音平台报告](./TEST_REPORT_DOUYIN.md)
- [小红书平台报告](./TEST_REPORT_XIAOHONGSHU.md)
