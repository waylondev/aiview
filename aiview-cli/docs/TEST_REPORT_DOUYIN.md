# Aiview CLI — Douyin 平台测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录
- Go 版本: 1.24+

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 11 |
| ✅ 通过 | 2 |
| ⚠️ 需认证 | 9 |
| ❌ 失败 | 0 |

## 一、公开命令（无需认证 - 2个）

### hot — 热搜
**命令：** `aiview douyin hot --json`
**状态：** ✅ PASS
**输出：** 返回热搜数据，含 word、hotValue

### status — 登录状态
**命令：** `aiview douyin status --json`
**状态：** ✅ PASS
**输出：** {"logged_in": false, "platform": "douyin"}

## 二、需认证命令（9个）

### trending — 趋势
**命令：** `aiview douyin trending --json`
**状态：** ⚠️ AUTH - 需要登录

### search — 搜索
**命令：** `aiview douyin search <关键词> --json`
**状态：** ⚠️ AUTH - 需要登录

### video — 视频详情
**命令：** `aiview douyin video <url/id> --json`
**状态：** ⚠️ AUTH - 需要登录

### user — 用户信息
**命令：** `aiview douyin user <uid> --json`
**状态：** ⚠️ AUTH - 需要登录

### user-posts — 用户作品
**命令：** `aiview douyin user-posts <uid> --json`
**状态：** ⚠️ AUTH - 需要登录

### comment — 评论
**命令：** `aiview douyin comment <video_id> --json`
**状态：** ⚠️ AUTH - 需要登录

### collect — 采集
**命令：** `aiview douyin collect --json`
**状态：** ⚠️ AUTH - 需要登录

### login — 登录
**命令：** `aiview douyin login --cookie <cookie>`
**状态：** ⚠️ AUTH - 未提供cookie

### logout — 登出
**命令：** `aiview douyin logout`
**状态：** ⚠️ AUTH - 未登录，无cookie可清除

## 总结
- 11个命令全部可正常注册和访问 --help
- Douyin 平台大部分操作需要 Cookie 认证（仅 hot 和 status 公开）
- 所有认证错误均返回标准化 aiverr 错误类型
- 无崩溃或 panic