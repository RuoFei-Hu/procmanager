# procmanager - 通用进程管理工具

`procmanager` 是一个功能强大的命令行工具，用于管理程序的完整生命周期。它提供了启动、停止、重启和查看程序状态的功能，支持守护进程模式和日志管理。

## 功能特点

- 启动程序，支持普通模式和守护进程模式
- 停止正在运行的程序，支持优雅终止和强制终止
- 重启程序，保持原有的配置和参数
- 查看程序的运行状态
- 通过PID文件跟踪进程
- 支持日志文件配置
- 支持自定义超时时间

## 安装

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/RuoFei-Hu/procmanager.git
cd procmanager

# 编译
go build -o procmanager

# 安装到系统路径（可选）
sudo mv procmanager /usr/local/bin/
```

## 基本用法

```bash
# 查看帮助信息
procmanager --help

# 查看特定命令的帮助信息
procmanager start --help
procmanager stop --help
procmanager restart --help
procmanager status --help
```

## 命令详解

### 启动程序 (start)

```bash
procmanager start [程序路径] [参数...] --pid [PID文件路径] [选项]
```

#### 选项

- `-p, --pid <文件路径>`: 指定PID文件路径，用于跟踪进程
- `-d, --daemon`: 以守护进程方式在后台运行
- `-l, --log <文件路径>`: 指定日志文件路径，用于记录程序输出
- `-c, --config <文件路径>`: 指定配置文件路径（全局选项）

#### 示例

```bash
# 启动一个程序
procmanager start /path/to/your/program --pid /var/run/program.pid

# 启动程序并传递参数
procmanager start /path/to/your/program arg1 arg2 --pid /var/run/program.pid

# 以守护进程方式启动并指定日志文件
procmanager start /path/to/your/program --daemon --log /var/log/program.log --pid /var/run/program.pid
```

### 停止程序 (stop)

```bash
procmanager stop --pid [PID文件路径] [选项]
```

#### 选项

- `-p, --pid <文件路径>`: 指定PID文件路径（必需）
- `-f, --force`: 强制终止程序（发送SIGKILL信号）
- `-t, --timeout <秒数>`: 等待程序终止的超时时间，默认为10秒
- `-c, --config <文件路径>`: 指定配置文件路径（全局选项）

#### 示例

```bash
# 停止程序
procmanager stop --pid /var/run/program.pid

# 强制停止程序
procmanager stop --pid /var/run/program.pid --force

# 设置较长的超时时间
procmanager stop --pid /var/run/program.pid --timeout 30
```

### 重启程序 (restart)

```bash
procmanager restart [程序路径] [参数...] --pid [PID文件路径] [选项]
```

#### 选项

- `-p, --pid <文件路径>`: 指定PID文件路径（必需）
- `-d, --daemon`: 以守护进程方式在后台运行
- `-l, --log <文件路径>`: 指定日志文件路径
- `-t, --timeout <秒数>`: 等待程序终止的超时时间，默认为10秒
- `-c, --config <文件路径>`: 指定配置文件路径（全局选项）

#### 示例

```bash
# 重启程序
procmanager restart /path/to/your/program --pid /var/run/program.pid

# 重启程序并传递参数
procmanager restart /path/to/your/program arg1 arg2 --pid /var/run/program.pid

# 以守护进程方式重启并指定日志文件
procmanager restart /path/to/your/program --daemon --log /var/log/program.log --pid /var/run/program.pid
```

### 查看程序状态 (status)

```bash
procmanager status --pid [PID文件路径] [选项]
```

#### 选项

- `-p, --pid <文件路径>`: 指定PID文件路径（必需）
- `-c, --config <文件路径>`: 指定配置文件路径（全局选项）

#### 示例

```bash
# 查看程序状态
procmanager status --pid /var/run/program.pid
```

## 最佳实践

### PID文件管理

PID文件是跟踪进程的关键。建议将PID文件存放在以下位置：

- 系统级服务：`/var/run/<程序名>.pid`
- 用户级服务：`$HOME/.local/run/<程序名>.pid` 或 `$HOME/.cache/<程序名>.pid`

### 日志文件管理

日志文件对于排查问题至关重要。建议将日志文件存放在以下位置：

- 系统级服务：`/var/log/<程序名>.log`
- 用户级服务：`$HOME/.local/log/<程序名>.log` 或 `$HOME/.cache/<程序名>.log`

### 守护进程模式

使用守护进程模式（`--daemon`）时，程序将在后台运行，不会受到终端关闭的影响。这对于长期运行的服务特别有用。

## 常见问题

### Q: 如何确认程序是否正在运行？

A: 使用 `status` 命令检查：

```bash
procmanager status --pid /path/to/your/program.pid
```

### Q: 如何优雅地停止程序？

A: 默认情况下，`stop` 命令会先发送 SIGTERM 信号，给程序一个清理资源的机会，然后等待程序终止。如果在超时时间内程序未终止，可以使用 `--force` 选项强制终止。

### Q: PID文件已存在但程序未运行怎么办？

A: 使用 `status` 命令会自动清理过时的PID文件。或者手动删除PID文件后重新启动程序。

### Q: 如何在系统启动时自动启动程序？

A: 可以创建一个系统服务（如systemd服务）或添加到启动脚本中。例如，创建systemd服务：

```ini
[Unit]
Description=My Program Service
After=network.target

[Service]
ExecStart=/usr/local/bin/procmanager start /path/to/your/program --pid /var/run/program.pid
ExecStop=/usr/local/bin/procmanager stop --pid /var/run/program.pid
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交问题报告和拉取请求！