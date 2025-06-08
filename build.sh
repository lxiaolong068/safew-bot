#!/bin/bash

# SafeW Bot 本地编译脚本
# 支持交叉编译为Linux服务器架构

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_NAME="safew-bot"
BUILD_DIR="build"
VERSION=$(date +"%Y%m%d_%H%M%S")

echo -e "${BLUE}🔨 开始编译 SafeW Bot...${NC}"

# 清理之前的编译文件
if [ -d "$BUILD_DIR" ]; then
    echo -e "${YELLOW}清理旧的编译文件...${NC}"
    rm -rf "$BUILD_DIR"
fi

# 创建编译目录
mkdir -p "$BUILD_DIR"

# 获取Go版本信息
GO_VERSION=$(go version)
echo -e "${BLUE}Go版本: ${GO_VERSION}${NC}"

# 编译配置
GOOS_TARGET="linux"
GOARCH_TARGET="amd64"

echo -e "${BLUE}目标平台: ${GOOS_TARGET}/${GOARCH_TARGET}${NC}"

# 设置编译参数
LDFLAGS="-w -s -X main.Version=${VERSION} -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"

echo -e "${YELLOW}开始交叉编译...${NC}"

# 交叉编译
GOOS=$GOOS_TARGET GOARCH=$GOARCH_TARGET go build \
    -ldflags "$LDFLAGS" \
    -o "$BUILD_DIR/${PROJECT_NAME}" \
    .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 编译成功！${NC}"
else
    echo -e "${RED}❌ 编译失败！${NC}"
    exit 1
fi

# 复制配置文件
echo -e "${YELLOW}复制配置文件...${NC}"
cp .env.example "$BUILD_DIR/"
cp README.md "$BUILD_DIR/"

# 复制部署脚本
cp -r deploy "$BUILD_DIR/"

# 创建启动脚本
cat > "$BUILD_DIR/start.sh" << 'EOF'
#!/bin/bash

# SafeW Bot 启动脚本

# 设置权限
chmod +x safew-bot

# 检查配置文件
if [ ! -f ".env" ]; then
    echo "❌ 配置文件 .env 不存在，请复制 .env.example 并修改配置"
    exit 1
fi

# 启动Bot
echo "🚀 启动 SafeW Bot..."
./safew-bot
EOF

chmod +x "$BUILD_DIR/start.sh"

# 创建停止脚本
cat > "$BUILD_DIR/stop.sh" << 'EOF'
#!/bin/bash

# SafeW Bot 停止脚本

echo "🛑 停止 SafeW Bot..."

# 查找并停止进程
PID=$(pgrep -f safew-bot)
if [ -n "$PID" ]; then
    kill -TERM $PID
    sleep 2
    
    # 如果进程仍在运行，强制杀死
    if pgrep -f safew-bot > /dev/null; then
        kill -9 $PID
        echo "✅ 强制停止 SafeW Bot (PID: $PID)"
    else
        echo "✅ 优雅停止 SafeW Bot (PID: $PID)"
    fi
else
    echo "ℹ️  SafeW Bot 未运行"
fi
EOF

chmod +x "$BUILD_DIR/stop.sh"

# 获取文件信息
FILE_SIZE=$(du -sh "$BUILD_DIR/${PROJECT_NAME}" | cut -f1)
echo -e "${GREEN}📦 编译完成！${NC}"
echo -e "${BLUE}文件大小: ${FILE_SIZE}${NC}"
echo -e "${BLUE}输出目录: ${BUILD_DIR}/${NC}"

# 显示文件列表
echo -e "${YELLOW}编译输出文件:${NC}"
ls -la "$BUILD_DIR/"

echo ""
echo -e "${GREEN}🎉 编译完成！可以将 ${BUILD_DIR} 目录上传到服务器${NC}"
echo -e "${BLUE}💡 使用说明:${NC}"
echo -e "1. 将整个 ${BUILD_DIR} 目录上传到服务器"
echo -e "2. 复制 .env.example 为 .env 并配置"
echo -e "3. 运行 ./start.sh 启动Bot"
echo -e "4. 运行 ./stop.sh 停止Bot" 