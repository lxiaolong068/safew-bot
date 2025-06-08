package bot

import "encoding/json"

// ApiResponse API响应基础结构
type ApiResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Result      json.RawMessage `json:"result,omitempty"`
}

// Update 更新结构
type Update struct {
	UpdateID         int               `json:"update_id"`
	Message          *Message          `json:"message,omitempty"`
	EditedMessage    *Message          `json:"edited_message,omitempty"`
	CallbackQuery    *CallbackQuery    `json:"callback_query,omitempty"`
	ChatJoinRequest  *ChatJoinRequest  `json:"chat_join_request,omitempty"`
}

// Message 消息结构
type Message struct {
	MessageID       int              `json:"message_id"`
	From            *User            `json:"from,omitempty"`
	Date            int64            `json:"date"`
	Chat            *Chat            `json:"chat"`
	ForwardFrom     *User            `json:"forward_from,omitempty"`
	ForwardDate     int64            `json:"forward_date,omitempty"`
	ReplyToMessage  *Message         `json:"reply_to_message,omitempty"`
	Text            string           `json:"text,omitempty"`
	Entities        []MessageEntity  `json:"entities,omitempty"`
	Photo           []PhotoSize      `json:"photo,omitempty"`
	Video           *Video           `json:"video,omitempty"`
	Document        *Document        `json:"document,omitempty"`
	Audio           *Audio           `json:"audio,omitempty"`
	Voice           *Voice           `json:"voice,omitempty"`
	Caption         string           `json:"caption,omitempty"`
	Contact         *Contact         `json:"contact,omitempty"`
	Location        *Location        `json:"location,omitempty"`
}

// User 用户结构
type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

// Chat 聊天结构
type Chat struct {
	ID                          int64            `json:"id"`
	Type                        string           `json:"type"`
	Title                       string           `json:"title,omitempty"`
	Username                    string           `json:"username,omitempty"`
	FirstName                   string           `json:"first_name,omitempty"`
	LastName                    string           `json:"last_name,omitempty"`
	AllMembersAreAdministrators bool             `json:"all_members_are_administrators,omitempty"`
	Description                 string           `json:"description,omitempty"`
	InviteLink                  string           `json:"invite_link,omitempty"`
	Permissions                 *ChatPermissions `json:"permissions,omitempty"`
}

// ChatMember 聊天成员结构
type ChatMember struct {
	User   *User  `json:"user"`
	Status string `json:"status"`
}

// ChatPermissions 聊天权限结构
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendMediaMessages  bool `json:"can_send_media_messages,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
}

// MessageEntity 消息实体结构
type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url,omitempty"`
	User   *User  `json:"user,omitempty"`
}

// PhotoSize 图片尺寸结构
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size,omitempty"`
}

// Video 视频结构
type Video struct {
	FileID   string     `json:"file_id"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Duration int        `json:"duration"`
	Thumb    *PhotoSize `json:"thumb,omitempty"`
	MimeType string     `json:"mime_type,omitempty"`
	FileSize int        `json:"file_size,omitempty"`
}

// Document 文档结构
type Document struct {
	FileID   string     `json:"file_id"`
	Thumb    *PhotoSize `json:"thumb,omitempty"`
	FileName string     `json:"file_name,omitempty"`
	MimeType string     `json:"mime_type,omitempty"`
	FileSize int        `json:"file_size,omitempty"`
}

// Audio 音频结构
type Audio struct {
	FileID    string     `json:"file_id"`
	Duration  int        `json:"duration"`
	Performer string     `json:"performer,omitempty"`
	Title     string     `json:"title,omitempty"`
	MimeType  string     `json:"mime_type,omitempty"`
	FileSize  int        `json:"file_size,omitempty"`
	Thumb     *PhotoSize `json:"thumb,omitempty"`
}

// Voice 语音结构
type Voice struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type,omitempty"`
	FileSize int    `json:"file_size,omitempty"`
}

// Contact 联系人结构
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
}

// Location 位置结构
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// CallbackQuery 回调查询结构
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message,omitempty"`
	InlineMessageID string   `json:"inline_message_id,omitempty"`
	Data            string   `json:"data,omitempty"`
}

// ChatJoinRequest 加群请求结构
type ChatJoinRequest struct {
	Chat       *Chat `json:"chat"`
	From       *User `json:"from"`
	Date       int64 `json:"date"`
	Bio        string `json:"bio,omitempty"`
	InviteLink string `json:"invite_link,omitempty"`
}

// InlineKeyboardMarkup 内联键盘标记
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton 内联键盘按钮
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

// ReplyKeyboardMarkup 回复键盘标记
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
	Selective       bool               `json:"selective,omitempty"`
}

// KeyboardButton 键盘按钮
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
} 