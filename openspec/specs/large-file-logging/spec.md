# large-file-logging Specification

## Purpose
TBD - created by archiving change downgrade-large-file-log-level. Update Purpose after archive.
## Requirements
### Requirement: Push Failure Log Levels by Error Type
系统在 Git push 操作失败时 SHALL 根据错误类型使用适当的日志级别:
- 对于可自动恢复的大文件超限错误，系统 SHALL 使用 WARN 级别
- 对于其他真正的致命错误，系统 SHALL 使用 ERROR 级别

**修改位置**: `pkg/git/git.go` 的 `Push()` 函数（第 307-310 行）

**修改内容**:
在错误处理中添加条件判断，使用 `IsExceededLimitError(output)` 函数区分错误类型:
- 如果是 `exceeded limit` 错误 → WARN 级别 + 说明文字
- 如果是其他错误 → ERROR 级别（保持原有行为）

**修改理由**:
1. 大文件超限错误会被 `migrate.go` 自动处理（通过 git lfs migrate）
2. 第一次 push 失败是预期的，不应视为错误
3. ERROR 日志会误导用户和监控系统
4. 其他错误（网络、权限等）仍需保持 ERROR 级别以便及时发现

#### Scenario: 大文件超限 push 失败使用 WARN 日志
- **GIVEN** 系统正在执行 `git push` 操作
- **WHEN** push 失败且错误输出包含 "exceeded limit"
- **AND** 启用了 `use_lfs_migrate` 配置
- **THEN** 系统 SHALL 记录 WARN 级别日志: "xxx 裸仓push失败(文件超过大小限制，系统将自动处理): ..."
- **AND** 不记录 ERROR 级别日志
- **AND** 返回错误以触发自动修复流程

#### Scenario: 网络错误 push 失败使用 ERROR 日志
- **GIVEN** 系统正在执行 `git push` 操作
- **WHEN** push 因网络问题失败
- **AND** 错误输出不包含 "exceeded limit"
- **THEN** 系统 SHALL 记录 ERROR 级别日志: "xxx 裸仓push失败: ..."
- **AND** 不记录 WARN 日志
- **AND** 返回错误导致迁移失败

#### Scenario: 权限错误 push 失败使用 ERROR 日志
- **GIVEN** 系统正在执行 `git push` 操作
- **WHEN** push 因权限问题失败（如 token 无效）
- **AND** 错误输出不包含 "exceeded limit"
- **THEN** 系统 SHALL 记录 ERROR 级别日志: "xxx 裸仓push失败: ..."
- **AND** 不记录 WARN 日志
- **AND** 返回错误导致迁移失败

#### Scenario: 大文件自动修复后重新 push 成功
- **GIVEN** 第一次 push 因大文件超限失败（记录 WARN 日志）
- **WHEN** 系统执行 `git lfs migrate import` 处理大文件
- **AND** 处理成功后重新执行 push
- **AND** 第二次 push 成功
- **THEN** 系统 SHALL 记录 INFO 日志: "xxx 裸仓push成功"
- **AND** 迁移继续正常进行
- **AND** 整个过程无 ERROR 日志

#### Scenario: 大文件修复失败仍报告错误
- **GIVEN** 第一次 push 因大文件超限失败（记录 WARN 日志）
- **WHEN** 系统执行 `git lfs migrate import` 但失败
- **OR** 第二次 push 仍然失败
- **THEN** 系统 SHALL 返回错误导致迁移失败
- **AND** 在 `migrate.go` 层面记录最终失败的 ERROR

