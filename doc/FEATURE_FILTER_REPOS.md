# 按指定仓库迁移功能实现总结

## 功能概述

实现了支持工蜂、GitLab、GitHub、Gitee 等平台按指定仓库列表进行精确迁移的功能。当 `PLUGIN_SOURCE_REPO` 环境变量（或配置文件中的 `source.repo`）不为空时，只迁移指定的仓库，其他仓库将被跳过。

## 实现内容

### 1. 代码修改

#### 1.1 配置层修改 (`pkg/config/config.go`)

**文件**: `pkg/config/config.go`  
**修改位置**: 第155-161行

**修改内容**:
```go
// 调整 source.repo 的校验逻辑
// repo 模式下，只有 common 和 local 平台需要强制要求 source.repo
// 其他平台（gitlab、github、gitee、gongfeng等）可以通过 source.repo 过滤仓库
if config.Migrate.Type == "repo" && (platform == "common" || platform == "local") && (len(config.Source.Repo) == 0 || config.Source.Repo[0] == "") {
    return fmt.Errorf("when migrate.type is repo and platform is common or local, source.repo is required")
}
```

**变更说明**:
- 原逻辑: 所有平台在 repo 模式下都强制要求 `source.repo`
- 新逻辑: 只有 common 和 local 平台强制要求，其他平台可选
- 目的: 允许其他平台使用 `source.repo` 作为过滤条件

#### 1.2 迁移逻辑层修改 (`pkg/migrate/migrate.go`)

**文件**: `pkg/migrate/migrate.go`

**新增函数**: `filterReposByConfigList` (第133-179行)

```go
// filterReposByConfigList 根据 source.repo 配置过滤仓库列表
// 当 source.repo 不为空时，只保留在配置列表中的仓库
func filterReposByConfigList(depotList []vcs.VCS) []vcs.VCS {
    configRepos := config.Cfg.GetStringSlice("source.repo")
    
    // 如果配置为空，返回完整列表
    if len(configRepos) == 0 || (len(configRepos) == 1 && configRepos[0] == "") {
        return depotList
    }

    // 构建仓库路径映射表，提高查找效率
    repoMap := make(map[string]bool, len(configRepos))
    for _, repoPath := range configRepos {
        trimmedPath := strings.TrimSpace(repoPath)
        if trimmedPath != "" {
            repoMap[trimmedPath] = true
        }
    }

    // 如果配置的仓库列表为空（全是空字符串），返回完整列表
    if len(repoMap) == 0 {
        return depotList
    }

    logger.Logger.Infof("检测到 source.repo 配置，将只迁移指定的 %d 个仓库", len(repoMap))
    
    // 过滤仓库列表
    filteredDepotList := make([]vcs.VCS, 0, len(repoMap))
    matchedRepos := make(map[string]bool)
    
    for _, depot := range depotList {
        repoPath := depot.GetRepoPath()
        if repoMap[repoPath] {
            filteredDepotList = append(filteredDepotList, depot)
            matchedRepos[repoPath] = true
            logger.Logger.Infof("匹配到配置仓库: %s", repoPath)
        }
    }

    // 检查是否有配置的仓库未找到
    for repoPath := range repoMap {
        if !matchedRepos[repoPath] {
            logger.Logger.Warnf("配置的仓库 %s 在源平台未找到", repoPath)
        }
    }

    logger.Logger.Infof("根据 source.repo 配置过滤后，待迁移仓库数: %d", len(filteredDepotList))
    return filteredDepotList
}
```

**修改位置**: `Run` 函数，第184-196行

```go
// 获取并过滤仓库列表
depotList := sourceVcsList

// 先根据 source.repo 配置过滤（如果配置了的话）
depotList = filterReposByConfigList(depotList)

// 再根据 repo-path.txt 文件过滤（如果启用了仓库选择功能）
depotList, err = filterReposBySelection(depotList)
if err != nil {
    logger.Logger.Errorf("%s", err)
    return 1
}

logger.Logger.Infof("待迁移仓库总数%d", len(depotList))
```

### 2. 测试用例

**文件**: `pkg/migrate/migrate_test.go` (新建)

**测试覆盖**:
- ✅ 空配置测试
- ✅ 单个仓库过滤测试
- ✅ 多个仓库过滤测试
- ✅ 带空格的配置测试
- ✅ 空字符串配置测试
- ✅ 无匹配仓库测试
- ✅ GitLab 格式测试
- ✅ GitHub 格式测试
- ✅ 工蜂格式测试
- ✅ 混合格式测试
- ✅ 精确匹配测试
- ✅ 大小写敏感测试
- ✅ 空仓库列表测试
- ✅ 重复配置测试

**测试结果**: 所有 14 个测试用例全部通过 ✅

```
PASS
ok  	ccrctl/pkg/migrate	0.233s
```

### 3. 文档

#### 3.1 使用说明文档
**文件**: `doc/filter-repos-by-config.md`

内容包括:
- 功能概述
- 使用场景
- 配置方式（环境变量和配置文件）
- 不同平台的仓库路径格式
- 使用示例
- 行为说明
- 日志输出示例
- 注意事项
- 常见问题

#### 3.2 示例脚本
**文件**: `doc/filter-repos-example.sh`

包含 5 个平台的实际使用示例:
- GitLab 单个仓库迁移
- GitLab 多个仓库迁移
- GitHub 仓库迁移
- 工蜂仓库迁移
- Gitee 仓库迁移

## 功能特性

### 核心特性

1. **平台支持广泛**
   - ✅ GitLab
   - ✅ GitHub
   - ✅ Gitee
   - ✅ 工蜂（Gongfeng）
   - ✅ CODING

2. **灵活的配置方式**
   - 环境变量: `PLUGIN_SOURCE_REPO`
   - 配置文件: `source.repo`
   - 支持单个或多个仓库（逗号分隔）

3. **智能过滤**
   - 精确匹配（区分大小写）
   - 自动去除前后空格
   - 未匹配仓库警告提示
   - 空配置自动跳过过滤

4. **友好的日志**
   - 匹配成功提示
   - 未找到仓库警告
   - 过滤后数量统计

5. **向后兼容**
   - 不影响现有功能
   - 可与 `allow_select_repos` 配合使用
   - 支持所有迁移类型（team/project/repo）

### 执行流程

```
1. 获取源平台仓库列表
   ↓
2. filterReposByConfigList (source.repo 配置过滤)
   ↓
3. filterReposBySelection (repo-path.txt 文件过滤)
   ↓
4. 执行迁移
```

## 使用示例

### 环境变量方式

```bash
# 迁移单个仓库
export PLUGIN_SOURCE_REPO="owner/repo1"

# 迁移多个仓库
export PLUGIN_SOURCE_REPO="owner/repo1,group/subgroup/repo2,org/team/project/repo3"

./ccrctl migrate
```

### 配置文件方式

```yaml
source:
  platform: gitlab
  url: https://gitlab.example.com
  token: glpat-xxxxxxxxxxxx
  repo:
    - group1/subgroup1/repo1
    - group1/subgroup2/repo2
    - group2/repo3
```

## 代码规范检查

### Linter 检查结果

✅ 没有新增 linter 错误  
✅ 遵循 Go 编码规范  
✅ 符合项目代码风格要求

### 预先存在的提示（非本次修改引入）

- INFO: `codingTokenLength` 常量未使用
- INFO: `rebaseBranchesMap` 变量未使用
- HINT: `io/ioutil` 已废弃（建议使用 `os` 包）

## 编译验证

```bash
cd /Users/vincentliu/git_dir/code-import
go build -o ccrctl main.go
```

✅ 编译成功，无错误

## 适用场景

1. **精确迁移**: 只需要迁移特定的几个仓库
2. **批量选择**: 从大量仓库中选择部分进行迁移
3. **分批迁移**: 将迁移任务分成多批次执行
4. **测试验证**: 先迁移少量仓库进行测试

## 兼容性说明

- ✅ 与现有 `migrate.type` (team/project/repo) 兼容
- ✅ 与 `allow_select_repos` 功能兼容
- ✅ 不影响其他迁移配置项
- ✅ 向后兼容，不配置则保持原有行为

## 性能优化

1. 使用 map 进行仓库路径查找，时间复杂度 O(1)
2. 预先构建映射表，避免重复查找
3. 及时提供过滤进度反馈

## 总结

本次实现完成了按指定仓库迁移的功能，主要亮点：

1. **功能完整**: 支持所有主流 Git 平台
2. **代码质量**: 遵循 Go 规范，通过所有测试
3. **文档齐全**: 提供详细的使用文档和示例
4. **向后兼容**: 不影响现有功能
5. **用户友好**: 清晰的日志提示和错误处理

该功能已完全实现并通过测试，可以投入使用。
