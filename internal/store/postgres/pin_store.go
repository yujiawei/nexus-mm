package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type PinStore struct {
	db *sqlx.DB
}

func NewPinStore(db *sqlx.DB) *PinStore {
	return &PinStore{db: db}
}

func (s *PinStore) Create(ctx context.Context, pin *model.PinnedMessage) error {
	query := `INSERT INTO pinned_messages (id, channel_id, message_id, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (channel_id, message_id) DO NOTHING`
	_, err := s.db.ExecContext(ctx, query,
		pin.ID, pin.ChannelID, pin.MessageID, pin.UserID, pin.CreatedAt)
	return err
}

func (s *PinStore) Delete(ctx context.Context, channelID, messageID string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM pinned_messages WHERE channel_id = $1 AND message_id = $2",
		channelID, messageID)
	return err
}

func (s *PinStore) ListByChannel(ctx context.Context, channelID string) ([]*model.PinnedMessage, error) {
	var pins []*model.PinnedMessage
	err := s.db.SelectContext(ctx, &pins,
		"SELECT * FROM pinned_messages WHERE channel_id = $1 ORDER BY created_at DESC", channelID)
	return pins, err
}
