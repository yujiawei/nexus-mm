package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
	"github.com/yujiawei/nexus-mm/internal/wkim"
)

type BotService struct {
	userStore      *postgres.UserStore
	botUpdateStore *postgres.BotUpdateStore
	teamStore      *postgres.TeamStore
	channelStore   *postgres.ChannelStore
	msgStore       *postgres.MessageStore
	wk             *wkim.Client
}

func NewBotService(
	userStore *postgres.UserStore,
	botUpdateStore *postgres.BotUpdateStore,
	teamStore *postgres.TeamStore,
	channelStore *postgres.ChannelStore,
	msgStore *postgres.MessageStore,
	wk *wkim.Client,
) *BotService {
	return &BotService{
		userStore:      userStore,
		botUpdateStore: botUpdateStore,
		teamStore:      teamStore,
		channelStore:   channelStore,
		msgStore:       msgStore,
		wk:             wk,
	}
}

func (s *BotService) Register(ctx context.Context, req *model.BotRegisterRequest) (*model.BotRegisterResponse, error) {
	// Find or create owner by email.
	owner, err := s.userStore.GetByEmail(ctx, req.Email)
	if err != nil {
		// Create user with auto-generated password.
		autoPass := generateBotToken()[:16]
		hash, err := bcrypt.GenerateFromPassword([]byte(autoPass), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		now := time.Now().UTC()
		owner = &model.User{
			ID:           ulid.Make().String(),
			Username:     req.Email,
			Email:        req.Email,
			PasswordHash: string(hash),
			Nickname:     req.Email,
			Role:         "member",
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := s.userStore.Create(ctx, owner); err != nil {
			return nil, fmt.Errorf("create owner: %w", err)
		}
		// Register in WuKongIM.
		wkToken := randomHex(16)
		if err := s.wk.RegisterUser(ctx, owner.ID, wkToken); err != nil {
			log.Warn().Err(err).Msg("register owner in wukongim")
		}
	}

	// Create bot user.
	botToken := generateBotToken()
	now := time.Now().UTC()
	botUser := &model.User{
		ID:             ulid.Make().String(),
		Username:       fmt.Sprintf("bot_%s", req.Name),
		Email:          fmt.Sprintf("bot_%s@nexus.local", ulid.Make().String()),
		PasswordHash:   "",
		Nickname:       req.Name,
		Role:           "member",
		IsBot:          true,
		BotToken:       &botToken,
		BotOwnerID:     &owner.ID,
		BotDescription: "",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.userStore.Create(ctx, botUser); err != nil {
		return nil, fmt.Errorf("create bot user: %w", err)
	}

	// Register bot in WuKongIM.
	wkToken := randomHex(16)
	if err := s.wk.RegisterUser(ctx, botUser.ID, wkToken); err != nil {
		log.Warn().Err(err).Msg("register bot in wukongim")
	}

	// Auto-add bot to owner's teams.
	teams, err := s.teamStore.List(ctx, owner.ID)
	if err == nil {
		for _, team := range teams {
			member := &model.TeamMember{
				TeamID:    team.ID,
				UserID:    botUser.ID,
				Role:      "member",
				CreatedAt: now,
			}
			if err := s.teamStore.AddMember(ctx, member); err != nil {
				log.Warn().Err(err).Str("team", team.ID).Msg("add bot to team")
			}
		}
	}

	return &model.BotRegisterResponse{
		BotUserID:   botUser.ID,
		BotToken:    botToken,
		OwnerUserID: owner.ID,
		Message:     "Bot created and bound to owner. Use bot_token for Bot API.",
	}, nil
}

func (s *BotService) BindToOwner(ctx context.Context, bot *model.User, email string) error {
	owner, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("owner not found: %w", err)
	}
	bot.BotOwnerID = &owner.ID
	return s.userStore.UpdateBotWebhook(ctx, bot.ID, bot.BotWebhookURL) // triggers updated_at
}

func (s *BotService) GetByToken(ctx context.Context, token string) (*model.User, error) {
	return s.userStore.GetByBotToken(ctx, token)
}

func (s *BotService) SendMessage(ctx context.Context, bot *model.User, req *model.BotSendMessageRequest) (*model.Message, error) {
	// Check bot is a member of the channel.
	isMember, err := s.channelStore.IsMember(ctx, req.ChannelID, bot.ID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		// Auto-join: bots automatically join channels when they try to send
		if err := s.channelStore.AddMember(ctx, &model.ChannelMember{
			ChannelID: req.ChannelID,
			UserID:    bot.ID,
			CreatedAt: time.Now().UTC(),
		}); err != nil {
			return nil, fmt.Errorf("auto-join channel: %w", err)
		}
	}

	now := time.Now().UTC()
	msg := &model.Message{
		ID:        ulid.Make().String(),
		ChannelID: req.ChannelID,
		UserID:    bot.ID,
		Content:   req.Content,
		Type:      "text",
		RootID:    req.RootID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.msgStore.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("store message: %w", err)
	}

	if req.RootID != "" {
		if err := s.msgStore.IncrementReplyCount(ctx, req.RootID); err != nil {
			log.Warn().Err(err).Str("root_id", req.RootID).Msg("increment reply count")
		}
	}

	// Send through WuKongIM.
	payload, _ := json.Marshal(map[string]string{
		"type":    "text",
		"content": req.Content,
		"msg_id":  msg.ID,
	})
	if err := s.wk.SendMessage(ctx, &wkim.SendMsgReq{
		FromUID:     bot.ID,
		ChannelID:   req.ChannelID,
		ChannelType: wkim.ChannelTypeGroup,
		Payload:     payload,
	}); err != nil {
		log.Warn().Err(err).Msg("send wukongim message from bot")
	}

	return msg, nil
}

func (s *BotService) GetUpdates(ctx context.Context, botUserID string, offset int64, limit int) (*model.BotUpdateResponse, error) {
	updates, err := s.botUpdateStore.ListSince(ctx, botUserID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("list updates: %w", err)
	}

	result := make([]model.BotUpdateMessage, 0, len(updates))
	for _, u := range updates {
		msg, err := s.msgStore.GetByID(ctx, u.MessageID)
		if err != nil {
			continue
		}
		result = append(result, model.BotUpdateMessage{
			UpdateID:  u.ID,
			ChannelID: u.ChannelID,
			MessageID: u.MessageID,
			UserID:    msg.UserID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format(time.RFC3339),
		})
	}

	return &model.BotUpdateResponse{OK: true, Result: result}, nil
}

func (s *BotService) SetWebhook(ctx context.Context, botID, url string) error {
	return s.userStore.UpdateBotWebhook(ctx, botID, url)
}

func (s *BotService) CreateBotForUser(ctx context.Context, ownerID string, req *model.CreateBotRequest) (*model.User, error) {
	botToken := generateBotToken()
	now := time.Now().UTC()
	botUser := &model.User{
		ID:             ulid.Make().String(),
		Username:       fmt.Sprintf("bot_%s", req.Name),
		Email:          fmt.Sprintf("bot_%s@nexus.local", ulid.Make().String()),
		PasswordHash:   "",
		Nickname:       req.Name,
		Role:           "member",
		IsBot:          true,
		BotToken:       &botToken,
		BotOwnerID:     &ownerID,
		BotDescription: req.Description,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.userStore.Create(ctx, botUser); err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	// Register in WuKongIM.
	wkToken := randomHex(16)
	if err := s.wk.RegisterUser(ctx, botUser.ID, wkToken); err != nil {
		log.Warn().Err(err).Msg("register bot in wukongim")
	}

	// Add bot to owner's teams.
	teams, err := s.teamStore.List(ctx, ownerID)
	if err == nil {
		for _, team := range teams {
			member := &model.TeamMember{
				TeamID:    team.ID,
				UserID:    botUser.ID,
				Role:      "member",
				CreatedAt: now,
			}
			if err := s.teamStore.AddMember(ctx, member); err != nil {
				log.Warn().Err(err).Str("team", team.ID).Msg("add bot to team")
			}
		}
	}

	return botUser, nil
}

func (s *BotService) ListByOwner(ctx context.Context, ownerID string) ([]*model.User, error) {
	return s.userStore.ListBotsByOwner(ctx, ownerID)
}

func (s *BotService) RegenerateToken(ctx context.Context, botID, ownerID string) (string, error) {
	bot, err := s.userStore.GetByID(ctx, botID)
	if err != nil {
		return "", fmt.Errorf("bot not found: %w", err)
	}
	if !bot.IsBot || bot.BotOwnerID == nil || *bot.BotOwnerID != ownerID {
		return "", fmt.Errorf("not authorized")
	}
	newToken := generateBotToken()
	if err := s.userStore.UpdateBotToken(ctx, botID, newToken); err != nil {
		return "", fmt.Errorf("update token: %w", err)
	}
	return newToken, nil
}

// DeliverToBot is called when a message is sent in a channel with bot members.
func (s *BotService) DeliverToBot(ctx context.Context, channelID string, msg *model.Message) {
	botIDs, err := s.botUpdateStore.ListBotIDsByChannel(ctx, channelID)
	if err != nil {
		return
	}

	for _, botID := range botIDs {
		if botID == msg.UserID {
			continue // Don't deliver bot's own messages.
		}

		bot, err := s.userStore.GetByID(ctx, botID)
		if err != nil {
			continue
		}

		if bot.BotWebhookURL != "" {
			// POST to webhook.
			go s.postToWebhook(bot.BotWebhookURL, channelID, msg)
		} else {
			// Store for getUpdates polling.
			if err := s.botUpdateStore.Create(ctx, botID, msg.ID, channelID); err != nil {
				log.Warn().Err(err).Str("bot", botID).Msg("create bot update")
			}
		}
	}
}

func (s *BotService) postToWebhook(url, channelID string, msg *model.Message) {
	payload := map[string]interface{}{
		"update_id":  time.Now().UnixNano(),
		"channel_id": channelID,
		"message": map[string]interface{}{
			"id":         msg.ID,
			"user_id":    msg.UserID,
			"content":    msg.Content,
			"created_at": msg.CreatedAt.Format(time.RFC3339),
		},
	}
	data, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Warn().Err(err).Str("url", url).Msg("bot webhook callback")
		return
	}
	resp.Body.Close()
}

func generateBotToken() string {
	b := make([]byte, 24)
	rand.Read(b)
	return "bf_" + hex.EncodeToString(b)
}
