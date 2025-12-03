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
# 优雅退出标志
GRACEFUL_SHUTDOWN=0

# 优雅退出信号处理函数
graceful_exit() {
  echo "收到退出信号，设置优雅退出标志..."
  GRACEFUL_SHUTDOWN=1
}

# 注册信号处理函数
trap graceful_exit SIGTERM SIGINT

# 主循环：持续执行迁移任务
while true; do
  # 检查是否需要退出
  if [[ $GRACEFUL_SHUTDOWN -eq 1 ]]; then
    echo "优雅退出标志已设置，等待当前迁移任务完成后退出"
    break
  fi

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
  # 创建新的工作目录并进入
  mkdir ${workDir}
  cd ${workDir}
  
  # TODO: 在这里添加具体的迁移命令
  # 补充迁移的 docker run 命令（请将 xxx 替换为实际参数值）
  # 
  # 从 CODING 迁移示例：
  # docker run --rm \
  #   -e PLUGIN_SOURCE_TOKEN="xxx" \
  #   -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  #   -e PLUGIN_CNB_TOKEN="xxx" \
  #   -v $(pwd):$(pwd) -w $(pwd) \
  #   cnbcool/code-import
  #
  # 注意：命令会阻塞等待执行完成，确保迁移任务完成后才继续下一次循环
  
  # 清理旧的工作目录
  # 查找并删除120分钟前创建的所有workdir-*目录
  cd ${ROOT_WORK_DIR}
  find . -maxdepth 1 -type d -name 'workdir-*' -mmin +120 -exec rm -rf {} \;
  
  # 可选：添加休眠时间
  #sleep 300  # 休眠5分钟
done

echo "迁移任务已完成，脚本正常退出"