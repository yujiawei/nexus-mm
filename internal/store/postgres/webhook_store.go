package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type WebhookStore struct {
	db *sqlx.DB
}

func NewWebhookStore(db *sqlx.DB) *WebhookStore {
	return &WebhookStore{db: db}
}

type OutgoingWebhookInfo struct {
	ID           string
	ChannelID    string
	TeamID       string
	Token        string
	TriggerWords []string
	CallbackURLs []string
}

func (s *WebhookStore) GetOutgoingByChannel(ctx context.Context, channelID string) ([]*OutgoingWebhookInfo, error) {
	type whRow struct {
		ID        string `db:"id"`
		ChannelID string `db:"channel_id"`
		TeamID    string `db:"team_id"`
		Token     string `db:"token"`
	}
	var rows []whRow
	err := s.db.SelectContext(ctx, &rows,
		"SELECT id, channel_id, team_id, token FROM outgoing_webhooks WHERE channel_id = $1", channelID)
	if err != nil {
		return nil, err
	}

	var result []*OutgoingWebhookInfo
	for _, row := range rows {
		info := &OutgoingWebhookInfo{
			ID:        row.ID,
			ChannelID: row.ChannelID,
			TeamID:    row.TeamID,
			Token:     row.Token,
		}

		var triggers []string
		s.db.SelectContext(ctx, &triggers,
			"SELECT trigger_word FROM outgoing_webhook_triggers WHERE webhook_id = $1", row.ID)
		info.TriggerWords = triggers

		var urls []string
		s.db.SelectContext(ctx, &urls,
			"SELECT callback_url FROM outgoing_webhook_urls WHERE webhook_id = $1", row.ID)
		info.CallbackURLs = urls

		result = append(result, info)
	}
	return result, nil
}
