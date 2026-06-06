# Tasks

## 1. 项目骨架初始化
- [x] 1.1 初始化 Go 模块 (`go mod init github.com/jackwener/aiview`)
- [x] 1.2 安装依赖 (cobra, viper)
- [x] 1.3 创建目录结构 (`cmd/aiview/`, `internal/`, `pkg/`)
- [x] 1.4 创建 `cmd/aiview/main.go` 入口文件

## 2. 基础设施层
- [x] 2.1 实现 `internal/config/config.go` — 配置加载（viper，支持 config.yaml 和环境变量）
- [x] 2.2 实现 `internal/auth/store.go` — 通用凭据存储（JSON 文件读写、TTL 检查）
- [x] 2.3 实现 `internal/output/formatter.go` — 输出格式化（JSON/YAML 信封 + 表格输出）
- [x] 2.4 实现 `internal/cache/cache.go` — 简单内存缓存（TTL 支持）

## 3. 平台抽象层
- [x] 3.1 实现 `internal/platform/platform.go` — Platform 接口定义
- [x] 3.2 实现 `internal/platform/registry.go` — 平台注册表
- [x] 3.3 实现 `internal/cli/root.go` — cobra 根命令（集成平台注册、全局 flags）

## 4. Bilibili 平台实现
- [x] 4.1 实现 `internal/platform/bilibili/models.go` — 数据模型（Video, User, Comment, Dynamic 等）
- [x] 4.2 实现 `internal/platform/bilibili/client.go` — API 客户端（HTTP 请求封装、错误映射）
- [x] 4.3 实现 `internal/platform/bilibili/auth.go` — 认证（QR 登录、Cookie 提取、凭证持久化）
- [x] 4.4 实现 `internal/platform/bilibili/commands/video.go` — 视频命令（详情、字幕、AI 总结、评论、相关推荐）
- [x] 4.5 实现 `internal/platform/bilibili/commands/search.go` — 搜索命令（视频搜索、用户搜索）
- [x] 4.6 实现 `internal/platform/bilibili/commands/user.go` — 用户命令（信息、视频列表）
- [x] 4.7 实现 `internal/platform/bilibili/commands/account.go` — 账户命令（login, logout, status, whoami）
- [x] 4.8 实现 `internal/platform/bilibili/commands/collections.go` — 收藏命令（收藏夹、关注、历史、稍后再看）
- [x] 4.9 实现 `internal/platform/bilibili/commands/discovery.go` — 发现命令（热门、排行榜、动态）
- [x] 4.10 实现 `internal/platform/bilibili/commands/interactions.go` — 互动命令（点赞、投币、三连、取关）
- [x] 4.11 实现 `internal/platform/bilibili/commands/audio.go` — 音频命令（下载、切分）
- [x] 4.12 实现 `internal/platform/bilibili/bilibili.go` — 平台入口（实现 Platform 接口，注册所有命令）

## 5. 集成与验证
- [x] 5.1 在 `main.go` 中注册 bilibili 平台
- [x] 5.2 编译验证 (`go build ./...`)
- [x] 5.3 验证命令执行 (`aiview bilibili --help`)

# Task Dependencies
- Task 2 (基础设施) 可以并行执行
- Task 3 (平台抽象) 依赖 Task 2
- Task 4 (Bilibili) 依赖 Task 3
- Task 4.2-4.11 内部可部分并行（4.2, 4.3, 4.4 先完成，其余并行）
- Task 5 (集成) 依赖 Task 4