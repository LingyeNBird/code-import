#!/bin/bash
# 正确闭合引号并转义变量
new="-e PLUGIN_SOURCE_REPO=\"$1,"
escaped_new=$(sed 's/[\/&]/\\&/g' <<< "$new")

# 使用双引号包裹 sed 命令以解析变量
sed -i "s|-e PLUGIN_SOURCE_REPO=\"|$escaped_new|" single/sync-single-repo.sh
sed -i "s|-e PLUGIN_SOURCE_REPO=\"|$escaped_new|" cron-sync-all-gitlab.sh