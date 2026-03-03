package model

import "time"

type Channel struct {
	ID          string    `json:"id" db:"id"`
	TeamID      string    `json:"team_id" db:"team_id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Type        string    `json:"type" db:"type"` // "open", "private", "direct"
	Purpose     string    `json:"purpose,omitempty" db:"purpose"`
	CreatorID   string    `json:"creator_id" db:"creator_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ChannelMember struct {
	ChannelID string    `json:"channel_id" db:"channel_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"` // "admin", "member"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateChannelRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=32"`
	DisplayName string `json:"display_name" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=open private"`
	Purpose     string `json:"purpose"`
}
