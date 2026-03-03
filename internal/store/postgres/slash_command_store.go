package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type SlashCommandStore struct {
	db *sqlx.DB
}

func NewSlashCommandStore(db *sqlx.DB) *SlashCommandStore {
	return &SlashCommandStore{db: db}
}

func (s *SlashCommandStore) Create(ctx context.Context, cmd *model.SlashCommand) error {
	query := `INSERT INTO slash_commands (id, team_id, trigger, url, method, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query,
		cmd.ID, cmd.TeamID, cmd.Trigger, cmd.URL, cmd.Method,
		cmd.CreatorID, cmd.CreatedAt, cmd.UpdatedAt)
	return err
}

func (s *SlashCommandStore) GetByTrigger(ctx context.Context, teamID, trigger string) (*model.SlashCommand, error) {
	var cmd model.SlashCommand
	err := s.db.GetContext(ctx, &cmd,
		"SELECT * FROM slash_commands WHERE team_id = $1 AND trigger = $2", teamID, trigger)
	if err != nil {
		return nil, fmt.Errorf("command not found: %w", err)
	}
	return &cmd, nil
}

func (s *SlashCommandStore) ListByTeam(ctx context.Context, teamID string) ([]*model.SlashCommand, error) {
	var cmds []*model.SlashCommand
	err := s.db.SelectContext(ctx, &cmds,
		"SELECT * FROM slash_commands WHERE team_id = $1 ORDER BY trigger ASC", teamID)
	return cmds, err
}
