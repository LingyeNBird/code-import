# CNB Code Import(CNB代码仓库批量迁移工具)

![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/git-clone-yyds)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/status/push)

## 功能介绍
1. 支持CODING、Github、Gitlab、Gitee、阿里云以及通用第三方代码托管平台的代码仓库批量迁移至CNB
2. 自动跳过迁移成功的仓库(⚠️依赖工作目录下的`successful.log`文件)


## 在Docker上使用

### 注意事项
1. ⚠️开始迁移前，请确保CNB根组织已存在。  
2. `xxx`为需要用户自行替换的字段，具体含义详见参数介绍-核心参数。  
3. 云原生开发自带 docker 命令，内网运行更快速

### 从 Coding 迁移

迁移之后的效果：原 CODING 项目会在 CNB 中创建一个同名的子组织，并将原项目下的仓库迁移至该子组织下面

迁移团队下所有仓库
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

迁移指定项目仓库  
PLUGIN_SOURCE_PROJECT 字段根据需要自行替换，详见参数介绍
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_PROJECT="project1,project2" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_TYPE="project" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```


迁移指定仓库  
PLUGIN_SOURCE_REPO 字段根据需要自行替换，详见参数介绍
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_REPO="<TEAM-NAME>/<PROJECT-NAME>/<REPO-NAME>,test-team/project1/repoA,test-team/project2/repoB" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_TYPE="repo" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### 从 Github 迁移

迁移之后的效果：原 Github 账号下有权限的所有组织，会在 CNB 中创建同名的子组织，并将原组织下的仓库迁移至该子组织下面


```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://github.com" \
  -e PLUGIN_SOURCE_PLATFORM="github" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### 从 Gitlab 迁移

迁移之后的效果：原 GitLab 账号下有权限的所有 group，会在 CNB 中创建同名的子组织，并将原group下的仓库迁移至该子组织下面（如果 gitlab 是多级的group，迁移至 CNB 子组织和仓库仍会保留原有的多层级结构）

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://gitlab.com" \
  -e PLUGIN_SOURCE_PLATFORM="gitlab" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### 从 Gitee 迁移

迁移之后的效果：原 Gitee 账号下有权限的所有组织，会在 CNB 中创建同名的子组织，并将原组织下的仓库迁移至该子组织下面（如果 Gitee 是多级的组织/仓库组，迁移至 CNB 子组织和仓库仍会保留原有的多层级结构）

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://gitee.com" \
  -e PLUGIN_SOURCE_PLATFORM="gitee" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### 从阿里云迁移

迁移之后的效果：原云效账号下有权限的所有组织，会在 CNB 中创建同名的子组织，并将原仓库组下的仓库迁移至该子组织下面（如果云效是多级的仓库组，迁移至 CNB 子组织和仓库仍会保留原有的多层级结构）

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_USERNAME="xxx"  \
  -e PLUGIN_SOURCE_PASSWORD="xxx"  \
  -e PLUGIN_SOURCE_PLATFORM="aliyun" \
  -e PLUGIN_SOURCE_AK="xxx" \
  -e PLUGIN_SOURCE_AS="xxx" \
  -e PLUGIN_SOURCE_ORGANIZATIONID="xxx" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### 从通用第三方代码平台迁移

PLUGIN_SOURCE_REPO 字段中，group字段会映射为子组织，如果 cnb 的根组织下，没有该命名的子组织将会自动创建，如果有该命名的子组织会将仓库创建在已有的同名子组织下面

http协议
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_USERNAME="xxx"  \
  -e PLUGIN_SOURCE_PASSWORD="xxx"  \
  -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
  -e PLUGIN_SOURCE_URL="https://common.example.com" \
  -e PLUGIN_SOURCE_PLATFORM="common" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```
ssh协议

⚠️使用ssh协议时请在当前工作目录确保有对应的私钥文件，文件名固定为`ssh.key`
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
  -e PLUGIN_SOURCE_URL="https://common.example.com" \
  -e PLUGIN_SOURCE_PLATFORM="common" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e GIT_SSH_COMMAND='ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null' \  
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### 从工蜂迁移

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://git.woa.com" \
  -e PLUGIN_SOURCE_PLATFORM="gongfeng" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### 从 CNB 迁移

迁移指定根组织下所有仓库
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_GROUP="xxx" \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://cnb.example1.com" \
  -e PLUGIN_SOURCE_PLATFORM="cnb" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_CNB_URL="https://cnb.example2.com" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### 只迁移部分仓库（仓库选择功能）

首次运行，生成仓库列表文件 `repo-path.txt`  


这里以CODING为例，其他平台只需在原有迁移命令基础上增加`-e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \`参数即可：

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

**首次运行后，工具会在当前目录生成 `repo-path.txt`，请手动编辑该文件，仅保留需要迁移的仓库路径。**

编辑完成后，再次运行同样的命令即可只迁移你选择的仓库：

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

**如需重新选择仓库，只需删除 `repo-path.txt`，重新运行上述命令即可。**


### 仅下载仓库模式

如果只需要将仓库下载到本地而不推送到 CNB 平台，可以使用仅下载模式。该模式下会跳过所有 CNB 相关操作，仅执行仓库克隆，无需提供 CNB 相关配置信息。

这里以 CODING 为例，其他平台只需在原有迁移命令基础上增加 `-e PLUGIN_MIGRATE_DOWNLOAD_ONLY="true"` 参数即可：

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_DOWNLOAD_ONLY="true" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

**仅下载模式特点：**
- 仅执行仓库克隆操作，不推送到 CNB 平台
- 无需提供 CNB 相关配置信息（如 CNB_TOKEN、CNB_ROOT_ORGANIZATION 等）
- 跳过所有 CNB 相关操作（如创建子组织、创建仓库等）
- 下载完成后保留工作目录（格式：`source_git_dir_时间戳`）

### 迁移完成后，增量更新原平台最新内容
清空原工作目录下的 successful.log

效果：重新同步原平台的所有仓库，已迁移至 CNB 的仓库，如在原平台有更新，会将内容增量同步至 CNB 平台。



## 参数介绍

### 核心参数

- PLUGIN_SOURCE_URL
    - 类型：字符串
    - 必填：是
    - 默认值：https://e.coding.net
    - 说明：源仓库代码托管平台URL
    - Ex:
       - github: https://github.com
       - gitlab: https://gitlab.com
       - gitee: https://gitee.com
       - common: https://common.com

- PLUGIN_SOURCE_TOKEN
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：调用源代码托管平台 API 的 token (当 source_platform 不为 common 时必填)
        - CODING权限：**仅限团队所有者或团队管理员token**,用户信息-只读、项目信息-只读、代码仓库-只读 https://e.coding.net/user/account/setting/tokens
        - Github权限：repo:all、read:org https://github.com/settings/tokens
        - Gitlab权限：read_api https://gitlab.com/-/user_settings/personal_access_tokens
        - Gitee权限：user_info、projects https://gitee.com/profile/personal_access_tokens
        - CNB权限：account-engage:r、group-resource:r https://cnb.cool/profile/token
        - 工蜂权限: api、read_repository https://git.woa.com/profile/account

- PLUGIN_SOURCE_PLATFORM
    - 类型：字符串
    - 必填：是
    - 默认值：coding
    - 说明：迁移的平台名称，支持 coding/gitlab/github/gitee，其他通用平台填写 common

- PLUGIN_SOURCE_REPO
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：仓库路径，需要与source_url拼接成完整的源仓库http克隆地址,如https://common.com/group1/repo1  
  当source_platform 为 common 或者 source_platform 为 coding 且 **migrate_type 为 repo** 时必填，多个代码仓库以英文逗号隔开
    - Ex: group1/repo1,group1/repo2,group2/repo3

- PLUGIN_SOURCE_USERNAME
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：当 source_platform 为 aliyun或common 时必填，clone 代码仓库时要用到的用户名，需要确保能够clone所有仓库。
      [阿里云帮助文档](https://help.aliyun.com/zh/yunxiao/user-guide/configure-https-clone-account-password?spm=a2c4g.11186623.0.0.78b240cdASV98n)

- PLUGIN_SOURCE_PASSWORD
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：当 source_platform 为 aliyun或common 时必填，clone 代码仓库时要用到的密码

- PLUGIN_SOURCE_AK
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：AccessKey ID,当 source_platform 为 aliyun时必填，需要有AliyunRDCReadOnlyAccess权限，如果是RAM用户，需要关联至云效账号，并授权管理员角色。
      [阿里云帮助文档](https://help.aliyun.com/zh/yunxiao/user-guide/add-a-ram-user?spm=5176.28366559.console-base_help.dexternal.211e336a7R37d8&scm=20140722.S_help%40%40%E6%96%87%E6%A1%A3%40%40203014.S_RQW%40ag0%2BBB2%40ag0%2BBB1%40ag0%2Bos0.ID_203014-RL_ram%E7%94%A8%E6%88%B7%E5%A6%82%E4%BD%95%E5%85%B3%E8%81%94%E8%87%B3%E4%BA%91%E6%95%88-LOC_console~UND~help-OR_ser-V_4-P0_0-P1_0)

- PLUGIN_SOURCE_AS
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：AccessKey Secret,当 source_platform 为 aliyun时必填。

- PLUGIN_SOURCE_ENDPOINT
    - 类型：字符串
    - 必填：否
    - 默认值：devops.cn-hangzhou.aliyuncs.com
    - 说明：AccessKey 请求的地址。

- PLUGIN_SOURCE_ORGANIZATIONID
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：阿里云云效代码仓库企业ID，可在云效访问链接中获取，如https://devops.aliyun.com/organization/【OrganizationId】

- PLUGIN_CNB_URL
    - 类型：字符串
    - 必填：是
    - 默认值：https://cnb.cool
    - 说明：CNB访问URL

- PLUGIN_CNB_TOKEN
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：CNB 授权令牌，个人令牌-访问令牌创建 https://cnb.cool/profile/token。
    - 权限要求：
        - repo-code 读写
        - repo-basic-info 只读
        - account-profile 只读
        - account-engage 只读
        - group-resource 读写
        - group-manage 读写
        - repo-content 读写

- PLUGIN_CNB_ROOT_ORGANIZATION
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：迁移后，CNB对应的根组织名称，请确保根组织已提前创建,不需要带/
    - Ex: root-group


## 其他参数
- PLUGIN_SOURCE_GROUP
  - 类型：字符串
  - 必填：否
  - 说明：当从CNB迁移至CNB时，指定迁移根组织下仓库

- PLUGIN_MIGRATE_TYPE
    - 类型：字符串
    - 必填：否
    - 默认值：team
    - 说明：要迁移的类型，支持项目(project)、仓库(repo)、团队(team)多维度迁移，只支持 coding 平台

- PLUGIN_MIGRATE_CONCURRENCY
    - 类型：数值
    - 必填：否
    - 默认值：10
    - 说明：仓库迁移并发数，最大10

- PLUGIN_MIGRATE_FORCE_PUSH
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：⚠️强制push到CNB仓库,CNB侧仓库会被强制覆盖

- PLUGIN_MIGRATE_SKIP_EXISTS_REPO
    - 类型：布尔值
    - 必填：否
    - 默认值：true
    - 说明：跳过CNB已存在的仓库

- PLUGIN_MIGRATE_USE_LFS_MIGRATE
    - 类型：字符串
    - 必填：否
    - 默认值：false
    - 说明：是否使用lfs migrate处理历史提交中超过CNB单文件最大限制错误
      ⚠️如开启该配置，迁移后commit ID会与源仓库不一致

- PLUGIN_MIGRATE_ORGANIZATION_MAPPING_LEVEL
    - 类型：字符串
    - 必填：否
    - 默认值：1
    - 说明：CODING与CNB组织映射关系，仅支持 Coding 平台
      1: CODING项目映射为CNB子组织，仓库在子组织下面
      2: CODING项目不会映射为CNB子组织，仓库直接在CNB根组织下面

- PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH
    - 类型：字符串
    - 必填：否
    - 默认值：false
    - 说明：⚠️针对LFS源文件丢失的仓库，忽略LFS文件下载报错，LFS推送时忽略丢失的对象报错，继续推送

- PLUGIN_MIGRATE_LOG_LEVEL
    - 类型：字符串
    - 必填：否
    - 默认值：info
    - 说明：日志级别(debug/info/warn/error)

- PLUGIN_MIGRATE_RELEASE
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：迁移release（暂时只支持 gitlab/github/gitee/coding release迁移）

- PLUGIN_MIGRATE_FILE_LIMIT_SIZE
    - 类型：数值
    - 必填：否
    - 默认值：500
    - 说明：CNB最大文件大小限制，单位Mib

- PLUGIN_MIGRATE_CODE
    - 类型：布尔值
    - 必填：是
    - 默认值：true
    - 说明：迁移代码

- PLUGIN_MIGRATE_SSH
    - 类型：布尔值
    - 必填：是
    - 默认值：false
    - 说明：使用ssh协议克隆通用第三方平台代码仓库

- PLUGIN_MIGRATE_REBASE
    - 类型：布尔值
    - 必填：是
    - 默认值：false
    - 说明：在源和目标都有变更的情况下,且CNB侧仓库对应分支根目录有`.cnb.yml`文件，使用git rebase保证代码同步，CNB侧提交的流水线配置代码不会被覆盖
      ⚠️如开启该配置，将启用强制推送（PLUGIN_MIGRATE_FORCE_PUSH="true"），并在迁移工具执行的工作目录备份CNB侧代码仓库

- PLUGIN_SOURCE_PROJECT
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：要迁移的 CODING 项目名称 (当 source_platform 为 coding 且 **migrate_type 为 project** 时必填)，多个项目以英文逗号隔开

- PLUGIN_MIGRATE_ALLOW_SELECT_REPOS
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：是否允许用户选择迁移仓库。为 true 时启用 repo-path.txt 选择功能。

- PLUGIN_MIGRATE_DOWNLOAD_ONLY
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：是否只执行仓库下载操作，不执行迁移。为 true 时仅克隆仓库到本地，不推送到 CNB 平台。该模式下无需提供 CNB 相关配置信息。

## 常见问题
1. 超过了单个文件大小限制 500 MiB
可以开启`PLUGIN_MIGRATE_USE_LFS_MIGRATE`参数，详见参数介绍
2. 获取仓库列表失败: The current scope does not support access to this API
检查PLUGIN_SOURCE_TOKEN权限是否符合要求，如源平台为 CODING，确保token属于团队所有者或团队管理员。
3. LFS: Repository or object not found
可以开`PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH`详见参数介绍
4. push失败：git pull before pushing again
可根据实际情况开启`PLUGIN_MIGRATE_FORCE_PUSH`，详见参数介绍
5. 只迁移部分仓库怎么操作？  
设置 `PLUGIN_MIGRATE_ALLOW_SELECT_REPOS=true`，首次运行后编辑 `repo-path.txt`，只保留需要迁移的仓库路径即可。
6. 如何重新选择迁移仓库？  
删除 `repo-path.txt` 文件，重新运行迁移命令即可。
7. repo-path.txt 没有生成？  
   请确认 `PLUGIN_MIGRATE_ALLOW_SELECT_REPOS=true`，并确保有写入权限。
8. 如何只下载仓库而不推送到 CNB？  
   设置 `PLUGIN_MIGRATE_DOWNLOAD_ONLY=true`，该模式下仅执行仓库克隆操作，无需提供 CNB 相关配置，下载完成后会保留工作目录。