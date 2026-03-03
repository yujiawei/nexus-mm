package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type AgentHandler struct {
	botSvc *service.BotService
}

func NewAgentHandler(botSvc *service.BotService) *AgentHandler {
	return &AgentHandler{botSvc: botSvc}
}

// Register creates a new bot and binds it to an owner by email.
func (h *AgentHandler) Register(c *gin.Context) {
	var req model.BotRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.botSvc.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Bind binds an existing bot to an owner by email.
func (h *AgentHandler) Bind(c *gin.Context) {
	botToken := c.GetHeader("X-Bot-Token")
	if botToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-Bot-Token header required"})
		return
	}

	bot, err := h.botSvc.GetByToken(c.Request.Context(), botToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid bot token"})
		return
	}

	var req model.BotBindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.botSvc.BindToOwner(c.Request.Context(), bot, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Bot bound to owner"})
}

// CreateBot creates a bot for the authenticated user (JWT auth required).
func (h *AgentHandler) CreateBot(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.CreateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bot, err := h.botSvc.CreateBotForUser(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.BotInfo{
		ID:          bot.ID,
		Username:    bot.Username,
		Nickname:    bot.Nickname,
		Description: bot.BotDescription,
		WebhookURL:  bot.BotWebhookURL,
		Token:       *bot.BotToken,
		CreatedAt:   bot.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// ListBots lists bots owned by the authenticated user.
func (h *AgentHandler) ListBots(c *gin.Context) {
	userID := c.GetString("user_id")

	bots, err := h.botSvc.ListByOwner(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]model.BotInfo, len(bots))
	for i, bot := range bots {
		token := ""
		if bot.BotToken != nil {
			token = *bot.BotToken
		}
		result[i] = model.BotInfo{
			ID:          bot.ID,
			Username:    bot.Username,
			Nickname:    bot.Nickname,
			Description: bot.BotDescription,
			WebhookURL:  bot.BotWebhookURL,
			Token:       token,
			CreatedAt:   bot.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, result)
}

// RegenerateToken regenerates a bot's token.
func (h *AgentHandler) RegenerateToken(c *gin.Context) {
	userID := c.GetString("user_id")
	botID := c.Param("id")

	newToken, err := h.botSvc.RegenerateToken(c.Request.Context(), botID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bot_token": newToken})
}

// UpdateWebhook updates a bot's webhook URL.
func (h *AgentHandler) UpdateWebhook(c *gin.Context) {
	userID := c.GetString("user_id")
	botID := c.Param("id")

	// Verify ownership.
	bots, err := h.botSvc.ListByOwner(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	owned := false
	for _, bot := range bots {
		if bot.ID == botID {
			owned = true
			break
		}
	}
	if !owned {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req model.BotSetWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.botSvc.SetWebhook(c.Request.Context(), botID, req.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
