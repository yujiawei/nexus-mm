package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type ReactionStore struct {
	db *sqlx.DB
}

func NewReactionStore(db *sqlx.DB) *ReactionStore {
	return &ReactionStore{db: db}
}

func (s *ReactionStore) Create(ctx context.Context, r *model.Reaction) error {
	query := `INSERT INTO reactions (id, message_id, user_id, emoji_name, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (message_id, user_id, emoji_name) DO NOTHING`
	_, err := s.db.ExecContext(ctx, query,
		r.ID, r.MessageID, r.UserID, r.EmojiName, r.CreatedAt)
	return err
}

func (s *ReactionStore) Delete(ctx context.Context, messageID, userID, emojiName string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM reactions WHERE message_id = $1 AND user_id = $2 AND emoji_name = $3",
		messageID, userID, emojiName)
	return err
}

func (s *ReactionStore) ListByMessage(ctx context.Context, messageID string) ([]*model.Reaction, error) {
	var reactions []*model.Reaction
	err := s.db.SelectContext(ctx, &reactions,
		"SELECT * FROM reactions WHERE message_id = $1 ORDER BY created_at ASC", messageID)
	return reactions, err
}
