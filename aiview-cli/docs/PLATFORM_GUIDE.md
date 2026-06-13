# 平台接入指南

本文档详细说明各平台的认证流程、Cookie 获取方法、API 接口说明及常见问题。

---

## 目录

- [Bilibili 平台](#bilibili-平台)
- [抖音平台](#抖音平台)
- [小红书平台](#小红书平台)
- [通用问题](#通用问题)

---

## Bilibili 平台

### 认证方式

Bilibili 支持两种认证方式：

#### 1. Cookie 登录（推荐）

需要提供以下 Cookie 值：
- `SESSDATA` — 会话标识
- `bili_jct` — CSRF Token（用于写操作）

**获取步骤**：

1. 在浏览器中访问 [bilibili.com](https://www.bilibili.com) 并登录
2. 按 `F12` 打开开发者工具
3. 切换到 **Application**（应用）标签
4. 在左侧导航中找到 **Cookies** → `https://www.bilibili.com`
5. 找到以下 Cookie 值并复制：
   - `SESSDATA`
   - `bili_jct`
6. 使用命令登录：

```bash
aiview bilibili login --sessdata "<SESSDATA>" --bili-jct "<BILI_JCT>"
```

**可选参数**（高级用户）：
- `--ac_time_value` — 加速验证
- `--buvid3` / `--buvid4` — 设备标识
- `--dedeuserid` — 用户 ID

#### 2. 扫码登录

```bash
aiview bilibili login
```

命令会生成二维码，使用 Bilibili APP 扫码完成登录。

### 凭证存储

登录凭证自动保存在 `~/.aiview/bilibili_credential.json`，包含以下字段：

```json
{
  "sessdata": "xxx",
  "bili_jct": "xxx",
  "ac_time_value": "xxx",
  "buvid3": "xxx",
  "buvid4": "xxx",
  "dede_userid": "xxx",
  "saved_at": "2024-01-01T12:00:00Z"
}
```

### 权限分级

Bilibili 凭证分为两种权限级别：

- **读权限** — 查看视频、用户、评论等（仅需 `SESSDATA`）
- **写权限** — 点赞、投币、发评论、发弹幕等（需要 `SESSDATA` + `bili_jct`）

写操作时会验证 `bili_jct` 是否存在，缺失则返回 `permission_denied` 错误。

### 凭证过期

凭证默认 7 天后过期。过期后需要重新登录：

```bash
aiview bilibili login --sessdata "<NEW_SESSDATA>" --bili-jct "<NEW_BILI_JCT>"
```

### 常用 API 接口

| 接口 | 说明 | 认证要求 |
|------|------|----------|
| 视频详情 | `/x/web-interface/view` | 无需登录 |
| 热门视频 | `/x/web-interface/popular` | 无需登录 |
| 搜索视频 | `/x/web-interface/search/type` | 无需登录 |
| 用户信息 | `/x/web-interface/card` | 无需登录 |
| 发送弹幕 | `/x/v2/dm/post` | 需登录 + 写权限 |
| 发送评论 | `/x/v2/reply/add` | 需登录 + 写权限 |
| 点赞 | `/x/web-interface/archive/like` | 需登录 + 写权限 |

### WBI 签名

Bilibili 部分接口需要 WBI 签名验证。Aiview 已内置 WBI 签名实现，自动处理签名逻辑，无需用户干预。

### 常见问题

#### Q: 登录后仍然提示 "not_authenticated"？

**原因**：
1. 凭证已过期（超过 7 天）
2. Cookie 值不完整（缺少必要字段）

**解决**：
```bash
# 检查登录状态
aiview bilibili status

# 重新登录
aiview bilibili login --sessdata "<SESSDATA>" --bili-jct "<BILI_JCT>"
```

#### Q: 写操作提示 "permission_denied"？

**原因**：登录时未提供 `bili_jct`。

**解决**：重新登录，确保提供完整的 `bili_jct` 值。

#### Q: 搜索结果为空或不完整？

**原因**：
1. 关键词过于热门，触发风控
2. 未登录导致结果受限

**解决**：
1. 尝试更换关键词
2. 登录后重试

#### Q: 弹幕导出失败？

**原因**：视频不存在或无弹幕。

**解决**：
```bash
# 检查视频是否存在
aiview bilibili video BV1GJ411x7Rq

# 确认视频有弹幕后再导出
aiview bilibili danmaku BV1GJ411x7Rq -o ./data/
```

---

## 抖音平台

### 认证方式

抖音使用 Cookie 认证，需要提供完整的 Cookie 字符串。

**获取步骤**：

1. 在浏览器中访问 [douyin.com](https://www.douyin.com) 并登录
2. 按 `F12` 打开开发者工具
3. 切换到 **Network**（网络）标签
4. 刷新页面，找到任意请求（如 `hot` 或 `trending`）
5. 点击请求，在 **Headers**（标头）中找到 **Cookie** 字段
6. 复制完整的 Cookie 字符串
7. 使用命令登录：

```bash
aiview douyin login --cookie "<完整Cookie字符串>"
```

**Cookie 示例**：
```
ttwid=xxx; passport_csrf_token=xxx; msToken=xxx; ...
```

### 凭证存储

登录凭证自动保存在 `~/.aiview/douyin_credential.json`：

```json
{
  "cookie": "ttwid=xxx; passport_csrf_token=xxx; ..."
}
```

### 接口权限

| 接口 | 说明 | 认证要求 |
|------|------|----------|
| 热榜 | `/hot/search` | 无需登录 |
| 趋势话题 | `/trending` | 无需登录 |
| 搜索 | `/search` | 无需登录（登录结果更全） |
| 视频详情 | `/video/detail` | 需登录 |
| 用户信息 | `/user/info` | 需登录 |
| 评论列表 | `/comment/list` | 需登录 |
| 用户作品 | `/user/posts` | 需登录 |

### 常见问题

#### Q: 搜索结果不完整？

**原因**：未登录状态下，抖音返回结果有限。

**解决**：
```bash
# 先登录
aiview douyin login --cookie "<COOKIE>"

# 再搜索
aiview douyin search "AI技术"
```

#### Q: 视频详情返回 "video not found"？

**原因**：
1. 视频 ID 不正确
2. 视频已被删除或设为私密
3. 未登录

**解决**：
1. 确认视频 ID 正确（从分享链接中提取）
2. 登录后重试

#### Q: 如何从分享链接提取视频 ID？

抖音分享链接格式：
```
https://www.douyin.com/video/7123456789012345678
```

视频 ID 为 URL 末尾的数字部分：`7123456789012345678`

---

## 小红书平台

### 认证方式

小红书使用 Cookie 认证，与抖音类似。

**获取步骤**：

1. 在浏览器中访问 [xiaohongshu.com](https://www.xiaohongshu.com) 并登录
2. 按 `F12` 打开开发者工具
3. 切换到 **Network**（网络）标签
4. 刷新页面，找到任意请求
5. 复制完整的 Cookie 字符串
6. 使用命令登录：

```bash
aiview xiaohongshu login --cookie "<完整Cookie字符串>"
```

### 凭证存储

登录凭证自动保存在 `~/.aiview/xiaohongshu_credential.json`：

```json
{
  "cookie": "a1=xxx; webId=xxx; ..."
}
```

### 接口权限

| 接口 | 说明 | 认证要求 |
|------|------|----------|
| 热门笔记 | `/explore/hot` | 无需登录 |
| 搜索笔记 | `/search/notes` | 无需登录 |
| 笔记详情 | `/note/detail` | 无需登录 |
| 用户信息 | `/user/info` | 无需登录 |

### 常见问题

#### Q: 热门笔记返回为空？

**原因**：小红书接口可能触发风控。

**解决**：
1. 等待几分钟后重试
2. 登录后重试

#### Q: 如何获取笔记 ID？

小红书笔记链接格式：
```
https://www.xiaohongshu.com/explore/64a1b2c3d4e5f60000000001
```

笔记 ID 为 URL 末尾的字符串：`64a1b2c3d4e5f60000000001`

---

## 通用问题

### 凭证管理

#### 查看所有平台登录状态

```bash
aiview bilibili status
aiview douyin status
aiview xiaohongshu status
```

#### 登出

```bash
aiview bilibili logout
aiview douyin logout
aiview xiaohongshu logout
```

登出会删除本地凭证文件，但不会影响浏览器中的登录状态。

### 网络问题

#### Q: 请求超时？

**原因**：网络不稳定或平台接口响应慢。

**解决**：
1. 检查网络连接
2. 增加超时时间（配置文件）：

```yaml
# ~/.aiview/config.yaml
platforms:
  bilibili:
    timeout: 60  # 增加到 60 秒
```

#### Q: 返回 403/429 错误？

**原因**：触发平台风控或请求频率过高。

**解决**：
1. 等待几分钟后重试
2. 降低请求频率
3. 登录后重试（登录用户风控阈值更高）

### 数据格式

#### JSON 输出

所有命令均支持 `--json` 参数，输出标准 JSON 格式：

```bash
aiview bilibili hot --json
```

输出示例：
```json
[
  {
    "id": "BV1GJ411x7Rq",
    "title": "视频标题",
    "author": "UP主",
    "play": 1234567,
    "like": 12345
  }
]
```

#### YAML 输出

```bash
aiview bilibili hot --yaml
```

输出示例：
```yaml
- id: BV1GJ411x7Rq
  title: 视频标题
  author: UP主
  play: 1234567
  like: 12345
```

### 调试模式

启用详细日志输出，便于排查问题：

```bash
aiview bilibili hot -v
```

`-v` 或 `--verbose` 参数会输出请求 URL、响应状态码等调试信息。

---

## 安全建议

1. **保护凭证** — 不要将 `~/.aiview/` 目录下的凭证文件提交到 Git 仓库
2. **定期更新** — Cookie 有效期有限，过期后及时更新
3. **避免频繁请求** — 过高的请求频率可能触发平台风控
4. **使用官方接口** — 不要使用第三方破解接口，存在账号封禁风险

---

## 相关文档

- [README](../README.md) — 项目总览
- [API 文档](API.md) — Go Library API 参考
- [贡献指南](../CONTRIBUTING.md) — 代码规范与提交流程
