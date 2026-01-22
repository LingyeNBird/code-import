# Spec: CODING API 缓存机制

## ADDED Requirements

### Requirement: 项目元数据缓存
系统 SHALL 实现 CODING 项目元数据缓存机制,减少重复 API 调用次数。

#### Scenario: 首次获取项目元数据
- **当** 首次调用 `coding.GetProjectByName(projectName)` 且缓存中不存在该项目
- **则** 系统应调用 CODING API 获取项目元数据并存入缓存
- **并且** 返回获取到的项目元数据

#### Scenario: 缓存命中
- **当** 调用 `coding.GetProjectByName(projectName)` 且缓存中已存在该项目
- **则** 系统应直接从缓存返回项目元数据
- **并且** 不调用 CODING API
- **并且** 日志记录 "从缓存获取项目 {projectName} 信息"

#### Scenario: 并发访问缓存
- **假设** 10 个 goroutine 并发调用 `GetProjectByName("project-A")`
- **当** 缓存中不存在 "project-A"
- **则** 系统应确保缓存的并发安全性(使用 sync.Map)
- **并且** CODING API 调用次数应远小于 10 次(理想情况为 1 次,实际可能为 2-3 次)

#### Scenario: API 失败不缓存错误结果
- **当** 调用 `coding.GetProjectByName(projectName)` 时 CODING API 返回错误
- **则** 系统应返回错误给调用方
- **并且** 不将错误结果存入缓存
- **并且** 下次调用时仍会重新尝试 API 请求

### Requirement: 缓存生命周期
缓存 MUST 在程序生命周期内有效,无需实现过期和持久化机制。

#### Scenario: 程序启动时初始化缓存
- **当** 程序启动时
- **则** 系统应初始化空的项目缓存 `projectCache sync.Map`
- **并且** 缓存在程序运行期间持续有效

#### Scenario: 缓存不持久化
- **当** 程序退出后
- **则** 缓存数据应被丢弃
- **并且** 下次程序启动时重新初始化空缓存

### Requirement: 缓存透明性
缓存机制 MUST 对调用方完全透明,不改变 API 函数签名和行为。

#### Scenario: GetProjectByName 函数签名保持不变
- **当** 实现缓存后
- **则** `GetProjectByName(url, token, projectName string) (Project, error)` 函数签名应保持不变
- **并且** 返回值类型和错误处理逻辑应保持不变

#### Scenario: GetSubGroup 调用方无需修改
- **当** 实现缓存后
- **则** `CodingVcs.GetSubGroup()` 方法(pkg/vcs/coding.go:40-88)应无需任何修改
- **并且** 重试机制应继续正常工作

#### Scenario: 缓存失败时降级为直接 API 调用
- **当** 缓存读写发生异常(理论上不会发生,但作为保护机制)
- **则** 系统应降级为直接调用 CODING API
- **并且** 确保迁移流程不受缓存故障影响

### Requirement: 性能优化效果
缓存机制 MUST 显著减少 API 调用次数,提升迁移性能。

#### Scenario: 单项目多仓库场景
- **假设** 一个 CODING 项目包含 100 个仓库
- **当** 迁移这 100 个仓库时
- **则** `GetProjectByName()` 应仅调用 CODING API 1 次
- **并且** 其余 99 次调用应从缓存读取
- **并且** API 调用次数减少率应达到 99%

#### Scenario: 多项目混合场景
- **假设** 迁移涉及 10 个 CODING 项目,每个项目 10 个仓库
- **当** 迁移这 100 个仓库时
- **则** `GetProjectByName()` 应调用 CODING API 10 次(每个项目 1 次)
- **并且** 其余 90 次调用应从缓存读取
- **并且** API 调用次数减少率应达到 90%
