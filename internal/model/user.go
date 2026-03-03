package model

import "time"

type User struct {
	ID             string    `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	Nickname       string    `json:"nickname" db:"nickname"`
	AvatarURL      string    `json:"avatar_url,omitempty" db:"avatar_url"`
	WkToken        string    `json:"wk_token,omitempty" db:"wk_token"`
	Role           string    `json:"role" db:"role"` // "admin", "member"
	IsBot          bool      `json:"is_bot" db:"is_bot"`
	BotToken       *string   `json:"bot_token,omitempty" db:"bot_token"`
	BotOwnerID     *string   `json:"bot_owner_id,omitempty" db:"bot_owner_id"`
	BotDescription string    `json:"bot_description,omitempty" db:"bot_description"`
	BotWebhookURL  string    `json:"bot_webhook_url,omitempty" db:"bot_webhook_url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
	WsURL string `json:"ws_url,omitempty"`
}

type BotRegisterRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=32"`
	Email string `json:"email" binding:"required,email"`
}

type BotRegisterResponse struct {
	BotUserID   string `json:"bot_user_id"`
	BotToken    string `json:"bot_token"`
	OwnerUserID string `json:"owner_user_id"`
	Message     string `json:"message"`
}

type BotBindRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type BotSendMessageRequest struct {
	ChannelID string `json:"channel_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	RootID    string `json:"root_id,omitempty"`
}

type BotSetWebhookRequest struct {
	URL string `json:"url"`
}

type BotSendReactionRequest struct {
	ChannelID string `json:"channel_id" binding:"required"`
	MessageID string `json:"message_id" binding:"required"`
	EmojiName string `json:"emoji_name" binding:"required"`
}

type BotUpdate struct {
	ID        int64    `json:"id" db:"id"`
	BotUserID string   `json:"bot_user_id" db:"bot_user_id"`
	MessageID string   `json:"message_id" db:"message_id"`
	ChannelID string   `json:"channel_id" db:"channel_id"`
	CreatedAt string   `json:"created_at" db:"created_at"`
	Message   *Message `json:"message,omitempty"`
}

type BotUpdateResponse struct {
	OK     bool               `json:"ok"`
	Result []BotUpdateMessage `json:"result"`
}

type BotUpdateMessage struct {
	UpdateID  int64  `json:"update_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CreateBotRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=32"`
	Description string `json:"description"`
}

type BotInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	WebhookURL  string `json:"webhook_url"`
	Token       string `json:"token,omitempty"`
	CreatedAt   string `json:"created_at"`
}
