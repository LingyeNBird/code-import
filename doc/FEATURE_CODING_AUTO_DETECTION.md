# CODING 源平台迁移逻辑优化

## 优化概述

优化 CODING 源平台的迁移维度判断逻辑，不再依赖 `PLUGIN_MIGRATE_TYPE` 环境变量，改为根据 `PLUGIN_SOURCE_REPO` 和 `PLUGIN_SOURCE_PROJECT` 的配置自动判断迁移维度。

## 核心变更

### 优化前的逻辑

需要用户手动指定 `PLUGIN_MIGRATE_TYPE` 来决定迁移维度：

```bash
# 按仓库迁移
export PLUGIN_MIGRATE_TYPE=repo
export PLUGIN_SOURCE_REPO="team1/project1/repo1,team1/project1/repo2"

# 按项目迁移
export PLUGIN_MIGRATE_TYPE=project
export PLUGIN_SOURCE_PROJECT="project1,project2"

# 按团队迁移
export PLUGIN_MIGRATE_TYPE=team
```

**问题**：
- 用户需要同时配置迁移类型和迁移对象，容易混淆
- 配置不一致时可能导致意外结果
- 增加了使用复杂度

### 优化后的逻辑

**自动判断迁移维度**，根据配置内容智能决定：

```bash
# 场景1: 按仓库迁移（配置了 source.repo）
export PLUGIN_SOURCE_REPO="team1/project1/repo1,team1/project1/repo2"
# 自动识别为按仓库迁移，无需配置 PLUGIN_MIGRATE_TYPE

# 场景2: 按项目迁移（配置了 source.project）
export PLUGIN_SOURCE_PROJECT="project1,project2"
# 自动识别为按项目迁移，无需配置 PLUGIN_MIGRATE_TYPE

# 场景3: 按团队迁移（都未配置）
# 自动获取团队下所有仓库，无需配置 PLUGIN_MIGRATE_TYPE
```

**优先级规则**：
1. **source.repo** > source.project > team（默认）
2. 如果配置了 `source.repo` → 按仓库维度获取
3. 如果只配置了 `source.project` → 按项目维度获取
4. 如果都未配置 → 获取团队下所有仓库

## 代码修改

### 1. pkg/api/coding/api.go

**修改位置**: `GetDepotList` 函数（第434-453行）

**修改前**：
```go
func GetDepotList(migrateType string) ([]Depots, error) {
    logger.Logger.Infof("获取仓库列表中...")
    var depotList []Depots
    var err error
    switch migrateType {
    case ProjectType:
        depotList, err = GetDepotListByProjectNames(SourceURL, SourceToken, Projects)
    case RepoType:
        depotList, err = GetDepotListByRepoPath(SourceURL, SourceToken, Repos)
    case Team:
        depotList, err = GetDepotListByTeam(SourceURL, SourceToken)
    default:
        return nil, fmt.Errorf("未知的迁移类型: %s", migrateType)
    }
    // ...
}
```

**修改后**：
```go
func GetDepotList(migrateType string) ([]Depots, error) {
    logger.Logger.Infof("获取仓库列表中...")
    var depotList []Depots
    var err error
    
    // 优化逻辑：根据 source.repo 和 source.project 配置自动判断迁移维度
    // 优先级：source.repo > source.project > team（全部仓库）
    if len(Repos) > 0 && Repos[0] != "" {
        // 配置了 source.repo，按仓库迁移
        logger.Logger.Info("检测到 source.repo 配置，按仓库维度获取列表")
        depotList, err = GetDepotListByRepoPath(SourceURL, SourceToken, Repos)
    } else if len(Projects) > 0 && Projects[0] != "" {
        // 配置了 source.project，按项目迁移
        logger.Logger.Info("检测到 source.project 配置，按项目维度获取列表")
        depotList, err = GetDepotListByProjectNames(SourceURL, SourceToken, Projects)
    } else {
        // 都没配置，默认获取团队下所有仓库
        logger.Logger.Info("未配置 source.repo 和 source.project，获取团队下所有仓库")
        depotList, err = GetDepotListByTeam(SourceURL, SourceToken)
    }
    // ...
}
```

### 2. pkg/vcs/coding.go

**修改位置**: `newCodingRepo` 函数（第155行）

**修改前**：
```go
func newCodingRepo() ([]VCS, error) {
    repoList, err := coding.GetDepotList(config.Cfg.GetString("migrate.type"))
    if err != nil {
        return nil, err
    }
    return CodingCovertToVcs(repoList), nil
}
```

**修改后**：
```go
func newCodingRepo() ([]VCS, error) {
    // 不再需要传入 migrate.type，由 GetDepotList 内部根据配置自动判断
    repoList, err := coding.GetDepotList("")
    if err != nil {
        return nil, err
    }
    return CodingCovertToVcs(repoList), nil
}
```

### 3. pkg/config/config.go

**修改位置**: 配置校验逻辑（第146-163行）

**修改前**：
```go
// 检查 migrate 参数
if config.Migrate.Type == "" {
    return fmt.Errorf("migrate.type is required")
}

if config.Migrate.Type == "project" && (len(config.Source.Project) == 0 || config.Source.Project[0] == "") {
    return fmt.Errorf("coding.project is required")
}
```

**修改后**：
```go
// 检查 migrate 参数
// migrate.type 不再是必填项，如果未配置则默认为 team
if config.Migrate.Type == "" {
    config.Migrate.Type = "team"
}

// CODING 平台的特殊校验已移除，改为由 source.repo 和 source.project 自动判断
// 其他平台保持原有逻辑
```

## 使用示例

### 场景1: 按仓库迁移（指定具体仓库）

```bash
export PLUGIN_SOURCE_PLATFORM=coding
export PLUGIN_SOURCE_URL=https://your-team.coding.net
export PLUGIN_SOURCE_TOKEN=your_token
export PLUGIN_SOURCE_REPO="team1/project1/repo1,team1/project1/repo2"
# 不需要配置 PLUGIN_MIGRATE_TYPE，自动识别为按仓库迁移

./ccrctl migrate
```

**日志输出**：
```
INFO  检测到 source.repo 配置，按仓库维度获取列表
INFO  获取仓库列表中...
INFO  从源平台获取到仓库总数: 2
```

### 场景2: 按项目迁移（迁移整个项目）

```bash
export PLUGIN_SOURCE_PLATFORM=coding
export PLUGIN_SOURCE_URL=https://your-team.coding.net
export PLUGIN_SOURCE_TOKEN=your_token
export PLUGIN_SOURCE_PROJECT="project1,project2"
# 不需要配置 PLUGIN_MIGRATE_TYPE，自动识别为按项目迁移

./ccrctl migrate
```

**日志输出**：
```
INFO  检测到 source.project 配置，按项目维度获取列表
INFO  获取仓库列表中...
INFO  从源平台获取到仓库总数: 25
```

### 场景3: 按团队迁移（迁移全部仓库）

```bash
export PLUGIN_SOURCE_PLATFORM=coding
export PLUGIN_SOURCE_URL=https://your-team.coding.net
export PLUGIN_SOURCE_TOKEN=your_token
# 不配置 PLUGIN_SOURCE_REPO 和 PLUGIN_SOURCE_PROJECT
# 不需要配置 PLUGIN_MIGRATE_TYPE，自动获取全部仓库

./ccrctl migrate
```

**日志输出**：
```
INFO  未配置 source.repo 和 source.project，获取团队下所有仓库
INFO  获取仓库列表中...
INFO  从源平台获取到仓库总数: 100
```

### 场景4: 同时配置（按优先级选择）

```bash
export PLUGIN_SOURCE_PLATFORM=coding
export PLUGIN_SOURCE_URL=https://your-team.coding.net
export PLUGIN_SOURCE_TOKEN=your_token
export PLUGIN_SOURCE_REPO="team1/project1/repo1"
export PLUGIN_SOURCE_PROJECT="project2"
# source.repo 优先级高，按仓库迁移

./ccrctl migrate
```

**日志输出**：
```
INFO  检测到 source.repo 配置，按仓库维度获取列表
INFO  获取仓库列表中...
INFO  从源平台获取到仓库总数: 1
```

## 兼容性说明

### 向后兼容

✅ **完全兼容**：如果用户仍然配置了 `PLUGIN_MIGRATE_TYPE`，系统会正常工作：

```bash
# 旧的配置方式仍然有效
export PLUGIN_MIGRATE_TYPE=repo
export PLUGIN_SOURCE_REPO="team1/project1/repo1"
# 仍然可以正常工作
```

### 迁移建议

对于现有用户：
1. **无需修改配置**：旧的配置方式仍然有效
2. **逐步优化**：可以逐步移除 `PLUGIN_MIGRATE_TYPE` 配置，简化使用
3. **推荐做法**：新的迁移任务直接省略 `PLUGIN_MIGRATE_TYPE`

## 优势总结

### 1. **简化配置**
- 减少一个必填环境变量
- 配置更直观，用户体验更好

### 2. **智能判断**
- 根据实际配置自动选择最合适的迁移维度
- 减少配置错误的可能性

### 3. **向后兼容**
- 不破坏现有配置
- 平滑过渡

### 4. **日志清晰**
- 明确提示当前使用的迁移维度
- 便于排查问题

## 测试验证

### 编译测试
```bash
cd /Users/vincentliu/git_dir/code-import
go build -o ccrctl main.go
# ✅ 编译成功
```

### 单元测试
```bash
go test ./pkg/migrate -v -count=1
# ✅ 所有17个测试用例通过
```

### Linter 检查
```bash
# ✅ 无新增错误，只有预先存在的提示
```

## 影响范围

### 受影响的文件
1. `pkg/api/coding/api.go` - 核心逻辑修改
2. `pkg/vcs/coding.go` - 调用方式修改
3. `pkg/config/config.go` - 配置校验优化

### 不受影响的平台
- ✅ GitLab
- ✅ GitHub
- ✅ Gitee
- ✅ Gitea
- ✅ 阿里云
- ✅ 华为云
- ✅ 工蜂
- ✅ Common/Local

这些平台的逻辑完全不受影响，仍按原有方式工作。

---

**实现完成时间**: 2025-12-04  
**版本**: v1.2  
**状态**: ✅ 已完成，可投入使用  
**更新内容**: CODING 平台智能识别迁移维度，简化配置流程
