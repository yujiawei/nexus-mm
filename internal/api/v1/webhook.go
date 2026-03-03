package v1

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type WebhookHandler struct {
	db     *sqlx.DB
	msgSvc *service.MessageService
}

func NewWebhookHandler(db *sqlx.DB, msgSvc *service.MessageService) *WebhookHandler {
	return &WebhookHandler{db: db, msgSvc: msgSvc}
}

func (h *WebhookHandler) CreateIncoming(c *gin.Context) {
	var req model.CreateIncomingWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	now := time.Now().UTC()

	// Look up channel to get team_id.
	var teamID string
	err := h.db.GetContext(c.Request.Context(), &teamID,
		"SELECT team_id FROM channels WHERE id = $1", req.ChannelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel not found"})
		return
	}

	hook := &model.IncomingWebhook{
		ID:          ulid.Make().String(),
		ChannelID:   req.ChannelID,
		TeamID:      teamID,
		CreatorID:   userID,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Token:       generateToken(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = h.db.ExecContext(c.Request.Context(),
		`INSERT INTO incoming_webhooks (id, channel_id, team_id, creator_id, display_name, description, token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		hook.ID, hook.ChannelID, hook.TeamID, hook.CreatorID,
		hook.DisplayName, hook.Description, hook.Token, hook.CreatedAt, hook.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hook)
}

func (h *WebhookHandler) CreateOutgoing(c *gin.Context) {
	var req model.CreateOutgoingWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	now := time.Now().UTC()

	var teamID string
	err := h.db.GetContext(c.Request.Context(), &teamID,
		"SELECT team_id FROM channels WHERE id = $1", req.ChannelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel not found"})
		return
	}

	hook := &model.OutgoingWebhook{
		ID:           ulid.Make().String(),
		ChannelID:    req.ChannelID,
		TeamID:       teamID,
		CreatorID:    userID,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		TriggerWords: req.TriggerWords,
		CallbackURLs: req.CallbackURLs,
		Token:        generateToken(),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tx, err := h.db.BeginTxx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(c.Request.Context(),
		`INSERT INTO outgoing_webhooks (id, channel_id, team_id, creator_id, display_name, description, token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		hook.ID, hook.ChannelID, hook.TeamID, hook.CreatorID,
		hook.DisplayName, hook.Description, hook.Token, hook.CreatedAt, hook.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, word := range req.TriggerWords {
		_, err = tx.ExecContext(c.Request.Context(),
			"INSERT INTO outgoing_webhook_triggers (webhook_id, trigger_word) VALUES ($1, $2)",
			hook.ID, word)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	for _, url := range req.CallbackURLs {
		_, err = tx.ExecContext(c.Request.Context(),
			"INSERT INTO outgoing_webhook_urls (webhook_id, callback_url) VALUES ($1, $2)",
			hook.ID, url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hook)
}

// PostIncoming handles external incoming webhook payloads (no auth required, token in URL).
func (h *WebhookHandler) PostIncoming(c *gin.Context) {
	hookID := c.Param("id")
	token := c.Query("token")

	var hook model.IncomingWebhook
	err := h.db.GetContext(c.Request.Context(), &hook,
		"SELECT * FROM incoming_webhooks WHERE id = $1", hookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "webhook not found"})
		return
	}

	if hook.Token != token {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var payload model.IncomingWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.msgSvc.Send(c.Request.Context(), hook.ChannelID, hook.CreatorID, &model.SendMessageRequest{
		Content: payload.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, msg)
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
