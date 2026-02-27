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
        - Gitea permissions: read:organization, read:repository, read:user http(s)://<YOUR-GITEA-HOST>/user/settings/applications (Settings-Applications-Generate New Token)
        - Huawei Cloud CodeArts Repo permissions: Repository read/write https://devcloud.{YOUR-REGION}.huaweicloud.com/codehub/tokens, replace YOUR-REGION with corresponding region, e.g. cn-north-4 for North China-Beijing4, see docs https://support.huaweicloud.com/api-codeartsrepo/codeartsrepo_05_0001.html#section0
- **PLUGIN_SOURCE_AK**
    - Type: string
    - Required: No
    - Default: -
    - Description: Access key for source code hosting platform API (required when source_platform is huaweicloud), create at: Console-Personal Info-My Credentials-Access Keys
- **PLUGIN_SOURCE_SK**
    - Type: string
    - Required: No
    - Default: -
    - Description: Secret key for source code hosting platform API (required when source_platform is huaweicloud), create at: Console-Personal Info-My Credentials-Access Keys
- **PLUGIN_SOURCE_REGION**
    - Type: string
    - Required: No
    - Default: cn-north-4
    - Description: Huawei Cloud CodeArts Repo region code, see https://support.huaweicloud.com/api-codeartsrepo/codeartsrepo_05_0001.html#section0

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
    - Description: Migration platform name, supports coding/gitlab/github/gitee/aliyun/cnb/gongfeng/gitea/huaweicloud, other common platforms use common; local bare repositories use local

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
    - Description: ⚠️Force push to CNB repositories, will overwrite existing CNB repositories. When enabled, a WARN-level warning message will be logged to alert users about the risks of this operation.

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
    - Description: Source repository path to CNB organization mapping relationship, generally keep default.  
      1: Migrated repository path will be `<CNB root org>/<source repo path>`, e.g. if source repo path is `group1/repo1`, migrated CNB repo path will be `<CNB root org>/group1/repo1`, will auto-create sub-organizations.  
      2: Migrated repository path will be `<CNB root org>/<source repo name>`, e.g. if source repo path is `group1/repo1`, migrated CNB repo path will be `<CNB root org>/repo1`, will not auto-create sub-organizations. ⚠️Cannot have repositories with same names, otherwise will conflict and error

- **PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH**
    - Type: string
    - Required: No
    - Default: false
    - Description: ⚠️For repositories with missing LFS source files, ignore LFS file download errors and missing object errors during LFS push, continue pushing

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

- **PLUGIN_MIGRATE_RELEASE_TAG**
    - Type: string
    - Required: No
    - Default: empty
    - Description: Sync only the release with this tag (for example `v1.0.1`). If not set, only the latest release is synced.

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
    - Description: Use git rebase to ensure code sync, CNB-side committed pipeline config code will not be overwritten, conflicts need manual resolution.  
      ⚠️If enabled, will enable force push (PLUGIN_MIGRATE_FORCE_PUSH="true") and backup CNB-side code repository in migration tool working directory

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
