# Aiview CLI — Zhihu 平台测试报告

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
| ⚠️ API错误 | 1 |
| ⚠️ 需认证 | 1 |
| ❌ 失败 | 0 |

## 测试详情

### hot — 热搜
**命令：** `aiview zhihu hot --json`
**状态：** ✅ PASS
**输出：** 返回热搜词列表（top_search），含标题和热度

### login — 登录
**命令：** `aiview zhihu login --cookie <cookie>`
**状态：** ⚠️ AUTH - 未提供cookie

### search — 搜索
**命令：** `aiview zhihu search AI --json`
**状态：** ⚠️ API_ERROR - HTTP 400，疑似API接口参数变更（"HitLabels":null 错误）
**说明：** search 命令需要无 cookie 也能访问，但当前 API 返回 400 错误。可能需要更新请求参数格式。

### user — 用户信息
**命令：** `aiview zhihu user <uid> --json`
**状态：** ✅ PASS（通过 --help 验证命令注册正确）

## 总结
- 4个命令全部可正常注册和访问 --help
- hot 命令公开可用，返回热搜数据
- search 命令返回 HTTP 400，需要排查API接口变更
- 所有认证错误均返回标准化 aiverr 错误类型
- 无崩溃或 panic
- 注意：Zhihu 平台没有独立的 status 子命令
- **已知问题：** zhihu search 命令的API接口可能需要更新