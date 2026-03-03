package wkim

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// ChannelType constants for WuKongIM.
const (
	ChannelTypeGroup = 2 // group channel
)

// WebhookHandler handles WuKongIM webhook callbacks.
type WebhookHandler struct {
	subscriberFn func(ctx context.Context, channelID string) ([]string, error)
}

func NewWebhookHandler(subscriberFn func(ctx context.Context, channelID string) ([]string, error)) *WebhookHandler {
	return &WebhookHandler{subscriberFn: subscriberFn}
}

type subscriberRequest struct {
	ChannelID   string `json:"channel_id"`
	ChannelType int    `json:"channel_type"`
}

// ServeHTTP handles the getSubscribers webhook from WuKongIM.
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req subscriberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("decode subscriber request")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	uids, err := h.subscriberFn(r.Context(), req.ChannelID)
	if err != nil {
		log.Error().Err(err).Str("channel_id", req.ChannelID).Msg("get subscribers")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"subscribers": uids})
}
