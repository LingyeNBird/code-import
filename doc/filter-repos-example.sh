#!/bin/bash

# 按指定仓库迁移功能示例脚本（Docker 环境变量方式）
# 此脚本演示如何使用 Docker 容器和环境变量进行精确仓库迁移

set -e

echo "========================================="
echo "按指定仓库迁移功能示例（Docker）"
echo "========================================="
echo ""

# 示例 1: GitLab 迁移单个仓库
echo "示例 1: GitLab 迁移单个仓库"
echo "----------------------------"
cat <<EOF
docker run --rm \\
  -e PLUGIN_SOURCE_PLATFORM=gitlab \\
  -e PLUGIN_SOURCE_URL=https://gitlab.example.com \\
  -e PLUGIN_SOURCE_TOKEN=glpat-xxxxxxxxxxxx \\
  -e PLUGIN_SOURCE_REPO="group1/subgroup1/repo1" \\
  -e PLUGIN_MIGRATE_TYPE=team \\
  -e PLUGIN_CNB_URL=https://cnb.example.com \\
  -e PLUGIN_CNB_TOKEN=your-cnb-token \\
  -e PLUGIN_CNB_ROOT_ORGANIZATION=your-org \\
  cnbcool/code-import:latest
EOF
echo ""

# 示例 2: GitLab 迁移多个仓库
echo "示例 2: GitLab 迁移多个仓库"
echo "----------------------------"
cat <<EOF
docker run --rm \\
  -e PLUGIN_SOURCE_PLATFORM=gitlab \\
  -e PLUGIN_SOURCE_URL=https://gitlab.example.com \\
  -e PLUGIN_SOURCE_TOKEN=glpat-xxxxxxxxxxxx \\
  -e PLUGIN_SOURCE_REPO="group1/repo1,group2/subgroup/repo2,group3/repo3" \\
  -e PLUGIN_MIGRATE_TYPE=team \\
  -e PLUGIN_CNB_URL=https://cnb.example.com \\
  -e PLUGIN_CNB_TOKEN=your-cnb-token \\
  -e PLUGIN_CNB_ROOT_ORGANIZATION=your-org \\
  cnbcool/code-import:latest
EOF
echo ""

# 示例 3: GitHub 迁移指定仓库
echo "示例 3: GitHub 迁移指定仓库"
echo "----------------------------"
cat <<EOF
docker run --rm \\
  -e PLUGIN_SOURCE_PLATFORM=github \\
  -e PLUGIN_SOURCE_URL=https://github.com \\
  -e PLUGIN_SOURCE_TOKEN=ghp_xxxxxxxxxxxx \\
  -e PLUGIN_SOURCE_REPO="owner1/repo1,owner2/repo2" \\
  -e PLUGIN_MIGRATE_TYPE=team \\
  -e PLUGIN_CNB_URL=https://cnb.example.com \\
  -e PLUGIN_CNB_TOKEN=your-cnb-token \\
  -e PLUGIN_CNB_ROOT_ORGANIZATION=your-org \\
  cnbcool/code-import:latest
EOF
echo ""

# 示例 4: 工蜂迁移指定仓库
echo "示例 4: 工蜂迁移指定仓库"
echo "----------------------------"
cat <<EOF
docker run --rm \\
  -e PLUGIN_SOURCE_PLATFORM=gongfeng \\
  -e PLUGIN_SOURCE_URL=https://git.code.tencent.com \\
  -e PLUGIN_SOURCE_TOKEN=your-token \\
  -e PLUGIN_SOURCE_REPO="tencent/WXG/project1/repo1,tencent/PCG/repo2" \\
  -e PLUGIN_MIGRATE_TYPE=team \\
  -e PLUGIN_CNB_URL=https://cnb.example.com \\
  -e PLUGIN_CNB_TOKEN=your-cnb-token \\
  -e PLUGIN_CNB_ROOT_ORGANIZATION=your-org \\
  cnbcool/code-import:latest
EOF
echo ""

# 示例 5: Gitee 迁移指定仓库
echo "示例 5: Gitee 迁移指定仓库"
echo "----------------------------"
cat <<EOF
docker run --rm \\
  -e PLUGIN_SOURCE_PLATFORM=gitee \\
  -e PLUGIN_SOURCE_URL=https://gitee.com \\
  -e PLUGIN_SOURCE_TOKEN=your-token \\
  -e PLUGIN_SOURCE_REPO="owner1/repo1,owner2/repo2" \\
  -e PLUGIN_MIGRATE_TYPE=team \\
  -e PLUGIN_CNB_URL=https://cnb.example.com \\
  -e PLUGIN_CNB_TOKEN=your-cnb-token \\
  -e PLUGIN_CNB_ROOT_ORGANIZATION=your-org \\
  cnbcool/code-import:latest
EOF
echo ""



echo "========================================="
echo "功能说明"
echo "========================================="
echo "1. 使用 Docker 容器运行，通过环境变量配置"
echo "2. 镜像名称: cnbcool/code-import:latest"
echo "3. 当 PLUGIN_SOURCE_REPO 为空时，迁移所有仓库"
echo "4. 当 PLUGIN_SOURCE_REPO 不为空时，只迁移指定的仓库"
echo "5. 支持平台: GitLab、GitHub、Gitee、工蜂、CODING"
echo "6. 仓库路径格式:"
echo "   - GitHub/Gitee: owner/repo"
echo "   - GitLab: group/subgroup/repo"
echo "   - 工蜂: org/team/project/repo"
echo "   - CODING: team/project/repo"
echo "7. 多个仓库用逗号分隔"
echo "8. 仓库路径区分大小写"
echo ""
