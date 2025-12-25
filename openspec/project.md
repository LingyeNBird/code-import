# 项目上下文

## 目的
CNB Code Import 是一个批量代码仓库迁移工具,旨在促进从多个源平台向 CNB(集中式代码托管平台)进行大规模仓库迁移。

**核心目标:**
- 自动化从 10 多个 VCS 平台迁移(CODING、GitHub、GitLab、Gitee、Gitea、阿里云 Codeup、华为云 CodeArts Repo、CNB、腾讯工蜂以及通用平台)
- 使用 `<CNB 根组织>/<源仓库路径>` 层级结构保留仓库结构
- 自动处理大文件(>256 MiB),将其转换为 Git LFS 对象
- 支持 CODING 特定功能:将项目显示名称映射到 CNB 子组织别名,将项目描述映射到组织描述
- 跟踪迁移状态以跳过已迁移的仓库(通过 `successful.log`)
- 支持并发迁移(最多 10 个仓库同时进行)以提高效率

## 技术栈
- **语言:** Go 1.23.0+ (工具链 1.23.4)
- **CLI 框架:** Cobra v1.8.0 - 命令行接口构建器
- **配置管理:** Viper v1.18.2 - 多源配置管理(YAML、环境变量、命令行参数)
- **VCS 集成:**
  - `google/go-github` v66 - GitHub API v3 客户端
  - `xanzy/go-gitlab` v0.108.0 - GitLab API 客户端
  - `huaweicloud/huaweicloud-sdk-go-v3` - 华为云 CodeArts Repo SDK
- **日志:** Uber Zap v1.21.0 - 高性能结构化日志
- **并发控制:** `golang.org/x/sync` v0.11.0 - 信号量限流
- **HTTP/认证:** `golang.org/x/oauth2` v0.27.0 - OAuth2 认证
- **序列化:** `gopkg.in/yaml.v3` v3.0.1 - YAML 解析
- **构建/运行时:** Docker(多阶段构建:Golang 构建器 + Alpine 运行时)
- **版本控制:** Git + Git LFS 处理大文件

## 项目约定

### 代码风格
- **语言:** 使用 Go 惯用模式;内部代码使用中文注释,全局/导出 API 使用英文
- **命名约定:**
  - 包级变量:PascalCase 或 camelCase
  - 常量:UPPER_CASE 使用下划线
  - 导出类型/函数:PascalCase
  - 私有类型/函数:camelCase
- **错误处理:** 显式错误返回,包含丰富的上下文信息;错误使用结构化字段记录
- **日志记录:** 使用 `uber/zap`,采用适当级别(Info、Debug、Warn、Error);在日志中屏蔽敏感数据(令牌、密码)
- **注释:** 内部文档主要使用中文;面向公众的文档使用英文
- **文件组织:** 按功能分组到包中(`cmd/` 用于 CLI,`pkg/` 用于核心逻辑)
- **无正式 Linter:** 项目非正式遵循 Go 惯用法;考虑添加 `.golangci.yml` 以保持一致性

### 架构模式
- **插件架构:** VCS 平台实现通用 `VCS` 接口(`pkg/vcs/interface.go`)以实现可扩展性
- **工厂模式:** `NewVcs(platform)` 返回平台特定实现(coding.go、github.go、gitlab.go 等)
- **仓库模式:** `pkg/api/*` 中的 API 客户端抽象平台特定的 REST/SDK 调用
- **并发模型:** 基于信号量的速率限制(最多 10 个并发迁移)以避免 API 过载
- **命令模式:** `cmd/` 中的 Cobra CLI 命令编排高级操作
- **包装器模式:** `pkg/git/git.go` 包装 Git CLI 命令,提供重试逻辑和错误处理
- **配置层级:** YAML 文件 → 环境变量 → CLI 参数(优先级递增)
- **原子操作:** 使用原子计数器进行线程安全的统计跟踪(成功/失败/跳过计数)

### 测试策略
- **框架:** Go 标准 `testing` 包(不使用 testify 等外部框架)
- **测试覆盖:** 8 个测试文件覆盖关键路径:
  - 配置解析(`config_test.go`、`url_test.go`)
  - VCS 平台实现(`coding_test.go`、`gitee_test.go`)
  - 迁移逻辑(`migrate_test.go`)
  - Git 操作(`git_test.go`)
  - API 客户端(`api_test.go`)
  - 工具函数(`util_test.go`)
- **测试模式:** 表驱动测试处理多种场景;模拟外部依赖(API、Git 命令)
- **CI 集成:** 在 CNB CI/CD 流水线中对拉取请求自动运行测试
- **未来改进:** 考虑添加端到端迁移场景的集成测试

### Git 工作流
- **分支策略:**
  - `main` 分支:受保护,需要合并请求
  - 功能分支:描述性命名(如 `vincent-patch-1`、`fix-gitlab-patch0`、`openspec`)
  - 不直接提交到 `main`
- **提交消息约定:** 主要使用中文;遵循模式:`类型: 描述`
  - **类型:** `feat:`(功能)、`fix:`(修复)、`docs:`(文档)、`merge:`(合并请求)、`chore:`(维护)
  - **示例:** `feat: 优化手动流水线,设置CODING迁移按钮为默认选项`
- **版本管理:** 使用 Git 标签的语义化版本控制(如 `v1.60.1`、`v1.61.0`)
- **CI/CD 触发器:**
  - **推送到 main:** Docker 构建 + 推送、自动标记、知识库构建
  - **拉取请求:** AI 驱动的代码审查
  - **标签推送:** 变更日志生成、发布上传
- **代码审查:** 所有更改在合并前都需要合并请求批准

## 领域上下文
**代码仓库迁移:**
- **源平台:** 各种 VCS 平台上的组织/项目/仓库
- **目标平台:** CNB(类似 GitLab/GitHub 的中国代码托管平台)
- **迁移范围:**
  - 代码(带完整历史的 Git 仓库)
  - Release(标签、发布说明、资产) - 可选
  - 大文件(如 >256 MiB 自动转换为 Git LFS)
- **组织层级:** CNB 使用嵌套组织(`root_org/sub_org/repo`);工具根据源仓库路径自动创建子组织
- **CODING 平台特性:**
  - 将 CODING 项目显示名称映射到 CNB 子组织别名
  - 将 CODING 项目描述映射到 CNB 子组织描述
  - 处理 CODING 特定的 API 细节(团队 ID 要求、仓库结构)
- **迁移状态跟踪:** `successful.log` 文件记录已迁移的仓库,以实现增量迁移并跳过重复处理
- **并发控制:** 限制并发操作为 10 个,以避免 API 速率限制和资源耗尽
- **凭证管理:** 使用平台特定的认证(令牌、OAuth2、AK/SK);日志中屏蔽凭证

## 重要约束
- **API 速率限制:** 必须遵守平台特定的速率限制(如 GitHub:认证用户 5000 请求/小时)
- **大文件处理:** >256 MiB 的文件需要 Git LFS;标准 Git 推送会失败
- **CNB 根组织:** 迁移前必须存在;工具无法创建根组织
- **并发迁移限制:** 最多 10 个并发仓库迁移以避免资源耗尽
- **Docker 执行超时:** CNB Web 触发器对长时间运行的迁移有 10 小时超时限制
- **网络连接:** 需要访问源平台 API 和 CNB API;防火墙可能阻止某些平台
- **磁盘空间:** 推送到 CNB 之前将整个仓库克隆到 `source_git_dir/`;确保有足够的磁盘空间
- **认证令牌:** 需要适当的权限范围:
  - 源平台:读取仓库、Release
  - CNB:创建组织、仓库、推送代码
- **平台可用性:** 工具行为取决于源平台 API 可用性和 CNB 服务健康状况
- **Git 版本:** 执行环境需要安装 Git 2.x+ 和 Git LFS

## 外部依赖
**必需服务:**
1. **CNB 平台** (https://cnb.cool 或自定义实例)
   - API:仓库创建、组织管理、代码推送
   - 认证:具有 `api`、`write_repository` 范围的个人访问令牌
2. **源 VCS 平台:**
   - **CODING** (https://coding.net):OAuth2 或个人访问令牌
   - **GitHub** (https://github.com):具有 `repo`、`read:org` 范围的个人访问令牌
   - **GitLab** (https://gitlab.com 或自托管):具有 `api`、`read_repository` 范围的个人访问令牌
   - **Gitee** (https://gitee.com):个人访问令牌
   - **Gitea** (自托管):个人访问令牌
   - **阿里云 Codeup** (https://codeup.aliyun.com):AK/SK 凭证
   - **华为云 CodeArts Repo** (https://devcloud.huaweicloud.com):AK/SK 凭证
   - **腾讯工蜂(Gongfeng)**:平台特定认证
   - **通用平台**:HTTP Basic Auth 或基于令牌的认证
3. **Git + Git LFS:** 用于仓库操作的本地 Git 安装
4. **Docker Registry**(用于容器化部署):
   - 镜像:`cnbcool/code-import:latest`、`cnbcool/code-import:v{VERSION}`

**配置文件:**
- `config.yaml`:主要配置(源/目标平台、迁移选项)
- `successful.log`:跟踪已迁移的仓库(自动生成)
- `migrate.log`:迁移执行日志(自动生成)
- `repo-path.txt`:用于选择性仓库迁移的可选白名单

**环境变量(Docker 模式):**
- 所有配置键都可以通过 `PLUGIN_*` 前缀的变量设置(如 `PLUGIN_SOURCE_PLATFORM`、`PLUGIN_CNB_TOKEN`)

**CI/CD 集成:**
- `.cnb.yml`:CNB CI/CD 流水线配置(Docker 构建、标记、变更日志生成)
- 使用 CNB 提供的 Docker 服务和容器镜像仓库
