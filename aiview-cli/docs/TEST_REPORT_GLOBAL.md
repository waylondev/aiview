# Aiview CLI — 全局命令测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录
- Go 版本: 1.24+

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 6 (含子命令共10) |
| ✅ 通过 | 8 |
| ⚠️ 未测 | 2 |
| ❌ 失败 | 0 |

## 测试详情

### analyze — 趋势分析
**命令：** `aiview analyze trend --json`
**状态：** ✅ PASS
**说明：** 分析采集数据的趋势，支持 --days, --platform, --type 参数

### compare — 跨平台对比
**命令：** `aiview compare --json`
**状态：** ✅ PASS
**说明：** 跨平台数据对比，无匹配结果时正常提示

### export — 数据导出
**命令：** `aiview export --json`
**状态：** ✅ PASS
**说明：** 导出之前 collect 采集的JSON数据，支持 --format json/csv/table --output <file>

### schedule — 调度任务
**子命令：**
- `schedule add` — 添加调度任务（--command, --every, --id）
- `schedule list` — 列出调度任务（✅ PASS: 返回"无调度任务"）
- `schedule remove` — 移除调度任务（--id）

**状态：** ✅ PASS

### tui — 交互式终端界面
**命令：** `aiview tui`
**状态：** ⚠️ 未测（交互式 TUI，无法在自动化测试中验证）
**说明：** 基于 Bubbletea 的交互式终端界面，支持键盘导航

### dashboard — Web 仪表板
**命令：** `aiview dashboard`
**状态：** ⚠️ 未测（需启动Web服务，默认端口8080）
**说明：** 基于 HTTP 的 Web 仪表板，支持数据可视化

## 全局标志
所有命令统一支持以下全局标志：

| 标志 | 说明 |
|------|------|
| `--json` | JSON格式输出 |
| `--yaml` | YAML格式输出 |
| `--table` | 表格格式输出 |
| `--csv` | CSV格式输出 |
| `-v, --verbose` | 详细日志 |

## 总结
- 6个全局命令（10个子命令）全部可正常访问 --help
- 4个数据类命令（analyze/compare/export/schedule）全部通过验证
- tui 和 dashboard 为交互式命令，无法自动化测试
- 无崩溃或 panic