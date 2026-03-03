package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
)

type SlashCommandService struct {
	store      *postgres.SlashCommandStore
	httpClient *http.Client
}

func NewSlashCommandService(store *postgres.SlashCommandStore) *SlashCommandService {
	return &SlashCommandService{
		store:      store,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SlashCommandService) Create(ctx context.Context, teamID, creatorID string, req *model.CreateSlashCommandRequest) (*model.SlashCommand, error) {
	now := time.Now().UTC()
	method := req.Method
	if method == "" {
		method = "POST"
	}

	cmd := &model.SlashCommand{
		ID:        ulid.Make().String(),
		TeamID:    teamID,
		Trigger:   req.Trigger,
		URL:       req.URL,
		Method:    method,
		CreatorID: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.Create(ctx, cmd); err != nil {
		return nil, fmt.Errorf("create command: %w", err)
	}
	return cmd, nil
}

func (s *SlashCommandService) ListByTeam(ctx context.Context, teamID string) ([]*model.SlashCommand, error) {
	return s.store.ListByTeam(ctx, teamID)
}

// Execute finds and executes a slash command, returning the response text.
func (s *SlashCommandService) Execute(ctx context.Context, teamID, channelID, userID, trigger, text string) (string, error) {
	cmd, err := s.store.GetByTrigger(ctx, teamID, trigger)
	if err != nil {
		return "", err
	}

	payload, _ := json.Marshal(map[string]string{
		"team_id":    teamID,
		"channel_id": channelID,
		"user_id":    userID,
		"command":    "/" + trigger,
		"text":       text,
	})

	req, err := http.NewRequestWithContext(ctx, cmd.Method, cmd.URL, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create command request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute command: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Text string `json:"text"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Text, nil
}
