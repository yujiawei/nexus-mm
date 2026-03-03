package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
)

type TeamService struct {
	store      *postgres.TeamStore
	inviteStore *postgres.InviteLinkStore
}

func NewTeamService(store *postgres.TeamStore, inviteStore *postgres.InviteLinkStore) *TeamService {
	return &TeamService{store: store, inviteStore: inviteStore}
}

func (s *TeamService) Create(ctx context.Context, req *model.CreateTeamRequest, creatorID string) (*model.Team, error) {
	now := time.Now().UTC()
	team := &model.Team{
		ID:          ulid.Make().String(),
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatorID:   creatorID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.store.Create(ctx, team); err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}

	// Add creator as owner.
	member := &model.TeamMember{
		TeamID:    team.ID,
		UserID:    creatorID,
		Role:      "owner",
		CreatedAt: now,
	}
	if err := s.store.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("add team owner: %w", err)
	}

	return team, nil
}

func (s *TeamService) GetByID(ctx context.Context, id string) (*model.Team, error) {
	return s.store.GetByID(ctx, id)
}

func (s *TeamService) List(ctx context.Context, userID string) ([]*model.Team, error) {
	return s.store.List(ctx, userID)
}

func (s *TeamService) IsMember(ctx context.Context, teamID, userID string) (bool, error) {
	return s.store.IsMember(ctx, teamID, userID)
}

func (s *TeamService) SetRetention(ctx context.Context, teamID string, days int) error {
	return s.store.SetRetention(ctx, teamID, days)
}

func (s *TeamService) ListAll(ctx context.Context) ([]*model.Team, error) {
	return s.store.ListAll(ctx)
}

func (s *TeamService) AddMember(ctx context.Context, teamID, userID, role string) error {
	now := time.Now().UTC()
	member := &model.TeamMember{
		TeamID:    teamID,
		UserID:    userID,
		Role:      role,
		CreatedAt: now,
	}
	return s.store.AddMember(ctx, member)
}

func (s *TeamService) ListMembers(ctx context.Context, teamID string) ([]*model.TeamMember, error) {
	return s.store.ListMembers(ctx, teamID)
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, userID string) error {
	return s.store.RemoveMember(ctx, teamID, userID)
}

func (s *TeamService) GetMemberRole(ctx context.Context, teamID, userID string) (string, error) {
	return s.store.GetMemberRole(ctx, teamID, userID)
}

func (s *TeamService) CreateInviteLink(ctx context.Context, teamID, creatorID string, req *model.CreateInviteLinkRequest) (*model.InviteLink, error) {
	code := generateCode(8)
	now := time.Now().UTC()
	link := &model.InviteLink{
		ID:        ulid.Make().String(),
		TeamID:    teamID,
		Code:      code,
		CreatorID: creatorID,
		MaxUses:   req.MaxUses,
		UseCount:  0,
		CreatedAt: now,
	}
	if req.ExpireDay > 0 {
		exp := now.Add(time.Duration(req.ExpireDay) * 24 * time.Hour)
		link.ExpiresAt = &exp
	}
	if err := s.inviteStore.Create(ctx, link); err != nil {
		return nil, fmt.Errorf("create invite link: %w", err)
	}
	return link, nil
}

func (s *TeamService) JoinByCode(ctx context.Context, code, userID string) (*model.Team, error) {
	link, err := s.inviteStore.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("invalid invite code")
	}
	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return nil, fmt.Errorf("invite link has expired")
	}
	if link.MaxUses > 0 && link.UseCount >= link.MaxUses {
		return nil, fmt.Errorf("invite link has reached max uses")
	}

	team, err := s.store.GetByID(ctx, link.TeamID)
	if err != nil {
		return nil, fmt.Errorf("team not found")
	}

	now := time.Now().UTC()
	member := &model.TeamMember{
		TeamID:    link.TeamID,
		UserID:    userID,
		Role:      "member",
		CreatedAt: now,
	}
	if err := s.store.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("join team: %w", err)
	}

	_ = s.inviteStore.IncrementUseCount(ctx, link.ID)
	return team, nil
}

func generateCode(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)[:length]
}
