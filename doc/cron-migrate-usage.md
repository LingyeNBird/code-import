# cron-migrate.sh 使用说明

## 功能概述

`cron-migrate.sh` 是一个支持优雅退出的循环迁移脚本，主要特性：

1. **循环执行**：持续执行迁移任务
2. **优雅退出**：收到信号后等待当前迁移完成再退出
3. **热更新**：支持脚本自动热更新
4. **自动清理**：定期清理旧的工作目录

## 快速使用

### 1. 添加迁移命令

编辑 `scripts/cron-migrate.sh`，在 TODO 部分添加 docker run 迁移命令。

以从 CODING 迁移为例（将 `xxx` 替换为实际参数值）：

```bash
# 从 CODING 迁移团队下所有仓库
docker run --rm \
  -e PLUGIN_SOURCE_TOKEN="xxx" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

更多迁移平台示例请参考：[Docker使用文档](./docker-usage.md)

### 2. 运行脚本

```bash
# 直接运行
./scripts/cron-migrate.sh

# 后台运行
nohup ./scripts/cron-migrate.sh > migrate.log 2>&1 &
```

### 3. 优雅停止脚本

```bash
# 方式1：发送 SIGTERM 信号
kill -TERM $(pgrep -f cron-migrate.sh)

# 方式2：使用 Ctrl+C（前台运行时）
# 在运行脚本的终端按 Ctrl+C
```

## 工作原理

### 执行流程

```
启动脚本
  ↓
注册信号处理函数
  ↓
进入主循环 ←─────┐
  ↓              │
检查退出标志      │
  ↓              │
检查脚本更新      │
  ↓              │
创建工作目录      │
  ↓              │
执行迁移命令（阻塞等待完成）
  ↓              │
清理旧目录        │
  ↓              │
继续循环 ─────────┘
  ↓
退出
```

### 优雅退出机制

1. **信号捕获**：脚本捕获 `SIGTERM` 和 `SIGINT` 信号
2. **设置标志**：收到信号后设置 `GRACEFUL_SHUTDOWN=1`
3. **等待完成**：当前迁移命令继续执行直到完成
4. **检查退出**：下次循环开始时检查标志并退出
5. **正常清理**：执行清理逻辑后正常退出

## 关键参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `ROOT_WORK_DIR` | `/code-import` | 工作目录根路径 |
| `GRACEFUL_SHUTDOWN` | `0` | 优雅退出标志（0=运行中，1=准备退出） |
| 工作目录清理时间 | `120分钟` | 保留最近120分钟的工作目录 |

## 使用场景

### 场景1：一次性批量迁移

编辑脚本，添加迁移命令：

```bash
# 从 CODING 迁移
docker run --rm \
  -e PLUGIN_SOURCE_TOKEN="your_token" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="your_org" \
  -e PLUGIN_CNB_TOKEN="your_cnb_token" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

运行脚本：

```bash
./scripts/cron-migrate.sh

# 迁移完成后停止
kill -TERM $(pgrep -f cron-migrate.sh)
```

### 场景2：定期增量迁移

编辑脚本，添加迁移命令并设置休眠时间：

```bash
# 在 TODO 部分添加迁移命令
docker run --rm \
  -e PLUGIN_SOURCE_TOKEN="your_token" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="your_org" \
  -e PLUGIN_CNB_TOKEN="your_cnb_token" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import

# 在循环末尾取消注释休眠时间
sleep 3600  # 每小时执行一次
```

后台运行：

```bash
nohup ./scripts/cron-migrate.sh > migrate.log 2>&1 &
```

### 场景3：迁移指定仓库

```bash
# 从 CODING 迁移指定仓库
docker run --rm \
  -e PLUGIN_SOURCE_TOKEN="your_token" \
  -e PLUGIN_SOURCE_REPO="<TEAM-NAME>/<PROJECT-NAME>/<REPO-NAME>,test-team/project1/repoA" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="your_org" \
  -e PLUGIN_CNB_TOKEN="your_cnb_token" \
  -e PLUGIN_MIGRATE_TYPE="repo" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

## 日志查看

```bash
# 实时查看日志（如果使用了 nohup）
tail -f migrate.log

# 查看完整日志
cat migrate.log

# 搜索关键信息
grep "迁移完成" migrate.log
grep "收到退出信号" migrate.log
```

## 脚本热更新

脚本支持热更新，修改后会自动重新加载：

```bash
# 1. 脚本正在运行中
./scripts/cron-migrate.sh &

# 2. 修改脚本
vim scripts/cron-migrate.sh

# 3. 脚本会自动检测变化并重新加载
# 输出：检测到脚本更新，重新加载新版本...
```

**注意**：热更新会立即重启脚本，当前正在执行的迁移任务会被中断。建议等待迁移完成后再更新。

## 常见问题

### Q1: 如何确认脚本正在运行？

```bash
# 查看进程
ps aux | grep cron-migrate.sh

# 查看日志
tail -f migrate.log
```

### Q2: 脚本收到信号后多久退出？

脚本会等待当前迁移命令完成后退出，退出时间取决于：
- 单次迁移任务的耗时
- 仓库数量和大小
- 网络速度

### Q3: 如何立即强制退出？

```bash
# 强制终止（不推荐，可能导致数据不完整）
kill -9 $(pgrep -f cron-migrate.sh)
```

### Q4: 工作目录清理规则是什么？

脚本会在每次循环后清理 120 分钟前创建的 `workdir-*` 目录。修改清理时间：

```bash
# 修改为 60 分钟
find . -maxdepth 1 -type d -name 'workdir-*' -mmin +60 -exec rm -rf {} \;

# 修改为 240 分钟
find . -maxdepth 1 -type d -name 'workdir-*' -mmin +240 -exec rm -rf {} \;
```

### Q5: 支持哪些迁移平台？

支持以下平台的迁移（详见 [Docker使用文档](./docker-usage.md)）：
- CODING
- GitHub
- GitLab
- Gitee
- Gitea
- 阿里云 Codeup
- 华为云 CodeArts Repo
- CNB
- 腾讯工蜂
- 通用第三方代码托管平台

## 最佳实践

1. **监控日志**：使用日志收集工具（如 ELK、Loki）监控迁移进度
2. **测试验证**：先在测试环境验证脚本配置
3. **备份数据**：重要数据在迁移前先备份
4. **分批迁移**：对于大量仓库，建议分批进行
5. **错误处理**：迁移失败时查看日志，根据错误信息调整配置

## 完整配置示例

### 从 CODING 迁移所有仓库

```bash
#!/bin/bash
set -ex

SCRIPT_PATH=$(realpath "$0")
INITIAL_HASH=$(md5sum "$SCRIPT_PATH" | awk '{print $1}')
ROOT_WORK_DIR="/code-import"
GRACEFUL_SHUTDOWN=0

graceful_exit() {
  echo "收到退出信号，设置优雅退出标志..."
  GRACEFUL_SHUTDOWN=1
}

trap graceful_exit SIGTERM SIGINT

while true; do
  if [[ $GRACEFUL_SHUTDOWN -eq 1 ]]; then
    echo "优雅退出标志已设置，等待当前迁移任务完成后退出"
    break
  fi

  CURRENT_HASH=$(md5sum "$SCRIPT_PATH" | awk '{print $1}')
  if [[ "$CURRENT_HASH" != "$INITIAL_HASH" ]]; then
     echo "检测到脚本更新，重新加载新版本..."
     exec "$SCRIPT_PATH"
  fi

  cd ${ROOT_WORK_DIR}
  date=$(date +"%Y%m%d%H%M%S")
  workDir="workdir-${date}"
  mkdir ${workDir}
  cd ${workDir}
  
  # 执行迁移任务
  docker run --rm \
    -e PLUGIN_SOURCE_TOKEN="your_coding_token" \
    -e PLUGIN_CNB_ROOT_ORGANIZATION="your_cnb_org" \
    -e PLUGIN_CNB_TOKEN="your_cnb_token" \
    -v $(pwd):$(pwd) -w $(pwd) \
    cnbcool/code-import
  
  cd ${ROOT_WORK_DIR}
  find . -maxdepth 1 -type d -name 'workdir-*' -mmin +120 -exec rm -rf {} \;
  
  # 可选：每次迁移后休眠
  # sleep 300
done

echo "迁移任务已完成，脚本正常退出"
```

## 总结

`cron-migrate.sh` 提供了一个简单但强大的循环迁移解决方案：

- ✅ 支持优雅退出，确保数据完整性
- ✅ 自动清理旧文件，节省磁盘空间
- ✅ 支持热更新，方便调试和维护
- ✅ 支持多种迁移平台

通过合理配置 docker run 命令，可以满足各种迁移场景的需求。详细的迁移参数配置请参考 [Docker使用文档](./docker-usage.md) 和 [参数说明](./parameters.md)。


## 工作原理

### 执行流程

```
启动脚本
  ↓
注册信号处理函数
  ↓
进入主循环 ←─────┐
  ↓              │
检查退出标志      │
  ↓              │
检查脚本更新      │
  ↓              │
创建工作目录      │
  ↓              │
执行迁移命令（阻塞等待完成）
  ↓              │
清理旧目录        │
  ↓              │
继续循环 ─────────┘
  ↓
退出
```

### 优雅退出机制

1. **信号捕获**：脚本捕获 `SIGTERM` 和 `SIGINT` 信号
2. **设置标志**：收到信号后设置 `GRACEFUL_SHUTDOWN=1`
3. **等待完成**：当前迁移命令继续执行直到完成
4. **检查退出**：下次循环开始时检查标志并退出
5. **正常清理**：执行清理逻辑后正常退出

## 关键参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `ROOT_WORK_DIR` | `/code-import` | 工作目录根路径 |
| `GRACEFUL_SHUTDOWN` | `0` | 优雅退出标志（0=运行中，1=准备退出） |
| 工作目录清理时间 | `120分钟` | 保留最近120分钟的工作目录 |

## 使用场景

### 场景1：一次性批量迁移

编辑脚本，添加迁移命令：

```bash
# 从 CODING 迁移
docker run --rm \
  -e PLUGIN_SOURCE_TOKEN="your_token" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="your_org" \
  -e PLUGIN_CNB_TOKEN="your_cnb_token" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

运行脚本：

```bash
./scripts/cron-migrate.sh

# 迁移完成后停止
kill -TERM $(pgrep -f cron-migrate.sh)
```

### 场景2：定期增量迁移

编辑脚本，添加迁移命令并设置休眠时间：

```bash
# 在 TODO 部分添加迁移命令
docker run --rm \
  -e PLUGIN_SOURCE_TOKEN="your_token" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="your_org" \
  -e PLUGIN_CNB_TOKEN="your_cnb_token" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import

# 在循环末尾取消注释休眠时间
sleep 3600  # 每小时执行一次
```

后台运行：

```bash
nohup ./scripts/cron-migrate.sh > migrate.log 2>&1 &
```

## 日志查看

```bash
# 实时查看日志（如果使用了 nohup）
tail -f migrate.log

# 查看完整日志
cat migrate.log

# 搜索关键信息
grep "迁移完成" migrate.log
grep "收到退出信号" migrate.log
```

## 脚本热更新

脚本支持热更新，修改后会自动重新加载：

```bash
# 1. 脚本正在运行中
./scripts/cron-migrate.sh &

# 2. 修改脚本
vim scripts/cron-migrate.sh

# 3. 脚本会自动检测变化并重新加载
# 输出：检测到脚本更新，重新加载新版本...
```

**注意**：热更新会立即重启脚本，当前正在执行的迁移任务会被中断。建议等待迁移完成后再更新。

## 常见问题

### Q1: 如何确认脚本正在运行？

```bash
# 查看进程
ps aux | grep cron-migrate.sh

# 查看日志
tail -f migrate.log
```

### Q2: 脚本收到信号后多久退出？

脚本会等待当前迁移命令完成后退出，退出时间取决于：
- 单次迁移任务的耗时
- 仓库数量和大小
- 网络速度

### Q3: 如何立即强制退出？

```bash
# 强制终止（不推荐，可能导致数据不完整）
kill -9 $(pgrep -f cron-migrate.sh)
```

### Q4: 工作目录清理规则是什么？

脚本会在每次循环后清理 120 分钟前创建的 `workdir-*` 目录。修改清理时间：

```bash
# 修改为 60 分钟
find . -maxdepth 1 -type d -name 'workdir-*' -mmin +60 -exec rm -rf {} \;

# 修改为 240 分钟
find . -maxdepth 1 -type d -name 'workdir-*' -mmin +240 -exec rm -rf {} \;
```

### Q5: 支持哪些迁移平台？

支持以下平台的迁移（详见 [Docker使用文档](./docker-usage.md)）：
- CODING
- GitHub
- GitLab
- Gitee
- Gitea
- 阿里云 Codeup
- 华为云 CodeArts Repo
- CNB
- 腾讯工蜂
- 通用第三方代码托管平台

## 最佳实践

1. **监控日志**：使用日志收集工具（如 ELK、Loki）监控迁移进度
2. **资源限制**：在容器中运行时设置合理的资源限制
3. **错误处理**：迁移命令应有适当的错误处理和重试机制
4. **备份数据**：重要数据在迁移前先备份
5. **测试验证**：先在测试环境验证脚本配置

## 示例配置

### 完整的迁移命令示例

```bash
# 进入工作目录
cd ${ROOT_WORK_DIR}
date=$(date +"%Y%m%d%H%M%S")
workDir="workdir-${date}"
mkdir ${workDir}
cd ${workDir}

# 执行迁移（会等待完成）
docker run --rm \
  -v ${ROOT_WORK_DIR}/config.yaml:/app/config.yaml:ro \
  -v ${ROOT_WORK_DIR}/${workDir}:/app/workspace \
  -e TZ=Asia/Shanghai \
  your-migrate-image:latest

# 迁移完成后的处理
echo "迁移任务 ${workDir} 已完成"

# 清理旧目录
cd ${ROOT_WORK_DIR}
find . -maxdepth 1 -type d -name 'workdir-*' -mmin +120 -exec rm -rf {} \;
```

## 总结

`cron-migrate.sh` 提供了一个简单但强大的循环迁移解决方案：

- ✅ 支持优雅退出，确保数据完整性
- ✅ 自动清理旧文件，节省磁盘空间
- ✅ 支持热更新，方便调试和维护
- ✅ 适用于多种运行环境（本地、Docker、K8s）

通过合理配置，可以满足各种迁移场景的需求。
