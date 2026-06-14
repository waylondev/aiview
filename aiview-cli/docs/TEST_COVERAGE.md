# Aiview CLI — 测试覆盖率报告

生成日期：2026-06-14

## 测试环境

- OS: Windows amd64
- Go 版本: 1.26.4
- 编译检查: ✅ 全部通过（`go test -run=^$ ./...` 无编译错误）

## 总览

| 指标 | 数量 |
|------|------|
| 总包数 | 28 |
| 有测试的包 | 12 |
| 无测试的包 | 16 |
| 覆盖率 ≥ 80%（良好） | 2 |
| 覆盖率 60%~80%（中等） | 2 |
| 覆盖率 < 60%（高风险） | 8 |
| 测试文件总数 | 14 |

## 覆盖率分布

```
良好 (≥80%)  : ██ 2
中等 (60-80%): ██ 2
高风险 (<60%): ████████ 8
无测试 (0%)  : ████████████████ 16
```

---

## 一、Internal 包覆盖率

### ✅ 良好（≥ 80%）

| 包 | 覆盖率 | 测试文件 |
|----|--------|----------|
| `internal/ratelimit` | **96.4%** | `ratelimit_test.go` |
| `internal/config` | **91.3%** | `config_test.go` |

### ⚡ 中等（60% ~ 80%）

| 包 | 覆盖率 | 测试文件 |
|----|--------|----------|
| `internal/cache` | **77.8%** | `cache_test.go` |
| `internal/scheduler` | **71.7%** | `scheduler_test.go` |

### 🔴 高风险（< 60%）

| 包 | 覆盖率 | 测试文件 | 建议 |
|----|--------|----------|------|
| `internal/helper` | 51.4% | `helper_test.go` | 需补充边界条件和错误路径测试 |
| `internal/analyzer` | 41.3% | `analyzer_test.go` | 分析逻辑覆盖不足，需增加趋势/对比测试用例 |
| `internal/auth` | 39.6% | `store_test.go` | 认证存储测试不完整，需补充 Token 刷新和过期场景 |
| `internal/platform/douyin` | 24.4% | `client_test.go`, `auth_test.go` | API 客户端测试需大幅扩充 |
| `internal/output` | 21.3% | `formatter_test.go` | 格式化输出各模式覆盖不足 |
| `internal/platform/bilibili` | 14.4% | `client_test.go` | B站客户端测试覆盖极低 |

### ⬜ 无测试（0%）

| 包 | 说明 |
|----|------|
| `internal/browser` | 浏览器自动化模块，无测试 |
| `internal/errors` | 错误定义模块，无测试 |
| `internal/pipeline` | 数据管道模块，无测试 |
| `internal/platform` | 平台公共模块，无测试 |
| `internal/platform/base` | 平台基类，无测试 |
| `internal/platform/bilibili/bilibilitypes` | B站类型定义，无测试 |
| `internal/platform/kuaishou` | 快手平台，无测试 |
| `internal/platform/weibo` | 微博平台，无测试 |
| `internal/platform/xiaohongshu` | 小红书平台，无测试 |
| `internal/platform/zhihu` | 知乎平台，无测试 |
| `internal/storage` | 存储模块，无测试 |

---

## 二、Commands 包覆盖率

| 包 | 覆盖率 | 测试文件 |
|----|--------|----------|
| `commands/bilibili` | 1.3% | `hot_test.go`, `types_test.go` |
| `commands/douyin` | 0.8% | `hot_test.go` |
| `commands` | 0.0% | 无测试 |
| `commands/kuaishou` | 0.0% | 无测试 |
| `commands/weibo` | 0.0% | 无测试 |
| `commands/xiaohongshu` | 0.0% | 无测试 |
| `commands/zhihu` | 0.0% | 无测试 |

---

## 三、测试文件清单

共 **14** 个测试文件：

| 文件路径 | 所属模块 |
|----------|----------|
| `internal/ratelimit/ratelimit_test.go` | 速率限制 |
| `internal/config/config_test.go` | 配置管理 |
| `internal/cache/cache_test.go` | 缓存 |
| `internal/scheduler/scheduler_test.go` | 调度器 |
| `internal/helper/helper_test.go` | 工具函数 |
| `internal/analyzer/analyzer_test.go` | 数据分析 |
| `internal/auth/store_test.go` | 认证存储 |
| `internal/platform/douyin/client_test.go` | 抖音客户端 |
| `internal/platform/douyin/auth_test.go` | 抖音认证 |
| `internal/output/formatter_test.go` | 输出格式化 |
| `internal/platform/bilibili/client_test.go` | B站客户端 |
| `commands/bilibili/hot_test.go` | B站命令 |
| `commands/bilibili/types_test.go` | B站类型 |
| `commands/douyin/hot_test.go` | 抖音命令 |

---

## 四、需要补充测试的模块（按优先级排序）

### 🔴 紧急（覆盖率 < 20%）

1. **`internal/platform/bilibili`** — 14.4%，需补充 API 调用测试
2. **`internal/platform/douyin`** — 24.4%，需补充 API 调用测试
3. **`internal/output`** — 21.3%，需补充 JSON/YAML/CSV/Table 各格式测试

### 🟡 重要（覆盖率 20% ~ 60%）

4. **`internal/auth`** — 39.6%，需补充 Token 刷新、过期、多平台场景
5. **`internal/analyzer`** — 41.3%，需补充趋势分析、跨平台对比测试
6. **`internal/helper`** — 51.4%，需补充边界条件和错误路径

### 🟢 建议（无测试文件）

7. **`internal/errors`** — 错误类型定义测试
8. **`internal/pipeline`** — 数据管线流程测试
9. **`internal/platform/zhihu`** — 知乎平台客户端测试
10. **`internal/platform/weibo`** — 微博平台客户端测试
11. **`internal/platform/kuaishou`** — 快手平台客户端测试

---

## 五、现有测试报告完整性

`docs/` 目录下共有 **7** 个 TEST_REPORT_*.md 文件，**全部完整**：

| 报告文件 | 内容 | 状态 |
|----------|------|------|
| `TEST_REPORT_GLOBAL.md` | 全局命令测试（6个命令/10个子命令） | ✅ |
| `TEST_REPORT_DOUYIN.md` | 抖音平台（11个命令） | ✅ |
| `TEST_REPORT_BILIBILI.md` | B站平台（42个命令） | ✅ |
| `TEST_REPORT_ZHIHU.md` | 知乎平台（4个命令） | ✅ |
| `TEST_REPORT_XIAOHONGSHU.md` | 小红书平台（5个命令） | ✅ |
| `TEST_REPORT_WEIBO.md` | 微博平台（4个命令） | ✅ |
| `TEST_REPORT_KUAISHOU.md` | 快手平台（4个命令） | ✅ |

> 所有报告均以 `2026-06-13` 为生成日期，格式统一，内容完整。

---

## 六、总结

- 编译检查：**全部通过**，无编译错误
- 单元测试覆盖率整体偏低：28 个包中仅 2 个达到 80% 以上
- 最佳模块：`ratelimit`（96.4%）和 `config`（91.3%）
- 最大缺口：6 个平台包（kuaishou/weibo/xiaohongshu/zhihu）完全没有单元测试
- 现有 7 个 CLI 测试报告完整，覆盖所有平台命令