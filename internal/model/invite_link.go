package model

import "time"

type InviteLink struct {
	ID        string     `json:"id" db:"id"`
	TeamID    string     `json:"team_id" db:"team_id"`
	Code      string     `json:"code" db:"code"`
	CreatorID string     `json:"creator_id" db:"creator_id"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	MaxUses   int        `json:"max_uses" db:"max_uses"`
	UseCount  int        `json:"use_count" db:"use_count"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type CreateInviteLinkRequest struct {
	MaxUses   int `json:"max_uses"`
	ExpireDay int `json:"expire_days"` // 0 = never
}

type InviteLinkResponse struct {
	Code string `json:"code"`
	Link string `json:"link"`
}
