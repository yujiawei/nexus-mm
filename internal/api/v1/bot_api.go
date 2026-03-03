package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type BotAPIHandler struct {
	botSvc      *service.BotService
	reactionSvc *service.ReactionService
}

func NewBotAPIHandler(botSvc *service.BotService, reactionSvc *service.ReactionService) *BotAPIHandler {
	return &BotAPIHandler{botSvc: botSvc, reactionSvc: reactionSvc}
}

// getBot extracts and validates the bot from the token in the URL.
func (h *BotAPIHandler) getBot(c *gin.Context) *model.User {
	token := c.Param("token")
	bot, err := h.botSvc.GetByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "invalid bot token"})
		return nil
	}
	return bot
}

func (h *BotAPIHandler) GetMe(c *gin.Context) {
	bot := h.getBot(c)
	if bot == nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"result": gin.H{
			"id":          bot.ID,
			"username":    bot.Username,
			"nickname":    bot.Nickname,
			"is_bot":      true,
			"description": bot.BotDescription,
		},
	})
}

func (h *BotAPIHandler) SendMessage(c *gin.Context) {
	bot := h.getBot(c)
	if bot == nil {
		return
	}

	var req model.BotSendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	msg, err := h.botSvc.SendMessage(c.Request.Context(), bot, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "result": msg})
}

func (h *BotAPIHandler) GetUpdates(c *gin.Context) {
	bot := h.getBot(c)
	if bot == nil {
		return
	}

	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	resp, err := h.botSvc.GetUpdates(c.Request.Context(), bot.ID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *BotAPIHandler) SetWebhook(c *gin.Context) {
	bot := h.getBot(c)
	if bot == nil {
		return
	}

	var req model.BotSetWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	if err := h.botSvc.SetWebhook(c.Request.Context(), bot.ID, req.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "description": "Webhook was set"})
}

func (h *BotAPIHandler) SendReaction(c *gin.Context) {
	bot := h.getBot(c)
	if bot == nil {
		return
	}

	var req model.BotSendReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	reaction, err := h.reactionSvc.Add(c.Request.Context(), req.MessageID, bot.ID, &model.CreateReactionRequest{EmojiName: req.EmojiName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "result": reaction})
}
