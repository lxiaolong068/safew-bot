# SafeW Bot 本地编译部署指南

本指南介绍如何在本地编译SafeW Bot并部署到宝塔服务器。

## 🚀 快速开始

### 1. 本地编译

```bash
# 克隆项目
git clone <repository-url>
cd safew-bot

# 安装依赖
go mod tidy

# 编译Linux版本
./build.sh
```

### 2. 上传到服务器

```bash
# 方式1: 使用脚本上传
./upload.sh 服务器IP root

# 方式2: 使用Makefile一键部署
make deploy SERVER=服务器IP USER=root

# 方式3: 手动上传（scp方式）
scp -r build/* root@服务器IP:/www/wwwroot/safew-bot/
```

### 3. 服务器配置

```bash
# SSH连接服务器
ssh root@服务器IP

# 进入目录
cd /www/wwwroot/safew-bot

# 配置环境变量
cp .env.example .env
vi .env

# 启动服务
./start.sh
```

## 📋 详细步骤

### 环境准备

**本地环境:**
- Go 1.21+
- Git
- SSH客户端

**服务器环境:**
- Linux服务器
- SSH访问权限
- 已安装宝塔面板（可选）

### 编译配置

#### 使用脚本编译（推荐）
```bash
./build.sh
```

编译脚本会：
- 交叉编译为Linux版本
- 复制配置文件模板
- 创建启动/停止脚本
- 打包所有必要文件

#### 使用Makefile
```bash
# 显示所有可用命令
make help

# 编译
make build

# 清理
make clean

# 一键部署
make deploy SERVER=192.168.1.100 USER=root
```

#### 手动编译
```bash
GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=$(date +%Y%m%d_%H%M%S)" \
    -o build/safew-bot .
```

### 上传方式

#### 1. 脚本上传（推荐）
```bash
# 基本用法
./upload.sh <服务器IP> <用户名> [远程目录]

# 示例
./upload.sh 192.168.1.100 root
./upload.sh example.com root /opt/safew-bot

# 脚本会自动：
# - 检查编译文件
# - 验证服务器连接
# - 停止旧服务
# - 上传文件
# - 设置权限
# - 可选配置和启动
```

#### 2. Makefile上传
```bash
# 上传到默认目录
make upload SERVER=192.168.1.100 USER=root

# 指定目录
make upload SERVER=192.168.1.100 USER=root DIR=/opt/safew-bot

# 编译+上传
make deploy SERVER=192.168.1.100 USER=root
```

#### 3. 手动上传
```bash
# 使用rsync（推荐）
rsync -avz --progress build/ root@服务器IP:/www/wwwroot/safew-bot/

# 使用scp
scp -r build/* root@服务器IP:/www/wwwroot/safew-bot/

# 使用宝塔面板
# 1. 压缩build目录：tar -czf safew-bot.tar.gz build/
# 2. 通过宝塔文件管理器上传
# 3. 在服务器解压：tar -xzf safew-bot.tar.gz
```

### 服务器配置

#### 1. 环境变量配置
```bash
# 复制配置模板
cp .env.example .env

# 编辑配置文件
vi .env
```

必要配置项：
```env
# Bot Token (必填)
SAFEW_BOT_TOKEN=your_bot_token_here

# 日志级别 (可选，默认：INFO)
LOG_LEVEL=INFO

# 长轮询超时时间 (可选，默认：30秒)
POLL_TIMEOUT=30

# 默认转发目标群组ID (可选)
# FORWARD_TARGET_CHAT=-1001234567890

# 超级管理员用户ID列表 (可选)
# SUPER_ADMINS=123456789,987654321
```

**重要说明**：
- 本Bot使用**长轮询模式**，不运行HTTP服务器，因此**不需要端口配置**
- 在宝塔进程守护器中**无需填写端口号**
- Bot通过主动请求SafeW API来获取消息，而不是被动接收

#### 2. 启动服务
```bash
# 方式1: 使用启动脚本
./start.sh

# 方式2: 直接运行
./safew-bot

# 方式3: 后台运行
nohup ./safew-bot > safew-bot.log 2>&1 &

# 方式4: 使用宝塔进程守护器
# 在宝塔面板中添加进程监控
```

#### 3. 服务管理
```bash
# 启动
./start.sh

# 停止
./stop.sh

# 查看状态
ps aux | grep safew-bot

# 查看日志
tail -f logs/safew-bot.log
```

## 🔧 高级配置

### 系统服务配置
```bash
# 创建系统服务文件
sudo tee /etc/systemd/system/safew-bot.service << EOF
[Unit]
Description=SafeW Bot
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/www/wwwroot/safew-bot
ExecStart=/www/wwwroot/safew-bot/safew-bot
Restart=always
RestartSec=5
Environment=PATH=/usr/local/bin:/usr/bin:/bin

[Install]
WantedBy=multi-user.target
EOF

# 启用服务
sudo systemctl enable safew-bot
sudo systemctl start safew-bot
```

### 宝塔进程守护器配置

#### 安装进程守护器
1. 登录宝塔面板
2. 软件商店 → 系统工具 → 搜索"进程守护器"
3. 点击安装

#### 添加SafeW Bot进程
1. 进入进程守护器管理页面
2. 点击"添加守护进程"
3. 填写以下配置：

| 配置项 | 值 | 说明 |
|--------|-----|------|
| **名称** | `SafeW Bot` | 进程显示名称 |
| **启动文件** | `/www/wwwroot/safew-bot/safew-bot` | 可执行文件完整路径 |
| **运行目录** | `/www/wwwroot/safew-bot` | 工作目录 |
| **运行用户** | `www` | 运行用户（建议非root） |
| **端口** | **留空或填0** | ⚠️ 重要：本Bot不使用端口 |
| **启动参数** | 留空 | 无需额外参数 |
| **开机启动** | 开启 | 服务器重启自动启动 |

#### 重要提示
- ⚠️ **端口字段必须留空或填写0**，因为SafeW Bot使用长轮询模式，不监听任何端口
- ✅ 如果宝塔要求必须填写端口，可以填写`0`
- 🔍 添加后可在进程守护器中查看运行状态、启动/停止服务、查看日志

#### 常见宝塔配置问题
1. **端口验证错误**：宝塔可能提示端口无效，这是正常的，直接保存即可
2. **权限问题**：确保`/www/wwwroot/safew-bot/safew-bot`有执行权限
3. **路径错误**：检查启动文件路径是否正确
4. **用户权限**：`www`用户需要对项目目录有读取和执行权限

### 自动更新脚本
```bash
# 创建更新脚本
cat > update.sh << 'EOF'
#!/bin/bash
cd /www/wwwroot/safew-bot
./stop.sh
# 这里可以添加从Git拉取或下载新版本的逻辑
./start.sh
EOF

chmod +x update.sh
```

## 🔍 故障排除

### 常见问题

1. **编译失败**
   ```bash
   # 检查Go版本
   go version
   
   # 重新安装依赖
   go mod tidy
   ```

2. **上传失败**
   ```bash
   # 检查SSH连接
   ssh -v root@服务器IP
   
   # 检查rsync是否安装
   which rsync
   ```

3. **启动失败**
   ```bash
   # 检查配置文件
   cat .env
   
   # 检查权限
   ls -la safew-bot
   chmod +x safew-bot
   
   # 查看错误日志
   ./safew-bot 2>&1 | tee error.log
   ```

4. **配置错误**
   ```bash
   # 检查配置文件语法
   cat .env
   
   # 验证Token格式
   echo $SAFEW_BOT_TOKEN
   
   # 测试Bot连接
   ./safew-bot -test || ./safew-bot --version
   ```

### 调试模式
```bash
# 启用调试日志
echo "LOG_LEVEL=debug" >> .env

# 前台运行查看详细输出
./safew-bot
```

## 📝 部署检查清单

### 编译前检查
- [ ] Go环境版本 >= 1.21
- [ ] 项目依赖已安装 (`go mod tidy`)
- [ ] 代码格式检查 (`make check`)

### 上传前检查
- [ ] 编译成功，build目录存在
- [ ] SSH密钥配置或密码准备就绪
- [ ] 服务器网络连接正常

### 部署后检查
- [ ] 配置文件 .env 已正确配置
- [ ] Bot Token 有效
- [ ] 程序能正常启动
- [ ] 日志文件正常生成
- [ ] Bot响应测试消息

### 生产环境检查
- [ ] 系统服务配置
- [ ] 自动重启配置
- [ ] 日志轮转配置
- [ ] 监控和告警配置
- [ ] 备份策略制定

## 🎯 最佳实践

1. **版本管理**：每次部署记录版本信息
2. **配置管理**：敏感配置使用环境变量
3. **日志管理**：定期清理和归档日志
4. **监控告警**：配置进程和资源监控
5. **安全加固**：最小权限原则，定期更新

---

更多详细信息请参考 [README.md](README.md) 