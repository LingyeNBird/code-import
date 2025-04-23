# CNB Code Import (Batch Migration Tool for CNB Repositories)

![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/git-clone-yyds)
![badge](https://cnb.cool/cnb/plugins/cnbcool/code-import/-/badge/git/latest/ci/status/push)

## Features
1. Supports batch migration of repositories from CODING, GitHub, GitLab, Gitee, Alibaba Cloud, and other third-party platforms to CNB
2. Automatically skips successfully migrated repositories (⚠️ depends on `successful.log` file in working directory)

## Using with Docker

### Notes
1. ⚠️ Ensure the CNB root organization exists before migration  
2. Replace `xxx` fields with your actual values (see Core Parameters section for details)  
3. Cloud Native Development includes docker commands for faster internal execution

### Migrating from CODING

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
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
        - CODING permissions: Read-only for user info, project info, and repositories
        - GitHub permissions: repo:all, read:org
        - GitLab permissions: read_api
        - Gitee permissions: user_info, projects

[Additional parameters and troubleshooting guide continue...]