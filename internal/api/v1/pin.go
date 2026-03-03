package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type PinHandler struct {
	pinSvc     *service.PinService
	channelSvc *service.ChannelService
}

func NewPinHandler(pinSvc *service.PinService, channelSvc *service.ChannelService) *PinHandler {
	return &PinHandler{pinSvc: pinSvc, channelSvc: channelSvc}
}

func (h *PinHandler) Pin(c *gin.Context) {
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

	pin, err := h.pinSvc.Pin(c.Request.Context(), channelID, messageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pin)
}

func (h *PinHandler) Unpin(c *gin.Context) {
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

	if err := h.pinSvc.Unpin(c.Request.Context(), channelID, messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *PinHandler) List(c *gin.Context) {
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

	pins, err := h.pinSvc.ListByChannel(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pins)
}
