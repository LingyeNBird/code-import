# 设计文档:CODING 项目元数据缓存优化

## Context

CNB Code Import 工具从 CODING 平台迁移仓库到 CNB。在迁移过程中,每个仓库都会调用 `CodingVcs.GetSubGroup()` 来获取项目元数据(DisplayName、Description),进而调用 `coding.GetProjectByName()` API。

**现状**:
- 位置:`pkg/vcs/coding.go:49`
- 调用链:`migrateDo() -> depot.GetSubGroup() -> coding.GetProjectByName(projectName)`
- 问题:同一个项目的多个仓库会重复调用 GetProjectByName API

**示例场景**:
```
项目 "backend-services" 包含 100 个仓库
当前行为:GetProjectByName("backend-services") 被调用 100 次
期望行为:GetProjectByName("backend-services") 被调用 1 次,其余 99 次从缓存读取
```

**约束**:
- 迁移时最多 10 个 goroutine 并发执行(`MaxConcurrency = 10`)
- 缓存必须线程安全
- API 失败时必须保持现有重试机制和降级行为
- 不能改变 VCS 接口或调用方代码

## Goals / Non-Goals

**Goals**:
- 减少 CODING GetProjectByName API 调用次数至每个项目最多 1 次
- 缓存对调用方完全透明,不改变 VCS 接口
- 保持现有的错误处理和重试机制
- 确保并发安全(最多 10 个 goroutine)

**Non-Goals**:
- 不实现缓存过期机制(迁移过程通常在数小时内完成,项目元数据不会变化)
- 不实现缓存持久化(每次运行工具时重新初始化缓存)
- 不优化其他 CODING API 调用(范围限制为 GetProjectByName)
- 不改变 GetSubGroup() 的重试逻辑

## Decisions

### 决策 1:缓存位置和数据结构

**选择**:在 `pkg/api/coding/api.go` 中添加包级全局变量 `projectCache sync.Map`

**理由**:
- `sync.Map` 是 Go 标准库提供的并发安全 map,适合读多写少的场景
- 包级变量在 `GetDepotList()` 首次调用时初始化,在程序生命周期内有效
- API 层缓存对 VCS 层完全透明,不需要修改 `CodingVcs` 结构体

**替代方案**:
1. 在 `CodingVcs` 结构体中添加缓存字段
   - ❌ 每个 VCS 实例都有独立缓存,无法共享
   - ❌ 需要修改 `CodingCovertToVcs()` 函数传递缓存引用
2. 使用普通 `map[string]Project` + `sync.RWMutex`
   - ❌ 需要手动管理锁,增加复杂度
   - ✅ `sync.Map` 在读多写少场景性能更优

### 决策 2:缓存键和值

**缓存键**:项目名称(string)
**缓存值**:`coding.Project` 结构体

```go
type Project struct {
    Name        string `json:"Name"`
    Id          int    `json:"Id"`
    Type        int    `json:"Type"`
    DisplayName string `json:"DisplayName"`
    Icon        string `json:"Icon"`
    Description string `json:"Description"`
    // ... 其他字段
}
```

**理由**:
- 项目名称是 GetProjectByName() 的唯一输入参数,是自然的缓存键
- 缓存完整的 Project 结构体而非仅 DisplayName/Description,以备未来扩展
- CODING 项目名称在团队内唯一,不会发生键冲突

### 决策 3:缓存填充策略

**选择**:懒加载(Lazy Loading)

在 `GetProjectByName()` 函数中:
1. 检查缓存是否命中
2. 命中:直接返回缓存值
3. 未命中:调用 API 获取数据,存入缓存后返回

**理由**:
- 简单直接,不需要预加载所有项目元数据
- 仅缓存实际需要的项目,节省内存
- 对调用方完全透明

**替代方案**:
- 在 `GetDepotList()` 中预加载所有项目元数据
  - ❌ 需要额外 API 调用获取项目列表
  - ❌ 如果只迁移少数仓库,会缓存大量无用数据
  - ❌ 增加初始化时间

### 决策 4:缓存生命周期

**选择**:程序生命周期内有效,不实现过期机制

**理由**:
- 迁移过程通常在数小时内完成,项目元数据在此期间不会变化
- 不需要缓存持久化,每次运行工具时重新初始化
- 简化实现,避免引入定时器等复杂逻辑

**风险与缓解**:
- 风险:如果迁移过程中项目元数据被修改,缓存会过期
- 缓解:极低风险,实际场景中迁移期间项目元数据通常不会变化;如需强制刷新,重启工具即可

### 决策 5:错误处理和降级

**选择**:保持现有的重试机制和降级逻辑

- `GetProjectByName()` 失败时返回错误(不缓存错误结果)
- `GetSubGroup()` 中的重试机制保持不变(1秒、5秒、10秒间隔)
- 所有重试失败后返回最小 SubGroup(仅包含 Name 字段)

**理由**:
- 缓存仅优化成功路径,不改变失败处理逻辑
- 保持与现有错误处理行为的一致性
- 避免缓存污染(不缓存错误或部分数据)

## Implementation Details

### 代码修改点

#### 1. `pkg/api/coding/api.go`

添加包级缓存变量:
```go
var (
    SourceToken = config.Cfg.GetString("source.token")
    SourceURL   = config.Cfg.GetString("source.url")
    Projects    = config.Cfg.GetStringSlice("source.project")
    Repos       = config.Cfg.GetStringSlice("source.repo")
    c           = http_client.NewClient(SourceURL)
    c2          = http_client.NewCodingClient()
    
    // 新增:项目元数据缓存
    projectCache sync.Map // map[string]Project
)
```

修改 `GetProjectByName()` 函数(第250-271行):
```go
func GetProjectByName(url, token, projectName string) (project Project, err error) {
    // 1. 先检查缓存
    if cachedProject, ok := projectCache.Load(projectName); ok {
        logger.Logger.Debugf("从缓存获取项目 %s 信息", projectName)
        return cachedProject.(Project), nil
    }
    
    // 2. 缓存未命中,调用 API
    body := &DescribeProjectByNameRequest{
        Action:      "DescribeProjectByName",
        ProjectName: projectName,
    }
    resp, err := c.Request("POST", endpoint, token, body)
    if err != nil {
        return project, err
    }
    err = checkResponse(resp)
    if err != nil {
        return project, err
    }
    var projectInfo DescribeProjectByNameResponse
    err = c.Unmarshal(resp, &projectInfo)
    if err != nil {
        return project, err
    }
    logger.Logger.Debugf("%s项目ID: %d", projectName, projectInfo.Response.Project.Id)
    
    // 3. 存入缓存
    project = projectInfo.Response.Project
    projectCache.Store(projectName, project)
    
    return project, nil
}
```

#### 2. `pkg/vcs/coding.go`

**无需修改**,`GetSubGroup()` 方法(第40-88行)继续调用 `coding.GetProjectByName()`,缓存优化对其完全透明。

### 并发安全性

`sync.Map` 提供以下并发安全保证:
- `Load()`:多个 goroutine 可以并发读取,无锁
- `Store()`:写入时使用内部锁,确保原子性
- 读多写少场景性能优于 `map + sync.RWMutex`

**测试场景**:
- 10 个 goroutine 并发调用 `GetProjectByName("project-A")`
- 第一个调用触发 API 请求并缓存结果
- 后续 9 个调用从缓存读取(可能部分调用在缓存前发起,会导致少量重复 API 调用,但数量远小于无缓存情况)

### 缓存命中率分析

**场景 1:单项目 100 个仓库**
- 第 1 个仓库:缓存未命中,调用 API
- 第 2-100 个仓库:缓存命中
- 命中率:99%

**场景 2:10 个项目,每个 10 个仓库**
- 每个项目的第 1 个仓库:缓存未命中(共 10 次 API 调用)
- 其余 90 个仓库:缓存命中
- 命中率:90%

**场景 3:并发竞争(10 个 goroutine 同时处理同一项目的仓库)**
- 最坏情况:10 个仓库同时调用 GetProjectByName,缓存尚未建立
- 可能触发 2-3 次 API 调用(取决于 API 响应时间和调度顺序)
- 后续仓库仍从缓存读取,整体命中率仍远高于无缓存情况

## Risks / Trade-offs

### 风险 1:缓存不一致

**描述**:迁移过程中项目元数据被修改,缓存值过期

**影响**:低
- 实际场景中迁移期间项目元数据通常不会变化
- 即使变化,仅影响 CNB 子组织的 Description 和 Remark 字段(非核心功能)

**缓解**:
- 文档说明:迁移前确保项目元数据稳定
- 如需强制刷新,重启工具即可

### 风险 2:并发竞争导致少量重复 API 调用

**描述**:10 个 goroutine 同时处理同一项目的仓库,可能在缓存建立前触发 2-3 次 API 调用

**影响**:极低
- 仅在迁移开始的短时间窗口内发生
- 即使发生,API 调用次数仍远小于无缓存情况(3 次 vs 100 次)

**缓解**:
- 可接受的 trade-off,避免引入额外的锁或初始化逻辑
- `sync.Map` 的内部实现已尽量减少并发写入时的竞争

### 权衡:内存占用 vs 性能

**内存占用**:
- 每个 Project 结构体约 200 字节
- 100 个项目:约 20 KB
- 1000 个项目:约 200 KB

**结论**:内存占用可忽略不计,性能收益远大于内存成本

## Migration Plan

### 部署步骤

1. 合并代码到主分支
2. 构建新版本 Docker 镜像
3. 用户更新到最新版本

### 回滚计划

如果缓存导致问题(虽然概率极低):
1. 回滚到上一版本
2. 缓存逻辑完全在 `GetProjectByName()` 函数内部,回滚后自动恢复原行为

### 验证方法

**功能验证**:
1. 迁移包含 50+ 仓库的 CODING 项目
2. 检查日志,确认 "从缓存获取项目 xxx 信息" 出现次数符合预期
3. 验证 CNB 子组织的 Description 和 Remark 字段正确映射

**性能验证**:
1. 对比优化前后迁移耗时(100 个仓库的项目)
2. 预期节省时间:约 (99 次 × 1秒) + (失败重试时间) = 100+ 秒

**压测验证**:
1. 模拟 10 个 goroutine 并发迁移同一项目的仓库
2. 验证缓存并发安全性
3. 确认 API 调用次数符合预期(≤ 10 次)

## Open Questions

无
