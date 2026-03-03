package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type CategoryStore struct {
	db *sqlx.DB
}

func NewCategoryStore(db *sqlx.DB) *CategoryStore {
	return &CategoryStore{db: db}
}

func (s *CategoryStore) Create(ctx context.Context, cat *model.ChannelCategory) error {
	query := `INSERT INTO channel_categories (id, user_id, team_id, display_name, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, query,
		cat.ID, cat.UserID, cat.TeamID, cat.DisplayName,
		cat.SortOrder, cat.CreatedAt, cat.UpdatedAt)
	return err
}

func (s *CategoryStore) Update(ctx context.Context, id, displayName string, sortOrder int) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE channel_categories SET display_name = $1, sort_order = $2, updated_at = NOW() WHERE id = $3",
		displayName, sortOrder, id)
	return err
}

func (s *CategoryStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM channel_categories WHERE id = $1", id)
	return err
}

func (s *CategoryStore) GetByID(ctx context.Context, id string) (*model.ChannelCategory, error) {
	var cat model.ChannelCategory
	err := s.db.GetContext(ctx, &cat, "SELECT * FROM channel_categories WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *CategoryStore) ListByUserTeam(ctx context.Context, userID, teamID string) ([]*model.ChannelCategory, error) {
	var cats []*model.ChannelCategory
	err := s.db.SelectContext(ctx, &cats,
		"SELECT * FROM channel_categories WHERE user_id = $1 AND team_id = $2 ORDER BY sort_order ASC",
		userID, teamID)
	return cats, err
}

func (s *CategoryStore) AddEntry(ctx context.Context, entry *model.CategoryEntry) error {
	query := `INSERT INTO channel_category_entries (category_id, channel_id, sort_order)
		VALUES ($1, $2, $3)
		ON CONFLICT (category_id, channel_id) DO UPDATE SET sort_order = $3`
	_, err := s.db.ExecContext(ctx, query,
		entry.CategoryID, entry.ChannelID, entry.SortOrder)
	return err
}

func (s *CategoryStore) RemoveEntry(ctx context.Context, categoryID, channelID string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM channel_category_entries WHERE category_id = $1 AND channel_id = $2",
		categoryID, channelID)
	return err
}

func (s *CategoryStore) ListEntries(ctx context.Context, categoryID string) ([]*model.CategoryEntry, error) {
	var entries []*model.CategoryEntry
	err := s.db.SelectContext(ctx, &entries,
		"SELECT * FROM channel_category_entries WHERE category_id = $1 ORDER BY sort_order ASC",
		categoryID)
	return entries, err
}
