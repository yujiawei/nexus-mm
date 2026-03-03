package model

import "time"

type Message struct {
	ID         string    `json:"id" db:"id"`
	ChannelID  string    `json:"channel_id" db:"channel_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Content    string    `json:"content" db:"content"`
	Type       string    `json:"type" db:"type"` // "text", "system"
	RootID     string    `json:"root_id,omitempty" db:"root_id"`
	ReplyCount int       `json:"reply_count" db:"reply_count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
	RootID  string `json:"root_id"`
}

type MessageListRequest struct {
	Before string `form:"before"`
	Limit  int    `form:"limit,default=50"`
}
