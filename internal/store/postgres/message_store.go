package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type MessageStore struct {
	db *sqlx.DB
}

func NewMessageStore(db *sqlx.DB) *MessageStore {
	return &MessageStore{db: db}
}

func (s *MessageStore) Create(ctx context.Context, msg *model.Message) error {
	query := `INSERT INTO messages (id, channel_id, user_id, content, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, query,
		msg.ID, msg.ChannelID, msg.UserID, msg.Content,
		msg.Type, msg.CreatedAt, msg.UpdatedAt)
	return err
}

func (s *MessageStore) ListByChannel(ctx context.Context, channelID string, before string, limit int) ([]*model.Message, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var messages []*model.Message
	var err error

	if before != "" {
		query := `SELECT * FROM messages
			WHERE channel_id = $1 AND created_at < (SELECT created_at FROM messages WHERE id = $2)
			ORDER BY created_at DESC LIMIT $3`
		err = s.db.SelectContext(ctx, &messages, query, channelID, before, limit)
	} else {
		query := `SELECT * FROM messages
			WHERE channel_id = $1
			ORDER BY created_at DESC LIMIT $2`
		err = s.db.SelectContext(ctx, &messages, query, channelID, limit)
	}
	return messages, err
}
