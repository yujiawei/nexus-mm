package model

import "time"

type Team struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatorID   string    `json:"creator_id" db:"creator_id"`
	RetentionDays int       `json:"retention_days" db:"retention_days"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type SetRetentionRequest struct {
	RetentionDays int `json:"retention_days"`
}

type TeamMember struct {
	TeamID    string    `json:"team_id" db:"team_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"` // "owner", "admin", "member"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=32"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
}
