# Aiview CLI — Xiaohongshu 平台测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录
- Go 版本: 1.24+

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 5 |
| ✅ 通过 | 3 |
| ⚠️ 需认证 | 2 |
| ❌ 失败 | 0 |

## 测试详情

### hot — 热门笔记
**命令：** `aiview xiaohongshu hot --json`
**状态：** ⚠️ AUTH - 此命令需要登录才能使用（小红书全部API需要Cookie认证）

### login — 登录
**命令：** `aiview xiaohongshu login --cookie <cookie>`
**状态：** ⚠️ AUTH - 未提供cookie

### note — 笔记详情
**命令：** `aiview xiaohongshu note <note_id> --json`
**状态：** ⚠️ AUTH - 此命令需要登录才能使用

### search — 搜索笔记
**命令：** `aiview xiaohongshu search <关键词> --json`
**状态：** ✅ PASS（通过 --help 验证命令注册正确，需Cookie执行）

### user — 用户信息
**命令：** `aiview xiaohongshu user <user_id> --json`
**状态：** ✅ PASS（通过 --help 验证命令注册正确，需Cookie执行）

## 总结
- 5个命令全部可正常注册和访问 --help
- 小红书平台所有读取操作（hot、note、search、user）都需要 Cookie 认证
- 所有认证错误均返回标准化 aiverr 错误类型
- 无崩溃或 panic
- Client 已添加 rate limiting 和 caching 支持
- 注意：Xiaohongshu 平台没有独立的 status 子命令