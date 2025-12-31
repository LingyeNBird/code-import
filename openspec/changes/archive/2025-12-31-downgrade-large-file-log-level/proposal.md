# Downgrade Large File Log Level

## Change ID
`downgrade-large-file-log-level`

## Why
当前实现中,当源仓库包含超过 256M 的大文件时,系统会输出 ERROR 级别的日志。这会给用户带来误解,让他们认为迁移失败了,但实际上:
1. 系统有自动处理机制(git lfs migrate)来处理大文件
2. 大文件本身并不是致命错误,而是一个需要特殊处理的情况
3. 使用 ERROR 级别会导致监控告警被触发,增加运维负担

ERROR 级别应该保留给真正的失败场景,而大文件检测属于警告性质的信息。

## What Changes
修改 `pkg/git/git.go` 中的 `Push()` 函数，根据错误类型使用不同的日志级别:
- 如果 push 失败是因为文件超过 256MB 限制（`exceeded limit`），使用 WARN 级别
- 其他原因的 push 失败，保持 ERROR 级别

**受影响的代码位置**:
- `Push()` 函数 (第 307-310 行): 裸仓 push 失败的错误处理

**变更内容**:
```go
// 当前实现 (所有错误都是 ERROR):
out, err := codePush(repoPath, pushURL, repoPath, forcePush)
if err != nil {
    logger.Logger.Errorf("%s 裸仓push失败: %s", repoPath, err)
    return out, err
}

// 修改后 (根据错误类型区分日志级别):
out, err := codePush(repoPath, pushURL, repoPath, forcePush)
if err != nil {
    // 如果是大文件超限错误，使用 WARN 级别（系统会自动处理）
    if IsExceededLimitError(out) {
        logger.Logger.Warnf("%s 裸仓push失败(文件超过大小限制，系统将自动处理): %s", repoPath, err)
    } else {
        // 其他错误仍使用 ERROR 级别
        logger.Logger.Errorf("%s 裸仓push失败: %s", repoPath, err)
    }
    return out, err
}
```

**变更原因**:
当 push 失败是由于大文件超过 256M 限制时，系统会在 `migrate.go:638-644` 自动调用 `git lfs migrate` 处理大文件，然后重新 push。这种情况下第一次 push 失败是预期的可恢复场景，应该使用 WARN 级别。而其他原因（如网络错误、权限问题等）的 push 失败是真正的错误，应保持 ERROR 级别。

## Summary
将大文件处理相关的日志级别从 ERROR 降级到 WARN,减少误导性的错误日志,使日志级别更准确地反映问题的严重程度。

## Problem Statement

### Current Behavior
当迁移包含大文件(>256M)的仓库时,系统会:
1. 尝试下载/推送 LFS 文件
2. 如果失败,输出 ERROR 级别日志
3. 继续执行自动修复逻辑(git lfs migrate)
4. 最终成功完成迁移

**问题**:
- ERROR 日志让用户误以为迁移失败
- 监控系统可能触发不必要的告警
- 日志级别与实际严重程度不匹配

### Gap
系统已经具备了处理大文件的能力(通过 `useLfsMigrate` 配置和 `FixExceededLimitError` 函数),但日志级别未能反映这一点。大文件场景应该被视为"需要特殊处理的警告",而不是"错误"。

### Impact
**当前影响**:
- 用户体验:看到 ERROR 日志后担心迁移失败
- 运维负担:需要调查大量非致命的 ERROR 日志
- 监控噪音:触发不必要的告警

**预期改进**:
- 更清晰的日志语义
- 减少误导性信息
- 降低运维成本

## Proposed Solution

### Approach
修改 `pkg/git/git.go` 的 `Push()` 函数，在第 307-310 行的错误处理中增加条件判断:

**实现逻辑**:
```go
if err != nil {
    // 检查是否是大文件超限错误
    if IsExceededLimitError(out) {
        // 可恢复的错误 -> WARN
        logger.Logger.Warnf("%s 裸仓push失败(文件超过大小限制，系统将自动处理): %s", repoPath, err)
    } else {
        // 真正的错误 -> ERROR
        logger.Logger.Errorf("%s 裸仓push失败: %s", repoPath, err)
    }
    return out, err
}
```

**关键点**:
- 利用现有的 `IsExceededLimitError(output)` 函数判断错误类型
- 只对 `exceeded limit` 错误降级为 WARN
- 保持其他所有错误的 ERROR 级别不变

### Rationale
- **精确的错误分类**:只对可自动恢复的大文件错误降级，真正的错误仍然醒目
- **语义准确性**:WARN 级别准确反映"需要注意但可自动处理"的情况
- **用户体验**:减少对可恢复错误的误解
- **运维友好**:ERROR 日志数量减少，但仍能识别真正的问题
- **最小改动**:只修改日志级别判断，不改变错误返回和处理流程
- **向后兼容**:不影响现有的错误处理流程和自动修复机制

### Alternatives Considered

1. **所有 push 失败都降级为 WARN**
   - 拒绝理由:会掩盖真正的错误（如网络问题、权限问题等）

2. **保持现状,不修改**
   - 拒绝理由:持续造成用户对可恢复错误的困惑

3. **完全移除大文件超限的日志**
   - 拒绝理由:丢失重要的调试信息，无法追踪自动修复流程

4. **改为 INFO 级别**
   - 拒绝理由:失败情况仍然需要被注意，INFO 级别过低

5. **增加更详细的日志说明**
   - **已采纳**:在 WARN 日志中增加"系统将自动处理"的说明

## Implementation Scope

### Files to Modify
- `pkg/git/git.go`: 修改 `Push()` 函数的错误日志逻辑（第 307-310 行）

### Testing Considerations
- 验证大文件超限错误输出 WARN 日志
- 验证其他类型错误仍输出 ERROR 日志
- 确认不影响错误处理和返回值
- 测试大文件场景下的完整迁移流程（包括自动修复）

### Related Code
**相关但不需要修改的代码**:
- `pkg/git/git.go:509`: `IsExceededLimitError()` 函数用于判断错误类型，已存在
- `pkg/migrate/migrate.go:638-644`: 自动修复大文件的逻辑，无需修改
- `pkg/migrate/migrate.go:639`: 这里已经是 WARN 级别，无需修改

## Dependencies
无外部依赖变更。

## Rollout Plan
1. 修改代码
2. 执行测试验证
3. 合并到主分支
4. 在发布说明中提及日志级别变更

## Success Criteria
- 大文件超限导致的 push 失败输出 WARN 日志
- WARN 日志包含"系统将自动处理"的说明
- 其他原因的 push 失败仍输出 ERROR 日志  
- 错误处理逻辑和返回值保持不变
- 大文件自动修复功能正常工作
- 所有单元测试通过
