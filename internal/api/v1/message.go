package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type MessageHandler struct {
	msgSvc     *service.MessageService
	channelSvc *service.ChannelService
}

func NewMessageHandler(msgSvc *service.MessageService, channelSvc *service.ChannelService) *MessageHandler {
	return &MessageHandler{msgSvc: msgSvc, channelSvc: channelSvc}
}

func (h *MessageHandler) Send(c *gin.Context) {
	channelID := c.Param("id")
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

	var req model.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.msgSvc.Send(c.Request.Context(), channelID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *MessageHandler) List(c *gin.Context) {
	channelID := c.Param("id")
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

	var req model.MessageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages, err := h.msgSvc.ListByChannel(c.Request.Context(), channelID, req.Before, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
