# CNB Code Import(CNB代码仓库批量导入工具)

![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/git-clone-yyds)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/status/push)

## 功能介绍
1. 支持CODING、Github、Gitlab、Gitee、阿里云以及通用第三方代码托管平台的代码仓库批量导入至CNB
2. 自动跳过迁移成功的仓库(⚠️依赖工作目录下的`successful.log`文件)



## 在Docker上使用

### 注意事项
1. ⚠️开始导入前，请确保CNB根组织已存在。  
2. `xxx`为需要用户自行替换的字段，具体含义详见参数介绍。  
3. 云原生开发自带 docker 命令，内网运行更快速

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

### 从阿里云导入

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
  -e PLUGIN_CNB_URL="https://cnb.example.com" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从通用第三方代码平台导入

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
  -e PLUGIN_CNB_URL="https://cnb.example.com" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````
ssh协议

⚠️使用ssh协议时请在当前工作目录确保有对应的私钥文件，文件名固定为`ssh.key`
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
  -e PLUGIN_SOURCE_URL="https://common.example.com" \
  -e PLUGIN_SOURCE_PLATFORM="common" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_CNB_URL="https://cnb.example.com" \
  -e PLUGIN_MIGRATE_SSH="true" \
  -e GIT_SSH_COMMAND='ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null' \  
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````




## 参数介绍


| 参数名                                       | 类型  | 必填 | 默认值                             | 说明                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
|-------------------------------------------|-----|----|---------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| PLUGIN_SOURCE_URL                         | 字符串 | 是  | -                               | 代码托管平台URL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| PLUGIN_SOURCE_TOKEN                       | 字符串 | 是  | -                               | 调用源代码托管平台 API 的 token (当 source_platform 不为 common 时必填)<br>- CODING权限：用户信息-只读、项目信息-只读、代码仓库-只读<br>- Github权限：repo:all、read:org<br>- Gitlab权限：read_api<br>- Gitee权限：user_info、projects                                                                                                                                                                                                                                                                                                        |
| PLUGIN_SOURCE_PLATFORM                    | 字符串 | 是  | coding                          | 导入的平台名称，支持 coding/gitlab/github/gitee，其他通用平台填写 common                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| PLUGIN_SOURCE_PROJECT                     | 字符串 | 否  | -                               | 要迁移的 CODING 项目名称 (当 source_platform 为 coding 且 migrate_type 为 project 时必填)<br>多个项目以英文逗号隔开                                                                                                                                                                                                                                                                                                                                                                                                   |
| PLUGIN_SOURCE_REPO                        | 字符串 | 否  | -                               | 当source_platform 为 common 或者 source_platform 为 coding 且 migrate_type 为 repo 时必填<br>多个代码仓库以英文逗号隔开                                                                                                                                                                                                                                                                                                                                                                                            |
| PLUGIN_SOURCE_USERNAME                    | 字符串 | 否  | -                               | 当 source_platform 为 aliyun或common 时必填，clone 代码仓库时要用到的用户名，需要确保能够clone所有仓库。<br/>[阿里云帮助文档](https://help.aliyun.com/zh/yunxiao/user-guide/configure-https-clone-account-password?spm=a2c4g.11186623.0.0.78b240cdASV98n)                                                                                                                                                                                                                                                                         |
| PLUGIN_SOURCE_PASSWORD                    | 字符串 | 否  | -                               | 当 source_platform 为 aliyun或common 时必填，clone 代码仓库时要用到的密码                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| PLUGIN_SOURCE_AK                          | 字符串 | 否  | -                               | AccessKey ID,当 source_platform 为 aliyun时必填，需要有AliyunRDCReadOnlyAccess权限，如果是RAM用户，需要关联至云效账号，并授权管理员角色。<br/>[阿里云帮助文档](https://help.aliyun.com/zh/yunxiao/user-guide/add-a-ram-user?spm=5176.28366559.console-base_help.dexternal.211e336a7R37d8&scm=20140722.S_help%40%40%E6%96%87%E6%A1%A3%40%40203014.S_RQW%40ag0%2BBB2%40ag0%2BBB1%40ag0%2Bos0.ID_203014-RL_ram%E7%94%A8%E6%88%B7%E5%A6%82%E4%BD%95%E5%85%B3%E8%81%94%E8%87%B3%E4%BA%91%E6%95%88-LOC_console~UND~help-OR_ser-V_4-P0_0-P1_0) |
| PLUGIN_SOURCE_AS                          | 字符串 | 否  | -                               | AccessKey Secret,当 source_platform 为 aliyun时必填。                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| PLUGIN_SOURCE_ENDPOINT                    | 字符串 | 否  | devops.cn-hangzhou.aliyuncs.com | AccessKey 请求的地址。                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| PLUGIN_SOURCE_ORGANIZATIONID              | 字符串 | 否  | -                               | 阿里云云效代码仓库企业ID，可在云效访问链接中获取，如https://devops.aliyun.com/organization/【OrganizationId】                                                                                                                                                                                                                                                                                                                                                                                                          |
| PLUGIN_CNB_URL                            | 字符串 | 是  | -                               | CNB访问URL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| PLUGIN_CNB_TOKEN                          | 字符串 | 是  | -                               | CNB 授权令牌，权限要求：<br>- repo-code 读写<br>- repo-basic-info 只读<br>- account-profile 只读<br>- account-engage 只读<br>- group-resource 读写<br>- group-manage 读写<br>- repo-content 读写                                                                                                                                                                                                                                                                                                                    |
| PLUGIN_CNB_ROOT_ORGANIZATION              | 字符串 | 是  | -                               | 迁移后，CNB对应的根组织名称，请确保根组织已提前创建                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| PLUGIN_MIGRATE_TYPE                       | 字符串 | 否  | team                            | 要迁移的类型，支持项目(project)、仓库(repo)、团队(team)多维度迁移，只支持 coding 平台                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| PLUGIN_MIGRATE_CONCURRENCY                | 数值  | 否  | 10                              | 仓库迁移并发数，最大10                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| PLUGIN_MIGRATE_FORCE_PUSH                 | 布尔值 | 否  | false                           | 强制push到CNB仓库                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| PLUGIN_MIGRATE_SKIP_EXISTS_REPO           | 布尔值 | 否  | true                            | 跳过CNB已存在的仓库                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| PLUGIN_MIGRATE_USE_LFS_MIGRATE            | 字符串 | 否  | false                           | 是否使用lfs migrate处理历史提交中超过CNB单文件最大限制错误<br>⚠️如开启该配置，迁移后commit ID会与源仓库不一致                                                                                                                                                                                                                                                                                                                                                                                                                       |
| PLUGIN_MIGRATE_ORGANIZATION_MAPPING_LEVEL | 字符串 | 否  | 1                               | CODING与CNB组织映射关系，仅支持 Coding 平台<br>1: CODING项目映射为CNB子组织，仓库在子组织下面<br>2: CODING项目不会映射为CNB子组织，仓库直接在CNB根组织下面                                                                                                                                                                                                                                                                                                                                                                                     |
| PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH      | 字符串 | 否  | false                           | 针对LFS源文件丢失的仓库，忽略LFS文件下载报错，LFS推送时忽略丢失的对象报错，继续推送                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| PLUGIN_MIGRATE_LOG_LEVEL                  | 字符串 | 否  | info                            | 日志级别(debug/info/warn/error)                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| PLUGIN_MIGRATE_RELEASE                    | 布尔值 | 否  | false                           | 迁移release（暂时只支持 gitlab release迁移）                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| PLUGIN_MIGRATE_FILE_LIMIT_SIZE            | 数值  | 否  | 100                             | CNB最大文件大小限制，单位Mib                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| PLUGIN_MIGRATE_CODE                       | 布尔值 | 是  | true                            | 迁移代码                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| PLUGIN_MIGRATE_SSH                        | 布尔值 | 是  | false                           | 使用ssh协议克隆通用第三方平台代码仓库                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |

