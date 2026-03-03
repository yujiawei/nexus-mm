package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type ReactionHandler struct {
	reactionSvc *service.ReactionService
	channelSvc  *service.ChannelService
}

func NewReactionHandler(reactionSvc *service.ReactionService, channelSvc *service.ChannelService) *ReactionHandler {
	return &ReactionHandler{reactionSvc: reactionSvc, channelSvc: channelSvc}
}

func (h *ReactionHandler) Add(c *gin.Context) {
	channelID := c.Param("id")
	messageID := c.Param("msg_id")
	userID := c.GetString("user_id")

	isMember, err := h.channelSvc.IsMember(c.Request.Context(), channelID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a channel member"})
		return
	}

	var req model.CreateReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reaction, err := h.reactionSvc.Add(c.Request.Context(), messageID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reaction)
}

func (h *ReactionHandler) Remove(c *gin.Context) {
	channelID := c.Param("id")
	messageID := c.Param("msg_id")
	emojiName := c.Param("emoji")
	userID := c.GetString("user_id")

	isMember, err := h.channelSvc.IsMember(c.Request.Context(), channelID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a channel member"})
		return
	}

	if err := h.reactionSvc.Remove(c.Request.Context(), messageID, userID, emojiName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
