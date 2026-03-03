package postgres

import (
	"context"
	"fmt"

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
	query := `INSERT INTO messages (id, channel_id, user_id, content, type, root_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query,
		msg.ID, msg.ChannelID, msg.UserID, msg.Content,
		msg.Type, msg.RootID, msg.CreatedAt, msg.UpdatedAt)
	return err
}

func (s *MessageStore) GetByID(ctx context.Context, id string) (*model.Message, error) {
	var msg model.Message
	err := s.db.GetContext(ctx, &msg, "SELECT * FROM messages WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}
	return &msg, nil
}

func (s *MessageStore) GetThread(ctx context.Context, rootID string) ([]*model.Message, error) {
	var messages []*model.Message
	query := `SELECT * FROM messages WHERE id = $1 OR root_id = $1 ORDER BY created_at ASC`
	err := s.db.SelectContext(ctx, &messages, query, rootID)
	return messages, err
}

func (s *MessageStore) IncrementReplyCount(ctx context.Context, messageID string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE messages SET reply_count = reply_count + 1 WHERE id = $1", messageID)
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

func (s *MessageStore) DeleteExpiredByChannel(ctx context.Context) (int64, error) {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM messages WHERE id IN (
			SELECT m.id FROM messages m
			JOIN channels c ON m.channel_id = c.id
			WHERE c.retention_days > 0
			  AND m.created_at < NOW() - MAKE_INTERVAL(days => c.retention_days)
		)`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *MessageStore) DeleteExpiredByTeam(ctx context.Context) (int64, error) {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM messages WHERE id IN (
			SELECT m.id FROM messages m
			JOIN channels c ON m.channel_id = c.id
			JOIN teams t ON c.team_id = t.id
			WHERE t.retention_days > 0
			  AND c.retention_days = 0
			  AND m.created_at < NOW() - MAKE_INTERVAL(days => t.retention_days)
		)`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
