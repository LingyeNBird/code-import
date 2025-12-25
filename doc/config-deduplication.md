# 配置去重功能说明

## 功能概述

系统会自动对 `PLUGIN_SOURCE_REPO` 和 `PLUGIN_SOURCE_PROJECT` 环境变量传入的列表进行去重处理，确保每个仓库或项目仅保留一个有效记录。

## 处理时机

去重发生在**业务逻辑执行阶段**：
- `source.repo` 在 `migrate` 包的 `filterReposByConfigList()` 函数中去重
- `source.project` 在 `coding` API 包的 `GetDepotList()` 函数中去重

## 处理逻辑

### 触发条件

当处理以下配置项时自动触发去重：
- `source.repo` - 仓库路径列表（在仓库过滤时）
- `source.project` - 项目名称列表（在获取仓库列表时）

### 处理流程

```
环境变量 PLUGIN_SOURCE_REPO / PLUGIN_SOURCE_PROJECT
    ↓
配置读取 (config 包)
    ↓
字符串分割（按逗号）
    ↓
业务逻辑使用（migrate 包 / coding API）
    ↓
去重处理 (util.DeduplicateStringSlice)
    ├─ 去除前后空格
    ├─ 过滤空字符串
    ├─ 移除重复项
    └─ 保留首次出现
    ↓
日志输出（logger）
    ├─ WARN: 重复项详情
    └─ INFO: 去重统计
    ↓
使用去重后的列表
```

### 去重规则

1. **空格处理**：自动去除每项的前后空格
2. **空值过滤**：跳过空字符串和纯空格字符串
3. **重复判断**：使用 Go map 进行 O(n) 时间复杂度去重
4. **大小写敏感**：严格区分大小写
5. **保留顺序**：保留首次出现的项，维持原始顺序

## 日志输出

### 重复项警告

每发现一个重复项，输出 WARN 级别日志（最多显示 3 个）：

```
2025-12-24T15:12:09.506+0800  WARN  migrate/migrate.go:72  配置 source.repo 中发现重复项: org/project/repo1，已自动过滤
```

如果重复项超过 3 个，会额外输出：

```
2025-12-24T15:12:09.506+0800  WARN  migrate/migrate.go:75  配置 source.repo 中还有 5 个重复项未显示
```

### 去重统计

去重完成后，输出 INFO 级别统计日志（仅当有重复项时）：

```
2025-12-24T15:12:09.506+0800  INFO  migrate/migrate.go:78  配置 source.repo 去重完成：原始配置 5 项，去重后 3 项，移除重复项 2 个
```

## 使用示例

### 示例 1：仓库列表去重

**环境变量配置**：
```bash
export PLUGIN_SOURCE_REPO="org/project/repo1,org/project/repo2,org/project/repo1,org/project/repo3"
```

**去重结果**：
```
原始: ["org/project/repo1", "org/project/repo2", "org/project/repo1", "org/project/repo3"]
去重: ["org/project/repo1", "org/project/repo2", "org/project/repo3"]
```

**日志输出**：
```
2025-12-24T15:12:09.506+0800  WARN  migrate/migrate.go:72  配置 source.repo 中发现重复项: org/project/repo1，已自动过滤
2025-12-24T15:12:09.506+0800  INFO  migrate/migrate.go:78  配置 source.repo 去重完成：原始配置 4 项，去重后 3 项，移除重复项 1 个
```

### 示例 2：包含空格的配置

**环境变量配置**：
```bash
export PLUGIN_SOURCE_REPO=" org/project/repo1 , org/project/repo2,  org/project/repo1  , org/project/repo3"
```

**去重结果**：
```
原始: [" org/project/repo1 ", " org/project/repo2", "  org/project/repo1  ", " org/project/repo3"]
去除空格: ["org/project/repo1", "org/project/repo2", "org/project/repo1", "org/project/repo3"]
去重: ["org/project/repo1", "org/project/repo2", "org/project/repo3"]
```

### 示例 3：项目列表去重

**环境变量配置**：
```bash
export PLUGIN_SOURCE_PROJECT="project1,project2,project1,project3,project2"
```

**去重结果**：
```
原始: ["project1", "project2", "project1", "project3", "project2"]
去重: ["project1", "project2", "project3"]
```

**日志输出**：
```
2025-12-24T15:12:09.506+0800  WARN  coding/api.go:72  配置 source.project 中发现重复项: project1，已自动过滤
2025-12-24T15:12:09.506+0800  WARN  coding/api.go:72  配置 source.project 中发现重复项: project2，已自动过滤
2025-12-24T15:12:09.506+0800  INFO  coding/api.go:78  配置 source.project 去重完成：原始配置 5 项，去重后 3 项，移除重复项 2 个
```

### 示例 4：包含空字符串

**环境变量配置**：
```bash
export PLUGIN_SOURCE_REPO="org/project/repo1,,org/project/repo2,  ,org/project/repo3"
```

**去重结果**：
```
原始: ["org/project/repo1", "", "org/project/repo2", "  ", "org/project/repo3"]
过滤空值: ["org/project/repo1", "org/project/repo2", "org/project/repo3"]
```

## 特性

### ✅ 自动化

- 无需手动干预，配置解析时自动去重
- 对用户完全透明
- 不影响正常的配置使用

### ✅ 性能优化

- 使用 Go map 实现 O(n) 时间复杂度
- 内存占用小
- 只在真正使用配置时才执行去重
- 重复项日志输出限制在 3 个，避免日志过多

### ✅ 保持兼容性

- 不破坏现有功能
- 对无重复配置无任何影响（不输出日志）
- 只处理 `source.repo` 和 `source.project` 两个配置项
- 职责分离：config 负责读取，util 提供工具，migrate/API 负责业务逻辑

### ✅ 清晰反馈

- 每个重复项都有日志记录（最多显示 3 个，避免日志泛滥）
- 提供详细的统计信息（原始数、去重后数、重复数）
- 使用结构化日志（zap），便于问题排查
- 日志格式与项目其他模块统一

## 代码位置

### 核心去重逻辑
- **工具函数**：`pkg/util/util.go::DeduplicateStringSlice()`
  - 返回 `DeduplicationResult` 结构体
  - 不依赖 logger，避免循环依赖
  - 可被任意模块复用

### 日志封装
- **migrate 包**：`pkg/migrate/migrate.go::deduplicateWithLog()`
  - 封装日志输出逻辑
  - 在 `filterReposByConfigList()` 中对 `source.repo` 去重
  
- **coding API**：`pkg/api/coding/api.go::GetDepotList()`
  - 对 `source.project` 去重并输出日志

### 单元测试
- **工具函数测试**：`pkg/util/dedup_test.go`（12 个测试用例）
- **集成测试**：`pkg/migrate/dedup_test.go`（带日志输出）

## 测试覆盖

包含 12 个单元测试用例，覆盖以下场景：

- ✅ 无重复配置
- ✅ 包含重复配置
- ✅ 包含空格的配置
- ✅ 包含空字符串的配置
- ✅ 全是空字符串的配置
- ✅ 空切片
- ✅ 保留首次出现顺序
- ✅ 大小写敏感
- ✅ 项目配置去重
- ✅ 混合格式的仓库路径
- ✅ 单个配置项
- ✅ 仅包含空字符串

所有测试用例均通过。

## 注意事项

1. **大小写敏感**：`Org/Project/Repo` 和 `org/project/repo` 被视为不同项
2. **仅处理指定配置**：只对 `source.repo` 和 `source.project` 去重
3. **不修改源数据**：只处理内存中的配置，不修改环境变量
4. **日志级别**：重复项使用 WARN 级别，统计使用 INFO 级别
5. **无副作用**：去重操作不影响任何其他配置项
6. **去重时机**：在业务逻辑使用配置时去重，而非配置读取时
7. **日志限制**：最多显示 3 个重复项，超过部分仅显示数量统计
8. **架构设计**：采用分层设计，util 提供纯函数，migrate/API 负责日志

## 架构设计

### 分层架构

```
┌─────────────────────────────────────────────────┐
│  config 包 (配置读取)                           │
│  - 从环境变量读取配置                           │
│  - 字符串分割处理                               │
│  - 不做去重（避免循环依赖）                     │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│  util 包 (工具函数)                             │
│  - DeduplicateStringSlice() 核心去重逻辑        │
│  - 返回 DeduplicationResult 结构体              │
│  - 不依赖 logger 和 config                      │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│  migrate 包 (迁移业务逻辑)                      │
│  - deduplicateWithLog() 封装日志输出            │
│  - filterReposByConfigList() 调用去重           │
│  - 使用 logger.Logger 输出统一格式日志          │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  coding API 包 (平台 API)                       │
│  - GetDepotList() 对 source.project 去重        │
│  - 直接调用 util 包并输出日志                   │
└─────────────────────────────────────────────────┘
```

### 设计原则

1. **避免循环依赖**：util 包不依赖 logger 和 config
2. **职责分离**：config 读取、util 工具、migrate/API 业务逻辑
3. **复用性**：DeduplicateStringSlice 可被任意模块使用
4. **测试性**：核心逻辑与日志输出分离，易于单元测试

## 版本历史

- **2025-12-24**：重构配置去重功能
  - 将去重逻辑从 config 包移至 util 包
  - 在 migrate 和 coding API 包中调用去重
  - 使用结构化日志（zap）输出
  - 添加日志输出限制（最多 3 个重复项）
  - 解决循环依赖问题
  - 优化架构设计，职责更清晰

- **2025-12-23**：初始实现配置去重功能
  - 在 config 包中实现去重逻辑
  - 支持 `source.repo` 和 `source.project` 自动去重
  - 添加详细的日志记录
  - 包含完整的单元测试（12 个测试用例）
