package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
	"github.com/yujiawei/nexus-mm/internal/wkim"
)

type MessageService struct {
	store *postgres.MessageStore
	wk    *wkim.Client
}

func NewMessageService(store *postgres.MessageStore, wk *wkim.Client) *MessageService {
	return &MessageService{store: store, wk: wk}
}

func (s *MessageService) Send(ctx context.Context, channelID, userID string, req *model.SendMessageRequest) (*model.Message, error) {
	now := time.Now().UTC()
	msg := &model.Message{
		ID:        ulid.Make().String(),
		ChannelID: channelID,
		UserID:    userID,
		Content:   req.Content,
		Type:      "text",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("store message: %w", err)
	}

	// Send message through WuKongIM for real-time delivery.
	payload, _ := json.Marshal(map[string]string{
		"type":    "text",
		"content": req.Content,
		"msg_id":  msg.ID,
	})
	if err := s.wk.SendMessage(ctx, &wkim.SendMsgReq{
		FromUID:     userID,
		ChannelID:   channelID,
		ChannelType: wkim.ChannelTypeGroup,
		Payload:     payload,
	}); err != nil {
		fmt.Printf("warn: send wukongim message: %v\n", err)
	}

	return msg, nil
}

func (s *MessageService) ListByChannel(ctx context.Context, channelID, before string, limit int) ([]*model.Message, error) {
	return s.store.ListByChannel(ctx, channelID, before, limit)
}
