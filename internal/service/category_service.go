package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
)

type CategoryService struct {
	store *postgres.CategoryStore
}

func NewCategoryService(store *postgres.CategoryStore) *CategoryService {
	return &CategoryService{store: store}
}

func (s *CategoryService) Create(ctx context.Context, userID, teamID string, req *model.CreateCategoryRequest) (*model.ChannelCategory, error) {
	now := time.Now().UTC()
	cat := &model.ChannelCategory{
		ID:          ulid.Make().String(),
		UserID:      userID,
		TeamID:      teamID,
		DisplayName: req.DisplayName,
		SortOrder:   req.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.store.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return cat, nil
}

func (s *CategoryService) Update(ctx context.Context, id string, req *model.UpdateCategoryRequest) error {
	cat, err := s.store.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	displayName := cat.DisplayName
	sortOrder := cat.SortOrder
	if req.DisplayName != "" {
		displayName = req.DisplayName
	}
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	return s.store.Update(ctx, id, displayName, sortOrder)
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *CategoryService) List(ctx context.Context, userID, teamID string) ([]*model.ChannelCategory, error) {
	return s.store.ListByUserTeam(ctx, userID, teamID)
}

func (s *CategoryService) AddEntry(ctx context.Context, categoryID string, req *model.AddCategoryEntryRequest) (*model.CategoryEntry, error) {
	entry := &model.CategoryEntry{
		CategoryID: categoryID,
		ChannelID:  req.ChannelID,
		SortOrder:  req.SortOrder,
	}

	if err := s.store.AddEntry(ctx, entry); err != nil {
		return nil, fmt.Errorf("add category entry: %w", err)
	}
	return entry, nil
}

func (s *CategoryService) RemoveEntry(ctx context.Context, categoryID, channelID string) error {
	return s.store.RemoveEntry(ctx, categoryID, channelID)
}

func (s *CategoryService) ListEntries(ctx context.Context, categoryID string) ([]*model.CategoryEntry, error) {
	return s.store.ListEntries(ctx, categoryID)
}
