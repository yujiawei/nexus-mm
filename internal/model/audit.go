package model

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID         string          `json:"id" db:"id"`
	UserID     string          `json:"user_id" db:"user_id"`
	Action     string          `json:"action" db:"action"`
	EntityType string          `json:"entity_type" db:"entity_type"`
	EntityID   string          `json:"entity_id" db:"entity_id"`
	IPAddr     string          `json:"ip_addr" db:"ip_addr"`
	Details    json.RawMessage `json:"details" db:"details"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

type AuditLogQuery struct {
	Action string `form:"action"`
	UserID string `form:"user_id"`
	Limit  int    `form:"limit,default=50"`
}
