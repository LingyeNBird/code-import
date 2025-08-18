# CNB Code Import(CNB代码仓库批量迁移工具)

![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/git-clone-yyds)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/status/push)

## 📒功能介绍
1. 支持 CODING、GitHub、GitLab、Gitee、Codeup(云效)、CNB、腾讯工蜂、通用第三方代码托管平台以及本地裸仓库(Local)的代码仓库批量迁移至 CNB
2. 自动创建 CNB 子组织及仓库(迁移完后的仓库路径为`<CNB根组织>/<源仓库路径>`)
3. 自动跳过迁移成功的仓库(⚠️依赖工作目录下的`successful.log`文件)
4. SVN 仓库不支持迁移，请先自行转换为git仓库


## 🌟迁移前准备
1. 创建源平台访问令牌 
2. CNB 创建根组织  
3. 创建 CNB 访问令牌 

    [详细步骤](doc/ready.md)


## 🚀快速开始
- [在云原生构建中使用](doc/web-trigger.md) 
🔥推荐
- [在 Docker上使用](doc/docker-usage.md)
在 Docker 上使用支持用户自定义更多高级用法



## 🔖参数介绍

- [详细介绍](doc/parameters.md)


## ❓常见问题
1. 超过了单个文件大小限制 500 MiB
可以开启`PLUGIN_MIGRATE_USE_LFS_MIGRATE`参数，详见[更多参数](doc/parameters.md)
2. 获取仓库列表失败: The current scope does not support access to this API
检查 PLUGIN_SOURCE_TOKEN 权限是否符合要求，如源平台为 CODING，确保 token 属于团队所有者或团队管理员。
3. LFS: Repository or object not found
可以开`PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH`详见[更多参数](doc/parameters.md)
4. push 失败：git pull before pushing again
可根据实际情况开启`PLUGIN_MIGRATE_FORCE_PUSH`，详见[更多参数](doc/parameters.md)
5. 只迁移部分仓库怎么操作？  
设置 `PLUGIN_MIGRATE_ALLOW_SELECT_REPOS=true`，首次运行后编辑 `repo-path.txt`，只保留需要迁移的仓库路径即可。
6. 如何重新选择迁移仓库？  
删除 `repo-path.txt` 文件，重新运行迁移命令即可。
7. repo-path.txt 没有生成？  
请确认 `PLUGIN_MIGRATE_ALLOW_SELECT_REPOS=true`，并确保有写入权限。
8. 如何只下载仓库而不推送到 CNB？  
设置 `PLUGIN_MIGRATE_DOWNLOAD_ONLY=true`，该模式下仅执行仓库克隆操作，无需提供 CNB 相关配置，下载完成后会保留工作目录。
