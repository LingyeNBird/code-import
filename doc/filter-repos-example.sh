#!/bin/bash

# 按指定仓库迁移功能测试脚本
# 此脚本演示如何使用 source.repo 配置进行精确仓库迁移

set -e

echo "========================================="
echo "按指定仓库迁移功能测试"
echo "========================================="
echo ""

# 示例 1: GitLab 迁移单个仓库
echo "示例 1: GitLab 迁移单个仓库"
echo "----------------------------"
cat <<EOF
export PLUGIN_SOURCE_PLATFORM=gitlab
export PLUGIN_SOURCE_URL=https://gitlab.example.com
export PLUGIN_SOURCE_TOKEN=glpat-xxxxxxxxxxxx
export PLUGIN_SOURCE_REPO="group1/subgroup1/repo1"
export PLUGIN_MIGRATE_TYPE=team
export PLUGIN_CNB_URL=https://cnb.example.com
export PLUGIN_CNB_TOKEN=your-cnb-token
export PLUGIN_CNB_ROOT_ORGANIZATION=your-org

./ccrctl migrate
EOF
echo ""

# 示例 2: GitLab 迁移多个仓库
echo "示例 2: GitLab 迁移多个仓库"
echo "----------------------------"
cat <<EOF
export PLUGIN_SOURCE_PLATFORM=gitlab
export PLUGIN_SOURCE_URL=https://gitlab.example.com
export PLUGIN_SOURCE_TOKEN=glpat-xxxxxxxxxxxx
export PLUGIN_SOURCE_REPO="group1/repo1,group2/subgroup/repo2,group3/repo3"
export PLUGIN_MIGRATE_TYPE=team
export PLUGIN_CNB_URL=https://cnb.example.com
export PLUGIN_CNB_TOKEN=your-cnb-token
export PLUGIN_CNB_ROOT_ORGANIZATION=your-org

./ccrctl migrate
EOF
echo ""

# 示例 3: GitHub 迁移指定仓库
echo "示例 3: GitHub 迁移指定仓库"
echo "----------------------------"
cat <<EOF
export PLUGIN_SOURCE_PLATFORM=github
export PLUGIN_SOURCE_URL=https://github.com
export PLUGIN_SOURCE_TOKEN=ghp_xxxxxxxxxxxx
export PLUGIN_SOURCE_REPO="owner1/repo1,owner2/repo2"
export PLUGIN_MIGRATE_TYPE=team
export PLUGIN_CNB_URL=https://cnb.example.com
export PLUGIN_CNB_TOKEN=your-cnb-token
export PLUGIN_CNB_ROOT_ORGANIZATION=your-org

./ccrctl migrate
EOF
echo ""

# 示例 4: 工蜂迁移指定仓库
echo "示例 4: 工蜂迁移指定仓库"
echo "----------------------------"
cat <<EOF
export PLUGIN_SOURCE_PLATFORM=gongfeng
export PLUGIN_SOURCE_URL=https://git.code.tencent.com
export PLUGIN_SOURCE_TOKEN=your-token
export PLUGIN_SOURCE_REPO="tencent/WXG/project1/repo1,tencent/PCG/repo2"
export PLUGIN_MIGRATE_TYPE=team
export PLUGIN_CNB_URL=https://cnb.example.com
export PLUGIN_CNB_TOKEN=your-cnb-token
export PLUGIN_CNB_ROOT_ORGANIZATION=your-org

./ccrctl migrate
EOF
echo ""

# 示例 5: Gitee 迁移指定仓库
echo "示例 5: Gitee 迁移指定仓库"
echo "----------------------------"
cat <<EOF
export PLUGIN_SOURCE_PLATFORM=gitee
export PLUGIN_SOURCE_URL=https://gitee.com
export PLUGIN_SOURCE_TOKEN=your-token
export PLUGIN_SOURCE_REPO="owner1/repo1,owner2/repo2"
export PLUGIN_MIGRATE_TYPE=team
export PLUGIN_CNB_URL=https://cnb.example.com
export PLUGIN_CNB_TOKEN=your-cnb-token
export PLUGIN_CNB_ROOT_ORGANIZATION=your-org

./ccrctl migrate
EOF
echo ""

echo "========================================="
echo "配置文件示例 (config.yaml)"
echo "========================================="
cat <<EOF
source:
  platform: gitlab
  url: https://gitlab.example.com
  token: glpat-xxxxxxxxxxxx
  repo:
    - group1/subgroup1/repo1
    - group1/subgroup2/repo2
    - group2/repo3

cnb:
  url: https://cnb.example.com
  token: your-cnb-token
  root_organization: your-org

migrate:
  type: team
  concurrency: 5
  force_push: false
  ignore_lfs_not_found_error: true
  use_lfs_migrate: false
  organization_mapping_level: 1
EOF
echo ""

echo "========================================="
echo "功能说明"
echo "========================================="
echo "1. 当 source.repo 为空时，迁移所有仓库"
echo "2. 当 source.repo 不为空时，只迁移指定的仓库"
echo "3. 支持平台: GitLab、GitHub、Gitee、工蜂、CODING"
echo "4. 仓库路径格式:"
echo "   - GitHub/Gitee: owner/repo"
echo "   - GitLab: group/subgroup/repo"
echo "   - 工蜂: org/team/project/repo"
echo "   - CODING: team/project/repo"
echo "5. 多个仓库用逗号分隔"
echo "6. 仓库路径区分大小写"
echo ""
