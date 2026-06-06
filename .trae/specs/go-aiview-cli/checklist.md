# Checklist

## 项目结构
- [x] 目录结构符合 `cmd/` + `internal/` + `pkg/` 的 Go 标准布局
- [x] `go.mod` 已初始化，依赖声明正确
- [x] `go build ./...` 编译通过

## 平台抽象层
- [x] `Platform` 接口定义清晰，包含 `Name()`、`Commands()`、`NewClient()` 方法
- [x] `Registry` 机制支持 `Register()` 和 `GetPlatform()` 操作
- [x] 根命令正确集成平台注册，所有平台命令显示在 `--help` 中

## Bilibili 平台
- [x] `aiview bilibili video <BV>` 输出视频详情（标题、UP主、时长、播放量等）
- [x] `aiview bilibili video <BV> --subtitle` 输出字幕内容
- [x] `aiview bilibili video <BV> --ai` 输出 AI 总结
- [x] `aiview bilibili video <BV> --comments` 输出评论
- [x] `aiview bilibili video <BV> --related` 输出相关推荐
- [x] `aiview bilibili search "关键词" --type video` 搜索视频
- [x] `aiview bilibili search "关键词" --type user` 搜索用户
- [x] `aiview bilibili user <uid>` 查看用户信息
- [x] `aiview bilibili user-videos <uid>` 查看用户视频列表
- [x] `aiview bilibili hot` 热门视频
- [x] `aiview bilibili rank` 排行榜
- [x] `aiview bilibili feed` 动态时间线（需登录）
- [x] `aiview bilibili favorites` 收藏夹（需登录）
- [x] `aiview bilibili following` 关注列表（需登录）
- [x] `aiview bilibili history` 观看历史（需登录）
- [x] `aiview bilibili watch-later` 稍后再看（需登录）
- [x] `aiview bilibili like <BV>` 点赞（需登录）
- [x] `aiview bilibili coin <BV>` 投币（需登录）
- [x] `aiview bilibili triple <BV>` 一键三连（需登录）
- [x] `aiview bilibili unfollow <uid>` 取消关注（需登录）
- [x] `aiview bilibili audio <BV>` 音频下载
- [x] `aiview bilibili login` QR 码登录
- [x] `aiview bilibili logout` 退出登录
- [x] `aiview bilibili status` 登录状态检查
- [x] `aiview bilibili whoami` 当前用户信息

## 输出格式
- [x] `--json` 输出符合 agent 信封格式：`{ok, schema_version, data}`
- [x] `--yaml` 输出符合 agent 信封格式
- [x] 非 TTY 环境默认输出 YAML
- [x] 错误输出也符合 agent 信封格式

## 错误处理
- [x] 错误的 BV 号给出清晰的错误提示
- [x] 网络错误给出 `network_error` code
- [x] 未登录操作给出 `not_authenticated` code
- [x] 不存在的资源给出 `not_found` code

## 扩展性
- [x] 新增平台只需实现 `Platform` 接口并调用 `Register()`
- [x] 平台内部命令通过 `init()` 自动注册，无需修改核心代码
- [x] 配置系统支持平台级配置隔离