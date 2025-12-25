# Change: 在启用 Git Force Push 时添加 WARN 级别日志提醒

## Why
当用户启用 `force_push` 配置或在 rebase 模式下自动启用强制推送时,强制推送会覆盖目标仓库的历史记录,这是一个高风险操作。目前代码中没有显式的警告日志来提醒用户该操作正在执行,可能导致用户在不知情的情况下覆盖重要数据。

为了提高系统的可观测性和用户安全性,需要在执行强制推送前记录一条 WARN 级别的日志,明确告知用户此操作的风险。

## What Changes
- 在 `pkg/git/git.go` 中的 `codePush` 函数中,当 `force` 参数为 `true` 时,在执行 push 命令前添加 WARN 级别日志
- 日志内容应包含:仓库路径、操作类型(强制推送)、风险提示(将覆盖目标仓库历史)
- 确保日志在重试循环开始前只记录一次,避免重复警告

## Impact
- **Affected specs**: `migration-logging`(新增能力)
- **Affected code**: 
  - `pkg/git/git.go:392-413` - `codePush` 函数需要添加 WARN 日志
- **Breaking changes**: 无
- **User impact**: 用户在日志中会看到额外的警告信息,帮助他们识别高风险操作
