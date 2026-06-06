# Go Aiview CLI 多平台工具 Spec

## Why
参考 Python 版 bilibili-cli、rdt-cli、twitter-cli 的成熟架构，用 Go 语言构建一个高内聚、低耦合、高扩展的多平台 CLI 工具。先实现 Bilibili 支持，架构设计支持后续轻松扩展其他平台。

## What Changes
- 新建 Go 项目 `aiview`，使用 cobra + viper 技术栈
- 实现 Platform 接口抽象层，支持多平台插件化注册
- 实现 Bilibili 平台完整功能（视频、搜索、用户、收藏、热门、音频等）
- 实现 agent-friendly 结构化输出（JSON/YAML 统一信封）
- 实现认证系统（文件持久化、浏览器 Cookie 提取、QR 登录）
- **BREAKING**: 全新项目，无 breaking changes

## Impact
- Affected specs: 无（全新项目）
- Affected code: 全新代码，不影响现有文件

---

## ADDED Requirements

### Requirement: 项目结构遵循 Go 标准布局
系统 SHALL 采用 `cmd/` + `internal/` + `pkg/` 的 Go 标准项目布局。

#### Scenario: 项目目录结构
- **WHEN** 项目初始化完成
- **THEN** 目录结构包含 `cmd/aiview/main.go`、`internal/`（cli、platform、config、auth、output、cache）、`pkg/httputil/`

### Requirement: Platform 接口抽象
系统 SHALL 定义 `Platform` 接口，包含 `Name() string`、`Commands() []*cobra.Command`、`NewClient(cfg *config.Config) (Client, error)` 三个方法，支持通过 `Register()` 函数注册新平台。

#### Scenario: 注册新平台
- **WHEN** 调用 `platform.Register(p)` 注册一个平台实现
- **THEN** 该平台可从 `platform.GetPlatform(name)` 获取
- **AND** 该平台的命令自动注册到根命令的子命令中

### Requirement: Bilibili 平台实现
系统 SHALL 实现完整的 Bilibili 平台支持，包括视频详情、字幕、AI 总结、评论、搜索（用户/视频）、用户信息、用户视频列表、热门、排行榜、收藏夹、关注列表、观看历史、稍后再看、动态时间线、点赞/投币/三连、取消关注、音频下载等功能。

#### Scenario: 查看视频详情
- **WHEN** 执行 `aiview bilibili video BV1xxx`
- **THEN** 输出视频标题、UP主、时长、播放量、弹幕、点赞、投币、收藏、分享、简介、链接

#### Scenario: 查看视频带字幕
- **WHEN** 执行 `aiview bilibili video BV1xxx --subtitle`
- **THEN** 输出包含字幕内容（如可用）

#### Scenario: 搜索视频
- **WHEN** 执行 `aiview bilibili search "关键词" --type video --max 5`
- **THEN** 输出前 5 条搜索结果

#### Scenario: 搜索用户
- **WHEN** 执行 `aiview bilibili search "关键词" --type user`
- **THEN** 输出用户搜索结果

#### Scenario: 查看用户信息
- **WHEN** 执行 `aiview bilibili user 946974`
- **THEN** 输出用户昵称、等级、签名、粉丝数等

#### Scenario: 查看热门视频
- **WHEN** 执行 `aiview bilibili hot`
- **THEN** 输出热门视频列表

#### Scenario: 结构化输出
- **WHEN** 执行 `aiview bilibili video BV1xxx --json` 或 `--yaml`
- **THEN** 输出符合 agent 信封格式的结构化数据：`{ok: true, schema_version: "1", data: {...}}`

### Requirement: 认证系统
系统 SHALL 支持三种认证方式（按优先级）：从文件加载已保存凭证、从浏览器提取 Cookie、QR 码登录。凭证保存到 `~/.aiview/credential.json`。

#### Scenario: 无凭证时执行需要登录的命令
- **WHEN** 未登录状态下执行 `aiview bilibili favorites`
- **THEN** 提示"需要登录，使用 `aiview bilibili login` 登录"

#### Scenario: 检查登录状态
- **WHEN** 执行 `aiview bilibili status`
- **THEN** 输出认证状态（已登录/未登录）

#### Scenario: QR 码登录
- **WHEN** 执行 `aiview bilibili login`
- **THEN** 终端显示 QR 码，用户扫码后完成登录，凭证保存到本地

### Requirement: 配置管理
系统 SHALL 使用 viper 管理配置，支持 `config.yaml` 配置文件和环境变量覆盖。配置包含默认平台、输出格式、缓存 TTL、平台特定配置等。

#### Scenario: 加载配置文件
- **WHEN** 启动时 `~/.aiview/config.yaml` 存在
- **THEN** 自动加载配置，可通过环境变量 `AIVIEW_*` 覆盖

### Requirement: 输出格式化
系统 SHALL 支持三种输出格式：终端美化的表格输出（默认）、JSON 结构化输出、YAML 结构化输出。当 stdout 非 TTY 时默认使用 YAML。

#### Scenario: TTY 终端输出
- **WHEN** stdout 是 TTY 且未指定 --json/--yaml
- **THEN** 使用美化的表格格式输出

#### Scenario: 管道输出
- **WHEN** stdout 不是 TTY（如管道到文件）
- **THEN** 默认使用 YAML 格式输出

### Requirement: 错误处理
系统 SHALL 使用分层错误类型，通过 `fmt.Errorf` 的 `%w` 包装。错误码映射到稳定的 code 字符串（not_authenticated、invalid_input、network_error、not_found 等）。

#### Scenario: 网络错误
- **WHEN** API 请求因网络问题失败
- **THEN** 输出 `❌ 获取视频信息: network error` 并以 code `network_error` 退出