# CNB Code Importer(CNB代码仓库导入工具)

![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/git-clone-yyds)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/status/push)

主要包含以下特性：
1. 支持CODING、Github、Gitlab、Gitee以及通用第三方代码托管平台的代码仓库批量导入至CNB
2. 支持按CODING团队、项目、仓库多维度迁移
3. 支持Git LFS文件迁移
4. 自动跳过迁移成功的仓库(⚠️依赖工作目录下的`successful.log`文件)

⚠️开始导入前，请确保根组织已存在。

**Tips:云原生开发自带 docker 命令，内网运行更快速**
## 在Docker上使用

### 从 Coding 导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://coding.example.com" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_CNB_URL="https://cnb.example.com" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从 Github 导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://github.com" \
  -e PLUGIN_SOURCE_PLATFORM="github" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_CNB_URL="https://cnb.example.com" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从 Gitlab 导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://gitlab.example.com" \
  -e PLUGIN_SOURCE_PLATFORM="gitlab" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_CNB_URL="https://cnb.example.com" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从 Gitee 导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://gitlab.com" \
  -e PLUGIN_SOURCE_PLATFORM="gitee" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_CNB_URL="https://cnb.example.com" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从通用第三方代码平台导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_USERNAME="xxx"  \
  -e PLUGIN_SOURCE_PASSWORD="xxx"  \
  -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
  -e PLUGIN_SOURCE_URL="https://common.example.com" \
  -e PLUGIN_SOURCE_PLATFORM="common" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_CNB_URL="https://cnb.example.com" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

## 参数
如需使用，请加上plugin_前缀并全部转为大写。  

| 参数名 | 类型  | 必填 | 默认值    | 说明 |
|--------|-----|---|--------|------|
| source_url | 字符串 | 是 | -      | 代码托管平台URL |
| source_token | 字符串 | 是 | -      | 调用源代码托管平台 API 的 token (当 source_platform 不为 common 时必填)<br>- CODING权限：用户信息-只读、项目信息-只读、代码仓库-只读<br>- Github权限：repo:all、read:org<br>- Gitlab权限：read_api<br>- Gitee权限：user_info、projects |
| source_platform | 字符串 | 是 | coding | 导入的平台名称，支持 coding/gitlab/github/gitee，其他通用平台填写 common |
| source_project | 字符串 | 否 | -      | 要迁移的 CODING 项目名称 (当 source_platform 为 coding 且 migrate_type 为 project 时必填)<br>多个项目以英文逗号隔开 |
| source_repo | 字符串 | 否 | -      | 当source_platform 为 common 或者 source_platform 为 coding 且 migrate_type 为 repo 时必填<br>多个代码仓库以英文逗号隔开 |
| source_username | 字符串 | 否 | -      | 当 source_platform 为 common 时必填，clone 代码仓库时要用到的用户名 |
| source_password | 字符串 | 否 | -      | 当 source_platform 为 common 时必填，clone 代码仓库时要用到的密码 |
| cnb_url | 字符串 | 是 | -      | CNB访问URL |
| cnb_token | 字符串 | 是 | -      | CNB 授权令牌，权限要求：<br>- repo-code 读写<br>- repo-basic-info 只读<br>- account-profile 只读<br>- account-engage 只读<br>- group-resource 读写<br>- group-manage 读写<br>- repo-content 读写 |
| cnb_root_organization | 字符串 | 是 | -      | 迁移后，CNB对应的根组织名称，请确保根组织已提前创建 |
| migrate_type | 字符串 | 否 | team   | 要迁移的类型，支持项目(project)、仓库(repo)、团队(team)多维度迁移，只支持 coding 平台 |
| migrate_concurrency | 数值  | 否 | 5      | 仓库迁移并发数，最大10 |
| migrate_force_push | 布尔值 | 否 | false  | 强制push到CNB仓库 |
| migrate_skip_exists_repo | 布尔值 | 否 | true   | 跳过CNB已存在的仓库 |
| migrate_use_lfs_migrate | 字符串 | 否 | false  | 是否使用lfs migrate处理历史提交中超过CNB单文件最大限制错误<br>**注意：如开启该配置，迁移后commit ID会与源仓库不一致** |
| migrate_organization_mapping_level | 字符串 | 否 | 1      | CODING与CNB组织映射关系，仅支持 Coding 平台<br>1: CODING项目映射为CNB子组织，仓库在子组织下面<br>2: CODING项目不会映射为CNB子组织，仓库直接在CNB根组织下面 |
| migrate_allow_incomplete_push | 字符串 | 否 | false  | 针对LFS源文件丢失的仓库，忽略LFS文件下载报错，LFS推送时忽略丢失的对象报错，继续推送 |
| migrate_log_level | 字符串 | 否 | info   | 日志级别(debug/info/warn/error) |
| migrate_release | 布尔值 | 否 | false  | 迁移release（暂时只支持 gitlab release迁移） |
| migrate_file_limit_size | 数值  | 否 | 100    | CNB最大文件大小限制，单位Mib |
