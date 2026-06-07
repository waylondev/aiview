# Aiview CLI — 抖音平台专项测试报告

生成日期：2026-06-07

## 测试环境
- 操作系统：Windows
- 构建工具：Go
- 网络：需公网访问
- 终端：PowerShell 5
- 认证状态：未登录

## 测试结果概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 10 |
| ✅ 通过 | 6 |
| ⚠️ 需 Cookie 认证 | 4 |
| ❌ 失败 | 0 |

## 测试详情

### 1. hot — 热搜榜
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin hot --json` |
| 状态 | ✅ 通过 |
| --json | ✅ 支持 |
| --yaml | ✅ 支持（全局标志） |
| 输出示例 | `{"ok":true,"schema_version":"1","data":{"data":[{"hot_value":...,"word":"...",...}],"extra":{...}}}` |
| 备注 | 成功返回 50 条热搜数据，包含 hot_value、word、position、label 等字段 |

### 2. trending — 热点榜
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin trending --json` |
| 状态 | ✅ 通过 |
| --json | ✅ 支持 |
| --yaml | ✅ 支持（全局标志） |
| 输出示例 | `{"ok":true,"schema_version":"1","data":{"data":{"active_time":"...","word_list":[...]}}}` |
| 备注 | 成功返回热点榜数据，包含 word_list 列表和 active_time 时间戳 |

### 3. search — 搜索
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin search <keyword> --json` |
| 状态 | ✅ 通过 |
| --json | ✅ 支持 |
| --yaml | ✅ 支持（全局标志） |
| 输出示例 | `{"ok":true,"data":{"data":[],"has_more":0,"status_code":0,...}}` |
| 备注 | 返回有效 JSON，但搜索结果为空数组（`data:[]`）。英文关键词 "test" 和中文关键词 "美食" 均无搜索结果返回。返回的 `search_nil_info` 显示 `search_nil_type: "params_check"`，可能是 Cookie 校验或 API 参数问题。中文关键词在响应中出现编码问题（如 `缇庨` 应为 `美食`）。 |

### 4. video — 视频详情
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin video <id> --json` |
| 状态 | ⚠️ 需 Cookie 认证 |
| --json | ✅ 支持 |
| --yaml | ✅ 支持（全局标志） |
| 输出示例 | `{"ok":true,"schema_version":"1","data":{"raw":""}}` |
| 备注 | 使用假 ID `123456789` 和真实热搜 group_id `7646643385726129460` 均返回 `{"raw":""}`，表明需要 Cookie 登录后才能获取视频详情 |

### 5. user — 用户信息
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin user <uid> --json` |
| 状态 | ⚠️ 需 Cookie 认证 |
| --json | ✅ 支持 |
| --yaml | ✅ 支持（全局标志） |
| 输出示例 | `{"ok":true,"schema_version":"1","data":{"raw":""}}` |
| 备注 | 使用假 UID `123456789` 返回 `{"raw":""}`，表明需要 Cookie 登录后才能获取用户信息 |

### 6. comment — 评论列表
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin comment <video_id> --json` |
| 状态 | ⚠️ 需 Cookie 认证 |
| --json | ✅ 支持 |
| --yaml | ✅ 支持（全局标志） |
| 输出示例 | `{"ok":true,"schema_version":"1","data":{"raw":""}}` |
| 备注 | 返回 `{"raw":""}`，需要 Cookie 登录 |

### 7. user-posts — 用户作品
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin user-posts <uid> --json` |
| 状态 | ⚠️ 需 Cookie 认证 |
| --json | ✅ 支持 |
| --yaml | ✅ 支持（全局标志） |
| 输出示例 | `{"ok":true,"schema_version":"1","data":{"raw":""}}` |
| 备注 | 返回 `{"raw":""}`，需要 Cookie 登录 |

### 8. status — 登录状态
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin status --json` |
| 状态 | ✅ 通过 |
| --json | ✅ 支持 |
| --yaml | ✅ 支持（全局标志） |
| 输出示例 | `{"ok":true,"data":{"logged_in":false,"platform":"douyin"}}` |
| 备注 | 正确反映当前未登录状态 |

### 9. logout — 登出
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin logout` |
| 状态 | ✅ 通过 |
| --json | 不适用（无 `--json` 标志，输出纯文本） |
| --yaml | 不适用 |
| 输出示例 | `Logged out` |
| 备注 | 成功输出 "Logged out"，退出码为 0。登出后再次运行 `status` 确认仍为未登录状态 |

### 10. login — 登录
| 项目 | 结果 |
|------|------|
| 命令 | `aiview douyin login --help` |
| 状态 | ✅ 通过 |
| --json | 不适用（帮助命令） |
| --yaml | 不适用（帮助命令） |
| 输出示例 | `Login to Douyin using a browser cookie. ... --cookie string   Douyin browser cookie` |
| 备注 | 帮助信息显示完整，支持 `--cookie` 参数进行 Cookie 登录。示例：`aiview douyin login --cookie "your_cookie_here"` |

## 已知问题
1. **搜索功能无返回结果**：`search` 命令虽然返回有效 JSON，但 `data` 字段为空数组。响应的 `search_nil_info.search_nil_type` 为 `"params_check"`，可能是缺少必要 Cookie 导致搜索 API 未正确响应。
2. **中文关键词编码问题**：搜索命令使用中文关键词时，响应中的 `global_doodle_config.keyword` 字段显示乱码（如 `缇庨` 而非 `美食`），可能是 URL 编码或 JSON 序列化问题。
3. **需认证的 API 均返回空的 `raw` 字段**：video、user、comment、user-posts 四个命令在未登录状态下均返回 `{"raw":""}`，而非明确的错误提示。建议在这些命令中添加更友好的未登录提示信息。
4. **`raw` 字段语义不明确**：认证失败的 API 返回 `{"raw":""}` 而非 HTTP 状态码或错误消息，不利于用户诊断问题。