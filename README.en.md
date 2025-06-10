# CNB Code Import (Batch Migration Tool for CNB Repositories)

![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/git-clone-yyds)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/status/push)

## Features
1. Supports batch migration of repositories from CODING, GitHub, GitLab, Gitee, Alibaba Cloud, Gongfeng, CNB, and other third-party platforms to CNB
2. Automatically skips successfully migrated repositories (⚠️ depends on `successful.log` file in working directory)
3. Supports selective repository migration with repository selection feature
4. Supports incremental updates from source platform

## Using with Docker

### Notes
1. ⚠️ Ensure the CNB root organization exists before migration  
2. Replace `xxx` fields with your actual values (see Core Parameters section for details)  
3. Cloud Native Development includes docker commands for faster internal execution

### Migrating from CODING

Migrate all repositories under the team
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

Migrate specific project repositories
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

Migrate specific repositories
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_REPO="project1/repoA,project1/repoB,project2/repoC" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_TYPE="repo" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### Migrating from GitHub

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

### Migrating from GitLab

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

### Migrating from Gitee

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

### Migrating from Alibaba Cloud

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

### Migrating from Generic Third-party Platforms

HTTP Protocol
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

SSH Protocol

⚠️ When using SSH protocol, ensure the private key file exists in current working directory with filename `ssh.key`
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
  -e PLUGIN_SOURCE_URL="https://common.example.com" \
  -e PLUGIN_SOURCE_PLATFORM="common" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_SSH="true" \
  -e GIT_SSH_COMMAND='ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null' \  
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

### Migrating from Gongfeng

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

### Migrating from CNB

Migrate all repositories under specified root organization
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

### Selective Repository Migration

First run to generate repository list file `repo-path.txt`

Using CODING as an example (add `-e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \` parameter for other platforms):

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

**After first run, the tool will generate `repo-path.txt` in current directory. Please manually edit this file to keep only the repository paths you want to migrate.**

After editing, run the same command again to migrate only selected repositories:

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

**To reselect repositories, simply delete `repo-path.txt` and run the command again.**

### Incremental Updates from Source Platform
Clear the `successful.log` file in the working directory

Effect: Re-sync all repositories from source platform. For repositories already migrated to CNB, if there are updates in source platform, the changes will be incrementally synced to CNB platform.

## Parameter Reference

### Core Parameters

- PLUGIN_SOURCE_URL
    - Type: String
    - Required: Yes
    - Default: https://e.coding.net
    - Description: Source repository platform URL
    - Examples:
       - github: https://github.com
       - gitlab: https://gitlab.com
       - gitee: https://gitee.com
       - common: https://common.com

- PLUGIN_SOURCE_TOKEN
    - Type: String
    - Required: Yes (when source_platform ≠ common)
    - Default: -
    - Description: API token for source platform
        - CODING permissions: **Team owner or team admin token only**, Read-only for user info, project info, and repositories https://e.coding.net/user/account/setting/tokens
        - GitHub permissions: repo:all, read:org https://github.com/settings/tokens
        - GitLab permissions: read_api https://gitlab.com/-/user_settings/personal_access_tokens
        - Gitee permissions: user_info, projects https://gitee.com/profile/personal_access_tokens
        - CNB permissions: account-engage:r, group-resource:r https://cnb.cool/profile/token
        - Gongfeng permissions: api, read_repository https://git.woa.com/profile/account

- PLUGIN_SOURCE_PLATFORM
    - Type: String
    - Required: Yes
    - Default: coding
    - Description: Source platform name, supports coding/gitlab/github/gitee/gongfeng/cnb, use common for other platforms

- PLUGIN_SOURCE_REPO
    - Type: String
    - Required: No
    - Default: -
    - Description: Repository path, will be concatenated with source_url to form complete source repository http clone URL, e.g. https://common.com/group1/repo1
    Required when source_platform is common or when source_platform is coding and **migrate_type is repo**. Multiple repositories should be separated by commas.
    - Example: group1/repo1,group1/repo2,group2/repo3

- PLUGIN_SOURCE_USERNAME
    - Type: String
    - Required: No
    - Default: -
    - Description: Required when source_platform is aliyun or common. Username for cloning repositories, must have access to all repositories.
    [Alibaba Cloud Help Documentation](https://help.aliyun.com/zh/yunxiao/user-guide/configure-https-clone-account-password?spm=a2c4g.11186623.0.0.78b240cdASV98n)

- PLUGIN_SOURCE_PASSWORD
    - Type: String
    - Required: No
    - Default: -
    - Description: Required when source_platform is aliyun or common. Password for cloning repositories.

- PLUGIN_SOURCE_AK
    - Type: String
    - Required: No
    - Default: -
    - Description: AccessKey ID, required when source_platform is aliyun. Must have AliyunRDCReadOnlyAccess permission. For RAM users, must be associated with Yunxiao account and granted admin role.
    [Alibaba Cloud Help Documentation](https://help.aliyun.com/zh/yunxiao/user-guide/add-a-ram-user?spm=5176.28366559.console-base_help.dexternal.211e336a7R37d8&scm=20140722.S_help%40%40%E6%96%87%E6%A1%A3%40%40203014.S_RQW%40ag0%2BBB2%40ag0%2BBB1%40ag0%2Bos0.ID_203014-RL_ram%E7%94%A8%E6%88%B7%E5%A6%82%E4%BD%95%E5%85%B3%E8%81%94%E8%87%B3%E4%BA%91%E6%95%88-LOC_console~UND~help-OR_ser-V_4-P0_0-P1_0)

- PLUGIN_SOURCE_AS
    - Type: String
    - Required: No
    - Default: -
    - Description: AccessKey Secret, required when source_platform is aliyun.

- PLUGIN_SOURCE_ENDPOINT
    - Type: String
    - Required: No
    - Default: devops.cn-hangzhou.aliyuncs.com
    - Description: AccessKey request endpoint.

- PLUGIN_SOURCE_ORGANIZATIONID
    - Type: String
    - Required: No
    - Default: -
    - Description: Alibaba Cloud Yunxiao enterprise ID, can be found in Yunxiao access URL, e.g. https://devops.aliyun.com/organization/【OrganizationId】

- PLUGIN_CNB_URL
    - Type: String
    - Required: Yes
    - Default: https://cnb.cool
    - Description: CNB access URL

- PLUGIN_CNB_TOKEN
    - Type: String
    - Required: Yes
    - Default: -
    - Description: CNB authorization token, personal token - access token creation https://cnb.cool/profile/token.
    - Permission requirements:
        - repo-code read/write
        - repo-basic-info read-only
        - account-profile read-only
        - account-engage read-only
        - group-resource read/write
        - group-manage read/write
        - repo-content read/write

- PLUGIN_CNB_ROOT_ORGANIZATION
    - Type: String
    - Required: Yes
    - Default: -
    - Description: CNB root organization name after migration, ensure root organization is created in advance, do not include /
    - Example: root-group

## Additional Parameters

- PLUGIN_SOURCE_GROUP
    - Type: String
    - Required: No
    - Description: When migrating from CNB to CNB, specifies the root organization for migration

- PLUGIN_MIGRATE_TYPE
    - Type: String
    - Required: No
    - Default: team
    - Description: Migration type, supports project, repository, and team migration, only supported for CODING platform

- PLUGIN_MIGRATE_CONCURRENCY
    - Type: Number
    - Required: No
    - Default: 10
    - Description: Repository migration concurrency, maximum 10

- PLUGIN_MIGRATE_FORCE_PUSH
    - Type: Boolean
    - Required: No
    - Default: false
    - Description: ⚠️ Force push to CNB repository, CNB repository will be forcefully overwritten

- PLUGIN_MIGRATE_SKIP_EXISTS_REPO
    - Type: Boolean
    - Required: No
    - Default: true
    - Description: Skip repositories that already exist in CNB

- PLUGIN_MIGRATE_USE_LFS_MIGRATE
    - Type: String
    - Required: No
    - Default: false
    - Description: Whether to use lfs migrate to handle historical commits exceeding CNB single file size limit
    ⚠️ If enabled, commit IDs will not match source repository after migration

- PLUGIN_MIGRATE_ORGANIZATION_MAPPING_LEVEL
    - Type: String
    - Required: No
    - Default: 1
    - Description: CODING to CNB organization mapping relationship, only supported for CODING platform
    1: CODING projects map to CNB sub-organizations, repositories under sub-organizations
    2: CODING projects do not map to CNB sub-organizations, repositories directly under CNB root organization

- PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH
    - Type: String
    - Required: No
    - Default: false
    - Description: ⚠️ For repositories with missing LFS source files, ignore LFS file download errors and missing object errors during LFS push, continue pushing

- PLUGIN_MIGRATE_LOG_LEVEL
    - Type: String
    - Required: No
    - Default: info
    - Description: Log level (debug/info/warn/error)

- PLUGIN_MIGRATE_RELEASE
    - Type: Boolean
    - Required: No
    - Default: true
    - Description: Migrate releases (currently only supports gitlab/github/gitee release migration)

- PLUGIN_MIGRATE_FILE_LIMIT_SIZE
    - Type: Number
    - Required: No
    - Default: 500
    - Description: CNB maximum file size limit, unit MiB

- PLUGIN_MIGRATE_CODE
    - Type: Boolean
    - Required: Yes
    - Default: true
    - Description: Migrate code

- PLUGIN_MIGRATE_SSH
    - Type: Boolean
    - Required: Yes
    - Default: false
    - Description: Use SSH protocol to clone repositories from common third-party platforms

- PLUGIN_MIGRATE_REBASE
    - Type: Boolean
    - Required: Yes
    - Default: false
    - Description: When both source and target have changes, and CNB repository branch root directory has `.cnb.yml` file, use git rebase to ensure code synchronization, CNB pipeline configuration code will not be overwritten
    ⚠️ If enabled, force push will be enabled (PLUGIN_MIGRATE_FORCE_PUSH="true"), and CNB repository will be backed up in working directory

- PLUGIN_SOURCE_PROJECT
    - Type: String
    - Required: No
    - Default: -
    - Description: CODING project names to migrate (required when source_platform is coding and **migrate_type is project**), multiple projects separated by commas

- PLUGIN_MIGRATE_ALLOW_SELECT_REPOS
    - Type: Boolean
    - Required: No
    - Default: false
    - Description: Whether to allow user to select repositories for migration. When true, enables repo-path.txt selection feature.

## Common Issues

1. Exceeding single file size limit of 500 MiB
   Enable `PLUGIN_MIGRATE_USE_LFS_MIGRATE` parameter, see parameter description for details

2. Failed to get repository list: The current scope does not support access to this API
   Check if PLUGIN_SOURCE_TOKEN permissions meet requirements. For CODING platform, ensure token belongs to team owner or team admin.

3. LFS: Repository or object not found
   Enable `PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH`, see parameter description for details

4. Push failed: git pull before pushing again
   Enable `PLUGIN_MIGRATE_FORCE_PUSH` based on actual situation, see parameter description for details

5. How to migrate only specific repositories?
   Set `PLUGIN_MIGRATE_ALLOW_SELECT_REPOS=true`, after first run edit `repo-path.txt` to keep only desired repository paths.

6. How to reselect repositories for migration?
   Delete `repo-path.txt` file and run migration command again.

7. repo-path.txt not generated?
   Ensure `PLUGIN_MIGRATE_ALLOW_SELECT_REPOS=true` and write permissions are available.