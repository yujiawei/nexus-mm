package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type TeamStore struct {
	db *sqlx.DB
}

func NewTeamStore(db *sqlx.DB) *TeamStore {
	return &TeamStore{db: db}
}

func (s *TeamStore) Create(ctx context.Context, team *model.Team) error {
	query := `INSERT INTO teams (id, name, display_name, description, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, query,
		team.ID, team.Name, team.DisplayName, team.Description,
		team.CreatorID, team.CreatedAt, team.UpdatedAt)
	return err
}

func (s *TeamStore) GetByID(ctx context.Context, id string) (*model.Team, error) {
	var team model.Team
	err := s.db.GetContext(ctx, &team, "SELECT * FROM teams WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}
	return &team, nil
}

func (s *TeamStore) List(ctx context.Context, userID string) ([]*model.Team, error) {
	var teams []*model.Team
	query := `SELECT t.* FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY t.created_at DESC`
	err := s.db.SelectContext(ctx, &teams, query, userID)
	return teams, err
}

func (s *TeamStore) AddMember(ctx context.Context, member *model.TeamMember) error {
	query := `INSERT INTO team_members (team_id, user_id, role, created_at)
		VALUES ($1, $2, $3, $4) ON CONFLICT (team_id, user_id) DO NOTHING`
	_, err := s.db.ExecContext(ctx, query,
		member.TeamID, member.UserID, member.Role, member.CreatedAt)
	return err
}

func (s *TeamStore) IsMember(ctx context.Context, teamID, userID string) (bool, error) {
	var count int
	err := s.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM team_members WHERE team_id = $1 AND user_id = $2",
		teamID, userID)
	return count > 0, err
}

func (s *TeamStore) SetRetention(ctx context.Context, teamID string, days int) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE teams SET retention_days = $1, updated_at = NOW() WHERE id = $2", days, teamID)
	return err
}
