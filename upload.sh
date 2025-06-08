#!/bin/bash

# SafeW Bot 上传部署脚本
# 自动上传编译好的文件到宝塔服务器

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认配置
BUILD_DIR="build"
DEFAULT_REMOTE_DIR="/www/wwwroot/safew-bot"

# 显示使用说明
show_usage() {
    echo -e "${BLUE}SafeW Bot 上传部署脚本${NC}"
    echo -e "${YELLOW}使用方法:${NC}"
    echo "  $0 <服务器地址> <用户名> [远程目录]"
    echo ""
    echo -e "${YELLOW}参数说明:${NC}"
    echo "  服务器地址  : 服务器IP或域名"
    echo "  用户名      : SSH登录用户名"
    echo "  远程目录    : 目标目录 (默认: $DEFAULT_REMOTE_DIR)"
    echo ""
    echo -e "${YELLOW}示例:${NC}"
    echo "  $0 192.168.1.100 root"
    echo "  $0 example.com root /opt/safew-bot"
    echo ""
    echo -e "${YELLOW}注意事项:${NC}"
    echo "  1. 请先运行 ./build.sh 编译程序"
    echo "  2. 确保已配置SSH密钥认证或准备输入密码"
    echo "  3. 服务器需要安装rsync命令"
}

# 检查参数
if [ $# -lt 2 ] || [ $# -gt 3 ]; then
    show_usage
    exit 1
fi

SERVER="$1"
USERNAME="$2"
REMOTE_DIR="${3:-$DEFAULT_REMOTE_DIR}"

# 检查编译文件是否存在
if [ ! -d "$BUILD_DIR" ]; then
    echo -e "${RED}❌ 编译目录 $BUILD_DIR 不存在，请先运行 ./build.sh 编译程序${NC}"
    exit 1
fi

if [ ! -f "$BUILD_DIR/safew-bot" ]; then
    echo -e "${RED}❌ 编译文件不存在，请先运行 ./build.sh 编译程序${NC}"
    exit 1
fi

echo -e "${BLUE}🚀 开始上传 SafeW Bot 到服务器...${NC}"
echo -e "${YELLOW}服务器: ${SERVER}${NC}"
echo -e "${YELLOW}用户名: ${USERNAME}${NC}"
echo -e "${YELLOW}目标目录: ${REMOTE_DIR}${NC}"

# 检查服务器连接
echo -e "${YELLOW}检查服务器连接...${NC}"
if ! ssh -o ConnectTimeout=10 "${USERNAME}@${SERVER}" "echo '连接成功'" 2>/dev/null; then
    echo -e "${RED}❌ 无法连接到服务器，请检查:${NC}"
    echo "  1. 服务器地址是否正确"
    echo "  2. SSH服务是否启动"
    echo "  3. 网络连接是否正常"
    echo "  4. SSH密钥配置是否正确"
    exit 1
fi

# 创建远程目录
echo -e "${YELLOW}创建远程目录...${NC}"
ssh "${USERNAME}@${SERVER}" "mkdir -p $REMOTE_DIR"

# 停止远程服务（如果在运行）
echo -e "${YELLOW}停止远程服务...${NC}"
ssh "${USERNAME}@${SERVER}" "cd $REMOTE_DIR && if [ -f stop.sh ]; then ./stop.sh; fi" 2>/dev/null || true

# 上传文件
echo -e "${YELLOW}上传文件...${NC}"
rsync -avz --progress \
    --exclude='.git' \
    --exclude='.env' \
    "$BUILD_DIR/" \
    "${USERNAME}@${SERVER}:$REMOTE_DIR/"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 文件上传成功！${NC}"
else
    echo -e "${RED}❌ 文件上传失败！${NC}"
    exit 1
fi

# 设置权限
echo -e "${YELLOW}设置文件权限...${NC}"
ssh "${USERNAME}@${SERVER}" "cd $REMOTE_DIR && chmod +x safew-bot start.sh stop.sh deploy/*.sh"

# 检查配置文件
echo -e "${YELLOW}检查配置文件...${NC}"
if ssh "${USERNAME}@${SERVER}" "cd $REMOTE_DIR && [ ! -f .env ]"; then
    echo -e "${YELLOW}⚠️  配置文件 .env 不存在，需要手动配置${NC}"
    
    # 询问是否自动配置
    read -p "是否现在配置 .env 文件？(y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}请输入Bot配置信息:${NC}"
        
        # 读取配置
        read -p "Bot Token: " BOT_TOKEN
        read -p "管理员用户ID (多个用逗号分隔): " ADMIN_USERS
        read -p "Log级别 (debug/info/warn/error) [info]: " LOG_LEVEL
        LOG_LEVEL=${LOG_LEVEL:-info}
        
        # 创建配置文件
        ssh "${USERNAME}@${SERVER}" "cd $REMOTE_DIR && cat > .env << EOF
# SafeW Bot 配置文件

# Bot基础配置
BOT_TOKEN=$BOT_TOKEN
BOT_USERNAME=

# 管理员配置
ADMIN_USERS=$ADMIN_USERS

# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# 日志配置
LOG_LEVEL=$LOG_LEVEL
LOG_FILE=logs/safew-bot.log

# Webhook配置 (可选)
WEBHOOK_URL=
WEBHOOK_SECRET=

# 数据库配置 (可选)
DATABASE_TYPE=sqlite
DATABASE_PATH=data/safew-bot.db

# Redis配置 (可选)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# 功能开关
ENABLE_WEBHOOK=false
ENABLE_DATABASE=true
ENABLE_REDIS=false
ENABLE_METRICS=true
ENABLE_ADMIN_COMMANDS=true
ENABLE_MESSAGE_FORWARDING=true
ENABLE_GROUP_MANAGEMENT=true

# 限流配置
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=30
RATE_LIMIT_WINDOW=60

# 文件上传配置
MAX_FILE_SIZE=50
ALLOWED_FILE_TYPES=jpg,jpeg,png,gif,pdf,doc,docx,txt

# 安全配置
ALLOWED_ORIGINS=*
CORS_ENABLED=true
HTTPS_ONLY=false
EOF"
        
        echo -e "${GREEN}✅ 配置文件创建成功！${NC}"
    fi
fi

# 询问是否立即启动
read -p "是否立即启动 SafeW Bot？(y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}启动 SafeW Bot...${NC}"
    ssh "${USERNAME}@${SERVER}" "cd $REMOTE_DIR && nohup ./start.sh > startup.log 2>&1 &"
    
    # 等待启动
    sleep 3
    
    # 检查启动状态
    if ssh "${USERNAME}@${SERVER}" "pgrep -f safew-bot > /dev/null"; then
        echo -e "${GREEN}✅ SafeW Bot 启动成功！${NC}"
    else
        echo -e "${RED}❌ SafeW Bot 启动失败，请检查日志${NC}"
        echo -e "${YELLOW}查看启动日志: ssh ${USERNAME}@${SERVER} 'cd $REMOTE_DIR && cat startup.log'${NC}"
    fi
fi

echo ""
echo -e "${GREEN}🎉 部署完成！${NC}"
echo -e "${BLUE}服务器管理命令:${NC}"
echo -e "连接服务器: ${YELLOW}ssh ${USERNAME}@${SERVER}${NC}"
echo -e "进入目录: ${YELLOW}cd $REMOTE_DIR${NC}"
echo -e "启动服务: ${YELLOW}./start.sh${NC}"
echo -e "停止服务: ${YELLOW}./stop.sh${NC}"
echo -e "查看日志: ${YELLOW}tail -f logs/safew-bot.log${NC}"
echo -e "查看状态: ${YELLOW}ps aux | grep safew-bot${NC}" 