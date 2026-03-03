package model

import "time"

type ChannelCategory struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	TeamID      string    `json:"team_id" db:"team_id"`
	DisplayName string    `json:"display_name" db:"display_name"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CategoryEntry struct {
	CategoryID string `json:"category_id" db:"category_id"`
	ChannelID  string `json:"channel_id" db:"channel_id"`
	SortOrder  int    `json:"sort_order" db:"sort_order"`
}

type CreateCategoryRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	DisplayName string `json:"display_name"`
	SortOrder   *int   `json:"sort_order"`
}

type AddCategoryEntryRequest struct {
	ChannelID string `json:"channel_id" binding:"required"`
	SortOrder int    `json:"sort_order"`
}
