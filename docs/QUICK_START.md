# SafeW Bot 快速开始指南

## 🚀 3分钟部署到宝塔服务器

### 前置条件
- 本地安装 Go 1.21+
- 服务器已安装宝塔面板
- 已获取SafeW Bot Token

### 步骤1: 本地编译 (1分钟)
```bash
# 克隆项目
git clone <repository-url>
cd safew-bot

# 一键编译
./build.sh
```

### 步骤2: 上传部署 (1分钟)
```bash
# 一键上传（替换为你的服务器IP）
./upload.sh 192.168.1.100 root

# 或使用Makefile
make deploy SERVER=192.168.1.100 USER=root
```

### 步骤3: 配置启动 (1分钟)
```bash
# SSH连接服务器
ssh root@192.168.1.100

# 进入目录
cd /www/wwwroot/safew-bot

# 配置Bot Token
cp .env.example .env
echo "BOT_TOKEN=your_bot_token_here" >> .env
echo "ADMIN_USERS=your_user_id" >> .env

# 启动
./start.sh
```

## ✅ 完成！

Bot现在应该已经启动并运行。在SafeW中向Bot发送 `/start` 来测试。

## 📱 使用Bot

### 基础命令
- `/start` - 开始使用
- `/help` - 帮助信息
- `/info` - 群组信息

### 转发功能
1. 回复要转发的消息
2. 输入 `/forward 目标群ID`

### 管理功能（仅管理员）
- `/ban @用户名` - 禁言用户
- `/promote @用户名` - 提升管理员
- `/admins` - 查看管理员

## 🔧 管理命令

```bash
# 查看状态
ps aux | grep safew-bot

# 停止服务
./stop.sh

# 重启服务
./stop.sh && ./start.sh

# 查看日志
tail -f logs/safew-bot.log
```

## 🆘 故障排除

### Bot不响应
```bash
# 检查进程
ps aux | grep safew-bot

# 检查配置
cat .env

# 查看日志
tail logs/safew-bot.log
```

### 重新部署
```bash
# 本地重新编译
make clean && make build

# 上传
make upload SERVER=你的服务器 USER=root

# 服务器重启
ssh root@你的服务器 "cd /www/wwwroot/safew-bot && ./stop.sh && ./start.sh"
```

---

更多详细信息请参考：
- [部署指南](DEPLOY_GUIDE.md)
- [完整文档](README.md) 