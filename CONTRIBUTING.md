# 贡献指南

感谢你对 Aiview 项目的关注！我们欢迎各种形式的贡献，包括代码改进、文档完善、Bug 报告和新功能建议。

---

## 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发环境设置](#开发环境设置)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [Pull Request 流程](#pull-request-流程)
- [测试要求](#测试要求)
- [文档贡献](#文档贡献)

---

## 行为准则

参与本项目即表示你同意遵守以下准则：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

---

## 如何贡献

### 报告 Bug

使用 GitHub Issues 报告 Bug，并包含以下信息：

1. **清晰的标题** — 简明扼要地描述问题
2. **复现步骤** — 详细说明如何复现问题
3. **期望行为** — 你期望发生什么
4. **实际行为** — 实际发生了什么
5. **环境信息** — 操作系统、Go 版本、Aiview 版本
6. **日志输出** — 使用 `-v` 参数获取详细日志

**Bug 报告模板**：

```markdown
**描述**
简要描述 Bug

**复现步骤**
1. 运行命令 '...'
2. 看到错误 '...'

**期望行为**
描述期望的行为

**实际行为**
描述实际发生的行为

**环境**
- OS: [e.g. Windows 10]
- Go: [e.g. 1.21]
- Aiview: [e.g. v1.0.0]

**日志**
```bash
aiview bilibili hot -v
```
```

### 建议新功能

在提出新功能建议前，请先：

1. 检查 Issues 中是否已有类似建议
2. 考虑功能的可行性和维护成本
3. 思考功能与项目目标的一致性

**功能建议模板**：

```markdown
**功能描述**
简要描述你想要的功能

**使用场景**
描述这个功能的使用场景

**建议方案**
如果有具体的实现思路，请描述

**替代方案**
是否考虑过其他解决方案

**额外信息**
其他相关截图、链接等
```

---

## 开发环境设置

### 1. Fork 项目

访问 [Aiview GitHub 仓库](https://github.com/jackwener/aiview)，点击右上角的 "Fork" 按钮。

### 2. 克隆仓库

```bash
git clone https://github.com/YOUR_USERNAME/aiview.git
cd aiview/aiview-cli
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 构建项目

```bash
go build -o aiview-cli.exe .
```

### 5. 运行测试

```bash
go test ./...
```

### 6. 创建分支

```bash
git checkout -b feature/your-feature-name
# 或
git checkout -b fix/issue-description
```

---

## 代码规范

### Go 代码风格

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- 使用 `gofmt` 或 `goimports` 格式化代码
- 使用 `go vet` 检查常见问题

### 命名规范

- **包名** — 简短、小写、不使用下划线
- **变量/函数名** — 使用驼峰命名法
- **导出函数** — 首字母大写，必须添加注释
- **常量** — 使用 `const` 声明，不使用全大写

**示例**：

```go
// Good
package aiview

// GetHotVideos fetches popular videos from platform
func GetHotVideos(page, count int) ([]VideoInfo, error) {
    // ...
}

const defaultTimeout = 30

// Bad
package aiview_pkg

func get_hot_videos(Page, Count int) ([]VideoInfo, error) {
    // ...
}

const DEFAULT_TIMEOUT = 30
```

### 错误处理

- 始终检查错误返回值
- 使用有意义的错误消息
- 不要忽略错误（使用 `_` 丢弃）

**示例**：

```go
// Good
videos, err := client.GetHotVideos(1, 10)
if err != nil {
    return fmt.Errorf("failed to fetch hot videos: %w", err)
}

// Bad
videos, _ := client.GetHotVideos(1, 10)
```

### 注释规范

- 导出函数必须有注释
- 注释应说明"为什么"而非"做什么"
- 使用完整的句子，以函数名开头

**示例**：

```go
// GetHotVideos fetches popular videos from Bilibili.
// It returns a list of VideoInfo sorted by popularity.
func GetHotVideos(page, count int) ([]VideoInfo, error) {
    // ...
}
```

### 平台实现规范

添加新平台时，需实现以下接口：

```go
type Platform interface {
    Name() string
    Commands() []*cobra.Command
    NewClient(cfg *config.Config) (Client, error)
}
```

并在 `internal/platform/` 下创建对应目录，包含：

- `platform.go` — Platform 接口实现
- `client.go` — HTTP Client 实现
- `auth.go` — 认证逻辑
- `types.go` — 平台特定类型定义

---

## 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范提交代码。

### 提交消息格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- **feat** — 新功能
- **fix** — Bug 修复
- **docs** — 文档更新
- **style** — 代码格式（不影响代码运行的变动）
- **refactor** — 重构（既不是新增功能，也不是修改 Bug 的代码变动）
- **perf** — 性能优化
- **test** — 增加测试
- **chore** — 构建过程或辅助工具的变动

### Scope 范围

- `bilibili` — Bilibili 相关
- `douyin` — 抖音相关
- `xiaohongshu` — 小红书相关
- `core` — 核心功能
- `docs` — 文档
- `ci` — CI/CD

### 示例

```bash
# 新功能
feat(bilibili): add live room info command

# Bug 修复
fix(douyin): fix cookie expiration handling

# 文档
docs: update API documentation

# 重构
refactor(core): simplify error handling

# 测试
test(bilibili): add unit tests for video search
```

---

## Pull Request 流程

### 1. 保持分支同步

```bash
git remote add upstream https://github.com/jackwener/aiview.git
git fetch upstream
git checkout main
git merge upstream/main
```

### 2. 提交更改

```bash
git add .
git commit -m "feat(bilibili): add new feature"
```

### 3. 推送到 Fork

```bash
git push origin feature/your-feature-name
```

### 4. 创建 Pull Request

访问你的 Fork 仓库，点击 "Compare & pull request" 按钮。

**PR 描述模板**：

```markdown
## 描述
简要描述此 PR 的目的

## 相关 Issue
Closes #123

## 改动类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 文档更新
- [ ] 重构
- [ ] 其他

## 测试
描述你进行的测试

## 截图（如适用）
添加截图展示改动

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 已添加必要的注释
- [ ] 已添加/更新测试
- [ ] 所有测试通过
- [ ] 文档已更新（如需要）
```

### 5. 代码审查

- 维护者会审查你的 PR
- 根据反馈进行修改
- 审查通过后，PR 会被合并

---

## 测试要求

### 单元测试

为新功能添加单元测试，确保代码覆盖率：

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/platform/bilibili/...

# 查看覆盖率
go test -cover ./...
```

### 测试命名

使用描述性的测试名称：

```go
func TestGetHotVideos_Success(t *testing.T) {
    // ...
}

func TestGetHotVideos_InvalidPage(t *testing.T) {
    // ...
}
```

### 测试结构

使用表驱动测试：

```go
func TestSearchVideo(t *testing.T) {
    tests := []struct {
        name     string
        keyword  string
        page     int
        wantErr  bool
    }{
        {
            name:    "valid search",
            keyword: "Golang",
            page:    1,
            wantErr: false,
        },
        {
            name:    "empty keyword",
            keyword: "",
            page:    1,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

---

## 文档贡献

### 更新文档

- 修改代码后，同步更新相关文档
- 确保代码示例可以正常运行
- 使用清晰的中文描述

### 文档结构

- `README.md` — 项目总览和快速开始
- `docs/API.md` — Go Library API 文档
- `docs/PLATFORM_GUIDE.md` — 平台接入指南
- `CONTRIBUTING.md` — 贡献指南（本文档）

### 添加示例

在 `examples/` 目录下添加示例代码：

```
examples/
├── basic/
│   └── main.go       # 基础使用示例
├── advanced/
│   └── main.go       # 高级使用示例
└── batch/
    └── main.go       # 批量操作示例
```

---

## 开发建议

### 从小处着手

- 第一次贡献可以从修复文档错误、改进注释开始
- 逐步参与更复杂的功能开发

### 寻求帮助

- 遇到问题可以在 Issue 中提问
- 查看现有代码和文档
- 参考其他平台的实现

### 保持沟通

- 在开发前先在 Issue 中讨论你的想法
- 定期更新 PR 进度
- 及时回应审查意见

---

## 许可证

贡献的代码将遵循项目的 MIT 许可证。

---

## 致谢

感谢所有贡献者！你的努力让 Aiview 变得更好。
