# migration-logging Specification

## Purpose
定义代码仓库迁移过程中的日志记录规范,确保关键操作(特别是高风险操作如强制推送)有适当的日志警告,提高系统可观测性和用户安全性。

## Requirements
### Requirement: Force Push Warning
当系统执行 Git 强制推送操作时,系统 SHALL 在推送前记录一条 WARN 级别的日志,明确提醒用户此操作的高风险性质。

#### Scenario: 用户显式配置强制推送
- **WHEN** 用户在配置文件中设置 `migrate.force_push=true`
- **AND** 系统执行仓库推送操作
- **THEN** 系统在执行 `git push -f` 前记录 WARN 级别日志
- **AND** 日志内容包含:仓库路径、强制推送标识、覆盖目标历史的风险提示

#### Scenario: Rebase 模式自动启用强制推送
- **WHEN** 用户启用 rebase 迁移模式
- **AND** CNB 侧仓库存在(非空)
- **AND** 系统自动将 `isForcePush` 设置为 `true`
- **THEN** 系统在执行强制推送前记录 WARN 级别日志
- **AND** 日志内容说明因 rebase 模式自动启用强制推送

#### Scenario: 普通推送不输出警告
- **WHEN** 用户未配置 `force_push` 或配置为 `false`
- **AND** 系统执行普通推送(非 rebase 场景)
- **THEN** 系统不输出强制推送警告日志

#### Scenario: 警告日志只输出一次
- **WHEN** 强制推送因网络等原因重试
- **THEN** WARN 级别的警告日志只在第一次推送尝试前输出
- **AND** 重试时不再重复输出相同的警告

