# SafeW Bot

一个基于 Go 语言开发的 SafeW Bot，支持消息转发和群组管理功能。

## ⚡ 快速开始

### 👥 普通用户
想要了解如何使用Bot功能？
👉 **[用户使用指南](docs/USER_GUIDE.md)** - 功能介绍、使用方法、常见问题

### 🔧 技术人员部署

#### 3分钟部署到服务器

```bash
# 1. 本地编译 (1分钟)
git clone <repository-url> && cd safew-bot
./build.sh

# 2. 一键部署 (1分钟)
./upload.sh 服务器IP root

# 3. 配置启动 (1分钟)
ssh root@服务器IP
cd /www/wwwroot/safew-bot
cp .env.example .env
# 编辑 .env 添加 BOT_TOKEN 和 ADMIN_USERS
./start.sh
```

**更多方式**：
- 🔧 **使用Makefile**: `make deploy SERVER=服务器IP USER=root`
- 📖 **详细指南**: 参见[详细部署指南](docs/DEPLOY_GUIDE.md)
- ⚡ **快速指南**: 参见[快速开始指南](docs/QUICK_START.md)
- 🏠 **宝塔面板**: 参见下方[宝塔面板部署](#🏠-宝塔面板部署)章节

## 🚀 功能特性

### 📤 消息转发
- 支持转发文本、图片、视频、文档等多种类型的消息
- 通过回复消息使用 `/forward` 命令进行转发
- 支持指定目标群组进行精准转发

### 👮‍♂️ 群组管理
- 查看群组信息和管理员列表
- 禁言违规用户
- 提升用户为管理员
- 权限验证确保安全性

### 🔧 基础功能
- 友好的命令帮助系统
- 优雅的错误处理和日志记录
- 支持环境变量配置
- 优雅关闭机制

## 📋 环境要求

### 本地开发环境
- Go 1.21 或更高版本
- Git
- SSH客户端（用于部署）

### 服务器环境
- Linux服务器（支持本地编译部署，无需Go环境）
- SSH访问权限
- SafeW Bot Token

### 可选环境
- 宝塔面板（简化服务器管理）
- rsync（用于文件同步）

## 🛠️ 安装配置

### 1. 克隆项目
```bash
git clone <repository-url>
cd safew-bot
```

### 2. 获取 Bot Token
1. 在 SafeW 中联系 @BotFather
2. 使用 `/newbot` 命令创建新的 Bot
3. 获取 Token（格式如：`4839574812:AAFD39kkdpWt3ywyRZergyOLMaJhac60qc`）

### 3. 配置环境变量

项目支持两种配置方式：

#### 方式1: 使用 .env 文件（推荐）
```bash
# 复制示例配置文件
cp .env.example .env

# 编辑 .env 文件，填入您的实际配置
# 至少需要设置 SAFEW_BOT_TOKEN
```

#### 方式2: 系统环境变量
```bash
export SAFEW_BOT_TOKEN="your_bot_token_here"
export FORWARD_TARGET_CHAT="默认转发目标群组ID（可选）"
export LOG_LEVEL="INFO"  # DEBUG, INFO, WARN, ERROR
export POLL_TIMEOUT="30"  # 长轮询超时时间（秒）
```

> **注意**: .env 文件的优先级高于系统环境变量

### 4. 编译运行

#### 🏗️ 本地编译部署（推荐）

本项目支持本地编译后上传到服务器，无需在服务器安装Go环境，部署更高效。

**步骤1: 本地编译**
```bash
# 使用构建脚本（推荐）
chmod +x build.sh upload.sh
./build.sh

# 或使用 Makefile
make build

# 或手动交叉编译
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/safew-bot .
```

编译完成后，`build` 目录包含：
- `safew-bot` - Linux可执行文件
- `.env.example` - 配置模板
- `start.sh` / `stop.sh` - 启动/停止脚本
- `deploy/` - 部署相关脚本

**步骤2: 上传到服务器**
```bash
# 使用上传脚本（推荐）
./upload.sh 服务器IP 用户名 [目标目录]

# 示例
./upload.sh 192.168.1.100 root
./upload.sh example.com root /opt/safew-bot

# 或使用 Makefile 一键部署
make deploy SERVER=192.168.1.100 USER=root
make deploy SERVER=example.com USER=root DIR=/opt/safew-bot
```

**步骤3: 服务器启动**
```bash
# SSH连接到服务器
ssh root@your-server

# 进入目录
cd /www/wwwroot/safew-bot

# 配置环境（首次部署）
cp .env.example .env
vi .env  # 修改配置

# 启动服务
./start.sh

# 停止服务
./stop.sh
```

#### 📦 构建命令参考

```bash
# 显示帮助
make help

# 编译相关
make build          # 交叉编译Linux版本
make build-local    # 本地编译
make clean          # 清理编译文件

# 测试和检查
make test           # 运行测试
make fmt            # 格式化代码
make vet            # 代码检查
make check          # 完整检查

# 部署相关
make upload SERVER=ip USER=user    # 上传到服务器
make deploy SERVER=ip USER=user    # 编译+上传

# 依赖管理
make deps           # 安装依赖
make update         # 更新依赖

# 版本信息
make version        # 显示构建版本
./safew-bot -v      # 显示程序版本
```

#### 🖥️ 本地开发运行

```bash
# 本地编译
go build -o safew-bot

# 本地运行
./safew-bot

# 或直接运行
go run .

# 或使用 make
make run
```

## 📖 命令列表

### 🔧 基础命令
- `/start` - 开始使用Bot，显示欢迎信息
- `/help` - 显示所有可用命令和使用说明
- `/info` - 获取当前群组的详细信息

### 📤 转发功能
- `/forward <目标群ID>` - 转发回复的消息到指定群组
  - 使用方法：回复要转发的消息，然后输入命令
  - 示例：`/forward -1001234567890`

### 👮‍♂️ 管理命令（仅管理员）
- `/ban <@用户名> [原因]` - 禁言指定用户
- `/promote <@用户名>` - 提升用户为管理员
- `/admins` - 查看群组管理员列表

## 🔒 权限说明

- **普通用户**：可以使用基础命令和转发功能
- **群组管理员**：可以使用所有群组管理命令
- **Bot 权限**：需要在群组中给予 Bot 以下权限：
  - 读取消息
  - 发送消息
  - 删除消息
  - 管理聊天（用于管理功能）

## 📁 项目结构

```
safew-bot/
├── go.mod                  # Go 模块文件
├── go.sum                  # Go 依赖校验文件
├── main.go                 # 程序入口
├── config.go               # 配置管理
├── .env.example            # 环境变量配置示例
├── .env                    # 环境变量配置文件 (需手动创建)
├── .gitignore              # Git 忽略文件
├── build.sh                # 本地编译脚本
├── upload.sh               # 自动上传部署脚本
├── Makefile                # 构建工具和任务管理
├── bot/                    # Bot 核心包
│   ├── models.go           # API 数据结构
│   ├── api.go              # API 客户端
│   ├── bot.go              # Bot 主循环
│   └── handlers.go         # 消息处理器
├── docs/                   # 文档目录
│   └── development-plan.md # 开发计划
├── deploy/                 # 部署脚本目录
│   ├── README.md           # 部署脚本说明文档
│   ├── bt-deploy.sh        # 宝塔环境一键部署脚本
│   └── bt-update.sh        # 宝塔环境快速更新脚本
├── build/                  # 编译输出目录 (gitignore)
│   ├── safew-bot           # Linux可执行文件
│   ├── .env.example        # 配置模板
│   ├── start.sh            # 启动脚本
│   ├── stop.sh             # 停止脚本
│   └── deploy/             # 部署脚本副本
└── README.md               # 项目说明
```

## 🚦 使用示例

### 启动Bot
```bash
# 方式1: 使用 .env 文件
# 确保 .env 文件中已设置正确的 SAFEW_BOT_TOKEN
go run .

# 方式2: 临时设置环境变量
export SAFEW_BOT_TOKEN="4839574812:AAFD39kkdpWt3ywyRZergyOLMaJhac60qc"
go run .
```

### 转发消息
1. 在群组中回复要转发的消息
2. 输入：`/forward -1001234567890`
3. Bot 会将回复的消息转发到指定群组

### 管理群组
```
/info                    # 查看群组信息
/admins                  # 查看管理员列表
/ban @username 违反群规   # 禁言用户
/promote @username       # 提升管理员
```

## ⚠️ 注意事项

1. **Token 安全**：
   - 绝不要在代码中硬编码 Token
   - 使用 .env 文件或环境变量进行配置
   - .env 文件已添加到 .gitignore，不会被提交到版本控制
   - 定期更换 Token

2. **权限管理**：
   - 管理命令会验证用户权限
   - Bot 需要相应的群组权限才能执行操作

3. **错误处理**：
   - Bot 会自动重试失败的请求
   - 网络错误会记录日志并继续运行

4. **性能优化**：
   - 使用长轮询减少不必要的请求
   - 合理设置超时时间

## 🚀 生产环境部署

### Linux 服务器部署（推荐）

#### 1. 服务器环境准备
```bash
# 更新系统包
sudo apt update && sudo apt upgrade -y  # Ubuntu/Debian
# 或
sudo yum update -y  # CentOS/RHEL

# 安装 Go 环境 (如果未安装)
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. 项目部署
```bash
# 创建应用目录
sudo mkdir -p /opt/safew-bot
cd /opt/safew-bot

# 上传项目文件或克隆代码
git clone <your-repository-url> .
# 或上传项目文件到此目录

# 配置环境变量
cp .env.example .env
nano .env  # 编辑配置文件

# 编译生产版本
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o safew-bot

# 设置可执行权限
chmod +x safew-bot

# 创建专用用户（安全考虑）
sudo useradd -r -s /bin/false safew
sudo chown -R safew:safew /opt/safew-bot
```

#### 3. 系统服务配置
```bash
# 创建 systemd 服务文件
sudo nano /etc/systemd/system/safew-bot.service
```

服务文件内容：
```ini
[Unit]
Description=SafeW Bot Service
After=network.target
Wants=network.target

[Service]
Type=simple
User=safew
Group=safew
WorkingDirectory=/opt/safew-bot
ExecStart=/opt/safew-bot/safew-bot
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=safew-bot

# 环境变量（可选，如果不使用 .env 文件）
# Environment=SAFEW_BOT_TOKEN=your_token_here
# Environment=LOG_LEVEL=INFO

[Install]
WantedBy=multi-user.target
```

```bash
# 重载 systemd 配置
sudo systemctl daemon-reload

# 启用服务（开机自启）
sudo systemctl enable safew-bot

# 启动服务
sudo systemctl start safew-bot

# 查看服务状态
sudo systemctl status safew-bot
```

### 🏠 宝塔面板部署

#### 一键部署脚本（推荐）

我们提供了自动化部署脚本，简化宝塔环境的部署过程：

```bash
# 下载部署脚本
wget https://raw.githubusercontent.com/your-repo/safew-bot/main/deploy/bt-deploy.sh
# 或使用 curl
curl -O https://raw.githubusercontent.com/your-repo/safew-bot/main/deploy/bt-deploy.sh

# 赋予执行权限
chmod +x bt-deploy.sh

# 执行部署
./bt-deploy.sh
```

**脚本功能**：
- ✅ 自动检测宝塔环境和Go语言
- ✅ 创建项目目录和用户权限
- ✅ 下载/更新项目代码
- ✅ 配置环境变量文件
- ✅ 编译和部署应用
- ✅ 创建系统服务和自启动
- ✅ 启动服务并检查状态

**更新脚本**：
```bash
# 下载更新脚本
wget https://raw.githubusercontent.com/your-repo/safew-bot/main/deploy/bt-update.sh
chmod +x bt-update.sh

# 执行更新
./bt-update.sh
```

#### 手动部署方式

#### 1. 环境准备
1. **安装 Go 环境**：
   - 在宝塔面板中进入 "软件商店"
   - 搜索并安装 "Go" 运行环境
   - 或手动安装：软件商店 → 运行环境 → Go

2. **创建网站目录**：
   - 网站 → 添加站点 → 创建目录 `/www/wwwroot/safew-bot`
   - 或使用命令行：`mkdir -p /www/wwwroot/safew-bot`

#### 2. 上传和配置
```bash
# 进入项目目录
cd /www/wwwroot/safew-bot

# 上传项目文件（通过宝塔文件管理器或 Git）
# 方式1: 使用 Git
git clone <your-repository-url> .

# 方式2: 通过宝塔面板文件管理器上传zip包并解压

# 配置环境变量
cp .env.example .env
# 通过宝塔面板文件管理器编辑 .env 文件
```

#### 3. 编译和运行
```bash
# 编译项目
go build -o safew-bot

# 设置权限
chmod +x safew-bot
```

#### 4. 宝塔进程管理
1. **安装进程守护器**：
   - 软件商店 → 系统工具 → 安装 "进程守护器"

2. **添加进程**：
   - 进程守护器 → 添加守护进程
   - 名称：`SafeW Bot`
   - 启动文件：`/www/wwwroot/safew-bot/safew-bot`
   - 运行目录：`/www/wwwroot/safew-bot`
   - 运行用户：`www`
   - **端口：留空或填写0**（⚠️ 重要：本Bot不使用端口）
   - 启动参数：留空

3. **进程管理**：
   - 启动：点击 "启动" 按钮
   - 停止：点击 "停止" 按钮
   - 重启：点击 "重启" 按钮
   - 查看日志：点击 "日志" 按钮

**重要说明**：SafeW Bot使用长轮询模式，不需要监听端口，宝塔面板端口字段必须留空或填写0。

#### 5. 宝塔定时任务（可选）
添加定时重启任务（凌晨重启，清理内存）：
```bash
# 计划任务 → 添加任务
# 任务类型：Shell脚本
# 任务名称：重启SafeW Bot
# 执行周期：每天 03:00
# 脚本内容：
systemctl restart safew-bot
# 或者如果使用进程守护器：
# 通过API重启进程
```

### 📊 监控和日志

#### 1. 日志查看
```bash
# 查看系统日志
sudo journalctl -u safew-bot -f

# 查看最近日志
sudo journalctl -u safew-bot --since "1 hour ago"

# 查看错误日志
sudo journalctl -u safew-bot -p err
```

#### 2. 宝塔日志管理
- **实时日志**：进程守护器 → 点击对应进程的 "日志" 按钮
- **日志文件**：`/www/wwwroot/safew-bot/` 目录下的日志文件
- **系统日志**：安全 → 系统日志 → 筛选 "safew-bot"

#### 3. 性能监控
```bash
# 查看进程状态
ps aux | grep safew-bot

# 查看内存使用
top -p $(pgrep safew-bot)

# 查看网络连接
netstat -tlnp | grep safew-bot
```

### 🔄 更新维护

#### 1. 更新代码
```bash
cd /opt/safew-bot  # 或 /www/wwwroot/safew-bot

# 停止服务
sudo systemctl stop safew-bot

# 拉取最新代码
git pull origin main

# 重新编译
go build -o safew-bot

# 启动服务
sudo systemctl start safew-bot

# 检查状态
sudo systemctl status safew-bot
```

#### 2. 备份配置
```bash
# 备份配置文件
cp .env .env.backup.$(date +%Y%m%d)

# 备份整个项目（可选）
tar -czf safew-bot-backup-$(date +%Y%m%d).tar.gz /opt/safew-bot
```

#### 3. 故障排除
```bash
# 检查配置文件
cat .env

# 检查编译错误
go build -v -o safew-bot

# 检查端口占用
ss -tlnp | grep :443  # 如果使用HTTPS

# 测试配置
./safew-bot --test-config  # 需要添加此功能
```

### 🔒 安全配置

#### 1. 防火墙设置
```bash
# 如果Bot需要接收Webhook（可选）
sudo ufw allow 8443/tcp  # 或您设置的端口

# 宝塔面板安全设置
# 安全 → 防火墙 → 添加规则
```

#### 2. SSL证书（如果使用Webhook）
- 宝塔面板 → SSL → 申请免费证书
- 或使用 Let's Encrypt

#### 3. 进程权限
```bash
# 确保Bot以非root用户运行
sudo chown -R safew:safew /opt/safew-bot
sudo chmod 644 .env  # 限制配置文件权限
```

## 🔧 开发相关

### 编译生产版本
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o safew-bot
```

### 运行测试
```bash
go test ./...
```

### 代码格式化
```bash
go fmt ./...
```

## 📝 更新日志

### v1.1.0 (最新)
- ✅ 本地编译部署支持
- ✅ 自动化构建脚本 (`build.sh`)
- ✅ 一键上传部署脚本 (`upload.sh`)
- ✅ Makefile 构建工具支持
- ✅ 版本信息和构建时间显示
- ✅ 交叉编译Linux版本
- ✅ 服务器启动/停止脚本
- ✅ 完善的部署文档

### v1.0.0
- ✅ 基础框架搭建
- ✅ 消息转发功能
- ✅ 群组管理功能
- ✅ 命令处理系统
- ✅ 配置管理
- ✅ 优雅关闭
- ✅ 宝塔面板部署脚本

## ❓ 常见问题

### Q: 本地编译失败怎么办？
```bash
# 检查Go版本 (需要1.21+)
go version

# 重新安装依赖
go mod tidy

# 清理后重新编译
make clean && make build
```

### Q: 上传到服务器失败？
```bash
# 检查SSH连接
ssh -v root@服务器IP

# 检查rsync是否安装
which rsync

# 使用详细模式查看错误
./upload.sh 服务器IP root
```

### Q: Bot不响应消息？
1. 检查Token是否正确设置
2. 确认Bot已添加到目标群组
3. 检查Bot是否有发送消息权限
4. 查看运行日志：`tail -f logs/safew-bot.log`

### Q: 如何更新Bot到新版本？
```bash
# 本地重新编译
make clean && make build

# 重新部署
make deploy SERVER=服务器IP USER=root

# 或使用更新脚本
ssh root@服务器IP "cd /www/wwwroot/safew-bot && ./deploy/bt-update.sh"
```

### Q: 如何在宝塔面板中管理？
1. 安装"进程守护器"插件
2. 添加进程：
   - 启动文件：`/www/wwwroot/safew-bot/safew-bot`
   - **端口：留空或填写0**（Bot不使用端口）
3. 或使用一键部署脚本：`./deploy/bt-deploy.sh`

### Q: 宝塔面板提示端口错误怎么办？
SafeW Bot使用长轮询模式，不监听端口：
- 端口字段留空或填写`0`
- 忽略宝塔的端口验证提示
- 这是正常现象，不影响Bot运行

## 📚 文档中心

根据您的需求选择合适的文档：

| 文档 | 适用人群 | 内容 | 链接 |
|------|----------|------|------|
| 📱 **用户使用指南** | Bot使用者、群组管理员 | 功能介绍、使用方法、常见问题 | [USER_GUIDE.md](docs/USER_GUIDE.md) |
| ⚡ **快速开始指南** | 技术人员、服务器管理员 | 3分钟快速部署、基础配置 | [QUICK_START.md](docs/QUICK_START.md) |
| 🔧 **详细部署指南** | 开发者、系统管理员、运维人员 | 编译、部署、配置、故障排除 | [DEPLOY_GUIDE.md](docs/DEPLOY_GUIDE.md) |
| 🏗️ **开发计划文档** | 开发者、贡献者 | 技术选型、API分析、项目结构 | [development-plan.md](docs/development-plan.md) |

💡 **不知道看哪个？** 查看 [docs/README.md](docs/README.md) 获取详细的文档指引。

## 🤝 贡献指南

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 支持

如果您遇到问题或有功能建议，请：
1. 查看现有的 Issues
2. 创建新的 Issue 描述问题
3. 提供详细的错误信息和步骤

---

**注意**: 本项目仅用于学习和研究目的，请遵守相关法律法规和平台规则。 