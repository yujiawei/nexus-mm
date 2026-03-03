package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type ChannelStore struct {
	db *sqlx.DB
}

func NewChannelStore(db *sqlx.DB) *ChannelStore {
	return &ChannelStore{db: db}
}

func (s *ChannelStore) Create(ctx context.Context, ch *model.Channel) error {
	query := `INSERT INTO channels (id, team_id, name, display_name, type, purpose, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := s.db.ExecContext(ctx, query,
		ch.ID, ch.TeamID, ch.Name, ch.DisplayName, ch.Type,
		ch.Purpose, ch.CreatorID, ch.CreatedAt, ch.UpdatedAt)
	return err
}

func (s *ChannelStore) GetByID(ctx context.Context, id string) (*model.Channel, error) {
	var ch model.Channel
	err := s.db.GetContext(ctx, &ch, "SELECT * FROM channels WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}
	return &ch, nil
}

func (s *ChannelStore) ListByTeam(ctx context.Context, teamID, userID string) ([]*model.Channel, error) {
	var channels []*model.Channel
	query := `SELECT c.* FROM channels c
		LEFT JOIN channel_members cm ON c.id = cm.channel_id AND cm.user_id = $2
		WHERE c.team_id = $1 AND (c.type = 'open' OR cm.user_id IS NOT NULL)
		ORDER BY c.created_at ASC`
	err := s.db.SelectContext(ctx, &channels, query, teamID, userID)
	return channels, err
}

func (s *ChannelStore) AddMember(ctx context.Context, member *model.ChannelMember) error {
	query := `INSERT INTO channel_members (channel_id, user_id, role, created_at)
		VALUES ($1, $2, $3, $4) ON CONFLICT (channel_id, user_id) DO NOTHING`
	_, err := s.db.ExecContext(ctx, query,
		member.ChannelID, member.UserID, member.Role, member.CreatedAt)
	return err
}

func (s *ChannelStore) IsMember(ctx context.Context, channelID, userID string) (bool, error) {
	var count int
	err := s.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM channel_members WHERE channel_id = $1 AND user_id = $2",
		channelID, userID)
	return count > 0, err
}

func (s *ChannelStore) GetMembers(ctx context.Context, channelID string) ([]string, error) {
	var userIDs []string
	err := s.db.SelectContext(ctx, &userIDs,
		"SELECT user_id FROM channel_members WHERE channel_id = $1", channelID)
	return userIDs, err
}

func (s *ChannelStore) SetRetention(ctx context.Context, channelID string, days int) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE channels SET retention_days = $1, updated_at = NOW() WHERE id = $2", days, channelID)
	return err
}

func (s *ChannelStore) ListIDsByTeam(ctx context.Context, teamID string) ([]string, error) {
	var ids []string
	err := s.db.SelectContext(ctx, &ids,
		"SELECT id FROM channels WHERE team_id = $1", teamID)
	return ids, err
}

func (s *ChannelStore) ListMembers(ctx context.Context, channelID string) ([]*model.ChannelMember, error) {
	var members []*model.ChannelMember
	err := s.db.SelectContext(ctx, &members,
		"SELECT * FROM channel_members WHERE channel_id = $1 ORDER BY created_at ASC", channelID)
	return members, err
}

func (s *ChannelStore) RemoveMember(ctx context.Context, channelID, userID string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM channel_members WHERE channel_id = $1 AND user_id = $2", channelID, userID)
	return err
}
