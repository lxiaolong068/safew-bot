# 部署脚本说明

本目录包含 SafeW Bot 在不同环境下的部署脚本。

## 🚀 推荐方式：本地编译部署

**最新推荐方式**：在本地编译后上传到服务器，无需在服务器安装Go环境。

```bash
# 在本地执行
./build.sh                              # 编译Linux版本
./upload.sh 服务器IP root                # 一键上传部署

# 或使用Makefile
make deploy SERVER=服务器IP USER=root    # 编译+上传一体化
```

**优势**：
- ✅ 服务器无需Go环境
- ✅ 编译速度更快
- ✅ 版本控制更清晰
- ✅ 支持交叉编译
- ✅ 自动化程度高

## 📦 脚本列表

### `bt-deploy.sh` - 宝塔环境一键部署脚本

**功能**：
- 🔍 检测宝塔面板和Go环境
- 📁 创建项目目录结构
- 🔄 下载/克隆项目代码
- ⚙️ 配置环境变量
- 🔨 编译项目
- 🛠️ 创建系统服务
- 🚀 启动服务

**使用方法**：
```bash
# 下载脚本
wget https://raw.githubusercontent.com/your-repo/safew-bot/main/deploy/bt-deploy.sh

# 赋予执行权限
chmod +x bt-deploy.sh

# 执行部署
./bt-deploy.sh
```

### `bt-update.sh` - 宝塔环境快速更新脚本

**功能**：
- 📋 检查项目环境
- 💾 备份当前版本
- ⏹️ 停止服务
- 🔄 更新代码
- 🔨 重新编译
- ▶️ 启动服务
- ✅ 检查运行状态

**使用方法**：
```bash
# 下载脚本
wget https://raw.githubusercontent.com/your-repo/safew-bot/main/deploy/bt-update.sh

# 赋予执行权限
chmod +x bt-update.sh

# 执行更新
./bt-update.sh
```

## 🛠️ 部署配置

### 默认目录结构

```
/www/wwwroot/safew-bot/    # 项目主目录
├── safew-bot              # 可执行文件
├── .env                   # 环境配置
├── go.mod                 # Go模块文件
└── ...                    # 其他项目文件

/www/backup/safew-bot/     # 备份目录
├── 20240608_120000/       # 按时间戳命名的备份
├── 20240608_140000/
└── ...

/etc/systemd/system/       # 系统服务
└── safew-bot.service      # 服务配置文件
```

### 服务配置

**服务名称**: `safew-bot`
**运行用户**: `www`
**工作目录**: `/www/wwwroot/safew-bot`
**自启动**: 已启用

### 常用命令

```bash
# 服务管理
systemctl start safew-bot      # 启动服务
systemctl stop safew-bot       # 停止服务
systemctl restart safew-bot    # 重启服务
systemctl status safew-bot     # 查看状态
systemctl enable safew-bot     # 启用自启动
systemctl disable safew-bot    # 禁用自启动

# 使用项目脚本管理
./start.sh                     # 启动Bot
./stop.sh                      # 停止Bot

# 日志查看
journalctl -u safew-bot -f              # 实时查看日志
journalctl -u safew-bot --since "1h"    # 查看1小时内日志
journalctl -u safew-bot -p err           # 查看错误日志
tail -f logs/safew-bot.log              # 查看项目日志文件

# 进程管理
ps aux | grep safew-bot                  # 查看进程
kill -9 $(pgrep safew-bot)              # 强制停止进程

# 本地编译部署 (推荐)
# 在本地机器执行：
make deploy SERVER=服务器IP USER=root    # 一键编译部署
./upload.sh 服务器IP root               # 上传编译好的文件

# 版本信息
./safew-bot -v                          # 查看程序版本
```

## 🔧 自定义配置

### 修改部署目录

编辑脚本中的配置变量：

```bash
# 在脚本开头修改
PROJECT_NAME="safew-bot"
PROJECT_DIR="/your/custom/path/${PROJECT_NAME}"
SERVICE_USER="your_user"
BACKUP_DIR="/your/backup/path/${PROJECT_NAME}"
```

### 修改服务配置

编辑 `/etc/systemd/system/safew-bot.service`：

```ini
[Unit]
Description=SafeW Bot Service
After=network.target

[Service]
Type=simple
User=www
Group=www
WorkingDirectory=/www/wwwroot/safew-bot
ExecStart=/www/wwwroot/safew-bot/safew-bot
Restart=always
RestartSec=5

# 自定义环境变量（可选）
Environment=LOG_LEVEL=DEBUG
Environment=POLL_TIMEOUT=60

[Install]
WantedBy=multi-user.target
```

重载配置：
```bash
systemctl daemon-reload
systemctl restart safew-bot
```

## ⚠️ 注意事项

1. **权限要求**：
   - 脚本需要sudo权限来创建系统服务
   - 项目文件归属于www用户

2. **环境依赖**：
   - 已安装宝塔面板
   - 已安装Go运行环境（1.19+）
   - 网络连接正常

3. **安全考虑**：
   - Token等敏感信息存储在.env文件中
   - .env文件权限设置为644
   - 服务以非特权用户运行

4. **备份策略**：
   - 每次更新前自动备份
   - 备份保留在`/www/backup/safew-bot/`
   - 建议定期清理旧备份

## 🆘 故障排除

### 常见问题

1. **服务启动失败**：
   ```bash
   # 查看详细错误信息
   journalctl -u safew-bot --no-pager
   
   # 检查配置文件
   cat /www/wwwroot/safew-bot/.env
   
   # 手动运行测试
   cd /www/wwwroot/safew-bot
   ./safew-bot
   ```

2. **编译失败**：
   ```bash
   # 检查Go环境
   go version
   go env
   
   # 手动编译测试
   cd /www/wwwroot/safew-bot
   go build -v
   ```

3. **权限问题**：
   ```bash
   # 修复文件权限
   chown -R www:www /www/wwwroot/safew-bot
   chmod +x /www/wwwroot/safew-bot/safew-bot
   ```

4. **端口占用**：
   ```bash
   # 检查端口使用
   netstat -tlnp | grep :端口号
   
   # 查找占用进程
   lsof -i :端口号
   ```

### 恢复备份

如果更新失败，可以恢复到之前的版本：

```bash
# 查看可用备份
ls -la /www/backup/safew-bot/

# 恢复指定备份
BACKUP_DATE="20240608_120000"  # 替换为实际备份时间
cd /www/wwwroot/safew-bot
systemctl stop safew-bot
cp /www/backup/safew-bot/$BACKUP_DATE/safew-bot ./
cp /www/backup/safew-bot/$BACKUP_DATE/.env ./
systemctl start safew-bot
```

## 📞 支持

如果遇到部署问题，请：

1. 检查系统日志：`journalctl -u safew-bot -f`
2. 验证环境配置：Go版本、宝塔面板、网络连接
3. 查看错误信息并参考故障排除指南
4. 在项目Issues中反馈问题 