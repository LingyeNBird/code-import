# gitea-visibility-mapping Specification

## Purpose
定义 Gitea 仓库迁移时的可见性映射规则,确保源仓库的访问控制策略在目标 CNB 平台得到正确保持,特别是处理 Gitea 特有的 "Internal"(内部)可见性级别。

## Requirements

### Requirement: Map Gitea Internal to CNB Private
当从 Gitea 迁移仓库到 CNB 平台时,系统 SHALL 将 Gitea 的 "Internal"(内部)可见性映射为 CNB 的 "Private"(私有)可见性,确保非公开仓库的访问控制不被削弱。

#### Scenario: 迁移 Gitea Internal 仓库
- **GIVEN** 源平台为 Gitea
- **AND** 存在一个可见性为 Internal 的仓库(API 返回 `internal: true`)
- **WHEN** 系统执行仓库迁移
- **THEN** 系统调用 CNB API 创建仓库时设置 `visibility: "private"`
- **AND** 迁移日志记录该仓库被创建为 private

#### Scenario: 迁移 Gitea Private 仓库
- **GIVEN** 源平台为 Gitea  
- **AND** 存在一个可见性为 Private 的仓库(API 返回 `private: true`)
- **WHEN** 系统执行仓库迁移
- **THEN** 系统调用 CNB API 创建仓库时设置 `visibility: "private"`
- **AND** 迁移日志记录该仓库被创建为 private

#### Scenario: 迁移 Gitea Public 仓库
- **GIVEN** 源平台为 Gitea
- **AND** 存在一个可见性为 Public 的仓库(API 返回 `private: false, internal: false`)
- **WHEN** 系统执行仓库迁移
- **THEN** 系统调用 CNB API 创建仓库时设置 `visibility: "public"`
- **AND** 迁移日志记录该仓库被创建为 public

#### Scenario: 迁移混合可见性的仓库批次
- **GIVEN** 源平台为 Gitea
- **AND** 待迁移仓库列表包含:
  - 2 个 Internal 仓库
  - 1 个 Private 仓库
  - 2 个 Public 仓库
- **WHEN** 系统执行批量迁移
- **THEN** 3 个仓库(2 internal + 1 private)在 CNB 创建为 private
- **AND** 2 个仓库在 CNB 创建为 public
- **AND** 迁移日志正确记录每个仓库的可见性

### Requirement: Preserve Existing Visibility Logic
系统 SHALL 保持对其他 VCS 平台(GitHub, GitLab, Gitee 等)的可见性判断逻辑不变,确保向后兼容性。

#### Scenario: 迁移 GitHub 仓库不受影响
- **GIVEN** 源平台为 GitHub
- **AND** 存在一个 private 仓库
- **WHEN** 系统执行仓库迁移
- **THEN** 系统仍然使用 `repo.Private` 字段判断可见性
- **AND** 迁移行为与修改前完全一致

#### Scenario: 迁移 GitLab 仓库不受影响
- **GIVEN** 源平台为 GitLab
- **AND** 存在一个 private 仓库
- **WHEN** 系统执行仓库迁移
- **THEN** 系统仍然使用 GitLab 的可见性字段判断
- **AND** 迁移行为与修改前完全一致

### Requirement: Gitea API Integration
系统 SHALL 正确解析 Gitea API 返回的仓库信息,包括 `private` 和 `internal` 字段,确保可见性数据准确传递到迁移逻辑。

#### Scenario: 正确解析 Gitea API 响应
- **GIVEN** 系统调用 Gitea API `/api/v1/user/repos` 获取仓库列表
- **WHEN** API 返回的仓库 JSON 包含 `"internal": true`
- **THEN** 系统将 `internal` 字段映射到 VCS 对象的 `Internal` 属性
- **AND** 该属性可被后续迁移逻辑访问

#### Scenario: 处理 API 字段缺失情况
- **GIVEN** 系统调用 Gitea API 获取仓库信息
- **WHEN** API 响应中缺少 `internal` 字段(如旧版 Gitea)
- **THEN** 系统将 `Internal` 属性默认为 `false`
- **AND** 仍然根据 `Private` 字段正确判断可见性
- **AND** 迁移不会因字段缺失而失败

## ADDED Requirements

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

## Cross-References
- **Related Change**: N/A (独立变更)
- **Depends On**: N/A
- **Blocks**: N/A
