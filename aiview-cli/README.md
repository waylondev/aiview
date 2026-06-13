# Aiview CLI — 多平台内容分析命令行工具

Aiview 是一个基于 Go 的多平台内容管理 CLI 工具，支持 **6 大内容平台**、**70+ 命令**、**4 种输出格式**，提供从数据采集到趋势分析的一站式终端体验。

## 核心特性

- **6 平台覆盖** — Bilibili (42) / 抖音 (11) / 小红书 (5) / 微博 (4) / 快手 (4) / 知乎 (4)
- **统一错误处理** — 标准化 `aiverr` 错误类型，HTML 响应检测，HTTP 状态码分类
- **多格式输出** — JSON / YAML / Table / CSV，全局 `--json` / `--yaml` / `--csv` / `--table` 标志
- **交互式 TUI** — 基于 [Bubbletea](https://github.com/charmbracelet/bubbletea) 的终端界面
- **Web Dashboard** — 内置 HTTP 可视化面板，Chart.js 数据图表
- **数据管道** — SQLite + JSON 双后端，支持采集 → 存储 → 分析 → 导出全流程
- **速率控制** — 令牌桶限流 + TTL 缓存，保护 API 配额
- **Go Library** — `pkg/aiview` 可作为 Go 模块在其他项目中引用

## 快速开始

```bash
# 编译
cd aiview-cli && go build -o aiview.exe .

# 查看帮助
./aiview --help

# 查看平台热搜
./aiview bilibili hot --json
./aiview douyin hot --json
./aiview weibo hot --json
./aiview kuaishou hot --json
./aiview zhihu hot --json

# Cookie 登录（解锁更多功能）
./aiview bilibili login --sessdata "<SESSDATA>" --bili-jct "<BILI_JCT>"
./aiview douyin login --cookie "<cookie>"

# 搜索
./aiview bilibili search "AI" --json
./aiview douyin search "AI" --json

# 数据分析
./aiview analyze trend --days 7 --platform bilibili
./aiview compare --json

# 启动 TUI
./aiview tui

# 启动 Dashboard
./aiview dashboard
```

## 架构

```
aiview-cli/
├── main.go / root.go          # 入口 & 全局 Flag
├── commands/                  # 命令层（每个平台独立子包）
│   ├── bilibili/              #   42 个命令
│   ├── douyin/                #   11 个命令
│   ├── xiaohongshu/           #   5 个命令
│   ├── weibo/                 #   4 个命令
│   ├── kuaishou/              #   4 个命令
│   └── zhihu/                 #   4 个命令
├── internal/
│   ├── platform/              # 平台抽象层 & 注册器
│   │   ├── platform.go        #   Platform 接口定义
│   │   ├── registry.go        #   全局平台注册
│   │   ├── bilibili/          #   Bilibili 实现
│   │   ├── douyin/            #   抖音实现
│   │   ├── xiaohongshu/       #   小红书实现
│   │   ├── weibo/             #   微博实现
│   │   ├── kuaishou/          #   快手实现
│   │   └── zhihu/             #   知乎实现
│   ├── errors/                # 统一错误类型 (aiverr)
│   ├── helper/                # 通用工具函数
│   ├── config/                # 配置加载 (viper)
│   ├── output/                # 输出格式化
│   ├── storage/               # 数据持久化 (SQLite/JSON)
│   ├── cache/                 # TTL 缓存
│   ├── ratelimit/             # 令牌桶限流
│   ├── analyzer/              # 趋势分析
│   ├── pipeline/              # 数据管道
│   ├── scheduler/             # 定时任务
│   ├── auth/                  # 认证管理
│   └── browser/               # 浏览器自动化 (chromedp)
├── tui/                       # 交互式 TUI (bubbletea)
├── dashboard/                 # Web Dashboard (Chart.js)
├── pkg/aiview/                # Go Library API
└── docs/                      # 测试报告
```

### 请求流转

```
User → Cobra Parser → Platform Registry → Platform Commands
  → HTTP Client (带限流+缓存) → API 响应
  → 输出格式化 (JSON/YAML/Table/CSV) → Terminal
```

## 平台命令概览

| 平台 | 命令数 | 公开可用 | 需认证 | 测试报告 |
|------|--------|----------|--------|----------|
| Bilibili | 42 | 14 | 8 | [报告](docs/TEST_REPORT_BILIBILI.md) |
| 抖音 | 11 | 2 | 9 | [报告](docs/TEST_REPORT_DOUYIN.md) |
| 小红书 | 5 | 0 | 5 | [报告](docs/TEST_REPORT_XIAOHONGSHU.md) |
| 微博 | 4 | 2 | 2 | [报告](docs/TEST_REPORT_WEIBO.md) |
| 快手 | 4 | 3 | 1 | [报告](docs/TEST_REPORT_KUAISHOU.md) |
| 知乎 | 4 | 2 | 2 | [报告](docs/TEST_REPORT_ZHIHU.md) |
| 全局 | 6 | 6 | 0 | [报告](docs/TEST_REPORT_GLOBAL.md) |

## 全局命令

| 命令 | 说明 |
|------|------|
| `analyze trend` | 趋势分析（支持 `--days` / `--platform` / `--type`） |
| `compare` | 跨平台数据对比 |
| `export` | 数据导出（支持 `--format` / `--output`） |
| `schedule add/list/remove` | 定时任务管理 |
| `tui` | 启动交互式终端界面 |
| `dashboard` | 启动 Web 可视化面板 |

## 技术栈

| 组件 | 技术选型 |
|------|----------|
| CLI 框架 | [Cobra](https://github.com/spf13/cobra) |
| 配置管理 | [Viper](https://github.com/spf13/viper) |
| TUI 框架 | [Bubbletea](https://github.com/charmbracelet/bubbletea) + [Lipgloss](https://github.com/charmbracelet/lipgloss) |
| 数据可视化 | [Chart.js](https://www.chartjs.org/) |
| 浏览器自动化 | [chromedp](https://github.com/chromedp/chromedp) |
| 数据库 | SQLite (CGO) + JSON 文件 |
| CI/CD | GitHub Actions |

## 构建

```bash
go build -o aiview.exe .
```

**依赖**: Go 1.21+ / [ffmpeg](https://ffmpeg.org/)（仅 `audio --segment` 需要）

## 测试报告

每个平台独立测试报告，覆盖全部命令验证：

- [Bilibili](docs/TEST_REPORT_BILIBILI.md) · [抖音](docs/TEST_REPORT_DOUYIN.md) · [小红书](docs/TEST_REPORT_XIAOHONGSHU.md)
- [微博](docs/TEST_REPORT_WEIBO.md) · [快手](docs/TEST_REPORT_KUAISHOU.md) · [知乎](docs/TEST_REPORT_ZHIHU.md)
- [全局命令](docs/TEST_REPORT_GLOBAL.md)
