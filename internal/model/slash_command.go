package model

import "time"

type SlashCommand struct {
	ID        string    `json:"id" db:"id"`
	TeamID    string    `json:"team_id" db:"team_id"`
	Trigger   string    `json:"trigger" db:"trigger"`
	URL       string    `json:"url" db:"url"`
	Method    string    `json:"method" db:"method"`
	CreatorID string    `json:"creator_id" db:"creator_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSlashCommandRequest struct {
	Trigger string `json:"trigger" binding:"required"`
	URL     string `json:"url" binding:"required"`
	Method  string `json:"method"`
}
