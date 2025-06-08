# SafeW Bot Makefile
# 简化编译和部署流程

# 项目配置
PROJECT_NAME := safew-bot
BUILD_DIR := build
VERSION := $(shell date +"%Y%m%d_%H%M%S")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Go编译配置
GOOS := linux
GOARCH := amd64
LDFLAGS := -w -s -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

# 默认目标
.PHONY: all
all: clean build

# 清理编译文件
.PHONY: clean
clean:
	@echo "🧹 清理编译文件..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(PROJECT_NAME)

# 本地编译
.PHONY: build-local
build-local:
	@echo "🔨 本地编译..."
	@go build -ldflags "$(LDFLAGS)" -o $(PROJECT_NAME) .
	@echo "✅ 本地编译完成: $(PROJECT_NAME)"

# 交叉编译
.PHONY: build
build:
	@echo "🔨 开始交叉编译 ($(GOOS)/$(GOARCH))..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(PROJECT_NAME) .
	@cp .env.example $(BUILD_DIR)/
	@cp README.md $(BUILD_DIR)/
	@cp -r deploy $(BUILD_DIR)/
	@./build.sh > /dev/null 2>&1 || true
	@echo "✅ 交叉编译完成！"

# 运行测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	@go test -v ./...

# 代码格式化
.PHONY: fmt
fmt:
	@echo "📝 格式化代码..."
	@go fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "🔍 代码检查..."
	@go vet ./...

# 完整检查
.PHONY: check
check: fmt vet test
	@echo "✅ 代码检查完成！"

# 本地运行
.PHONY: run
run: build-local
	@echo "🚀 启动 SafeW Bot..."
	@./$(PROJECT_NAME)

# 显示版本
.PHONY: version
version:
	@echo "SafeW Bot"
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"

# 上传到服务器 (需要参数: SERVER, USER, [DIR])
.PHONY: upload
upload:
	@if [ -z "$(SERVER)" ] || [ -z "$(USER)" ]; then \
		echo "❌ 缺少参数。使用方法: make upload SERVER=your_server USER=your_user [DIR=remote_dir]"; \
		exit 1; \
	fi
	@./upload.sh $(SERVER) $(USER) $(DIR)

# 一键部署 (编译+上传)
.PHONY: deploy
deploy: build upload

# 显示帮助
.PHONY: help
help:
	@echo "SafeW Bot 构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make build        - 交叉编译为Linux版本"
	@echo "  make build-local  - 本地编译"
	@echo "  make clean        - 清理编译文件"
	@echo "  make test         - 运行测试"
	@echo "  make fmt          - 格式化代码"
	@echo "  make vet          - 代码检查"
	@echo "  make check        - 完整检查 (fmt+vet+test)"
	@echo "  make run          - 本地运行"
	@echo "  make version      - 显示版本信息"
	@echo "  make upload       - 上传到服务器"
	@echo "  make deploy       - 一键部署 (build+upload)"
	@echo "  make help         - 显示此帮助"
	@echo ""
	@echo "上传示例:"
	@echo "  make upload SERVER=192.168.1.100 USER=root"
	@echo "  make upload SERVER=example.com USER=root DIR=/opt/safew-bot"
	@echo ""
	@echo "一键部署示例:"
	@echo "  make deploy SERVER=192.168.1.100 USER=root"

# 安装依赖
.PHONY: deps
deps:
	@echo "📦 安装依赖..."
	@go mod tidy
	@go mod download

# 更新依赖
.PHONY: update
update:
	@echo "🔄 更新依赖..."
	@go get -u ./...
	@go mod tidy 