package service

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
)

type ReactionService struct {
	store *postgres.ReactionStore
}

func NewReactionService(store *postgres.ReactionStore) *ReactionService {
	return &ReactionService{store: store}
}

func (s *ReactionService) Add(ctx context.Context, messageID, userID string, req *model.CreateReactionRequest) (*model.Reaction, error) {
	r := &model.Reaction{
		ID:        ulid.Make().String(),
		MessageID: messageID,
		UserID:    userID,
		EmojiName: req.EmojiName,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.store.Create(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *ReactionService) Remove(ctx context.Context, messageID, userID, emojiName string) error {
	return s.store.Delete(ctx, messageID, userID, emojiName)
}

func (s *ReactionService) ListByMessage(ctx context.Context, messageID string) ([]*model.Reaction, error) {
	return s.store.ListByMessage(ctx, messageID)
}
