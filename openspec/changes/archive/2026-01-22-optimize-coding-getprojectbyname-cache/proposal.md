# Change: 优化 CODING 平台迁移时的 GetProjectByName API 调用次数

## Why

当前在批量迁移 CODING 平台仓库时,每个仓库都会调用一次 `coding.GetProjectByName()` API 来获取项目元数据(DisplayName、Description)用于映射到 CNB 子组织属性。当一个项目包含大量仓库时(如100个仓库),会对同一个项目名称重复调用100次 GetProjectByName API,导致:

1. **性能问题**:不必要的 API 调用增加迁移总时间,每次调用耗时约 0.5-2 秒
2. **API 压力**:大量重复请求增加 CODING 平台 API 压力,可能触发限流
3. **网络带宽浪费**:每次请求返回相同的项目元数据
4. **重试成本高**:当 API 失败时,每个仓库都会独立重试(1秒、5秒、10秒间隔),放大了失败成本

## What Changes

在 CODING VCS 实现中添加项目元数据缓存机制:

- 在 `pkg/api/coding/api.go` 中添加全局 `projectCache` (map[string]Project) 用于缓存项目元数据
- 修改 `GetProjectByName()` 函数,在调用 API 前先检查缓存,命中则直接返回
- 缓存在首次调用 `GetDepotList()` 时初始化,在程序生命周期内有效
- 缓存使用 `sync.Map` 确保并发安全(迁移时最多10个 goroutine 并发访问)
- 在 `pkg/vcs/coding.go` 的 `GetSubGroup()` 方法中继续使用重试机制,但缓存命中后无需重试

**优化效果预估**:
- 100个仓库的项目:API 调用从100次降至1次,节约约 99 次 × (0.5-2秒 + 重试时间)
- 1000个仓库的项目:API 调用从1000次降至1次,节约约 999 次 × (0.5-2秒 + 重试时间)

## Impact

**受影响的规格**:
- 新增规格:`coding-api-cache` - CODING API 缓存机制

**受影响的代码**:
- `pkg/api/coding/api.go` - 添加缓存变量和缓存逻辑
- `pkg/vcs/coding.go` - GetSubGroup() 方法保持不变,继续调用 GetProjectByName(缓存优化对调用方透明)

**向后兼容**:
- ✅ 完全向后兼容
- ✅ 不改变 VCS 接口
- ✅ 不改变 GetProjectByName() 函数签名
- ✅ 缓存对调用方透明,失败时自动降级为 API 调用

**风险**:
- 缓存在程序生命周期内有效,如果 CODING 项目元数据在迁移过程中被修改,不会生效(低风险,迁移期间项目元数据通常不变)
- 并发访问缓存使用 `sync.Map`,确保线程安全

**不改变的行为**:
- GetSubGroup() 的重试机制保持不变
- API 失败时的降级行为保持不变
- 日志记录行为保持不变
