# Aiview Go Library API 文档

Aiview 提供 Go 包 API（`github.com/jackwener/aiview/pkg/aiview`），可在 Go 项目中直接引用，访问 Bilibili、抖音、小红书三大平台的数据。

---

## 安装

```bash
go get github.com/jackwener/aiview
```

---

## 核心类型

### Platform

```go
type Platform string

const (
    PlatformBilibili    Platform = "bilibili"
    PlatformDouyin      Platform = "douyin"
    PlatformXiaohongshu Platform = "xiaohongshu"
)
```

平台标识符，用于创建客户端时指定目标平台。

### VideoInfo

视频/笔记信息结构体，各平台通用。

```go
type VideoInfo struct {
    ID       string `json:"id"`        // 内容 ID（BV号/视频ID/笔记ID）
    Title    string `json:"title"`     // 标题
    Author   string `json:"author"`    // 作者昵称
    AuthorID string `json:"author_id"` // 作者 ID
    Play     int    `json:"play"`      // 播放量
    Danmaku  int    `json:"danmaku"`   // 弹幕数
    Like     int    `json:"like"`      // 点赞数
    Coin     int    `json:"coin"`      // 投币数
    Favorite int    `json:"favorite"`  // 收藏数
    Share    int    `json:"share"`     // 分享数
    Duration string `json:"duration"`  // 时长
    URL      string `json:"url"`       // 链接
}
```

### UserInfo

用户信息结构体，各平台通用。

```go
type UserInfo struct {
    ID     string `json:"id"`      // 用户 ID
    Name   string `json:"name"`    // 昵称
    Avatar string `json:"avatar"`  // 头像 URL
    Sign   string `json:"sign"`    // 签名
    Fans   int    `json:"fans"`    // 粉丝数
    Follow int    `json:"follow"`  // 关注数
    Videos int    `json:"videos"`  // 作品数
}
```

### HotItem

热门/趋势项结构体。

```go
type HotItem struct {
    Keyword  string `json:"keyword"`    // 关键词/标题
    HotValue int    `json:"hot_value"`  // 热度值
    Position int    `json:"position"`   // 排名位置
    URL      string `json:"url"`        // 链接
}
```

### SearchItem

搜索结果项结构体。

```go
type SearchItem struct {
    ID     string `json:"id"`      // 内容 ID
    Title  string `json:"title"`   // 标题
    Author string `json:"author"`  // 作者
    URL    string `json:"url"`     // 链接
}
```

---

## 核心函数

### New

创建新的客户端实例。

```go
func New(platformName string) (*Client, error)
```

**参数**：
- `platformName` — 平台标识符，可选值：`"bilibili"`、`"douyin"`、`"xiaohongshu"`

**返回值**：
- `*Client` — 客户端实例
- `error` — 平台不支持时返回错误

**示例**：

```go
client, err := aiview.New("bilibili")
if err != nil {
    log.Fatal(err)
}
fmt.Println(client.PlatformName()) // 输出: bilibili
```

---

### Client.PlatformName

返回当前平台名称。

```go
func (c *Client) PlatformName() string
```

---

### Client.BilibiliClient

获取 Bilibili 专用客户端。仅当平台为 `bilibili` 时可用。

```go
func (c *Client) BilibiliClient() (*BilibiliClient, error)
```

**返回值**：
- `*BilibiliClient` — Bilibili 客户端
- `error` — 平台不匹配或初始化失败时返回错误

---

### Client.DouyinClient

获取抖音专用客户端。仅当平台为 `douyin` 时可用。

```go
func (c *Client) DouyinClient() (*DouyinClient, error)
```

---

### Client.XiaohongshuClient

获取小红书专用客户端。仅当平台为 `xiaohongshu` 时可用。

```go
func (c *Client) XiaohongshuClient() (*XiaohongshuClient, error)
```

---

## BilibiliClient 方法

### GetHotVideos

获取 Bilibili 热门视频列表。

```go
func (b *BilibiliClient) GetHotVideos(page, count int) ([]VideoInfo, error)
```

**参数**：
- `page` — 页码，从 1 开始
- `count` — 每页数量

**返回值**：
- `[]VideoInfo` — 视频列表
- `error` — 请求失败时返回错误

**示例**：

```go
videos, err := biliClient.GetHotVideos(1, 10)
if err != nil {
    log.Fatal(err)
}
for _, v := range videos {
    fmt.Printf("%s - %s (播放: %d)\n", v.Title, v.Author, v.Play)
}
```

---

### GetVideoInfo

获取视频详细信息。

```go
func (b *BilibiliClient) GetVideoInfo(bvid string) (*VideoInfo, error)
```

**参数**：
- `bvid` — 视频 BV 号（如 `"BV1GJ411x7Rq"`）

**返回值**：
- `*VideoInfo` — 视频详情
- `error` — 视频不存在或请求失败时返回错误

**示例**：

```go
info, err := biliClient.GetVideoInfo("BV1GJ411x7Rq")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("标题: %s\n", info.Title)
fmt.Printf("播放: %d, 点赞: %d\n", info.Play, info.Like)
```

---

### GetUserInfo

获取用户信息。

```go
func (b *BilibiliClient) GetUserInfo(uid int) (*UserInfo, error)
```

**参数**：
- `uid` — 用户 ID（整数）

**返回值**：
- `*UserInfo` — 用户信息
- `error` — 用户不存在或请求失败时返回错误

**示例**：

```go
user, err := biliClient.GetUserInfo(37737161)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("昵称: %s, 粉丝: %d\n", user.Name, user.Fans)
```

---

### SearchVideo

按关键词搜索视频。

```go
func (b *BilibiliClient) SearchVideo(keyword string, page int) ([]SearchItem, error)
```

**参数**：
- `keyword` — 搜索关键词
- `page` — 页码，从 1 开始

**返回值**：
- `[]SearchItem` — 搜索结果列表
- `error` — 请求失败时返回错误

**示例**：

```go
results, err := biliClient.SearchVideo("Golang", 1)
if err != nil {
    log.Fatal(err)
}
for _, item := range results {
    fmt.Printf("%s - %s\n", item.Title, item.Author)
}
```

---

## DouyinClient 方法

### GetHotSearch

获取抖音热搜榜。

```go
func (d *DouyinClient) GetHotSearch() ([]HotItem, error)
```

**返回值**：
- `[]HotItem` — 热搜列表
- `error` — 请求失败时返回错误

**示例**：

```go
hotItems, err := dyClient.GetHotSearch()
if err != nil {
    log.Fatal(err)
}
for i, item := range hotItems {
    fmt.Printf("%d. %s (热度: %d)\n", i+1, item.Keyword, item.HotValue)
}
```

---

### Search

搜索抖音视频。

```go
func (d *DouyinClient) Search(keyword string, page, count int) ([]SearchItem, error)
```

**参数**：
- `keyword` — 搜索关键词
- `page` — 页码，从 1 开始
- `count` — 每页数量

**返回值**：
- `[]SearchItem` — 搜索结果列表
- `error` — 请求失败时返回错误

**示例**：

```go
results, err := dyClient.Search("AI技术", 1, 10)
if err != nil {
    log.Fatal(err)
}
for _, item := range results {
    fmt.Printf("%s - %s\n", item.Title, item.Author)
}
```

---

### GetVideoDetail

获取抖音视频详情。

```go
func (d *DouyinClient) GetVideoDetail(videoID string) (*VideoInfo, error)
```

**参数**：
- `videoID` — 视频 ID

**返回值**：
- `*VideoInfo` — 视频详情
- `error` — 视频不存在或请求失败时返回错误

**示例**：

```go
info, err := dyClient.GetVideoDetail("7123456789012345678")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("标题: %s, 播放: %d\n", info.Title, info.Play)
```

---

### GetUserInfo

获取抖音用户信息。

```go
func (d *DouyinClient) GetUserInfo(uid string) (*UserInfo, error)
```

**参数**：
- `uid` — 用户 ID（字符串）

**返回值**：
- `*UserInfo` — 用户信息
- `error` — 用户不存在或请求失败时返回错误

**示例**：

```go
user, err := dyClient.GetUserInfo("123456789")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("昵称: %s, 粉丝: %d\n", user.Name, user.Fans)
```

---

## XiaohongshuClient 方法

### GetHotNotes

获取小红书热门笔记。

```go
func (x *XiaohongshuClient) GetHotNotes() ([]HotItem, error)
```

**返回值**：
- `[]HotItem` — 热门笔记列表
- `error` — 请求失败时返回错误

**示例**：

```go
hotNotes, err := xhsClient.GetHotNotes()
if err != nil {
    log.Fatal(err)
}
for i, item := range hotNotes {
    fmt.Printf("%d. %s\n", i+1, item.Keyword)
}
```

---

### SearchNotes

搜索小红书笔记。

```go
func (x *XiaohongshuClient) SearchNotes(keyword string, page int) ([]SearchItem, error)
```

**参数**：
- `keyword` — 搜索关键词
- `page` — 页码，从 1 开始

**返回值**：
- `[]SearchItem` — 搜索结果列表
- `error` — 请求失败时返回错误

**示例**：

```go
results, err := xhsClient.SearchNotes("旅行攻略", 1)
if err != nil {
    log.Fatal(err)
}
for _, item := range results {
    fmt.Printf("%s - %s\n", item.Title, item.Author)
}
```

---

### GetNoteDetail

获取小红书笔记详情。

```go
func (x *XiaohongshuClient) GetNoteDetail(noteID string) (*VideoInfo, error)
```

**参数**：
- `noteID` — 笔记 ID

**返回值**：
- `*VideoInfo` — 笔记详情（复用 VideoInfo 结构体）
- `error` — 笔记不存在或请求失败时返回错误

**示例**：

```go
note, err := xhsClient.GetNoteDetail("64a1b2c3d4e5f60000000001")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("标题: %s, 点赞: %d\n", note.Title, note.Like)
```

---

### GetUserInfo

获取小红书用户信息。

```go
func (x *XiaohongshuClient) GetUserInfo(userID string) (*UserInfo, error)
```

**参数**：
- `userID` — 用户 ID（字符串）

**返回值**：
- `*UserInfo` — 用户信息
- `error` — 用户不存在或请求失败时返回错误

**示例**：

```go
user, err := xhsClient.GetUserInfo("5f1a2b3c4d5e6f0000000001")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("昵称: %s, 粉丝: %d\n", user.Name, user.Fans)
```

---

## 错误处理

所有方法均返回 `error` 类型，常见错误场景：

### 平台不支持

```go
client, err := aiview.New("unknown_platform")
// err: platform "unknown_platform" not supported
```

### 平台类型不匹配

```go
client, _ := aiview.New("bilibili")
dyClient, err := client.DouyinClient()
// err: not a douyin client
```

### 内容不存在

```go
info, err := biliClient.GetVideoInfo("BV_NOT_EXIST")
// err: 视频不存在或请求失败
```

### 推荐错误处理模式

```go
func fetchVideos() {
    client, err := aiview.New("bilibili")
    if err != nil {
        log.Printf("创建客户端失败: %v", err)
        return
    }

    biliClient, err := client.BilibiliClient()
    if err != nil {
        log.Printf("获取 Bilibili 客户端失败: %v", err)
        return
    }

    videos, err := biliClient.GetHotVideos(1, 10)
    if err != nil {
        log.Printf("获取热门视频失败: %v", err)
        return
    }

    for _, v := range videos {
        fmt.Printf("%s - %s\n", v.Title, v.Author)
    }
}
```

---

## 完整示例

```go
package main

import (
    "fmt"
    "log"

    "github.com/jackwener/aiview/pkg/aiview"
)

func main() {
    // === Bilibili 示例 ===
    fmt.Println("=== Bilibili ===")
    biliClient, err := createBilibiliClient()
    if err != nil {
        log.Printf("Bilibili 初始化失败: %v", err)
    } else {
        // 获取热门视频
        videos, err := biliClient.GetHotVideos(1, 5)
        if err == nil {
            for _, v := range videos {
                fmt.Printf("  %s - %s (播放: %d)\n", v.Title, v.Author, v.Play)
            }
        }

        // 搜索视频
        results, err := biliClient.SearchVideo("Golang", 1)
        if err == nil && len(results) > 0 {
            fmt.Printf("  搜索结果: %s\n", results[0].Title)
        }
    }

    // === 抖音示例 ===
    fmt.Println("\n=== 抖音 ===")
    dyClient, err := createDouyinClient()
    if err != nil {
        log.Printf("抖音初始化失败: %v", err)
    } else {
        hotItems, err := dyClient.GetHotSearch()
        if err == nil {
            for i, item := range hotItems {
                if i >= 5 {
                    break
                }
                fmt.Printf("  %d. %s (热度: %d)\n", i+1, item.Keyword, item.HotValue)
            }
        }
    }

    // === 小红书示例 ===
    fmt.Println("\n=== 小红书 ===")
    xhsClient, err := createXiaohongshuClient()
    if err != nil {
        log.Printf("小红书初始化失败: %v", err)
    } else {
        hotNotes, err := xhsClient.GetHotNotes()
        if err == nil {
            for i, item := range hotNotes {
                if i >= 5 {
                    break
                }
                fmt.Printf("  %d. %s\n", i+1, item.Keyword)
            }
        }
    }
}

func createBilibiliClient() (*aiview.BilibiliClient, error) {
    client, err := aiview.New("bilibili")
    if err != nil {
        return nil, err
    }
    return client.BilibiliClient()
}

func createDouyinClient() (*aiview.DouyinClient, error) {
    client, err := aiview.New("douyin")
    if err != nil {
        return nil, err
    }
    return client.DouyinClient()
}

func createXiaohongshuClient() (*aiview.XiaohongshuClient, error) {
    client, err := aiview.New("xiaohongshu")
    if err != nil {
        return nil, err
    }
    return client.XiaohongshuClient()
}
```

---

## 相关文档

- [README](../README.md) — 项目总览
- [平台接入指南](PLATFORM_GUIDE.md) — 各平台认证说明
- [示例代码](../examples/) — 可直接运行的示例程序
