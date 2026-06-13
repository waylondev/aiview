# Aiview - 多平台内容数据采集与分析工具

基于 Go 语言开发的多平台内容数据采集与分析工具，支持 **Bilibili**、**抖音**、**小红书** 三大平台，覆盖视频、用户、互动、热门等 40+ 数据场景，并支持 AI 摘要分析。

## 核心能力

| 能力 | 说明 |
|------|------|
| 数据采集 | 视频详情/统计、弹幕、评论、用户信息、热门趋势、搜索等 40+ 命令 |
| 数据导出 | 弹幕导出 XML、音频下载分段、JSON/YAML 结构化输出 |
| AI 分析 | 调用平台 AI 接口生成视频摘要、提取字幕文本 |
| 社交互动 | 点赞、投币、三连、收藏、发评论/弹幕、发动态 |
| Go Library | 提供 Go 包 API，可在项目中直接引用 |

## 支持平台

| 平台 | 标识符 | 认证方式 | 状态 |
|------|--------|----------|------|
| Bilibili | `bilibili` | SESSDATA + BILI_JCT / 扫码登录 | 完整支持 |
| 抖音 | `douyin` | Cookie | 完整支持 |
| 小红书 | `xiaohongshu` | Cookie | 完整支持 |

---

## 5 分钟快速上手

### 1. 安装

```bash
# 从源码构建
cd aiview-cli && go build -o aiview-cli.exe .

# 或直接运行（开发模式）
cd aiview-cli && go run . <command>
```

### 2. 登录平台

```bash
# Bilibili - Cookie 登录
aiview bilibili login --sessdata "<SESSDATA>" --bili-jct "<BILI_JCT>"

# Bilibili - 扫码登录
aiview bilibili login

# 抖音 - Cookie 登录
aiview douyin login --cookie "<COOKIE>"

# 小红书 - Cookie 登录
aiview xiaohongshu login --cookie "<COOKIE>"
```

### 3. 开始使用

```bash
# 查看热门视频
aiview bilibili hot

# 获取视频详情 + AI 摘要
aiview bilibili video BV1GJ411x7Rq --ai --subtitle

# 抖音热搜
aiview douyin hot

# 小红书热门
aiview xiaohongshu hot
```

---

## 详细使用示例

### Bilibili 平台

#### 视频与数据

```bash
# 视频详情
aiview bilibili video BV1GJ411x7Rq

# 视频详情 + AI 摘要 + 字幕
aiview bilibili video BV1GJ411x7Rq --ai --subtitle

# 视频详情 + 热门评论 + 相关推荐
aiview bilibili video BV1GJ411x7Rq --comments --related

# 视频六维统计（播放/弹幕/点赞/投币/收藏/分享）
aiview bilibili video-status BV1GJ411x7Rq

# 视频标签
aiview bilibili tags BV1GJ411x7Rq

# 实时在线人数
aiview bilibili online BV1GJ411x7Rq
```

#### 弹幕与音频

```bash
# 查看弹幕
aiview bilibili danmaku BV1GJ411x7Rq

# 导出弹幕为 XML 文件
aiview bilibili danmaku BV1GJ411x7Rq -o ./data/

# 发送弹幕（需登录）
aiview bilibili danmaku-send BV1GJ411x7Rq "这是一条弹幕"

# 下载音频并分段切割为 WAV（需 ffmpeg）
aiview bilibili audio BV1GJ411x7Rq

# 下载完整音频（不切割）
aiview bilibili audio BV1GJ411x7Rq --no-split

# 自定义分段时长（秒）
aiview bilibili audio BV1GJ411x7Rq --segment 30
```

#### 评论

```bash
# 查看评论（按时间）
aiview bilibili comment BV1GJ411x7Rq

# 按热度排序
aiview bilibili comment BV1GJ411x7Rq --sort 2

# 分页查看
aiview bilibili comment BV1GJ411x7Rq --page 2

# 发布评论（需登录）
aiview bilibili comment BV1GJ411x7Rq -m "这是一条评论"

# 删除评论（需登录）
aiview bilibili comment-delete BV1GJ411x7Rq <RPID>
```

#### 用户

```bash
# 用户信息
aiview bilibili user 37737161

# 用户作品列表
aiview bilibili user-videos 37737161

# 按播放量排序，取 50 条
aiview bilibili user-videos 37737161 --order click --max 50

# 粉丝列表
aiview bilibili fans 37737161

# 关注列表
aiview bilibili following 37737161

# 关系统计
aiview bilibili relation 37737161

# 视频合集
aiview bilibili collection 37737161

# 当前登录用户信息
aiview bilibili whoami
```

#### 热门与发现

```bash
# 热门视频
aiview bilibili hot

# 控制数量
aiview bilibili hot --max 20

# 排行榜
aiview bilibili rank --day 3 --type all

# 热搜词
aiview bilibili trending --limit 20

# 首页推荐
aiview bilibili recommend

# 获取全新推荐
aiview bilibili recommend --fresh

# 分区视频
aiview bilibili region 1

# 每周必看
aiview bilibili weekly 200

# 入站必刷
aiview bilibili precious

# 动态 Feed（需登录）
aiview bilibili feed
```

#### 搜索

```bash
# 搜索视频
aiview bilibili search "AI技术"

# 搜索用户
aiview bilibili search "AI技术" --type user

# 按播放量排序
aiview bilibili search "AI技术" --order click

# 时长筛选：1=<5min 2=5-30min 3=>30min
aiview bilibili search "AI技术" --duration 1

# 搜索联想词
aiview bilibili suggest "AI技术"
```

#### 收藏与历史

```bash
# 收藏夹列表
aiview bilibili favorites <UID>

# 添加到收藏夹
aiview bilibili favorite BV1GJ411x7Rq --fid <ID>

# 观看历史
aiview bilibili history

# 稍后再看
aiview bilibili watch-later
```

#### 动态

```bash
# 用户动态
aiview bilibili dynamic <UID>

# 发布动态（需登录）
aiview bilibili dynamic-post "今天天气真好"

# 删除动态（需登录）
aiview bilibili dynamic-delete <id>
```

#### 直播

```bash
# 直播间信息（通过房间号）
aiview bilibili live --room 21484828

# 直播间信息（通过用户 ID）
aiview bilibili live --uid 37737161
```

#### 互动操作（需登录 + 写权限）

```bash
# 点赞
aiview bilibili like BV1GJ411x7Rq

# 取消点赞
aiview bilibili like BV1GJ411x7Rq --undo

# 投币（1-2 枚）
aiview bilibili coin BV1GJ411x7Rq -n 2

# 三连（点赞+投币+收藏）
aiview bilibili triple BV1GJ411x7Rq

# 取消关注
aiview bilibili unfollow <UID>
```

### 抖音平台

```bash
# 热榜（无需登录）
aiview douyin hot

# 趋势话题（无需登录）
aiview douyin trending

# 搜索（无需登录，有 Cookie 结果更全）
aiview douyin search "AI技术"

# 视频详情（需登录）
aiview douyin video <id|url>

# 用户信息（需登录）
aiview douyin user <uid>

# 用户作品（需登录，--cursor 分页）
aiview douyin user-posts <uid>

# 评论（需登录，--cursor 分页）
aiview douyin comment <video_id>

# 登录状态
aiview douyin status

# 登出
aiview douyin logout
```

### 小红书平台

```bash
# 热门笔记
aiview xiaohongshu hot

# 搜索笔记
aiview xiaohongshu search "旅行攻略"

# 笔记详情
aiview xiaohongshu note <note_id>

# 用户信息
aiview xiaohongshu user <user_id>

# 登录
aiview xiaohongshu login --cookie "<COOKIE>"

# 登录状态
aiview xiaohongshu status

# 登出
aiview xiaohongshu logout
```

---

## 全局输出格式

所有命令均支持以下输出格式（通用 Flag）：

| Flag | 说明 | 适用场景 |
|------|------|----------|
| `--json` | JSON 格式输出 | 脚本处理、数据管道对接 |
| `--yaml` | YAML 格式输出 | 层次清晰，便于阅读 |
| 默认 | Table 格式 | 终端友好，数字自动格式化（如 `1234567` → `123.5万`） |

```bash
# JSON 输出 - 适合 jq 处理
aiview bilibili hot --json | jq '.[0].title'

# YAML 输出
aiview bilibili video BV1GJ411x7Rq --yaml
```

---

## 配置说明

### 配置文件

Aiview 支持 YAML 配置文件，存放路径为 `~/.aiview/config.yaml` 或当前目录下的 `config.yaml`。

```yaml
# ~/.aiview/config.yaml
platform: bilibili          # 默认平台
cache_ttl: 300              # 缓存过期时间（秒）
output: auto                # 输出格式：auto / json / yaml

platforms:
  bilibili:
    timeout: 30             # HTTP 超时（秒）
  douyin:
    timeout: 30
```

### 环境变量

所有配置项均可通过环境变量覆盖，前缀为 `AIVIEW_`：

```bash
export AIVIEW_PLATFORM=bilibili
export AIVIEW_CACHE_TTL=600
export AIVIEW_PLATFORMS_BILIBILI_TIMEOUT=60
```

### 凭证存储

登录凭证自动保存在 `~/.aiview/` 目录下：

| 文件 | 说明 |
|------|------|
| `bilibili_credential.json` | Bilibili 登录凭证 |
| `douyin_credential.json` | 抖音登录凭证 |
| `xiaohongshu_credential.json` | 小红书登录凭证 |

---

## Go Library 使用

Aiview 同时提供 Go 包 API，可在 Go 项目中直接引用：

```go
package main

import (
    "fmt"
    "log"

    "github.com/jackwener/aiview/pkg/aiview"
)

func main() {
    // 创建 Bilibili 客户端
    client, err := aiview.New("bilibili")
    if err != nil {
        log.Fatal(err)
    }

    biliClient, err := client.BilibiliClient()
    if err != nil {
        log.Fatal(err)
    }

    // 获取热门视频
    videos, err := biliClient.GetHotVideos(1, 10)
    if err != nil {
        log.Fatal(err)
    }

    for _, v := range videos {
        fmt.Printf("%s - %s (播放: %d)\n", v.Title, v.Author, v.Play)
    }
}
```

详见 [API 文档](docs/API.md) 和 [示例代码](examples/)。

---

## 常见问题 FAQ

### Q: 如何获取 Bilibili SESSDATA？

1. 在浏览器中登录 Bilibili
2. 按 F12 打开开发者工具 → Application → Cookies
3. 找到 `SESSDATA` 和 `bili_jct` 的值
4. 使用 `aiview bilibili login --sessdata "xxx" --bili-jct "xxx"` 登录

详见 [平台接入指南](docs/PLATFORM_GUIDE.md)。

### Q: 如何获取抖音 Cookie？

1. 在浏览器中登录抖音网页版
2. 按 F12 打开开发者工具 → Network
3. 刷新页面，找到任意请求的 Cookie 头
4. 复制完整 Cookie 字符串，使用 `aiview douyin login --cookie "xxx"` 登录

详见 [平台接入指南](docs/PLATFORM_GUIDE.md)。

### Q: 音频下载提示需要 ffmpeg？

`aiview bilibili audio` 命令默认会将音频分段切割为 WAV 格式，需要安装 ffmpeg。如不需要切割，可使用 `--no-split` 参数下载完整 M4A 音频。

### Q: 命令返回 "not_authenticated" 错误？

部分接口需要登录后才能访问。请先使用对应平台的 `login` 命令完成认证。

### Q: 命令返回 "permission_denied" 错误？

写操作（如点赞、投币、发评论）需要具有写权限的凭证。请重新登录，确保 Cookie 包含完整的认证信息。

### Q: 搜索结果不完整？

抖音搜索在未登录状态下返回结果有限，建议先登录获取完整结果。

### Q: 如何查看命令帮助？

```bash
# 查看全局帮助
aiview --help

# 查看子命令帮助
aiview bilibili video --help

# 查看具体 Flag
aiview bilibili search --help
```

---

## 技术架构

```
aiview-cli/
├── main.go / root.go           # 入口 + 全局 Flag
├── commands/
│   ├── bilibili/               # Bilibili 命令层（40+ 命令）
│   ├── douyin/                 # 抖音命令层
│   └── xiaohongshu/            # 小红书命令层
├── internal/
│   ├── platform/               # 平台抽象层（Platform 接口 + 全局注册器）
│   │   ├── bilibili/           # Bilibili HTTP Client + WBI 签名
│   │   ├── douyin/             # 抖音 HTTP Client + Cookie 管理
│   │   └── xiaohongshu/        # 小红书 HTTP Client
│   ├── auth/                   # 凭证持久化（读写权限分级）
│   ├── output/                 # 输出格式化（JSON/YAML/Table）
│   ├── config/                 # 配置加载（Viper）
│   └── helper/                 # 嵌套 Map 安全取值
├── pkg/aiview/                 # Go Library 公开 API
├── examples/                   # 示例代码
└── docs/                       # 文档
```

**设计要点**：
- **平台注册器**：新平台实现 `Platform` 接口并 import 即可自动注册
- **Client 接口抽象**：命令层通过接口调用，与具体 API 解耦
- **Payload 模式**：`video` 命令可组合字幕/AI/评论/推荐，非致命错误以 Warning 输出

---

## 依赖

- Go 1.21+
- [Cobra](https://github.com/spf13/cobra) — CLI 框架
- [Viper](https://github.com/spf13/viper) — 配置管理
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI 框架
- [chromedp](https://github.com/chromedp/chromedp) — 浏览器自动化
- FFmpeg — 音频分段切割（可选，未安装时降级为完整下载）

---

## 文档导航

- [API 文档](docs/API.md) — Go Library API 完整参考
- [平台接入指南](docs/PLATFORM_GUIDE.md) — 各平台认证与 API 说明
- [贡献指南](CONTRIBUTING.md) — 代码规范与提交流程
- [示例代码](examples/) — 可直接运行的示例程序

---

## 许可证

MIT License
