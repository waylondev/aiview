# Aiview - 多平台内容数据采集与分析 CLI 工具

基于 Go 语言开发的命令行工具，支持从 Bilibili 和抖音平台采集视频、用户、互动、热门等全场景数据，并调用平台 AI 接口进行视频摘要分析。

## 核心能力

| 能力 | 说明 |
|------|------|
| 数据采集 | 视频详情/统计、弹幕、评论、用户信息、热门趋势、搜索等 40+ 命令 |
| 数据导出 | 弹幕导出 XML、音频下载分段、JSON/YAML 结构化输出 |
| AI 分析 | 调用平台 AI 接口生成视频摘要、提取字幕文本 |
| 社交互动 | 点赞、投币、三连、收藏、发评论/弹幕、发动态 |

## 快速开始

```bash
# 构建
cd aiview-cli && go build -o aiview-cli.exe .

# 登录
aiview bilibili login --sessdata "<SESSDATA>" --bili-jct "<BILI_JCT>"
aiview douyin login --cookie "<COOKIE>"

# 采集数据
aiview bilibili video BV1GJ411x7Rq --ai --subtitle    # 视频详情 + AI 摘要 + 字幕
aiview bilibili hot --json                             # 热门视频
aiview bilibili danmaku BV1GJ411x7Rq -o ./data/        # 弹幕导出为 XML
aiview bilibili search "AI技术" --order click           # 搜索并按播放量排序
```

**全局输出格式**（所有命令通用）：
- `--json` — JSON 格式，适合脚本处理和数据管道对接
- `--yaml` — YAML 格式，层次清晰便于阅读
- 默认 — Table 格式，终端友好，数字自动格式化（如 `1234567` → `123.5万`）

---

## Bilibili 命令

### 视频与数据

```bash
aiview bilibili video <BV|URL>              # 视频详情（标题/UP主/时长/播放量等）
aiview bilibili video <BV> --ai             # + AI 生成的内容摘要
aiview bilibili video <BV> --subtitle       # + 视频字幕文本
aiview bilibili video <BV> --comments       # + 热门评论
aiview bilibili video <BV> --related        # + 相关视频推荐
aiview bilibili video-status <BV>           # 六维统计：播放/弹幕/点赞/投币/收藏/分享
aiview bilibili tags <BV>                   # 视频标签
aiview bilibili online <BV>                 # 实时在线人数
```

### 弹幕与音频

```bash
aiview bilibili danmaku <BV>                # 查看弹幕
aiview bilibili danmaku <BV> -o ./data/     # 导出弹幕为 XML 文件
aiview bilibili danmaku-send <BV> "消息"     # 发送弹幕（需登录）
aiview bilibili audio <BV>                  # 下载音频，分段切割为 WAV（需 ffmpeg）
aiview bilibili audio <BV> --no-split       # 下载完整音频（M4A）
aiview bilibili audio <BV> --segment 30     # 自定义分段时长（秒）
```

### 评论

```bash
aiview bilibili comment <BV>                # 查看评论（按时间）
aiview bilibili comment <BV> --sort 2       # 按热度排序
aiview bilibili comment <BV> --page 2       # 分页
aiview bilibili comment <BV> -m "内容"       # 发布评论（需登录）
aiview bilibili comment-delete <BV> <RPID>  # 删除评论（需登录）
```

### 用户

```bash
aiview bilibili user <UID>                  # 用户信息
aiview bilibili user-videos <UID>           # 用户作品列表
aiview bilibili user-videos <UID> --order click --max 50  # 按播放量排序，取 50 条
aiview bilibili fans <UID>                  # 粉丝列表（支持 --page 分页）
aiview bilibili following <UID>             # 关注列表
aiview bilibili relation <UID>              # 关系统计
aiview bilibili collection <UID>            # 视频合集/系列
aiview bilibili whoami                      # 当前登录用户信息
```

### 热门与发现

```bash
aiview bilibili hot                         # 热门视频（--max N 控制数量）
aiview bilibili rank                        # 排行榜（--day 1/3/7/30, --type all/origin/rookie）
aiview bilibili trending                    # 热搜词（--limit N）
aiview bilibili recommend                   # 首页推荐（--fresh 获取全新推荐）
aiview bilibili region <rid>                # 分区视频
aiview bilibili weekly <number>             # 每周必看
aiview bilibili precious                    # 入站必刷
aiview bilibili feed                        # 动态 Feed（需登录）
```

### 搜索

```bash
aiview bilibili search <keyword>            # 搜索视频（默认）
aiview bilibili search <keyword> --type user            # 搜索用户
aiview bilibili search <keyword> --order click/pubdate/dm/score  # 排序方式
aiview bilibili search <keyword> --duration 1             # 时长筛选: 1=<5min 2=5-30min 3=>30min
aiview bilibili suggest <keyword>           # 搜索联想词
```

### 收藏与历史

```bash
aiview bilibili favorites <UID>             # 收藏夹列表
aiview bilibili favorite <BV> --fid <ID>    # 添加到收藏夹
aiview bilibili history                     # 观看历史
aiview bilibili watch-later                 # 稍后再看
```

### 动态

```bash
aiview bilibili dynamic <UID>               # 用户动态
aiview bilibili dynamic-post "文本"          # 发布动态（需登录）
aiview bilibili dynamic-delete <id>         # 删除动态（需登录）
```

### 直播

```bash
aiview bilibili live --room <ID>            # 直播间信息
aiview bilibili live --uid <UID>            # 通过用户 ID 查询
```

### 互动操作（需登录 + 写权限）

```bash
aiview bilibili like <BV>                   # 点赞（--undo 取消）
aiview bilibili coin <BV> -n 2              # 投币（1-2 枚）
aiview bilibili triple <BV>                 # 三连（点赞+投币+收藏）
aiview bilibili unfollow <UID>              # 取消关注
```

### 账号管理

```bash
aiview bilibili login                       # 扫码登录
aiview bilibili login --sessdata "..." --bili-jct "..."  # Cookie 登录
aiview bilibili status                      # 登录状态
aiview bilibili logout                      # 登出
```

---

## 抖音命令

```bash
aiview douyin hot                           # 热榜（无需登录）
aiview douyin trending                      # 趋势话题（无需登录）
aiview douyin search <keyword>              # 搜索（无需登录，有 Cookie 结果更全）
aiview douyin video <id|url>                # 视频详情（需登录）
aiview douyin user <uid>                    # 用户信息（需登录）
aiview douyin user-posts <uid>              # 用户作品（需登录，--cursor 分页）
aiview douyin comment <video_id>            # 评论（需登录，--cursor 分页）
aiview douyin login --cookie "<COOKIE>"     # 登录
aiview douyin status                        # 登录状态
aiview douyin logout                        # 登出
```

---

## 技术架构

```
aiview-cli/
├── main.go / root.go           # 入口 + 全局 Flag
├── commands/
│   ├── bilibili/               # Bilibili 命令层（40+ 命令）
│   └── douyin/                 # 抖音命令层
├── internal/
│   ├── platform/               # 平台抽象层（Platform 接口 + 全局注册器）
│   │   ├── bilibili/           # Bilibili HTTP Client + WBI 签名
│   │   └── douyin/             # 抖音 HTTP Client + Cookie 管理
│   ├── auth/                   # 凭证持久化（读写权限分级）
│   ├── output/                 # 输出格式化（JSON/YAML/Table）
│   ├── config/                 # 配置加载（Viper）
│   └── helper/                 # 嵌套 Map 安全取值
└── docs/                       # 测试报告
```

**设计要点**：
- **平台注册器**：新平台实现 `Platform` 接口并 import 即可自动注册
- **Client 接口抽象**：命令层通过接口调用，与具体 API 解耦
- **Payload 模式**：`video` 命令可组合字幕/AI/评论/推荐，非致命错误以 Warning 输出

## 依赖

- Go 1.21+
- [Cobra](https://github.com/spf13/cobra) — CLI 框架
- [Viper](https://github.com/spf13/viper) — 配置管理
- FFmpeg — 音频分段切割（可选，未安装时降级为完整下载）

## 许可证

MIT License
