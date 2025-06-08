package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// BaseURL SafeW Bot API 基础URL
	BaseURL = "https://api.safew.org/bot"
)

// ApiClient SafeW Bot API 客户端
type ApiClient struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

// NewApiClient 创建新的API客户端
func NewApiClient(token string) *ApiClient {
	return &ApiClient{
		token:   token,
		baseURL: BaseURL + token,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// makeRequest 发送HTTP请求的通用方法
func (client *ApiClient) makeRequest(ctx context.Context, method, endpoint string, params interface{}) (*ApiResponse, error) {
	url := fmt.Sprintf("%s/%s", client.baseURL, endpoint)

	var body io.Reader
	if params != nil {
		jsonData, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !apiResp.Ok {
		return &apiResp, fmt.Errorf("API error: %s (code: %d)", apiResp.Description, apiResp.ErrorCode)
	}

	return &apiResp, nil
}

// GetMe 获取Bot信息
func (client *ApiClient) GetMe(ctx context.Context) (*User, error) {
	resp, err := client.makeRequest(ctx, "GET", "getMe", nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp.Result, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetUpdatesParams getUpdates 方法的参数
type GetUpdatesParams struct {
	Offset  int `json:"offset,omitempty"`
	Limit   int `json:"limit,omitempty"`
	Timeout int `json:"timeout,omitempty"`
}

// GetUpdates 获取更新
func (client *ApiClient) GetUpdates(ctx context.Context, params GetUpdatesParams) ([]Update, error) {
	resp, err := client.makeRequest(ctx, "POST", "getUpdates", params)
	if err != nil {
		return nil, err
	}

	var updates []Update
	if err := json.Unmarshal(resp.Result, &updates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updates: %w", err)
	}

	return updates, nil
}

// SendMessageParams sendMessage 方法的参数
type SendMessageParams struct {
	ChatID                int64                `json:"chat_id"`
	Text                  string               `json:"text"`
	ParseMode             string               `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool                 `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool                 `json:"disable_notification,omitempty"`
	ReplyToMessageID      int                  `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           interface{}          `json:"reply_markup,omitempty"`
}

// SendMessage 发送消息
func (client *ApiClient) SendMessage(ctx context.Context, params SendMessageParams) (*Message, error) {
	resp, err := client.makeRequest(ctx, "POST", "sendMessage", params)
	if err != nil {
		return nil, err
	}

	var message Message
	if err := json.Unmarshal(resp.Result, &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &message, nil
}

// ForwardMessageParams forwardMessage 方法的参数
type ForwardMessageParams struct {
	ChatID              int64 `json:"chat_id"`
	FromChatID          int64 `json:"from_chat_id"`
	MessageID           int   `json:"message_id"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
}

// ForwardMessage 转发消息
func (client *ApiClient) ForwardMessage(ctx context.Context, params ForwardMessageParams) (*Message, error) {
	resp, err := client.makeRequest(ctx, "POST", "forwardMessage", params)
	if err != nil {
		return nil, err
	}

	var message Message
	if err := json.Unmarshal(resp.Result, &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &message, nil
}

// GetChat 获取聊天信息
func (client *ApiClient) GetChat(ctx context.Context, chatID int64) (*Chat, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	resp, err := client.makeRequest(ctx, "POST", "getChat", params)
	if err != nil {
		return nil, err
	}

	var chat Chat
	if err := json.Unmarshal(resp.Result, &chat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chat: %w", err)
	}

	return &chat, nil
}

// GetChatAdministrators 获取聊天管理员列表
func (client *ApiClient) GetChatAdministrators(ctx context.Context, chatID int64) ([]ChatMember, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	resp, err := client.makeRequest(ctx, "POST", "getChatAdministrators", params)
	if err != nil {
		return nil, err
	}

	var admins []ChatMember
	if err := json.Unmarshal(resp.Result, &admins); err != nil {
		return nil, fmt.Errorf("failed to unmarshal administrators: %w", err)
	}

	return admins, nil
}

// GetChatMember 获取聊天成员信息
func (client *ApiClient) GetChatMember(ctx context.Context, chatID, userID int64) (*ChatMember, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}

	resp, err := client.makeRequest(ctx, "POST", "getChatMember", params)
	if err != nil {
		return nil, err
	}

	var member ChatMember
	if err := json.Unmarshal(resp.Result, &member); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chat member: %w", err)
	}

	return &member, nil
}

// BanChatMemberParams banChatMember 方法的参数
type BanChatMemberParams struct {
	ChatID         int64 `json:"chat_id"`
	UserID         int64 `json:"user_id"`
	UntilDate      int64 `json:"until_date,omitempty"`
	RevokeMessages bool  `json:"revoke_messages,omitempty"`
}

// BanChatMember 禁言或踢出聊天成员
func (client *ApiClient) BanChatMember(ctx context.Context, params BanChatMemberParams) error {
	_, err := client.makeRequest(ctx, "POST", "banChatMember", params)
	return err
}

// PromoteChatMemberParams promoteChatMember 方法的参数
type PromoteChatMemberParams struct {
	ChatID              int64  `json:"chat_id"`
	UserID              int64  `json:"user_id"`
	IsAnonymous         bool   `json:"is_anonymous,omitempty"`
	CanManageChat       bool   `json:"can_manage_chat,omitempty"`
	CanPostMessages     bool   `json:"can_post_messages,omitempty"`
	CanEditMessages     bool   `json:"can_edit_messages,omitempty"`
	CanDeleteMessages   bool   `json:"can_delete_messages,omitempty"`
	CanManageVoiceChats bool   `json:"can_manage_voice_chats,omitempty"`
	CanRestrictMembers  bool   `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   bool   `json:"can_promote_members,omitempty"`
	CanChangeInfo       bool   `json:"can_change_info,omitempty"`
	CanInviteUsers      bool   `json:"can_invite_users,omitempty"`
	CanPinMessages      bool   `json:"can_pin_messages,omitempty"`
}

// PromoteChatMember 提升聊天成员为管理员
func (client *ApiClient) PromoteChatMember(ctx context.Context, params PromoteChatMemberParams) error {
	_, err := client.makeRequest(ctx, "POST", "promoteChatMember", params)
	return err
}

// RestrictChatMemberParams restrictChatMember 方法的参数
type RestrictChatMemberParams struct {
	ChatID      int64            `json:"chat_id"`
	UserID      int64            `json:"user_id"`
	Permissions *ChatPermissions `json:"permissions"`
	UntilDate   int64            `json:"until_date,omitempty"`
}

// RestrictChatMember 限制聊天成员权限
func (client *ApiClient) RestrictChatMember(ctx context.Context, params RestrictChatMemberParams) error {
	_, err := client.makeRequest(ctx, "POST", "restrictChatMember", params)
	return err
}

// DeleteMessage 删除消息
func (client *ApiClient) DeleteMessage(ctx context.Context, chatID int64, messageID int) error {
	params := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	}

	_, err := client.makeRequest(ctx, "POST", "deleteMessage", params)
	return err
} 