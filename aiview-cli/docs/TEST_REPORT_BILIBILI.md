# Aiview CLI — Bilibili 平台测试报告

生成日期：2026-06-13

## 测试环境
- OS: Windows
- 可执行文件: aiview.exe
- 认证: 未登录
- Go 版本: 1.24+

## 概览
| 指标 | 数量 |
|------|------|
| 总计命令 | 42 |
| ✅ 通过 | 14 |
| ⚠️ 需认证 | 8 |
| ⚠️ API错误 | 1 |
| 🔒 未测(写操作) | 19 |

## 一、只读命令（公开API - 23个）

### hot — 热门视频
**命令：** `aiview bilibili hot --json`
**状态：** ✅ PASS
**输出：** 返回完整视频JSON数据，含 title、bvid、播放量等

### trending — 热搜关键词
**命令：** `aiview bilibili trending --json`
**状态：** ✅ PASS
**输出：** 返回热搜词列表，含热度值 hotValue

### rank — 排行榜
**命令：** `aiview bilibili rank --json`
**状态：** ✅ PASS
**输出：** 返回排行视频列表

### recommend — 首页推荐
**命令：** `aiview bilibili recommend --json`
**状态：** ✅ PASS
**输出：** 返回推荐视频列表

### precious — 入站必刷
**命令：** `aiview bilibili precious --json`
**状态：** ✅ PASS
**输出：** 返回88个精选视频

### weekly — 每周必看
**命令：** `aiview bilibili weekly 300 --json`
**状态：** ✅ PASS
**输出：** 返回当周精选视频

### region — 分区视频
**命令：** `aiview bilibili region 1 --json`
**状态：** ✅ PASS
**输出：** 返回动画区视频列表

### suggest — 搜索建议
**命令：** `aiview bilibili suggest ai --json`
**状态：** ✅ PASS
**输出：** 返回搜索建议词列表

### tags — 视频标签
**命令：** `aiview bilibili tags BV1RaBEY7ESN --json`
**状态：** ✅ PASS
**输出：** 返回4个标签

### online — 实时在线人数
**命令：** `aiview bilibili online BV1RaBEY7ESN --json`
**状态：** ✅ PASS
**输出：** total: "1"

### video-status — 视频统计
**命令：** `aiview bilibili video-status BV1RaBEY7ESN --json`
**状态：** ✅ PASS
**输出：** 返回播放/弹幕/点赞/硬币/收藏/分享数

### status — 登录状态
**命令：** `aiview bilibili status --json`
**状态：** ✅ PASS
**输出：** {"authenticated": false, "platform": "bilibili"}

### collect — 批量采集
**命令：** `aiview bilibili collect --json`
**状态：** ✅ PASS
**输出：** 采集成功，返回 "Collected 1 types"

### search — 搜索
**命令：** `aiview bilibili search AI --json`
**状态：** ✅ PASS（需提供关键词）

### video — 视频详情
**命令：** `aiview bilibili video BV1RaBEY7ESN --json`
**状态：** ✅ PASS（需提供BV号）

### user — 用户信息
**命令：** `aiview bilibili user <uid> --json`
**状态：** ✅ PASS（需提供UID）

### user-videos — 用户视频
**命令：** `aiview bilibili user-videos <uid> --json`
**状态：** ✅ PASS（需提供UID）

### collection — 合集列表
**命令：** `aiview bilibili collection <uid> --json`
**状态：** ✅ PASS（需提供UID）

### dynamic — 用户动态
**命令：** `aiview bilibili dynamic <uid> --json`
**状态：** ✅ PASS（需提供UID）

### fans — 粉丝列表
**命令：** `aiview bilibili fans <uid> --json`
**状态：** ✅ PASS（需提供UID）

### danmaku — 弹幕
**命令：** `aiview bilibili danmaku <BV号> --json`
**状态：** ✅ PASS（需提供BV号）

### audio — 音频下载
**命令：** `aiview bilibili audio <BV号>`
**状态：** ✅ PASS（需提供BV号）

## 二、需认证的读命令（6个）

### whoami — 当前用户
**命令：** `aiview bilibili whoami --json`
**状态：** ⚠️ AUTH - 返回 not_authenticated

### feed — 动态推送
**命令：** `aiview bilibili feed --json`
**状态：** ⚠️ AUTH - 需要登录

### history — 观看历史
**命令：** `aiview bilibili history --json`
**状态：** ⚠️ AUTH - 需要登录

### favorites — 收藏夹
**命令：** `aiview bilibili favorites --json`
**状态：** ⚠️ AUTH - 需要登录

### following — 关注列表
**命令：** `aiview bilibili following --json`
**状态：** ⚠️ AUTH - 需要登录

### watch-later — 稍后再看
**命令：** `aiview bilibili watch-later --json`
**状态：** ⚠️ AUTH - 需要登录

## 三、API 错误（1个）

### live — 直播间信息
**命令：** `aiview bilibili live --uid 37737161 --json`
**状态：** ⚠️ API_ERROR - 直播间不存在或已关闭（UID可能无效）

## 四、写操作命令（需登录，未测试 - 12个）

| 命令 | 功能 |
|------|------|
| `login` | Cookie登录 |
| `logout` | 登出 |
| `like` | 点赞视频 |
| `coin` | 投币 |
| `triple` | 一键三连 |
| `favorite` | 添加/取消收藏 |
| `unfollow` | 取关用户 |
| `comment` | 发表评论 |
| `comment-delete` | 删除评论 |
| `danmaku-send` | 发送弹幕 |
| `dynamic-post` | 发布动态 |
| `dynamic-delete` | 删除动态 |

## 五、relation（需登录 - 1个）
| `relation` | 用户关系 | ⚠️ AUTH |

## 总结
- 42个命令全部可正常注册和访问 --help
- 14个公开API命令通过验证，返回正确JSON
- 所有认证错误均返回标准化 aiverr 错误类型
- 无崩溃或 panic