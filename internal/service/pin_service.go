package service

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
)

type PinService struct {
	store *postgres.PinStore
}

func NewPinService(store *postgres.PinStore) *PinService {
	return &PinService{store: store}
}

func (s *PinService) Pin(ctx context.Context, channelID, messageID, userID string) (*model.PinnedMessage, error) {
	pin := &model.PinnedMessage{
		ID:        ulid.Make().String(),
		ChannelID: channelID,
		MessageID: messageID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.store.Create(ctx, pin); err != nil {
		return nil, err
	}
	return pin, nil
}

func (s *PinService) Unpin(ctx context.Context, channelID, messageID string) error {
	return s.store.Delete(ctx, channelID, messageID)
}

func (s *PinService) ListByChannel(ctx context.Context, channelID string) ([]*model.PinnedMessage, error) {
	return s.store.ListByChannel(ctx, channelID)
}
