# CNB Code Import (Batch Repository Migration Tool)

![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/git-clone-yyds)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/status/push)

## 📒Features
1. Support batch migration from CODING, GitHub, GitLab, Gitee, Codeup(Alibaba Cloud), CNB, Tencent Git, Gitea, Huawei Cloud(CodeArts Repo), and other third-party code hosting platforms to CNB
2. Automatically create CNB sub-organizations and repositories (migrated repository path will be `<CNB root org>/<source repo path>`)
3. CODING source repositories will map project display names to CNB sub-organization aliases, and project descriptions to sub-organization descriptions
4. Automatically skip successfully migrated repositories (⚠️ depends on `successful.log` file in working directory)

## 💥Important Notes (Must Read)
1. SVN repositories are not supported, please convert to git repositories first
2. **After migration is complete, ensure code commits are only made on one platform, otherwise re-migration may cause conflicts**
3. CNB sub-organizations are visible to external members by default. To modify this, enable Root Organization - Organization Settings - Organization Control - Hide Sub-organizations

## 🌟Preparation
1. Create access token for source platform
2. Create root organization in CNB  
3. Create CNB access token  

    [Detailed steps](doc/ready.md)

## 🚀Quick Start
- [Use with Cloud Native Build](doc/web-trigger.md) 
🔥Recommended
- [Use with Docker](doc/docker-usage.md)
Docker usage supports more advanced customizations

## 🔖Parameters
- [Detailed introduction](doc/parameters.md)

## ❓FAQ
1. Exceeds single file size limit of 256 MiB
Enable `PLUGIN_MIGRATE_USE_LFS_MIGRATE` parameter, see [More parameters](doc/parameters.md)
2. Failed to get repository list/project info: `The current scope does not support access to this API`
Check if PLUGIN_SOURCE_TOKEN has required permissions. For CODING platform, ensure token belongs to team owner or admin, see [CODING Token requirements](doc/ready.md)
3. Failed to download LFS files `LFS: Repository or object not found`  
Enable `PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH`, see [More parameters](doc/parameters.md)
4. Push failed: `git pull before pushing again`  
Enable `PLUGIN_MIGRATE_FORCE_PUSH` based on actual situation, see [More parameters](doc/parameters.md)
5. How to migrate only specific repositories?  
Set `PLUGIN_MIGRATE_ALLOW_SELECT_REPOS=true`, after first run edit `repo-path.txt` to keep only the repositories you want to migrate.