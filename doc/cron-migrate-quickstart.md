# cron-migrate.sh 快速开始

## 三步完成配置

### 第一步：编辑脚本添加迁移命令

编辑 `scripts/cron-migrate.sh`，在 TODO 部分（第 52-63 行）添加 docker run 命令。

**从 CODING 迁移示例**（将 `xxx` 替换为实际值）：

```bash
docker run --rm \
  -e PLUGIN_SOURCE_TOKEN="xxx" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

> 💡 更多平台迁移示例请参考：[Docker使用文档](./docker-usage.md)

### 第二步：运行脚本

```bash
# 前台运行（推荐测试时使用）
./scripts/cron-migrate.sh

# 后台运行（推荐生产环境）
nohup ./scripts/cron-migrate.sh > migrate.log 2>&1 &
```

### 第三步：优雅停止

```bash
# 发送退出信号
kill -TERM $(pgrep -f cron-migrate.sh)

# 或使用 Ctrl+C（前台运行时）
```

## 核心特性

### ✅ 优雅退出
收到信号后会等待当前迁移任务完成，不会中断正在进行的迁移。

```
收到 SIGTERM/SIGINT
    ↓
设置退出标志
    ↓
当前迁移继续执行
    ↓
任务完成后退出
```

### ✅ 自动清理
每次循环后自动清理 120 分钟前创建的工作目录。

### ✅ 热更新
修改脚本后会自动检测并重新加载（注意：会中断当前任务）。

## 常用配置

### 1. 一次性迁移

```bash
# 编辑脚本添加迁移命令后，直接运行
./scripts/cron-migrate.sh

# 迁移完成后停止
kill -TERM $(pgrep -f cron-migrate.sh)
```

### 2. 定期增量迁移

在脚本末尾（第 71 行）取消注释休眠时间：

```bash
sleep 300  # 每5分钟执行一次
```

### 3. 迁移指定仓库

在迁移命令中添加仓库参数：

```bash
docker run --rm \
  -e PLUGIN_SOURCE_TOKEN="xxx" \
  -e PLUGIN_SOURCE_REPO="team/project/repo1,team/project/repo2" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx" \
  -e PLUGIN_MIGRATE_TYPE="repo" \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
```

## 查看日志

```bash
# 实时查看
tail -f migrate.log

# 搜索关键信息
grep "迁移完成" migrate.log
grep "收到退出信号" migrate.log
```

## 常见问题

**Q: 如何确认脚本在运行？**
```bash
ps aux | grep cron-migrate.sh
```

**Q: 退出需要多久？**

等待当前迁移任务完成，时间取决于仓库大小和网络速度。

**Q: 如何强制停止？**
```bash
kill -9 $(pgrep -f cron-migrate.sh)  # 不推荐
```

**Q: 支持哪些平台？**

CODING、GitHub、GitLab、Gitee、阿里云、华为云、CNB、工蜂等。详见[文档](./docker-usage.md)。

## 更多信息

- 详细文档：[cron-migrate-usage.md](./cron-migrate-usage.md)
- Docker 使用：[docker-usage.md](./docker-usage.md)
- 参数说明：[parameters.md](./parameters.md)
