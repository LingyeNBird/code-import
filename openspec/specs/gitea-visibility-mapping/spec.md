# gitea-visibility-mapping Specification

## Purpose
定义 Gitea 仓库迁移时的可见性映射规则,确保源仓库的访问控制策略在目标 CNB 平台得到正确保持,特别是处理 Gitea 特有的 "Internal"(内部)可见性级别。

## Requirements
### Requirement: Repository Visibility Determination
对于 Gitea 平台,系统 SHALL 综合考虑 `repo.Private` 和 `repo.Internal` 字段来判断仓库可见性:
- 如果 `Private == true` OR `Internal == true`,则目标仓库 SHALL 创建为 private
- 否则目标仓库 SHALL 创建为 public

**实现说明**: 
- 在 `GiteaVcs` 结构体中添加 `Internal bool` 字段
- 修改 `GetRepoPrivate()` 方法返回 `c.Private || c.Internal`
- 在 `GiteaCovertToVcs()` 函数中传递 `Internal: repo.Internal`

**设计理由**: 支持 Gitea 的 Internal 可见性级别,确保非公开仓库的访问控制不被削弱。Gitea 的 Internal 仓库表示只有已登录用户可访问,语义上应映射为 CNB 的 Private(私有)而非 Public(公开)。

#### Scenario: Gitea 可见性判断使用新逻辑
- **GIVEN** 源平台为 Gitea
- **AND** 系统调用 `depot.GetRepoPrivate()` 方法
- **WHEN** 仓库的 `Internal` 为 true 且 `Private` 为 false
- **THEN** 方法返回 `true`(表示应创建 private 仓库)

