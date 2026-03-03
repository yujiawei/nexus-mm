package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/search"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
	"github.com/yujiawei/nexus-mm/internal/wkim"
)

// BotDeliverer is the interface for delivering messages to bots.
type BotDeliverer interface {
	DeliverToBot(ctx context.Context, channelID string, msg *model.Message)
}

type MessageService struct {
	store        *postgres.MessageStore
	webhookStore *postgres.WebhookStore
	wk           *wkim.Client
	search       *search.MeiliClient
	botDeliverer BotDeliverer
}

func NewMessageService(store *postgres.MessageStore, wk *wkim.Client) *MessageService {
	return &MessageService{store: store, wk: wk}
}

// SetWebhookStore sets the webhook store for outgoing webhook triggers.
func (s *MessageService) SetWebhookStore(ws *postgres.WebhookStore) {
	s.webhookStore = ws
}

// SetSearch sets the MeiliSearch client for message indexing.
func (s *MessageService) SetSearch(sc *search.MeiliClient) {
	s.search = sc
}

// SetBotService sets the bot service for delivering messages to bots.
func (s *MessageService) SetBotService(bd BotDeliverer) {
	s.botDeliverer = bd
}

func (s *MessageService) Send(ctx context.Context, channelID, userID string, req *model.SendMessageRequest) (*model.Message, error) {
	now := time.Now().UTC()
	msg := &model.Message{
		ID:        ulid.Make().String(),
		ChannelID: channelID,
		UserID:    userID,
		Content:   req.Content,
		Type:      "text",
		RootID:    req.RootID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("store message: %w", err)
	}

	// Increment reply count on root message if this is a reply.
	if req.RootID != "" {
		if err := s.store.IncrementReplyCount(ctx, req.RootID); err != nil {
			log.Warn().Err(err).Str("root_id", req.RootID).Msg("increment reply count")
		}
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
		log.Warn().Err(err).Msg("send wukongim message")
	}

	// Index in MeiliSearch (async).
	if s.search != nil {
		go func() {
			doc := &search.MessageDocument{
				ID:        msg.ID,
				ChannelID: msg.ChannelID,
				UserID:    msg.UserID,
				Content:   msg.Content,
				CreatedAt: msg.CreatedAt.Unix(),
			}
			if err := s.search.IndexMessage(context.Background(), doc); err != nil {
				log.Warn().Err(err).Msg("index message in meilisearch")
			}
		}()
	}

	// Trigger outgoing webhooks (async).
	if s.webhookStore != nil {
		go s.triggerOutgoingWebhooks(channelID, userID, msg.Content)
	}

	// Deliver to bots in the channel (async).
	if s.botDeliverer != nil {
		go s.botDeliverer.DeliverToBot(context.Background(), channelID, msg)
	}

	return msg, nil
}

func (s *MessageService) GetByID(ctx context.Context, id string) (*model.Message, error) {
	return s.store.GetByID(ctx, id)
}

func (s *MessageService) GetThread(ctx context.Context, rootID string) ([]*model.Message, error) {
	return s.store.GetThread(ctx, rootID)
}

func (s *MessageService) ListByChannel(ctx context.Context, channelID, before string, limit int) ([]*model.Message, error) {
	return s.store.ListByChannel(ctx, channelID, before, limit)
}

func (s *MessageService) triggerOutgoingWebhooks(channelID, userID, content string) {
	ctx := context.Background()
	webhooks, err := s.webhookStore.GetOutgoingByChannel(ctx, channelID)
	if err != nil {
		return
	}

	for _, wh := range webhooks {
		for _, trigger := range wh.TriggerWords {
			if strings.HasPrefix(content, trigger) {
				s.fireWebhookCallbacks(wh, userID, content, trigger)
				break
			}
		}
	}
}

func (s *MessageService) fireWebhookCallbacks(wh *postgres.OutgoingWebhookInfo, userID, content, trigger string) {
	payload := &model.OutgoingWebhookPayload{
		Token:       wh.Token,
		TeamID:      wh.TeamID,
		ChannelID:   wh.ChannelID,
		UserID:      userID,
		Text:        content,
		TriggerWord: trigger,
		Timestamp:   time.Now().Unix(),
	}
	data, _ := json.Marshal(payload)

	for _, url := range wh.CallbackURLs {
		go func(u string) {
			resp, err := http.Post(u, "application/json", bytes.NewReader(data))
			if err != nil {
				log.Warn().Err(err).Str("url", u).Msg("outgoing webhook callback")
				return
			}
			resp.Body.Close()
		}(url)
	}
}
