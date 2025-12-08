# 按指定仓库迁移功能使用说明

## 功能概述

支持工蜂、GitLab、GitHub、Gitee 等平台按指定仓库列表进行迁移。当 `PLUGIN_SOURCE_REPO` 环境变量（或配置文件中的 `source.repo`）不为空时，只迁移指定的仓库，其他仓库将被跳过。

## 使用场景

1. **精确迁移**: 只需要迁移特定的几个仓库，而不是整个团队/组织的所有仓库
2. **批量迁移**: 从多个项目/组中选择特定仓库进行迁移
3. **分批迁移**: 将大批量仓库分成多批次进行迁移
4. **测试验证**: 先迁移少量仓库进行测试验证

## 配置方式

### 方式一：使用环境变量（推荐）

```bash
# 单个仓库
export PLUGIN_SOURCE_REPO="owner/repo1"

# 多个仓库（逗号分隔）
export PLUGIN_SOURCE_REPO="owner/repo1,owner/repo2,group/subgroup/repo3"

# 执行迁移
./ccrctl migrate
```

### 方式二：使用配置文件

编辑 `config.yaml` 文件：

```yaml
source:
  platform: gitlab  # 或 github、gitee、gongfeng 等
  url: https://gitlab.example.com
  token: your-token
  # 指定要迁移的仓库列表（逗号分隔）
  repo:
    - group1/subgroup1/repo1
    - group1/subgroup2/repo2
    - group2/repo3

migrate:
  type: team  # 或 project
  # 其他迁移配置...
```

## 仓库路径格式

不同平台的仓库路径格式：

### GitHub
```
owner/repository
```
示例: `octocat/Hello-World`

### GitLab
```
group/subgroup/repository
或
group/repository
```
示例: 
- `gitlab-org/gitlab-ce`
- `group/subgroup1/subgroup2/repo`

### Gitee
```
owner/repository
```
示例: `openeuler/kernel`

### 工蜂（Gongfeng）
```
org/team/project/repository
或
org/team/repository
```
示例:
- `tencent/WXG/project1/repo1`
- `tencent/PCG/repo2`

### CODING
```
team/project/repository
```
示例: `my-team/my-project/my-repo`

## 使用示例

### 示例 1: GitLab 迁移指定仓库

```bash
export PLUGIN_SOURCE_PLATFORM=gitlab
export PLUGIN_SOURCE_URL=https://gitlab.example.com
export PLUGIN_SOURCE_TOKEN=glpat-xxxxxxxxxxxx
export PLUGIN_SOURCE_REPO="group1/repo1,group2/subgroup/repo2"
export PLUGIN_MIGRATE_TYPE=team
export PLUGIN_CNB_URL=https://cnb.example.com
export PLUGIN_CNB_TOKEN=your-cnb-token
export PLUGIN_CNB_ROOT_ORGANIZATION=your-org

./ccrctl migrate
```

### 示例 2: GitHub 迁移指定仓库

```bash
export PLUGIN_SOURCE_PLATFORM=github
export PLUGIN_SOURCE_URL=https://github.com
export PLUGIN_SOURCE_TOKEN=ghp_xxxxxxxxxxxx
export PLUGIN_SOURCE_REPO="owner1/repo1,owner2/repo2"
export PLUGIN_MIGRATE_TYPE=team
export PLUGIN_CNB_URL=https://cnb.example.com
export PLUGIN_CNB_TOKEN=your-cnb-token
export PLUGIN_CNB_ROOT_ORGANIZATION=your-org

./ccrctl migrate
```

### 示例 3: 工蜂迁移指定仓库

```bash
export PLUGIN_SOURCE_PLATFORM=gongfeng
export PLUGIN_SOURCE_URL=https://git.code.tencent.com
export PLUGIN_SOURCE_TOKEN=your-token
export PLUGIN_SOURCE_REPO="tencent/WXG/project1/repo1,tencent/PCG/repo2"
export PLUGIN_MIGRATE_TYPE=team
export PLUGIN_CNB_URL=https://cnb.example.com
export PLUGIN_CNB_TOKEN=your-cnb-token
export PLUGIN_CNB_ROOT_ORGANIZATION=your-org

./ccrctl migrate
```

### 示例 4: 配置文件方式（GitLab）

`config.yaml`:
```yaml
source:
  platform: gitlab
  url: https://gitlab.example.com
  token: glpat-xxxxxxxxxxxx
  repo:
    - group1/subgroup1/repo1
    - group1/subgroup2/repo2
    - group2/repo3

cnb:
  url: https://cnb.example.com
  token: your-cnb-token
  root_organization: your-org

migrate:
  type: team
  concurrency: 5
  force_push: false
```

执行迁移：
```bash
./ccrctl migrate
```

## 行为说明

### 1. 配置为空时
当 `source.repo` 为空或未配置时，将迁移源平台的所有仓库（根据 `migrate.type` 决定范围）。

### 2. 配置不为空时
只迁移配置列表中指定的仓库，其他仓库将被跳过。

### 3. 仓库路径匹配
- **精确匹配**: 仓库路径必须完全匹配（区分大小写）
- **未找到提示**: 如果配置的仓库在源平台不存在，会输出警告日志

### 4. 与其他过滤功能的关系
- `source.repo` 过滤优先执行
- 如果同时启用了 `migrate.allow_select_repos`，会先应用 `source.repo` 过滤，再应用文件选择过滤

## 日志输出

### 匹配成功
```
INFO  检测到 source.repo 配置，将只迁移指定的 3 个仓库
INFO  匹配到配置仓库: group1/repo1
INFO  匹配到配置仓库: group2/repo2
INFO  匹配到配置仓库: group3/repo3
INFO  根据 source.repo 配置过滤后，待迁移仓库数: 3
```

### 部分仓库未找到
```
INFO  检测到 source.repo 配置，将只迁移指定的 3 个仓库
INFO  匹配到配置仓库: group1/repo1
WARN  配置的仓库 group2/not-exist 在源平台未找到
INFO  匹配到配置仓库: group3/repo3
INFO  根据 source.repo 配置过滤后，待迁移仓库数: 2
```

## 注意事项

1. **仓库路径格式**: 必须使用完整的仓库路径，不支持模糊匹配
2. **大小写敏感**: 仓库路径区分大小写
3. **空格处理**: 配置中的前后空格会被自动去除
4. **逗号分隔**: 多个仓库使用英文逗号分隔
5. **路径验证**: 配置的仓库路径必须在源平台存在，否则会输出警告

## 常见问题

### Q1: 如何知道仓库的完整路径？
**A**: 可以在源平台的仓库页面查看完整路径，或使用以下方法：
- GitLab: 仓库首页的项目路径
- GitHub: `owner/repository` 格式
- 工蜂: 仓库 URL 中的路径部分

### Q2: 配置了仓库但没有迁移？
**A**: 检查以下几点：
1. 仓库路径是否正确（包括大小写）
2. 查看日志中是否有 "未找到" 的警告
3. 确认 token 有访问该仓库的权限

### Q3: 可以同时使用 source.repo 和 allow_select_repos 吗？
**A**: 可以，执行顺序为：
1. 先根据 `source.repo` 过滤
2. 再根据 `repo-path.txt` 文件过滤

### Q4: 如何迁移所有仓库？
**A**: 不配置 `source.repo` 或将其设置为空即可。

## 兼容性

- ✅ 支持的平台: GitLab、GitHub、Gitee、工蜂、CODING
- ✅ 支持的迁移类型: team、project、repo
- ✅ 与现有功能兼容: 不影响其他配置项的使用

## 更新日志

- 2025-12-04: 新增按指定仓库迁移功能，支持所有主流 Git 平台
