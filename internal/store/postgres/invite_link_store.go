package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type InviteLinkStore struct {
	db *sqlx.DB
}

func NewInviteLinkStore(db *sqlx.DB) *InviteLinkStore {
	return &InviteLinkStore{db: db}
}

func (s *InviteLinkStore) Create(ctx context.Context, link *model.InviteLink) error {
	query := `INSERT INTO team_invite_links (id, team_id, code, creator_id, expires_at, max_uses, use_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query,
		link.ID, link.TeamID, link.Code, link.CreatorID,
		link.ExpiresAt, link.MaxUses, link.UseCount, link.CreatedAt)
	return err
}

func (s *InviteLinkStore) GetByCode(ctx context.Context, code string) (*model.InviteLink, error) {
	var link model.InviteLink
	err := s.db.GetContext(ctx, &link,
		"SELECT * FROM team_invite_links WHERE code = $1", code)
	if err != nil {
		return nil, fmt.Errorf("invite link not found: %w", err)
	}
	return &link, nil
}

func (s *InviteLinkStore) IncrementUseCount(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE team_invite_links SET use_count = use_count + 1 WHERE id = $1", id)
	return err
}
