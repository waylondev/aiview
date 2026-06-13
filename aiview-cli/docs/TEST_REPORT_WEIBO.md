# Aiview CLI — Weibo 平台测试报告

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
| ✅ 通过 | 2 |
| ⚠️ 需认证 | 2 |
| ❌ 失败 | 0 |

## 测试详情

### hot — 热搜榜
**命令：** `aiview weibo hot --json`
**状态：** ✅ PASS
**输出：** 返回热搜榜JSON数据，含实时热搜关键词

### login — 登录
**命令：** `aiview weibo login --cookie <cookie>`
**状态：** ⚠️ AUTH - 未提供cookie

### search — 搜索
**命令：** `aiview weibo search <关键词> --json`
**状态：** ⚠️ AUTH - 重定向到登录页（ok: -100）

### user — 用户信息
**命令：** `aiview weibo user <uid> --json`
**状态：** ✅ PASS（需提供UID，未登录可访问公开用户信息）

## 总结
- 4个命令全部可正常注册和访问 --help
- hot 命令公开可用，返回热搜数据
- search 需要登录认证
- 所有认证错误均返回标准化 aiverr 错误类型
- 无崩溃或 panic
- 注意：Weibo 平台没有独立的 status 子命令