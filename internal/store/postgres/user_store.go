package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (id, username, email, password_hash, nickname, avatar_url, wk_token, role,
		is_bot, bot_token, bot_owner_id, bot_description, bot_webhook_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := s.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.Nickname, user.AvatarURL, user.WkToken, user.Role,
		user.IsBot, user.BotToken, user.BotOwnerID, user.BotDescription, user.BotWebhookURL,
		user.CreatedAt, user.UpdatedAt)
	return err
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (s *UserStore) Count(ctx context.Context) (int, error) {
	var count int
	err := s.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users")
	return count, err
}

func (s *UserStore) UpdateRole(ctx context.Context, userID, role string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2", role, userID)
	return err
}

func (s *UserStore) GetByBotToken(ctx context.Context, token string) (*model.User, error) {
	var user model.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE bot_token = $1 AND is_bot = true", token)
	if err != nil {
		return nil, fmt.Errorf("bot not found: %w", err)
	}
	return &user, nil
}

func (s *UserStore) ListBotsByOwner(ctx context.Context, ownerID string) ([]*model.User, error) {
	var bots []*model.User
	err := s.db.SelectContext(ctx, &bots,
		"SELECT * FROM users WHERE bot_owner_id = $1 AND is_bot = true ORDER BY created_at DESC", ownerID)
	return bots, err
}

func (s *UserStore) UpdateBotToken(ctx context.Context, botID, token string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE users SET bot_token = $1, updated_at = NOW() WHERE id = $2 AND is_bot = true", token, botID)
	return err
}

func (s *UserStore) UpdateBotWebhook(ctx context.Context, botID, url string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE users SET bot_webhook_url = $1, updated_at = NOW() WHERE id = $2 AND is_bot = true", url, botID)
	return err
}
