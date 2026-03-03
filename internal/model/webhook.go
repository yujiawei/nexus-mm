package model

import "time"

type IncomingWebhook struct {
	ID          string    `json:"id" db:"id"`
	ChannelID   string    `json:"channel_id" db:"channel_id"`
	TeamID      string    `json:"team_id" db:"team_id"`
	CreatorID   string    `json:"creator_id" db:"creator_id"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Description string    `json:"description,omitempty" db:"description"`
	Token       string    `json:"token" db:"token"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type OutgoingWebhook struct {
	ID           string   `json:"id" db:"id"`
	ChannelID    string   `json:"channel_id" db:"channel_id"`
	TeamID       string   `json:"team_id" db:"team_id"`
	CreatorID    string   `json:"creator_id" db:"creator_id"`
	DisplayName  string   `json:"display_name" db:"display_name"`
	Description  string   `json:"description,omitempty" db:"description"`
	TriggerWords []string `json:"trigger_words" db:"-"`
	CallbackURLs []string `json:"callback_urls" db:"-"`
	Token        string   `json:"token" db:"token"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type CreateIncomingWebhookRequest struct {
	ChannelID   string `json:"channel_id" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
}

type CreateOutgoingWebhookRequest struct {
	ChannelID    string   `json:"channel_id" binding:"required"`
	DisplayName  string   `json:"display_name" binding:"required"`
	Description  string   `json:"description"`
	TriggerWords []string `json:"trigger_words" binding:"required"`
	CallbackURLs []string `json:"callback_urls" binding:"required"`
}

type IncomingWebhookPayload struct {
	Text     string `json:"text" binding:"required"`
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
}

type OutgoingWebhookPayload struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	ChannelID   string `json:"channel_id"`
	UserID      string `json:"user_id"`
	Text        string `json:"text"`
	TriggerWord string `json:"trigger_word"`
	Timestamp   int64  `json:"timestamp"`
}
