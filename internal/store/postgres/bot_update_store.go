package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type BotUpdateStore struct {
	db *sqlx.DB
}

func NewBotUpdateStore(db *sqlx.DB) *BotUpdateStore {
	return &BotUpdateStore{db: db}
}

func (s *BotUpdateStore) Create(ctx context.Context, botUserID, messageID, channelID string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO bot_updates (bot_user_id, message_id, channel_id) VALUES ($1, $2, $3)`,
		botUserID, messageID, channelID)
	return err
}

func (s *BotUpdateStore) ListSince(ctx context.Context, botUserID string, offset int64, limit int) ([]*model.BotUpdate, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	var updates []*model.BotUpdate
	err := s.db.SelectContext(ctx, &updates,
		`SELECT * FROM bot_updates WHERE bot_user_id = $1 AND id > $2 ORDER BY id ASC LIMIT $3`,
		botUserID, offset, limit)
	return updates, err
}

func (s *BotUpdateStore) ListBotIDsByChannel(ctx context.Context, channelID string) ([]string, error) {
	var botIDs []string
	err := s.db.SelectContext(ctx, &botIDs,
		`SELECT u.id FROM users u
		 JOIN channel_members cm ON cm.user_id = u.id
		 WHERE cm.channel_id = $1 AND u.is_bot = true`, channelID)
	return botIDs, err
}
