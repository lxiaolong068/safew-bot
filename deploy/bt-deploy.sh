#!/bin/bash

# SafeW Bot 宝塔环境一键部署脚本
# 使用方法: chmod +x bt-deploy.sh && ./bt-deploy.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_NAME="safew-bot"
PROJECT_DIR="/www/wwwroot/${PROJECT_NAME}"
SERVICE_USER="www"
BACKUP_DIR="/www/backup/${PROJECT_NAME}"

# 输出函数
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为root用户
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_warning "检测到以root用户运行，建议使用普通用户"
    fi
}

# 检查宝塔面板
check_bt_panel() {
    if ! command -v bt &> /dev/null; then
        print_error "未检测到宝塔面板，请先安装宝塔面板"
        exit 1
    fi
    print_success "宝塔面板检测通过"
}

# 检查Go环境
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "未检测到Go环境，请在宝塔面板中安装Go运行环境"
        print_info "安装路径：软件商店 → 运行环境 → Go"
        exit 1
    fi
    print_success "Go环境检测通过 ($(go version))"
}

# 创建项目目录
create_directories() {
    print_info "创建项目目录..."
    
    # 创建主目录
    mkdir -p "$PROJECT_DIR"
    mkdir -p "$BACKUP_DIR"
    
    # 设置权限
    chown -R $SERVICE_USER:$SERVICE_USER "$PROJECT_DIR"
    
    print_success "项目目录创建完成: $PROJECT_DIR"
}

# 下载项目代码
download_project() {
    print_info "准备项目代码..."
    
    cd "$PROJECT_DIR"
    
    # 如果是Git仓库，则拉取代码
    if [ -d ".git" ]; then
        print_info "检测到Git仓库，更新代码..."
        git pull origin main || print_warning "Git拉取失败，请手动更新"
    else
        print_warning "未检测到Git仓库"
        print_info "请手动上传项目文件到: $PROJECT_DIR"
        print_info "或使用以下命令克隆代码:"
        print_info "  cd $PROJECT_DIR"
        print_info "  git clone <your-repository-url> ."
        read -p "是否已上传项目文件? (y/N): " uploaded
        if [[ ! "$uploaded" =~ ^[Yy]$ ]]; then
            print_error "请上传项目文件后重新运行脚本"
            exit 1
        fi
    fi
}

# 配置环境变量
setup_config() {
    print_info "配置环境变量..."
    
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            print_info "已创建 .env 文件，基于 .env.example"
        else
            print_warning "未找到 .env.example 文件"
            cat > .env << EOF
# SafeW Bot 配置文件

# 必需配置
SAFEW_BOT_TOKEN=your_bot_token_here

# 可选配置
LOG_LEVEL=INFO
POLL_TIMEOUT=30

# 默认转发目标群组ID (可选)
# FORWARD_TARGET_CHAT=-1001234567890
EOF
            print_info "已创建默认 .env 文件"
        fi
        
        print_warning "请编辑 .env 文件，设置您的Bot Token"
        print_info "编辑命令: nano .env"
        print_info "或通过宝塔面板文件管理器编辑"
        
        read -p "是否现在编辑配置文件? (y/N): " edit_config
        if [[ "$edit_config" =~ ^[Yy]$ ]]; then
            nano .env || vim .env || print_warning "请手动编辑 .env 文件"
        fi
    else
        print_success "发现现有 .env 文件"
    fi
    
    # 设置权限
    chmod 644 .env
    chown $SERVICE_USER:$SERVICE_USER .env
}

# 编译项目
build_project() {
    print_info "编译项目..."
    
    cd "$PROJECT_DIR"
    
    # 下载依赖
    go mod tidy
    
    # 编译
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $PROJECT_NAME
    
    if [ ! -f "$PROJECT_NAME" ]; then
        print_error "编译失败"
        exit 1
    fi
    
    # 设置权限
    chmod +x "$PROJECT_NAME"
    chown $SERVICE_USER:$SERVICE_USER "$PROJECT_NAME"
    
    print_success "项目编译完成"
}

# 创建systemd服务
create_systemd_service() {
    print_info "创建系统服务..."
    
    cat > "/etc/systemd/system/${PROJECT_NAME}.service" << EOF
[Unit]
Description=SafeW Bot Service
After=network.target
Wants=network.target

[Service]
Type=simple
User=$SERVICE_USER
Group=$SERVICE_USER
WorkingDirectory=$PROJECT_DIR
ExecStart=$PROJECT_DIR/$PROJECT_NAME
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$PROJECT_NAME

[Install]
WantedBy=multi-user.target
EOF

    # 重载systemd
    systemctl daemon-reload
    systemctl enable "$PROJECT_NAME"
    
    print_success "系统服务创建完成"
}

# 备份现有服务
backup_service() {
    if systemctl is-active --quiet "$PROJECT_NAME"; then
        print_info "备份现有服务..."
        
        # 创建备份目录
        BACKUP_TIME=$(date +%Y%m%d_%H%M%S)
        CURRENT_BACKUP_DIR="$BACKUP_DIR/$BACKUP_TIME"
        mkdir -p "$CURRENT_BACKUP_DIR"
        
        # 停止服务
        systemctl stop "$PROJECT_NAME"
        
        # 备份文件
        cp "$PROJECT_DIR/$PROJECT_NAME" "$CURRENT_BACKUP_DIR/" 2>/dev/null || true
        cp "$PROJECT_DIR/.env" "$CURRENT_BACKUP_DIR/" 2>/dev/null || true
        
        print_success "服务已备份到: $CURRENT_BACKUP_DIR"
    fi
}

# 启动服务
start_service() {
    print_info "启动服务..."
    
    systemctl start "$PROJECT_NAME"
    
    if systemctl is-active --quiet "$PROJECT_NAME"; then
        print_success "服务启动成功"
        print_info "服务状态:"
        systemctl status "$PROJECT_NAME" --no-pager -l
    else
        print_error "服务启动失败"
        print_info "查看日志: journalctl -u $PROJECT_NAME -f"
        exit 1
    fi
}

# 显示部署信息
show_deploy_info() {
    print_success "🎉 SafeW Bot 部署完成!"
    echo
    print_info "部署信息:"
    echo "  项目目录: $PROJECT_DIR"
    echo "  配置文件: $PROJECT_DIR/.env"
    echo "  可执行文件: $PROJECT_DIR/$PROJECT_NAME"
    echo "  服务名称: $PROJECT_NAME"
    echo
    print_info "常用命令:"
    echo "  启动服务: systemctl start $PROJECT_NAME"
    echo "  停止服务: systemctl stop $PROJECT_NAME"
    echo "  重启服务: systemctl restart $PROJECT_NAME"
    echo "  查看状态: systemctl status $PROJECT_NAME"
    echo "  查看日志: journalctl -u $PROJECT_NAME -f"
    echo
    print_info "宝塔面板管理:"
    echo "  安装进程守护器并添加进程管理"
    echo "  启动文件: $PROJECT_DIR/$PROJECT_NAME"
    echo "  运行目录: $PROJECT_DIR"
    echo "  运行用户: $SERVICE_USER"
    echo
    print_warning "重要提醒:"
    echo "  1. 请确保已正确配置 .env 文件中的 SAFEW_BOT_TOKEN"
    echo "  2. 建议定期备份配置文件和项目目录"
    echo "  3. 可通过宝塔面板的进程守护器进行可视化管理"
}

# 主函数
main() {
    echo "==============================================="
    echo "          SafeW Bot 宝塔环境部署脚本"
    echo "==============================================="
    echo
    
    check_root
    check_bt_panel
    check_go
    
    print_info "开始部署 SafeW Bot..."
    
    backup_service
    create_directories
    download_project
    setup_config
    build_project
    create_systemd_service
    start_service
    show_deploy_info
    
    echo
    print_success "部署完成! 🚀"
}

# 错误处理
trap 'print_error "部署过程中发生错误，退出码: $?"' ERR

# 运行主函数
main "$@" 