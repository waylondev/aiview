# Aiview CLI — Kuaishou 平台测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录
- Go 版本: 1.24+

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 4 |
| ✅ 通过 | 3 |
| ⚠️ 需认证 | 1 |
| ❌ 失败 | 0 |

## 测试详情

### hot — 热搜
**命令：** `aiview kuaishou hot --json`
**状态：** ✅ PASS
**输出：** 返回热搜JSON数据

### login — 登录
**命令：** `aiview kuaishou login --cookie <cookie>`
**状态：** ⚠️ AUTH - 未提供cookie

### search — 搜索
**命令：** `aiview kuaishou search <关键词> --json`
**状态：** ✅ PASS
**输出：** 返回 result: 2（无匹配结果但API正常响应JSON）

### user — 用户信息
**命令：** `aiview kuaishou user <uid> --json`
**状态：** ✅ PASS（需提供UID，未登录可通过GraphQL访问公开用户信息）

## 总结
- 4个命令全部可正常注册和访问 --help
- hot、search、user 均公开可用，通过 GraphQL API
- 所有认证错误均返回标准化 aiverr 错误类型
- 无崩溃或 panic
- 注意：Kuaishou 平台没有独立的 status 子命令