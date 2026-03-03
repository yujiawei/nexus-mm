package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type MessageHandler struct {
	msgSvc     *service.MessageService
	channelSvc *service.ChannelService
	commandSvc *service.SlashCommandService
}

func NewMessageHandler(msgSvc *service.MessageService, channelSvc *service.ChannelService, commandSvc *service.SlashCommandService) *MessageHandler {
	return &MessageHandler{msgSvc: msgSvc, channelSvc: channelSvc, commandSvc: commandSvc}
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

	// Check for slash command.
	if strings.HasPrefix(req.Content, "/") && h.commandSvc != nil {
		ch, _ := h.channelSvc.GetByID(c.Request.Context(), channelID)
		if ch != nil {
			parts := strings.SplitN(req.Content[1:], " ", 2)
			trigger := parts[0]
			text := ""
			if len(parts) > 1 {
				text = parts[1]
			}
			resp, err := h.commandSvc.Execute(c.Request.Context(), ch.TeamID, channelID, userID, trigger, text)
			if err == nil && resp != "" {
				// Post command response to channel.
				respMsg, err := h.msgSvc.Send(c.Request.Context(), channelID, userID, &model.SendMessageRequest{Content: resp})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, respMsg)
				return
			}
			// If command not found, fall through to normal send.
		}
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

func (h *MessageHandler) GetThread(c *gin.Context) {
	channelID := c.Param("id")
	msgID := c.Param("msg_id")
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

	messages, err := h.msgSvc.GetThread(c.Request.Context(), msgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
