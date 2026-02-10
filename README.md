# CNB Code Import(CNB代码仓库批量迁移工具)

![badge](/-/badge/git/latest/ci/pipeline-as-code)
![badge](/-/badge/git/latest/ci/git-clone-yyds)
![badge](/-/badge/git/latest/ci/status/push)

## 📒功能介绍
1. 支持 CODING、GitHub、GitLab、Gitee、Gitea、阿里云(Codeup)、华为云(CodeArts Repo)、CNB、腾讯工蜂、通用第三方代码托管平台代码仓库批量迁移至 CNB
2. 自动创建 CNB 子组织及仓库，迁移完后的仓库路径为`<CNB根组织>/<源仓库路径>`
3. 自动处理超过 256 MiB 的大文件，转为 LFS 对象
4. CODING 源仓库会将项目显示名称映射为 CNB 子组织别名，项目简介映射为子组织简介
5. 自动跳过迁移成功的仓库(⚠️依赖工作目录下的`successful.log`文件，云原生构建方法不支持)

## 💥注意事项（必读）
1. SVN 仓库不支持迁移，请先自行转换为git仓库
2. **迁移完成后，请确保只在一侧平台提交代码，否则再次迁移可能会冲突报错**
3. CNB 子组织默认外部成员可查看，如需修改请开启`根组织-组织设置-组织管控-隐藏子组织`


## 🌟迁移前准备
1. 源平台创建访问令牌 
2. CNB 创建根组织  
3. CNB 创建访问令牌 

    [详细步骤](doc/ready.md)


## 🚀快速开始
- [在云原生构建中使用](doc/web-trigger.md) 
🔥推荐
- [在 Docker上使用](doc/docker-usage.md)
在 Docker 上使用支持用户自定义更多高级用法



## 🔖参数介绍

- [详细介绍](doc/parameters.md)


## ❓常见问题
1. 可请点击仓库详情页上方![知识库](img/zhishiku-logo.png)按钮(迁移按钮旁边），直接提问，或者按`/`键直接输入问题，结尾带上`?`，如果仍然解决不了欢迎提issue。
1. 超过了单个文件大小限制 256 MiB
可以开启`PLUGIN_MIGRATE_USE_LFS_MIGRATE`参数，详见[更多参数](doc/parameters.md)
2. 获取仓库列表失败/获取项目信息失败: `The current scope does not support access to this API`
检查 PLUGIN_SOURCE_TOKEN 权限是否符合要求，如源平台为 CODING，确保 token 属于团队所有者或团队管理员，详见[CODDING Token要求](doc/ready.md)
3. 下载 LFS 文件失败 `LFS: Repository or object not found`  
 谨慎开启`PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH`详见[更多参数](doc/parameters.md#更多参数)
4. push 失败：`git pull before pushing again`  
可根据实际情况开启`PLUGIN_MIGRATE_FORCE_PUSH`，详见[更多参数](doc/parameters.md#更多参数)
5. 只迁移部分仓库怎么操作？  
设置 `PLUGIN_MIGRATE_ALLOW_SELECT_REPOS=true`，首次运行后编辑 `repo-path.txt`，只保留需要迁移的仓库路径即可。
