package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	client *ApiClient
}

// NewMessageHandler 创建新的消息处理器
func NewMessageHandler(client *ApiClient) *MessageHandler {
	return &MessageHandler{
		client: client,
	}
}

// HandleMessage 处理普通消息
func (h *MessageHandler) HandleMessage(ctx context.Context, message *Message) error {
	// 忽略空消息
	if message == nil {
		return nil
	}

	log.Printf("收到消息: [%s] %s: %s", message.Chat.Type, getUserName(message.From), message.Text)

	// 检查是否为命令
	if strings.HasPrefix(message.Text, "/") {
		return h.handleCommand(ctx, message)
	}

	// 如果不是命令，处理普通消息
	return h.handleNormalMessage(ctx, message)
}

// HandleEditedMessage 处理编辑的消息
func (h *MessageHandler) HandleEditedMessage(ctx context.Context, message *Message) error {
	if message == nil {
		return nil
	}

	log.Printf("收到编辑消息: [%s] %s: %s", message.Chat.Type, getUserName(message.From), message.Text)
	
	// 对于编辑的消息，暂时只记录日志
	// 后续可以根据需要添加特殊处理逻辑
	return nil
}

// HandleCallbackQuery 处理回调查询
func (h *MessageHandler) HandleCallbackQuery(ctx context.Context, query *CallbackQuery) error {
	if query == nil {
		return nil
	}

	log.Printf("收到回调查询: %s 点击了 %s", getUserName(query.From), query.Data)
	
	// 这里可以根据callback_data处理不同的按钮点击
	// 暂时只记录日志
	return nil
}

// HandleChatJoinRequest 处理加群请求
func (h *MessageHandler) HandleChatJoinRequest(ctx context.Context, request *ChatJoinRequest) error {
	if request == nil {
		return nil
	}

	log.Printf("收到加群请求: %s 想加入 %s", getUserName(request.From), request.Chat.Title)
	
	// 这里可以实现自动审批逻辑
	// 暂时只记录日志
	return nil
}

// handleCommand 处理命令
func (h *MessageHandler) handleCommand(ctx context.Context, message *Message) error {
	// 解析命令和参数
	parts := strings.Fields(message.Text)
	if len(parts) == 0 {
		return nil
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	switch command {
	case "/start":
		return h.handleStartCommand(ctx, message)
	case "/help":
		return h.handleHelpCommand(ctx, message)
	case "/info":
		return h.handleInfoCommand(ctx, message)
	case "/forward":
		return h.handleForwardCommand(ctx, message, args)
	case "/ban":
		return h.handleBanCommand(ctx, message, args)
	case "/promote":
		return h.handlePromoteCommand(ctx, message, args)
	case "/admins":
		return h.handleAdminsCommand(ctx, message)
	default:
		return h.handleUnknownCommand(ctx, message, command)
	}
}

// handleNormalMessage 处理普通消息
func (h *MessageHandler) handleNormalMessage(ctx context.Context, message *Message) error {
	// 这里可以实现自动转发或其他逻辑
	// 暂时只记录日志
	return nil
}

// handleStartCommand 处理 /start 命令
func (h *MessageHandler) handleStartCommand(ctx context.Context, message *Message) error {
	welcomeText := `🤖 欢迎使用 SafeW Bot！

我是一个功能强大的Bot，可以帮助您：
• 转发消息（支持图文和链接）
• 管理群组

使用 /help 查看所有可用命令。`

	return h.sendReply(ctx, message, welcomeText)
}

// handleHelpCommand 处理 /help 命令
func (h *MessageHandler) handleHelpCommand(ctx context.Context, message *Message) error {
	helpText := `📖 可用命令列表：

🔧 基础命令:
/start - 开始使用Bot
/help - 显示帮助信息
/info - 获取群组信息

📤 转发功能:
/forward <目标群ID> - 转发回复的消息到指定群组

👮‍♂️ 管理命令 (仅管理员):
/ban <@用户名> [原因] - 禁言用户
/promote <@用户名> - 提升用户为管理员
/admins - 查看管理员列表

💡 使用提示：
• 大部分管理命令需要管理员权限
• 转发功能支持图片、视频、文档等多种格式`

	return h.sendReply(ctx, message, helpText)
}

// handleInfoCommand 处理 /info 命令
func (h *MessageHandler) handleInfoCommand(ctx context.Context, message *Message) error {
	chat, err := h.client.GetChat(ctx, message.Chat.ID)
	if err != nil {
		return h.sendReply(ctx, message, "❌ 获取群组信息失败")
	}

	infoText := fmt.Sprintf(`📊 群组信息：

🏷️ 群组名称: %s
🆔 群组ID: %d
📝 类型: %s`,
		chat.Title,
		chat.ID,
		chat.Type)

	if chat.Description != "" {
		infoText += fmt.Sprintf("\n📄 描述: %s", chat.Description)
	}

	if chat.Username != "" {
		infoText += fmt.Sprintf("\n🔗 用户名: @%s", chat.Username)
	}

	return h.sendReply(ctx, message, infoText)
}

// handleForwardCommand 处理 /forward 命令
func (h *MessageHandler) handleForwardCommand(ctx context.Context, message *Message, args []string) error {
	if message.ReplyToMessage == nil {
		return h.sendReply(ctx, message, "❌ 请回复要转发的消息使用此命令")
	}

	if len(args) == 0 {
		return h.sendReply(ctx, message, "❌ 请指定目标群组ID\n用法: /forward <群组ID>")
	}

	targetChatID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return h.sendReply(ctx, message, "❌ 无效的群组ID")
	}

	// 转发消息
	params := ForwardMessageParams{
		ChatID:     targetChatID,
		FromChatID: message.Chat.ID,
		MessageID:  message.ReplyToMessage.MessageID,
	}

	_, err = h.client.ForwardMessage(ctx, params)
	if err != nil {
		return h.sendReply(ctx, message, "❌ 转发失败: "+err.Error())
	}

	return h.sendReply(ctx, message, "✅ 消息已成功转发")
}

// handleBanCommand 处理 /ban 命令
func (h *MessageHandler) handleBanCommand(ctx context.Context, message *Message, args []string) error {
	// 检查用户权限
	if !h.isUserAdmin(ctx, message.Chat.ID, message.From.ID) {
		return h.sendReply(ctx, message, "❌ 您没有管理员权限")
	}

	if len(args) == 0 {
		return h.sendReply(ctx, message, "❌ 请指定要禁言的用户\n用法: /ban <@用户名> [原因]")
	}

	// 解析用户ID (这里简化处理，实际应该支持@username)
	userIDStr := strings.TrimPrefix(args[0], "@")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return h.sendReply(ctx, message, "❌ 无效的用户ID")
	}

	// 执行禁言
	params := BanChatMemberParams{
		ChatID: message.Chat.ID,
		UserID: userID,
	}

	err = h.client.BanChatMember(ctx, params)
	if err != nil {
		return h.sendReply(ctx, message, "❌ 禁言失败: "+err.Error())
	}

	reason := "违反群规"
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}

	return h.sendReply(ctx, message, fmt.Sprintf("✅ 用户已被禁言\n原因: %s", reason))
}

// handlePromoteCommand 处理 /promote 命令
func (h *MessageHandler) handlePromoteCommand(ctx context.Context, message *Message, args []string) error {
	// 检查用户权限
	if !h.isUserAdmin(ctx, message.Chat.ID, message.From.ID) {
		return h.sendReply(ctx, message, "❌ 您没有管理员权限")
	}

	if len(args) == 0 {
		return h.sendReply(ctx, message, "❌ 请指定要提升的用户\n用法: /promote <@用户名>")
	}

	// 解析用户ID
	userIDStr := strings.TrimPrefix(args[0], "@")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return h.sendReply(ctx, message, "❌ 无效的用户ID")
	}

	// 提升为管理员
	params := PromoteChatMemberParams{
		ChatID:            message.Chat.ID,
		UserID:            userID,
		CanDeleteMessages: true,
		CanRestrictMembers: true,
		CanInviteUsers:    true,
		CanPinMessages:    true,
	}

	err = h.client.PromoteChatMember(ctx, params)
	if err != nil {
		return h.sendReply(ctx, message, "❌ 提升管理员失败: "+err.Error())
	}

	return h.sendReply(ctx, message, "✅ 用户已被提升为管理员")
}

// handleAdminsCommand 处理 /admins 命令
func (h *MessageHandler) handleAdminsCommand(ctx context.Context, message *Message) error {
	admins, err := h.client.GetChatAdministrators(ctx, message.Chat.ID)
	if err != nil {
		return h.sendReply(ctx, message, "❌ 获取管理员列表失败")
	}

	var adminList strings.Builder
	adminList.WriteString("👮‍♂️ 群组管理员列表:\n\n")

	for i, admin := range admins {
		adminList.WriteString(fmt.Sprintf("%d. %s", i+1, getUserName(admin.User)))
		if admin.Status == "creator" {
			adminList.WriteString(" 👑 (群主)")
		}
		adminList.WriteString("\n")
	}

	return h.sendReply(ctx, message, adminList.String())
}

// handleUnknownCommand 处理未知命令
func (h *MessageHandler) handleUnknownCommand(ctx context.Context, message *Message, command string) error {
	return h.sendReply(ctx, message, fmt.Sprintf("❓ 未知命令: %s\n使用 /help 查看可用命令", command))
}

// isUserAdmin 检查用户是否为管理员
func (h *MessageHandler) isUserAdmin(ctx context.Context, chatID, userID int64) bool {
	member, err := h.client.GetChatMember(ctx, chatID, userID)
	if err != nil {
		log.Printf("检查用户权限时出错: %v", err)
		return false
	}

	return member.Status == "administrator" || member.Status == "creator"
}

// sendReply 发送回复消息
func (h *MessageHandler) sendReply(ctx context.Context, originalMessage *Message, text string) error {
	params := SendMessageParams{
		ChatID: originalMessage.Chat.ID,
		Text:   text,
	}

	_, err := h.client.SendMessage(ctx, params)
	return err
}

// getUserName 获取用户显示名称
func getUserName(user *User) string {
	if user == nil {
		return "Unknown"
	}

	if user.Username != "" {
		return "@" + user.Username
	}

	name := user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}

	return name
} 