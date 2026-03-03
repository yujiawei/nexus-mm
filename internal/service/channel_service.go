package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
	"github.com/yujiawei/nexus-mm/internal/wkim"
)

type ChannelService struct {
	store *postgres.ChannelStore
	wk    *wkim.Client
}

func NewChannelService(store *postgres.ChannelStore, wk *wkim.Client) *ChannelService {
	return &ChannelService{store: store, wk: wk}
}

func (s *ChannelService) Create(ctx context.Context, teamID string, req *model.CreateChannelRequest, creatorID string) (*model.Channel, error) {
	now := time.Now().UTC()
	ch := &model.Channel{
		ID:          ulid.Make().String(),
		TeamID:      teamID,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Type:        req.Type,
		Purpose:     req.Purpose,
		CreatorID:   creatorID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.store.Create(ctx, ch); err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}

	// Add creator as channel admin.
	member := &model.ChannelMember{
		ChannelID: ch.ID,
		UserID:    creatorID,
		Role:      "admin",
		CreatedAt: now,
	}
	if err := s.store.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("add channel admin: %w", err)
	}

	// Create channel in WuKongIM and add the creator as subscriber.
	if err := s.wk.CreateChannel(ctx, ch.ID, wkim.ChannelTypeGroup, []string{creatorID}); err != nil {
		fmt.Printf("warn: create wukongim channel: %v\n", err)
	}

	return ch, nil
}

func (s *ChannelService) GetByID(ctx context.Context, id string) (*model.Channel, error) {
	return s.store.GetByID(ctx, id)
}

func (s *ChannelService) ListByTeam(ctx context.Context, teamID, userID string) ([]*model.Channel, error) {
	return s.store.ListByTeam(ctx, teamID, userID)
}

func (s *ChannelService) IsMember(ctx context.Context, channelID, userID string) (bool, error) {
	return s.store.IsMember(ctx, channelID, userID)
}

func (s *ChannelService) GetMembers(ctx context.Context, channelID string) ([]string, error) {
	return s.store.GetMembers(ctx, channelID)
}

func (s *ChannelService) SetRetention(ctx context.Context, channelID string, days int) error {
	return s.store.SetRetention(ctx, channelID, days)
}

func (s *ChannelService) ListIDsByTeam(ctx context.Context, teamID string) ([]string, error) {
	return s.store.ListIDsByTeam(ctx, teamID)
}

func (s *ChannelService) AddMember(ctx context.Context, channelID, userID, role string) error {
	now := time.Now().UTC()
	member := &model.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
		Role:      role,
		CreatedAt: now,
	}
	if err := s.store.AddMember(ctx, member); err != nil {
		return fmt.Errorf("add channel member: %w", err)
	}
	// Also add as subscriber in WuKongIM.
	if err := s.wk.AddSubscribers(ctx, channelID, wkim.ChannelTypeGroup, []string{userID}); err != nil {
		fmt.Printf("warn: add wukongim subscriber: %v\n", err)
	}
	return nil
}

func (s *ChannelService) ListMembers(ctx context.Context, channelID string) ([]*model.ChannelMember, error) {
	return s.store.ListMembers(ctx, channelID)
}

func (s *ChannelService) RemoveMember(ctx context.Context, channelID, userID string) error {
	if err := s.store.RemoveMember(ctx, channelID, userID); err != nil {
		return fmt.Errorf("remove channel member: %w", err)
	}
	// Also remove subscriber from WuKongIM.
	if err := s.wk.RemoveSubscribers(ctx, channelID, wkim.ChannelTypeGroup, []string{userID}); err != nil {
		fmt.Printf("warn: remove wukongim subscriber: %v\n", err)
	}
	return nil
}
