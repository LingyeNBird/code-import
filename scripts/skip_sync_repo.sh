#!/bin/bash
set -e

# 日志函数，同时输出到终端和日志文件
WORKDIR=$(ls -lt | grep workdir- | head -1 | awk '{print $9}')
WORKDIR_LOG_FILE=./$WORKDIR/successful.log
LOG_FILE="skip_sync_repo.log"
log() {
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "${timestamp} ${message}" | tee -a "${LOG_FILE}"
    echo "${timestamp} ${message}" | tee -a "${WORKDIR_LOG_FILE}"
     
}


# 检查参数是否存在
if [ $# -eq 0 ]; then
    log "错误：未提供 repoPath 参数，脚本终止。"
    exit 1
fi


for REPO_PATH in "$@";do
    log "${REPO_PATH}"
done



# 检查 sed 命令状态
if [ $? -eq 0 ]; then
    echo "OK"
else
    echo "错误：文件修改失败，请检查权限或语法。"
    exit 1
fi


echo  "========== 脚本执行完成 =========="