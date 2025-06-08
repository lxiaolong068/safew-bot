# SafeW Bot Go语言开发文档

## 1. 项目概述

本项目旨在使用 Go 语言创建一个 SafeW Bot，该 Bot 将具备消息转发和群组管理的核心功能。Bot 将通过 SafeW Bot API 与平台进行交互。

### 核心功能
1. **消息转发**: 能够识别并转发包含文本、图片或链接的消息
2. **群组管理**: 提供一系列命令来管理群聊，例如禁言、踢出成员、设置管理员等

### API 基础信息
- **Base URL**: `https://api.safew.org/bot<token>/`
- **认证方式**: Bot Token 认证
- **请求方式**: GET 和 POST 方法
- **响应格式**: JSON

## 2. 技术选型

### 架构方案
采用**轻量级混合方案**：以Go原生能力为核心，辅助少量实用库。

### 核心技术栈
- **编程语言**: Go
- **核心库**（Go 标准库）:
  - `net/http`: 用于向 SafeW API 发送 HTTP 请求
  - `encoding/json`: 用于解析 API 返回的 JSON 数据和构建请求体
  - `time`: 用于处理轮询 `getUpdates` 的间隔和超时控制
  - `context`: 用于请求上下文管理和优雅关闭
  - `log`: 用于基础日志记录
  - `os`: 用于读取环境变量（如 Bot Token）

### 辅助实用库（可选）
- `github.com/go-resty/resty/v2`: 简化HTTP客户端代码，提供更友好的API
- `github.com/sirupsen/logrus`: 结构化日志，提供更丰富的日志功能
- `github.com/spf13/viper`: 配置管理库，支持多种配置格式

### 方案优势
- **简单可控**: 主要依赖标准库，逻辑清晰
- **高性能**: Go原生HTTP客户端性能优秀
- **低依赖**: 避免过度工程化，减少潜在问题
- **易维护**: 代码易于理解、调试和扩展
- **快速开发**: 无需学习复杂框架，上手快

## 3. API 端点分析

根据 SafeW Bot API 文档，我们将主要使用以下 API 端点：

### 核心功能 API
- `getMe`: 获取 Bot 基本信息
- `getUpdates`: 以长轮询方式接收用户消息和事件
- `sendMessage`: 用于发送文本回复

### 消息转发相关 API
- `forwardMessage`: 将指定消息转发到另一个聊天
- `copyMessage`: 复制消息到另一个聊天
- `sendPhoto`: 发送图片消息
- `sendVideo`: 发送视频消息
- `sendDocument`: 发送文档消息
- `sendAudio`: 发送音频消息
- `sendVoice`: 发送语音消息
- `sendMediaGroup`: 发送媒体组

### 群组管理相关 API
- `banChatMember`: 禁言或踢出用户
- `promoteChatMember`: 提升用户为管理员
- `setChatPermissions`: 设置群组成员的权限
- `restrictChatMember`: 限制用户在群组中的权限
- `getChatAdministrators`: 获取群管理员列表
- `getChatMember`: 获取特定成员的信息
- `getChatMemberCount`: 获取群成员数量
- `getChat`: 获取聊天信息
- `setChatTitle`: 设置群组标题
- `setChatAdministratorCustomTitle`: 设置管理员自定义标题
- `leaveChat`: 离开群组

### 其他功能 API
- `setWebhook`: 设置 Webhook
- `deleteWebhook`: 删除 Webhook
- `getWebhookInfo`: 获取 Webhook 信息
- `deleteMessage`: 删除消息
- `deleteMessages`: 批量删除消息
- `editMessageText`: 编辑消息文本
- `editMessageReplyMarkup`: 编辑消息回复标记
- `setMyCommands`: 设置 Bot 命令
- `getMyCommands`: 获取 Bot 命令
- `createChatInviteLink`: 创建群组邀请链接
- `editChatInviteLink`: 编辑群组邀请链接
- `revokeChatInviteLink`: 撤销群组邀请链接
- `approveChatJoinRequest`: 批准加群请求
- `declineChatJoinRequest`: 拒绝加群请求
- `answerCallbackQuery`: 回答回调查询

## 4. 项目结构

建议采用以下模块化的项目结构：

```
safew-bot/
├── go.mod                  # Go 模块文件
├── go.sum                  # Go 依赖校验文件
├── main.go                 # 程序入口，初始化和启动 Bot
├── config.go               # 配置加载 (Bot Token 等)
├── docs/                   # 文档目录
│   └── development-plan.md # 开发计划文档
├── bot/                    # Bot 核心包
│   ├── bot.go              # Bot 核心逻辑，处理更新循环
│   ├── api.go              # 封装对 SafeW Bot API 的调用
│   ├── models.go           # 定义 API 的数据结构 (Update, Message, User 等)
│   └── handlers.go         # 处理具体的命令和消息
└── README.md               # 项目说明、配置和使用指南
```

## 5. 开发步骤

### 第一阶段：基础框架搭建

1. **初始化项目**: 创建 `safew-bot` 目录，并初始化 Go 模块
2. **定义数据结构 (`bot/models.go`)**: 根据 API 文档，定义 `Update`, `Message`, `Chat`, `User` 等核心 Go `struct`
3. **配置文件 (`config.go`)**: 实现从环境变量或配置文件加载 Bot Token 的功能
4. **API 客户端 (`bot/api.go`)**:
   - 创建一个 `ApiClient` 结构体，包含 `httpClient` 和 `baseURL`
   - 实现 `NewApiClient(token string)` 函数
   - 封装 `getUpdates`, `sendMessage`, `forwardMessage` 等基础 API 调用方法
   - 方法应处理 HTTP 请求、参数编码和 JSON 响应解析

### 第二阶段：核心逻辑实现

5. **Bot 主循环 (`bot/bot.go`)**:
   - 创建 `Bot` 结构体，包含 `ApiClient` 和 `updateOffset`
   - 实现 `Start()` 方法，该方法内含一个 `for` 循环，持续调用 `api.GetUpdates()`
   - 在循环中，遍历收到的 `Update`，并将其分发给 `handlers` 进行处理
   - 维护 `updateOffset` 以避免重复处理消息

6. **命令与消息处理 (`bot/handlers.go`)**:
   - 实现 `HandleUpdate(update models.Update)` 函数，作为更新分发的中枢
   - 该函数检查消息类型（命令、普通消息等）
   - 如果是命令（如 `/start`, `/help`），则调用相应的命令处理器
   - 如果是需要转发的消息，则调用转发处理器

### 第三阶段：功能实现 - 消息转发

7. **实现转发逻辑**:
   - 可以设计一个 `/forward` 命令或通过回复消息来触发转发
   - 处理器调用 `api.ForwardMessage()`，需要提供 `chat_id` (目标), `from_chat_id` (来源), 和 `message_id`
   - 需要考虑如何让用户指定转发的目标。一个简单的实现可以是转发到预设的 Channel 或 Chat ID
   - 支持不同类型的消息转发：文本、图片、视频、文档等

### 第四阶段：功能实现 - 群管理

8. **实现群管理命令**:
   - 为每个管理功能创建一个命令处理器，例如：
     - `/ban @user [reason]`: 调用 `api.BanChatMember()`
     - `/promote @user`: 调用 `api.PromoteChatMember()`
     - `/restrict @user`: 调用 `api.RestrictChatMember()`
     - `/info`: 获取群组信息
     - `/admins`: 获取管理员列表
   - **权限检查**: 在执行这些命令前，必须检查调用命令的用户是否是群管理员
   - 这需要先调用 `api.GetChatAdministrators()` 或 `api.GetChatMember()` 来验证权限
   - 处理器需要解析命令参数（如 `@username` 或 `user_id`）

### 第五阶段：完善与部署

9. **主程序入口 (`main.go`)**:
   - 调用 `config.Load()` 加载配置
   - 创建 `ApiClient` 实例
   - 创建 `Bot` 实例
   - 调用 `bot.Start()` 启动 Bot

10. **编写 `README.md`**: 详细说明如何获取 Bot Token，如何配置和运行此 Bot，以及所有可用的命令列表和用法

11. **构建与运行**:
    - 使用 `go build` 编译项目
    - 通过 `./safew-bot` 运行

## 6. 数据结构设计

### 核心数据结构

```go
// API 响应基础结构
type ApiResponse struct {
    Ok          bool   `json:"ok"`
    Description string `json:"description,omitempty"`
    ErrorCode   int    `json:"error_code,omitempty"`
}

// Update 结构
type Update struct {
    UpdateID int      `json:"update_id"`
    Message  *Message `json:"message,omitempty"`
    // 其他更新类型...
}

// Message 结构
type Message struct {
    MessageID int    `json:"message_id"`
    From      *User  `json:"from,omitempty"`
    Chat      *Chat  `json:"chat"`
    Date      int64  `json:"date"`
    Text      string `json:"text,omitempty"`
    // 媒体类型字段...
}

// User 结构
type User struct {
    ID        int64  `json:"id"`
    IsBot     bool   `json:"is_bot"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name,omitempty"`
    Username  string `json:"username,omitempty"`
}

// Chat 结构
type Chat struct {
    ID       int64  `json:"id"`
    Type     string `json:"type"`
    Title    string `json:"title,omitempty"`
    Username string `json:"username,omitempty"`
}
```

## 7. 配置管理

### 环境变量
- `SAFEW_BOT_TOKEN`: SafeW Bot 的认证 Token
- `FORWARD_TARGET_CHAT`: 默认转发目标聊天 ID
- `LOG_LEVEL`: 日志级别 (DEBUG, INFO, WARN, ERROR)

### 配置文件 (可选)
```json
{
  "bot_token": "YOUR_BOT_TOKEN",
  "forward_settings": {
    "default_target": "CHAT_ID",
    "auto_forward": false
  },
  "admin_settings": {
    "super_admins": ["USER_ID1", "USER_ID2"]
  }
}
```

## 8. 风险与注意事项

- **API 速率限制**: API 调用可能存在频率限制。需要设计合理的重试机制和错误处理逻辑
- **错误处理**: 完善的错误处理至关重要，特别是网络错误和 API 返回的 `{"ok": false, ...}` 错误
- **安全性**: Bot Token 必须保密，绝不能硬编码在代码中。推荐使用环境变量
- **长轮询超时**: `getUpdates` 的长轮询需要设置合理的超时时间
- **权限验证**: 群管理功能必须严格验证用户权限，防止误操作
- **消息类型处理**: 需要正确处理不同类型的消息（文本、图片、视频等）

## 9. 测试策略

- **单元测试**: 对 API 客户端和处理器进行单元测试
- **集成测试**: 测试完整的消息处理流程
- **手动测试**: 在实际 SafeW 环境中测试 Bot 功能

## 10. 后续扩展计划

- 添加更多媒体类型支持
- 实现 Webhook 模式
- 添加数据持久化功能
- 实现更复杂的权限管理系统
- 添加统计和监控功能 