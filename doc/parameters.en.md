# Core Parameters
- **PLUGIN_SOURCE_TOKEN**
    - Type: string
    - Required: Yes
    - Default: -
    - Description: API token for source code hosting platform (required when source_platform is not common)
        - CODING permissions: **Only team owner or admin token**, User info-read only, Project info-read only, Repository-read only https://e.coding.net/user/account/setting/tokens
        - Github permissions: repo:all, read:org https://github.com/settings/tokens (classic token)
        - Gitlab permissions: read_api https://gitlab.com/-/user_settings/personal_access_tokens
        - Gitee permissions: user_info, projects https://gitee.com/profile/personal_access_tokens
        - Aliyun Codeup permissions: Repository:read only https://account-devops.aliyun.com/settings/personalAccessToken
        - CNB permissions: account-engage:r, group-resource:r https://cnb.cool/profile/token
        - Tencent Git permissions: api, read_repository https://git.woa.com/profile/account

- **PLUGIN_CNB_ROOT_ORGANIZATION**
    - Type: string
    - Required: Yes
    - Default: -
    - Description: Root organization name in CNB after migration. Ensure it's created in advance, don't include /
    - Ex: cnb

- **PLUGIN_CNB_TOKEN**
    - Type: string
    - Required: Yes
    - Default: -
    - Description: CNB access token, create at https://cnb.cool/profile/token
    - Permission requirements: For common scenarios select `Migration Tool Credentials`
    - Scope: All repositories/artifacts

- **PLUGIN_SOURCE_URL**
    - Type: string
    - Required: Yes
    - Default: https://e.coding.net
    - Description: Source code hosting platform URL
    - Ex:
       - Github: https://github.com
       - GitLab: https://gitlab.com
       - Gitee: https://gitee.com
       - common: https://common.com

- **PLUGIN_SOURCE_PLATFORM**
    - Type: string
    - Required: Yes
    - Default: coding
    - Description: Migration platform name, supports coding/gitlab/github/gitee/aliyun/cnb/gongfeng, other common platforms use common; local bare repositories use local

- **PLUGIN_SOURCE_REPO**
    - Type: string
    - Required: No
    - Default: -
    - Description: Repository paths, needs to combine with source_url to form complete clone URLs like https://common.com/group1/repo1  
  Required when source_platform is common or source_platform is coding and **migrate_type is repo**, multiple repositories separated by commas
    - Ex: group1/repo1,group1/repo2,group2/repo3

- **PLUGIN_SOURCE_USERNAME**
    - Type: string
    - Required: No
    - Default: -
    - Description: Required when source_platform is common, username for cloning repositories (must have access to all repositories)

- **PLUGIN_SOURCE_PASSWORD**
    - Type: string
    - Required: No
    - Default: -
    - Description: Required when source_platform is common, password for cloning repositories

- **PLUGIN_SOURCE_ORGANIZATIONID**
    - Type: string
    - Required: No
    - Default: -
    - Description: Aliyun Codeup organization ID, can be obtained from management console. [Codeup docs](https://help.aliyun.com/zh/yunxiao/user-guide/enterprise-basic-operations?spm=a2c4g.11186623.0.0.b86bb43fPj59Ic#caf924e65dbad)

- **PLUGIN_CNB_URL**
    - Type: string
    - Required: Yes
    - Default: https://cnb.cool
    - Description: CNB access URL

# More Parameters
- **PLUGIN_SOURCE_GROUP**
  - Type: string
  - Required: No
  - Description: When migrating from CNB to CNB, specifies repositories under root organization to migrate

- **PLUGIN_MIGRATE_TYPE**
    - Type: string
    - Required: No
    - Default: team
    - Description: Migration type, supports project/repo/team dimensions, only for CODING platform

- **PLUGIN_MIGRATE_CONCURRENCY**
    - Type: number
    - Required: No
    - Default: 5
    - Description: Repository migration concurrency, max 10

- **PLUGIN_MIGRATE_FORCE_PUSH**
    - Type: boolean
    - Required: No
    - Default: false
    - Description: ⚠️Force push to CNB repositories, will overwrite existing CNB repositories

- **PLUGIN_MIGRATE_SKIP_EXISTS_REPO**
    - Type: boolean
    - Required: No
    - Default: false
    - Description: Skip repositories that already exist in CNB

- **PLUGIN_MIGRATE_USE_LFS_MIGRATE**
    - Type: string
    - Required: No
    - Default: true
    - Description: Whether to use lfs migrate for commits exceeding CNB single file size limit
      ⚠️If enabled, commit IDs will differ from source repository

- **PLUGIN_MIGRATE_ORGANIZATION_MAPPING_LEVEL**
    - Type: string
    - Required: No
    - Default: 1
    - Description: CODING to CNB organization mapping, only for CODING platform
      1: CODING projects map to CNB sub-orgs, repositories under sub-orgs
      2: CODING projects don't map to CNB sub-orgs, repositories directly under root org

- **PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH**
    - Type: string
    - Required: No
    - Default: true
    - Description: ⚠️For repositories with missing LFS source files, ignore LFS download errors and continue push

- **PLUGIN_MIGRATE_LOG_LEVEL**
    - Type: string
    - Required: No
    - Default: info
    - Description: Log level (debug/info/warn/error)

- **PLUGIN_MIGRATE_RELEASE**
    - Type: boolean
    - Required: No
    - Default: false
    - Description: Migrate releases (currently only supports gitlab/github/gitee/coding release migration)

- **PLUGIN_MIGRATE_FILE_LIMIT_SIZE**
    - Type: number
    - Required: No
    - Default: 256
    - Description: CNB maximum file size limit in Mib

- **PLUGIN_MIGRATE_CODE**
    - Type: boolean
    - Required: Yes
    - Default: true
    - Description: Migrate code

- **PLUGIN_MIGRATE_SSH**
    - Type: boolean
    - Required: Yes
    - Default: false
    - Description: Use SSH protocol to clone common third-party platform repositories

- **PLUGIN_MIGRATE_REBASE**
    - Type: boolean
    - Required: Yes
    - Default: false
    - Description: When both source and target have changes, and CNB repository has `.cnb.yml` at root, use git rebase to preserve CNB pipeline config
      ⚠️If enabled, will force push (PLUGIN_MIGRATE_FORCE_PUSH="true") and backup CNB repository in working directory

- **PLUGIN_SOURCE_PROJECT**
    - Type: string
    - Required: No
    - Default: -
    - Description: CODING project names to migrate (required when source_platform is coding and **migrate_type is project**), multiple projects separated by commas

- **PLUGIN_MIGRATE_ALLOW_SELECT_REPOS**
    - Type: boolean
    - Required: No
    - Default: false
    - Description: Whether to allow selecting repositories to migrate. When true, enables repo-path.txt selection feature.

- **PLUGIN_MIGRATE_DOWNLOAD_ONLY**
    - Type: boolean
    - Required: No
    - Default: false
    - Description: Whether to only download repositories without migration. When true, only clones repositories locally without pushing to CNB. No CNB config required in this mode.

- **PLUGIN_MIGRATE_EXCLUDE_GITHUB_FORK**
  - Type: boolean
  - Required: No
  - Default: false
  - Description: Exclude GitHub fork repositories from migration