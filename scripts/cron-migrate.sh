#!/bin/bash
# 设置脚本执行选项：
# -e: 遇到错误立即退出
# -x: 打印每个命令及其参数
set -ex

# 脚本自检机制：通过MD5哈希值检测脚本是否被更新
# 获取脚本的绝对路径
SCRIPT_PATH=$(realpath "$0")
# 记录脚本初始状态的MD5哈希值
INITIAL_HASH=$(md5sum "$SCRIPT_PATH" | awk '{print $1}')
# 定义工作目录根路径
ROOT_WORK_DIR="/code-import"

# 主循环：持续执行迁移任务
while true; do
  # 脚本热更新检测
  # 计算当前脚本的MD5哈希值
  CURRENT_HASH=$(md5sum "$SCRIPT_PATH" | awk '{print $1}')
  # 如果检测到脚本被更新，则重新执行新版本
  if [[ "$CURRENT_HASH" != "$INITIAL_HASH" ]]; then
     echo "检测到脚本更新，重新加载新版本..."
     exec "$SCRIPT_PATH"  # 使用exec替换当前进程，实现无缝热更新
  fi

  # 进入工作目录
  cd ${ROOT_WORK_DIR}
  # 生成基于时间戳的工作目录名称
  date=$(date +"%Y%m%d%H%M%S")
  workDir="workdir-${date}"
  # 创建新的工作目录
  mkdir ${workDir}
  
  # TODO: 在这里添加具体的迁移命令
  # 补充迁移的docker run 命令
  
  # 清理旧的工作目录
  # 查找并删除120分钟前创建的所有workdir-*目录
  cd ${ROOT_WORK_DIR}
  find . -maxdepth 1 -type d -name 'workdir-*' -mmin +120 -exec rm -rf {} \;
  
  # 可选：添加休眠时间
  #sleep 300  # 休眠5分钟
done