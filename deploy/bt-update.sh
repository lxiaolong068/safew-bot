#!/bin/bash

# SafeW Bot 快速更新脚本
# 使用方法: chmod +x bt-update.sh && ./bt-update.sh

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

# 检查项目目录
check_project_dir() {
    if [ ! -d "$PROJECT_DIR" ]; then
        print_error "项目目录不存在: $PROJECT_DIR"
        print_info "请先运行部署脚本: ./bt-deploy.sh"
        exit 1
    fi
    
    if [ ! -f "$PROJECT_DIR/.env" ]; then
        print_error "配置文件不存在: $PROJECT_DIR/.env"
        exit 1
    fi
    
    cd "$PROJECT_DIR"
    print_success "项目目录检查通过"
}

# 备份当前版本
backup_current() {
    print_info "备份当前版本..."
    
    BACKUP_TIME=$(date +%Y%m%d_%H%M%S)
    CURRENT_BACKUP_DIR="$BACKUP_DIR/$BACKUP_TIME"
    mkdir -p "$CURRENT_BACKUP_DIR"
    
    # 备份可执行文件和配置
    cp "$PROJECT_NAME" "$CURRENT_BACKUP_DIR/" 2>/dev/null || true
    cp ".env" "$CURRENT_BACKUP_DIR/" 2>/dev/null || true
    
    print_success "已备份到: $CURRENT_BACKUP_DIR"
}

# 停止服务
stop_service() {
    print_info "停止服务..."
    
    if systemctl is-active --quiet "$PROJECT_NAME"; then
        systemctl stop "$PROJECT_NAME"
        print_success "服务已停止"
    else
        print_warning "服务未在运行"
    fi
}

# 更新代码
update_code() {
    print_info "更新代码..."
    
    if [ -d ".git" ]; then
        git fetch origin
        LOCAL=$(git rev-parse HEAD)
        REMOTE=$(git rev-parse origin/main)
        
        if [ "$LOCAL" = "$REMOTE" ]; then
            print_info "代码已是最新版本"
        else
            print_info "发现新版本，开始更新..."
            git pull origin main
            print_success "代码更新完成"
        fi
    else
        print_warning "非Git仓库，请手动更新代码文件"
        read -p "是否已手动更新代码文件? (y/N): " updated
        if [[ ! "$updated" =~ ^[Yy]$ ]]; then
            print_error "请更新代码文件后重新运行"
            exit 1
        fi
    fi
}

# 更新依赖和编译
build_project() {
    print_info "更新依赖和编译..."
    
    # 更新依赖
    go mod tidy
    
    # 编译新版本
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o "${PROJECT_NAME}.new"
    
    if [ ! -f "${PROJECT_NAME}.new" ]; then
        print_error "编译失败"
        exit 1
    fi
    
    # 替换旧版本
    mv "${PROJECT_NAME}.new" "$PROJECT_NAME"
    chmod +x "$PROJECT_NAME"
    chown www:www "$PROJECT_NAME"
    
    print_success "编译完成"
}

# 启动服务
start_service() {
    print_info "启动服务..."
    
    systemctl start "$PROJECT_NAME"
    
    # 等待服务启动
    sleep 3
    
    if systemctl is-active --quiet "$PROJECT_NAME"; then
        print_success "服务启动成功"
    else
        print_error "服务启动失败"
        print_info "查看日志: journalctl -u $PROJECT_NAME -f"
        
        # 尝试恢复备份
        print_info "尝试恢复备份..."
        LATEST_BACKUP=$(ls -t "$BACKUP_DIR" | head -1)
        if [ -n "$LATEST_BACKUP" ] && [ -f "$BACKUP_DIR/$LATEST_BACKUP/$PROJECT_NAME" ]; then
            cp "$BACKUP_DIR/$LATEST_BACKUP/$PROJECT_NAME" "./"
            systemctl start "$PROJECT_NAME"
            print_warning "已恢复到备份版本"
        fi
        
        exit 1
    fi
}

# 检查服务状态
check_service_status() {
    print_info "检查服务状态..."
    
    if systemctl is-active --quiet "$PROJECT_NAME"; then
        print_success "✅ 服务运行正常"
        
        # 显示进程信息
        PID=$(pgrep -f "$PROJECT_NAME" | head -1)
        if [ -n "$PID" ]; then
            print_info "进程ID: $PID"
            print_info "内存使用: $(ps -p $PID -o rss= | awk '{print $1/1024 " MB"}')"
        fi
        
        # 显示最新日志
        print_info "最新日志:"
        journalctl -u "$PROJECT_NAME" --since "2 minutes ago" --no-pager | tail -5
    else
        print_error "❌ 服务未运行"
        exit 1
    fi
}

# 显示更新信息
show_update_info() {
    print_success "🎉 SafeW Bot 更新完成!"
    echo
    print_info "服务信息:"
    systemctl status "$PROJECT_NAME" --no-pager -l | head -10
    echo
    print_info "常用命令:"
    echo "  查看状态: systemctl status $PROJECT_NAME"
    echo "  查看日志: journalctl -u $PROJECT_NAME -f"
    echo "  重启服务: systemctl restart $PROJECT_NAME"
    echo
    print_info "备份目录: $BACKUP_DIR"
    echo "  最新备份: $(ls -t "$BACKUP_DIR" | head -1)"
}

# 主函数
main() {
    echo "==============================================="
    echo "          SafeW Bot 快速更新脚本"
    echo "==============================================="
    echo
    
    check_project_dir
    backup_current
    stop_service
    update_code
    build_project
    start_service
    check_service_status
    show_update_info
    
    echo
    print_success "更新完成! 🚀"
}

# 错误处理
trap 'print_error "更新过程中发生错误，退出码: $?"' ERR

# 运行主函数
main "$@" 