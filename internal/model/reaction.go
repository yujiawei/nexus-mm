package model

import "time"

type Reaction struct {
	ID        string    `json:"id" db:"id"`
	MessageID string    `json:"message_id" db:"message_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	EmojiName string    `json:"emoji_name" db:"emoji_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateReactionRequest struct {
	EmojiName string `json:"emoji_name" binding:"required"`
}
