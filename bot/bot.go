package bot

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Bot SafeW Bot 主结构
type Bot struct {
	client       *ApiClient
	updateOffset int
	pollTimeout  int
	handlers     *MessageHandler
}

// NewBot 创建新的Bot实例
func NewBot(token string, pollTimeout int) *Bot {
	client := NewApiClient(token)
	
	return &Bot{
		client:       client,
		updateOffset: 0,
		pollTimeout:  pollTimeout,
		handlers:     NewMessageHandler(client),
	}
}

// Start 启动Bot主循环
func (bot *Bot) Start(ctx context.Context) error {
	log.Println("正在启动SafeW Bot...")

	// 首先验证Bot Token
	user, err := bot.client.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("验证Bot Token失败: %w", err)
	}

	log.Printf("Bot已启动: %s (@%s)", user.FirstName, user.Username)

	// 开始长轮询循环
	for {
		select {
		case <-ctx.Done():
			log.Println("接收到停止信号，正在关闭Bot...")
			return ctx.Err()
		default:
			if err := bot.processUpdates(ctx); err != nil {
				log.Printf("处理更新时出错: %v", err)
				// 等待一段时间后重试
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// processUpdates 处理一轮更新
func (bot *Bot) processUpdates(ctx context.Context) error {
	params := GetUpdatesParams{
		Offset:  bot.updateOffset,
		Limit:   100,
		Timeout: bot.pollTimeout,
	}

	updates, err := bot.client.GetUpdates(ctx, params)
	if err != nil {
		return fmt.Errorf("获取更新失败: %w", err)
	}

	for _, update := range updates {
		// 更新offset
		bot.updateOffset = update.UpdateID + 1

		// 处理每个更新
		if err := bot.handleUpdate(ctx, update); err != nil {
			log.Printf("处理更新 %d 时出错: %v", update.UpdateID, err)
		}
	}

	return nil
}

// handleUpdate 处理单个更新
func (bot *Bot) handleUpdate(ctx context.Context, update Update) error {
	// 处理普通消息
	if update.Message != nil {
		return bot.handlers.HandleMessage(ctx, update.Message)
	}

	// 处理编辑消息
	if update.EditedMessage != nil {
		return bot.handlers.HandleEditedMessage(ctx, update.EditedMessage)
	}

	// 处理回调查询
	if update.CallbackQuery != nil {
		return bot.handlers.HandleCallbackQuery(ctx, update.CallbackQuery)
	}

	// 处理加群请求
	if update.ChatJoinRequest != nil {
		return bot.handlers.HandleChatJoinRequest(ctx, update.ChatJoinRequest)
	}

	// 如果没有处理任何类型的更新，记录日志
	log.Printf("收到未处理的更新类型: %+v", update)
	return nil
}

// Stop 停止Bot (优雅关闭)
func (bot *Bot) Stop() {
	log.Println("Bot正在关闭...")
} 