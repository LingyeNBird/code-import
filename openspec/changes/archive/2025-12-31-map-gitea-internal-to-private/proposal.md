# Map Gitea Internal Repositories to Private

## Change ID
`map-gitea-internal-to-private`

## Summary
当从 Gitea 迁移仓库时,如果源仓库的可见性为 "Internal"(内部仓库),则在 CNB 平台创建的仓库应设置为 "Private"(私有),而不是 "Public"(公开)。这确保了内部仓库的访问控制策略在迁移后得到保持。

## Problem Statement
### Current Behavior
目前代码在迁移 Gitea 仓库时,只检查 `repo.Private` 字段来决定 CNB 仓库的可见性:
- 如果 `repo.Private == true`,则创建为 private
- 如果 `repo.Private == false`,则创建为 public

### Gap
Gitea 支持三种仓库可见性级别:
1. **Public**: 公开仓库,任何人可访问
2. **Private**: 私有仓库,只有授权用户可访问  
3. **Internal**: 内部仓库,只有已登录的平台用户可访问

当前实现忽略了 `repo.Internal` 字段(Gitea API 返回的 `internal` 属性)。这导致 Gitea 内部仓库被迁移为 CNB 公开仓库,违背了原有的访问控制策略。

### Impact
- **安全风险**: 原本限制为内部访问的仓库变为公开可访问
- **合规问题**: 不符合组织的访问控制策略
- **数据泄露风险**: 敏感代码可能被未授权用户访问

## Proposed Solution
### Approach
修改 Gitea VCS 实现,在判断仓库可见性时同时考虑 `Private` 和 `Internal` 字段:
- 如果 `Private == true` OR `Internal == true`,则创建 CNB private 仓库
- 否则创建 CNB public 仓库

### Rationale
- **最小权限原则**: 将 internal 映射为 private 比映射为 public 更安全
- **语义一致性**: Internal 表示"非公开",与 Private 语义一致
- **向后兼容**: 不影响现有的 private/public 仓库迁移逻辑
- **简单实现**: 只需修改一处判断逻辑,无需新增配置项

### Alternatives Considered
1. **新增配置项控制映射策略**  
   - 拒绝理由: 增加复杂度,大多数用户期望的默认行为就是将 internal 视为 private

2. **将 internal 映射为新的 CNB 可见性级别**  
   - 拒绝理由: CNB 平台可能不支持 internal 级别,且语义可能不一致

3. **保持现状,依赖用户迁移后手动调整**  
   - 拒绝理由: 存在安全风险,用户体验差

## Implementation Scope
### Files to Modify
- `pkg/vcs/gitea.go`: 修改 `GetRepoPrivate()` 方法
- `pkg/api/gitea/api.go`: 已包含 `Internal` 字段定义,无需修改

### Testing Considerations
- 验证 private 仓库仍然创建为 private
- 验证 public 仓库仍然创建为 public
- 验证 internal 仓库创建为 private
- 验证混合场景(批量迁移包含不同可见性的仓库)

## Dependencies
无外部依赖变更。

## Rollout Plan
1. 修改代码并添加单元测试
2. 在测试环境验证迁移行为
3. 合并到主分支,随下一版本发布
4. 在发布说明中提及此行为变更

## Success Criteria
- Gitea internal 仓库迁移后在 CNB 为 private
- 现有 private/public 仓库迁移逻辑不受影响
- 无性能回归
- 代码通过所有单元测试
