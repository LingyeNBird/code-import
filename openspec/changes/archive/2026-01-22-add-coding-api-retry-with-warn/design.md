# 技术设计: Coding API 重试机制

## 背景

CNB Code Import 工具从 CODING 平台迁移仓库到 CNB。在迁移过程中,它调用 `coding.GetProjectByName()` 获取项目元数据(显示名称、描述)用于映射到 CNB 子组织属性。

**当前行为:**
- 在 `pkg/vcs/coding.go:41` 单次尝试 API 调用
- 失败时:记录 ERROR 日志,返回 `SubGroup`,错误信息嵌入到 `Desc` 字段,`Remark` 设为 "ERROR"
- 迁移继续,但使用了被污染的元数据

**约束条件:**
- API 失败时不能终止迁移(用户明确要求)
- 应遵循代码库中现有的重试模式(`pkg/git/git.go`)
- 日志必须使用结构化日志记录器(zap),并使用适当的级别
- 按项目约定使用中文注释

**相关方:**
- 执行批量迁移的用户(需要容错能力)
- 运维团队(需要清晰的警告日志用于监控)

## 目标 / 非目标

**目标:**
- 添加 3 次尝试的重试机制,使用指数退避(1秒、5秒、10秒)
- 使用 WARN 级别记录失败日志,包含完整错误上下文
- 即使所有重试失败后也继续迁移流程
- 失败时返回干净的 `SubGroup` 结构(仅填充 `Name`)
- 与现有重试模式保持一致

**非目标:**
- 为其他 API 调用添加重试(范围限制为用户请求的 `GetProjectByName`)
- 可配置的重试参数(像 `git.go` 中一样硬编码)
- 指数退避计算(使用简单的预定义间隔)
- 熔断器模式(批量迁移工具不需要)
- API 调用指标/遥测(没有现有基础设施)

## 决策

### 1. 重试间隔
**决策:** 使用固定间隔 [1秒, 5秒, 10秒]

**理由:**
- 与 `pkg/git/git.go:50` 中的现有模式匹配
- 总最大延迟约16秒对批量迁移可以接受
- 无需退避计算的简单实现

**考虑的替代方案:**
- 真正的指数退避 (2^n): 否决,对于3次尝试过度设计
- 通过 `config.yaml` 可配置: 否决,为边缘情况增加复杂性

### 2. 最终失败时的错误处理
**决策:** 返回 `SubGroup{Name: c.SubGroupName, Desc: "", Remark: ""}` 并记录 WARN 日志

**理由:**
- 调用者仍然可以通过名称识别子组织
- WARN 级别表示可恢复的问题(不是致命的 ERROR)
- 迁移继续处理剩余仓库
- 符合用户要求:"不退出主程序"
- 不污染业务数据字段

**考虑的替代方案:**
- 返回错误并跳过仓库: 否决,违背优雅降级原则
- Panic/退出: 否决,用户明确禁止
- 使用 ERROR 级别: 否决,对于非致命问题应使用 WARN

### 3. 日志消息格式
**决策:** 
```go
// 每次重试失败时
logger.Logger.Warnf("获取项目 %s 信息失败 (尝试 %d/3): %v", 
    c.SubGroupName, i+1, err)

// 最终失败时
logger.Logger.Warnf("获取项目 %s 信息最终失败,已重试3次,仅使用项目名称继续迁移: %v", 
    c.SubGroupName, err)
```

**理由:**
- 遵循 `pkg/git/git.go:66,117` 中的现有格式
- 包含尝试次数用于调试
- 包含完整错误详情用于故障排查
- 按项目约定使用中文

### 4. 范围限制
**决策:** 仅为 `GetSubGroup()` 方法添加重试,不包括 `GetReleases()`

**理由:**
- 用户明确指定"当前代码在 vcs/coding.go:43"(第43行在 GetSubGroup 中)
- `GetReleases()` 目前在错误时会 panic(第105行),需要单独讨论
- 保持变更范围最小化和聚焦

**未来考虑:** 
- 在单独的变更中对 `GetReleases()` 应用类似模式
- 考虑移除 panic,采用优雅错误处理

## 实现方法

### 代码位置
文件: `pkg/vcs/coding.go`  
方法: `CodingVcs.GetSubGroup()` (第40-67行)

### 伪代码
```go
func (c *CodingVcs) GetSubGroup() *SubGroup {
    // 重试间隔配置:第1次失败后等1秒,第2次失败后等5秒,第3次失败后等10秒
    retryIntervals := []time.Duration{1 * time.Second, 5 * time.Second, 10 * time.Second}
    
    var project coding.Project
    var err error
    
    // 重试循环:最多尝试3次
    for i, interval := range retryIntervals {
        project, err = coding.GetProjectByName(
            config.Cfg.GetString("source.url"), 
            c.GetToken(), 
            c.SubGroupName
        )
        
        if err == nil {
            // 成功,跳出重试循环
            break
        }
        
        // 记录警告日志
        logger.Logger.Warnf("获取项目 %s 信息失败 (尝试 %d/%d): %v", 
            c.SubGroupName, i+1, len(retryIntervals), err)
        
        // 如果不是最后一次尝试,则等待指定时间后重试
        if i < len(retryIntervals)-1 {
            time.Sleep(interval)
        }
    }
    
    // 所有重试都失败,记录最终警告并返回最小 SubGroup
    if err != nil {
        logger.Logger.Warnf("获取项目 %s 信息最终失败,已重试%d次,仅使用项目名称继续迁移: %v", 
            c.SubGroupName, len(retryIntervals), err)
        return &SubGroup{
            Name: c.SubGroupName,
        }
    }
    
    // 成功路径:使用获取的元数据填充 SubGroup
    var desc, remark string
    if config.Cfg.GetBool("migrate.map_coding_description") {
        desc = strings.TrimSpace(project.Description)
    }
    if config.Cfg.GetBool("migrate.map_coding_display_name") {
        remark = strings.TrimSpace(project.DisplayName)
    }
    return &SubGroup{
        Name:   c.SubGroupName,
        Desc:   desc,
        Remark: remark,
    }
}
```

### 需要的代码变更
1. 添加 `import "time"`(如果文件中尚不存在)
2. 用重试循环替换第41-50行
3. 更新成功路径(第52-66行基本保持不变)
4. 移除将错误嵌入 SubGroup 字段的旧错误处理

## 风险 / 权衡

### 风险1: 增加迁移时间
**描述:** 失败的 API 调用每个仓库增加最多16秒延迟  
**可能性:** 低(API 调用通常成功)  
**缓解:** 仅影响 API 失败的仓库;对批量迁移可接受

### 风险2: 掩盖配置问题
**描述:** 持久的 API 认证/配置错误可能被当作暂时性错误处理  
**可能性:** 中  
**影响:** 中 - 运维人员可能不会立即注意到配置错误  
**缓解:** 
- 生产环境应监控 WARN 日志
- 考虑在迁移结束时添加汇总统计
- 未来工作:区分持久性失败和暂时性失败(例如 401 vs 503)

### 风险3: WARN vs ERROR 级别的歧义
**描述:** 运维人员可能不够重视 WARN 日志  
**可能性:** 低  
**缓解:** 
- 日志消息明确说明"最终失败,已重试3次"
- 现有的 successful.log 机制跟踪成功的迁移
- 可以通过缺少元数据在 CNB 中识别失败的仓库

### 权衡: 代码重复
**描述:** 重试模式从 `git.go` 重复而不是抽象  
**理由:** 
- 代码库中目前只有2个地方使用此模式
- 抽象增加复杂性(通用重试包装器)
- 按照"简单、经过验证的模式"原则,简单重复是可接受的
- 未来工作:如果模式在3个以上位置使用则提取

## 迁移计划

**部署步骤:**
1. 将代码变更合并到 main 分支
2. 构建新的 Docker 镜像(`cnbcool/code-import:latest`)
3. 更新 CNB CI/CD 流水线使用新镜像
4. 无需配置变更(向后兼容)

**回滚计划:**
- 回退提交并重新构建 Docker 镜像
- 恢复之前的行为(无重试的立即失败)

**监控:**
- 检查 `migrate.log` 中的 WARN 级别消息
- 模式: `获取项目 .* 信息失败` 或 `获取项目 .* 信息最终失败`
- 这些警告的高频率表明需要调查的 API 问题

**生产测试:**
- 使用故意无效的 CODING token 运行迁移
- 验证发生3次重试尝试
- 验证 WARN 日志包含错误详情
- 验证迁移继续到下一个仓库

## 未决问题

1. ❓ 是否应该对第102行的 `GetReleases()` 应用类似的重试逻辑?
   - 目前在错误时会 panic(第105行)
   - 需要单独的变更讨论

2. ❓ 重试间隔是否应该通过 `config.yaml` 可配置?
   - 当前方法:为简单起见硬编码
   - 如果需要自定义,需要用户反馈

3. ❓ 是否应该跟踪 API 失败统计并在迁移结束时报告?
   - 示例:"成功迁移 95/100 个仓库,5个有元数据获取失败"
   - 超出此变更范围,潜在的未来增强
