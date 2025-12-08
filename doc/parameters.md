# 核心参数
- **PLUGIN_SOURCE_TOKEN**
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：调用源代码托管平台 API 的 token (当 source_platform 不为 common 时必填)
        - CODING权限：**仅限团队负责人或团队管理员token**,用户信息-只读、项目信息-只读、代码仓库-只读 https://e.coding.net/user/account/setting/tokens
        - Github权限：repo:all、read:org https://github.com/settings/tokens （classic token）
        - Gitlab权限：read_api https://gitlab.com/-/user_settings/personal_access_tokens
        - Gitee权限：user_info、projects https://gitee.com/profile/personal_access_tokens
        - 阿里云云效权限：代码仓库:只读 https://account-devops.aliyun.com/settings/personalAccessToken
        - CNB权限：account-engage:r、group-resource:r https://cnb.cool/profile/token
        - 工蜂权限: api、read_repository https://git.woa.com/profile/account
        - Gitea权限: read:organization、read:repository、read:user http(s)://<YOUR-GITEA-HOST>/user/settings/applications (设置-应用-生成新的令牌)
        - 华为云 Codearts repo 权限: 仓库读写 https://devcloud.{YOUR-REGION}.huaweicloud.com/codehub/tokens，YOUR-REGION 替换为对应的区域, 如华北-北京四对应值为 cn-north-4,参考文档 https://support.huaweicloud.com/api-codeartsrepo/codeartsrepo_05_0001.html#section0
- **PLUGIN_SOURCE_AK**
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：调用源代码托管平台 API 的 access key (当 source_platform 为 huaweicloud 时必填)，创建位置：控制台-个人信息-我的凭证-访问密钥
- **PLUGIN_SOURCE_SK**
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：调用源代码托管平台 API 的 secret key (当 source_platform 为 huaweicloud 时必填)，创建位置：控制台-个人信息-我的凭证-访问密钥
        - Gitea权限: read:organization、read:repository、read:user http(s)://<YOUR-GITEA-HOST>/user/settings/applications (设置-应用-生成新的令牌)
- **PLUGIN_SOURCE_REGION**
    - 类型：字符串
    - 必填：否
    - 默认值：cn-north-4
    - 说明: 华为云 CodeArts Repo 开通所在区域编号,详见https://support.huaweicloud.com/api-codeartsrepo/codeartsrepo_05_0001.html#section0

- **PLUGIN_CNB_ROOT_ORGANIZATION**
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：迁移后，CNB对应的根组织名称，请确保根组织已提前创建,不需要带/
    - Ex: cnb

- **PLUGIN_CNB_TOKEN**
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：CNB 授权令牌，个人令牌-访问令牌创建 https://cnb.cool/profile/token
    - 权限要求：常见场景选择`迁移工具凭据`
    - 授权范围：全部仓库/制品库

- **PLUGIN_SOURCE_URL**
    - 类型：字符串
    - 必填：是
    - 默认值：https://e.coding.net
    - 说明：源仓库代码托管平台URL
    - Ex:
       - Github: https://github.com
       - GitLab: https://gitlab.com
       - Gitee: https://gitee.com
       - common: https://common.com

- **PLUGIN_SOURCE_PLATFORM**
    - 类型：字符串
    - 必填：是
    - 默认值：coding
    - 说明：迁移的平台名称，支持 coding/gitlab/github/gitee/aliyun/cnb/gongfeng/huaweicloud，其他通用平台填写 common；

- **PLUGIN_SOURCE_REPO**
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：  
  仓库路径，当source_platform 为 common，需要与source_url拼接成完整的源仓库http克隆地址，如仓库地址是`https://common.git.com/group1/repo1.git`，PLUGIN_SOURCE_URL填写`https://common.git.com`，PLUGIN_SOURCE_REPO填写`group1/repo1`，多个以英文逗号分割。  
  当source_platform 不为 common 且 PLUGIN_SOURCE_REPO 不为空时，填写仓库路径，迁移指定仓库，多个仓库以英文逗号隔开。
    - Ex:  根据源平台不同，仓库路径格式不同，请参考下表：  
    <项目名>/<仓库名> （CODING、华为云）  
    <组织名>/<仓库名>（Gitlab、Gitee、Github、Gitea、Gongfeng）   
    <组织ID>/<组织名>/<仓库名> （阿里云云效）  
    <组织名>/<仓库名>（CNB，不用包含根组织）  
    如不清楚如何填写，可以开启`PLUGIN_MIGRATE_ALLOW_SELECT_REPOS`选项，通过查看生成的`repo_path.txt`文件来确认，详见参数介绍。


- **PLUGIN_SOURCE_USERNAME**
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：当 source_platform 为 common 时必填，clone 代码仓库时要用到的用户名，需要确保能够clone所有仓库。

- **PLUGIN_SOURCE_PASSWORD**
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：当 source_platform 为 common 时必填，clone 代码仓库时要用到的密码

- **PLUGIN_SOURCE_ORGANIZATIONID**
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明： 当 source_platform 为 aliyun 时必填，阿里云云效代码仓库组织ID，可在云效管理后台获取。  [云效帮助文档](https://help.aliyun.com/zh/yunxiao/user-guide/enterprise-basic-operations?spm=a2c4g.11186623.0.0.b86bb43fPj59Ic#caf924e65dbad)

- **PLUGIN_CNB_URL**
    - 类型：字符串
    - 必填：是
    - 默认值：https://cnb.cool
    - 说明：CNB 访问 URL
#  更多参数
- **PLUGIN_SOURCE_GROUP**
  - 类型：字符串
  - 必填：否
  - 说明：当从 CNB 迁移至 CNB 时，指定迁移根组织下仓库

- **PLUGIN_MIGRATE_TYPE**
    - 类型：字符串
    - 必填：否
    - 默认值：team
    - 说明：要迁移的类型，支持项目(project)、仓库(repo)、团队(team)多维度迁移，只支持 coding 平台

- **PLUGIN_MIGRATE_CONCURRENCY**
    - 类型：数值
    - 必填：否
    - 默认值：5
    - 说明：仓库迁移并发数，最大10

- **PLUGIN_MIGRATE_FORCE_PUSH**
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：⚠️强制push到CNB仓库,CNB侧仓库会被强制覆盖

- **PLUGIN_MIGRATE_SKIP_EXISTS_REPO**
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：跳过 CNB 已存在的仓库

- **PLUGIN_MIGRATE_USE_LFS_MIGRATE**
    - 类型：布尔值
    - 必填：否
    - 默认值：true
    - 说明：是否使用lfs migrate 处理历史提交中超过CNB单文件最大限制错误
      ⚠️如开启该配置，迁移后 commit ID会与源仓库不一致

- **PLUGIN_MIGRATE_ORGANIZATION_MAPPING_LEVEL**
    - 类型：字符串
    - 必填：否
    - 默认值：1
    - 说明：源仓库路径与CNB组织映射关系，一般保持默认即可。  
      1: 迁移完后的仓库路径为`<CNB根组织>/<源仓库路径>`，如源仓库路径为`group1/repo1`，迁移后CNB侧仓库路径为`<CNB根组织>/group1/repo1`，会自动创建子组织。  
      2: 迁移完后的仓库路径为`<CNB根组织>/<源仓库名>`， 如源仓库路径为`group1/repo1`，迁移后CNB侧仓库路径为`<CNB根组织>/repo1`，不会自动创建子组织。 ⚠️不能有同名仓库，否则会冲突报错

- **PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH**
    - 类型：字符串
    - 必填：否
    - 默认值：false
    - 说明：⚠️针对LFS源文件丢失的仓库，忽略LFS文件下载报错，LFS推送时忽略丢失的对象报错，继续推送

- **PLUGIN_MIGRATE_LOG_LEVEL**
    - 类型：字符串
    - 必填：否
    - 默认值：info
    - 说明：日志级别(debug/info/warn/error)

- **PLUGIN_MIGRATE_RELEASE**
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：迁移release（暂时只支持 gitlab/github/gitee/coding release迁移）

- **PLUGIN_MIGRATE_FILE_LIMIT_SIZE**
    - 类型：数值
    - 必填：否
    - 默认值：256
    - 说明：CNB最大文件大小限制，单位Mib

- **PLUGIN_MIGRATE_CODE**
    - 类型：布尔值
    - 必填：是
    - 默认值：true
    - 说明：迁移代码

- **PLUGIN_MIGRATE_SSH**
    - 类型：布尔值
    - 必填：是
    - 默认值：false
    - 说明：使用ssh协议克隆通用第三方平台代码仓库

- **PLUGIN_MIGRATE_REBASE**
    - 类型：布尔值
    - 必填：是
    - 默认值：false
    - 说明：使用git rebase保证代码同步，CNB侧提交的流水线配置代码不会被覆盖,如遇冲突需人工解决。  
      ⚠️如开启该配置，将启用强制推送（PLUGIN_MIGRATE_FORCE_PUSH="true"），并在迁移工具执行的工作目录备份CNB侧代码仓库

- **PLUGIN_SOURCE_PROJECT**
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：要迁移的 CODING 项目名称 (当 source_platform 为 coding 且 **migrate_type 为 project** 时必填)，多个项目以英文逗号隔开

- **PLUGIN_MIGRATE_ALLOW_SELECT_REPOS**
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：是否允许用户选择迁移指定仓库,为 true 时启用，将在工作目录生成 `repo-path.txt` ，编辑后再次运行迁移命令及只迁移 `repo-path.txt`中命中的仓库。

- **PLUGIN_MIGRATE_DOWNLOAD_ONLY**
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：是否只执行仓库下载操作，不执行迁移。为 true 时仅克隆仓库到本地，不推送到 CNB 平台。该模式下无需提供 CNB 相关配置信息。

- **PLUGIN_MIGRATE_EXCLUDE_GITHUB_FORK**
  - 类型：布尔值
  - 必填：否
  - 默认值：false
  - 说明：过滤掉 GitHub fork 的仓库不执行迁移 

- **PLUGIN_MIGRATE_GITLAB_PROJECTS_OWNED**
  - 类型：布尔值
  - 必填：否
  - 默认值：false
  - 说明：迁移 Gitlab 仓库时，仅限当前用户明确拥有的项目