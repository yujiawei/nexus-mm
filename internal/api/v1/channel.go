package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/api/middleware"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type ChannelHandler struct {
	channelSvc *service.ChannelService
	teamSvc    *service.TeamService
}

func NewChannelHandler(channelSvc *service.ChannelService, teamSvc *service.TeamService) *ChannelHandler {
	return &ChannelHandler{channelSvc: channelSvc, teamSvc: teamSvc}
}

func (h *ChannelHandler) Create(c *gin.Context) {
	teamID := c.Param("id")
	userID := c.GetString("user_id")

	isMember, err := h.teamSvc.IsMember(c.Request.Context(), teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a team member"})
		return
	}

	var req model.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !middleware.ValidateChannelName(req.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel name must be 3-50 chars, lowercase alphanumeric and hyphen only"})
		return
	}

	ch, err := h.channelSvc.Create(c.Request.Context(), teamID, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, ch)
}

func (h *ChannelHandler) ListByTeam(c *gin.Context) {
	teamID := c.Param("id")
	userID := c.GetString("user_id")

	isMember, err := h.teamSvc.IsMember(c.Request.Context(), teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a team member"})
		return
	}

	channels, err := h.channelSvc.ListByTeam(c.Request.Context(), teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channels)
}

func (h *ChannelHandler) Get(c *gin.Context) {
	channelID := c.Param("id")

	ch, err := h.channelSvc.GetByID(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
		return
	}

	c.JSON(http.StatusOK, ch)
}

func (h *ChannelHandler) JoinChannel(c *gin.Context) {
	channelID := c.Param("id")
	userID := c.GetString("user_id")

	// Get channel to find team_id.
	ch, err := h.channelSvc.GetByID(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
		return
	}

	// Must be team member first.
	isMember, err := h.teamSvc.IsMember(c.Request.Context(), ch.TeamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "must be a team member first"})
		return
	}

	if err := h.channelSvc.AddMember(c.Request.Context(), channelID, userID, "member"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ChannelHandler) AddMember(c *gin.Context) {
	channelID := c.Param("id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get channel to find team_id.
	ch, err := h.channelSvc.GetByID(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
		return
	}

	// Target must be team member.
	isMember, err := h.teamSvc.IsMember(c.Request.Context(), ch.TeamID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "user must be a team member first"})
		return
	}

	if err := h.channelSvc.AddMember(c.Request.Context(), channelID, req.UserID, "member"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ChannelHandler) ListMembers(c *gin.Context) {
	channelID := c.Param("id")

	members, err := h.channelSvc.ListMembers(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *ChannelHandler) RemoveMember(c *gin.Context) {
	channelID := c.Param("id")
	targetUserID := c.Param("user_id")
	callerID := c.GetString("user_id")

	// Allow self-removal (leave) or channel admin.
	if callerID != targetUserID {
		// Check if caller is channel admin.
		ch, err := h.channelSvc.GetByID(c.Request.Context(), channelID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		if ch.CreatorID != callerID {
			c.JSON(http.StatusForbidden, gin.H{"error": "only channel creator can remove others"})
			return
		}
	}

	if err := h.channelSvc.RemoveMember(c.Request.Context(), channelID, targetUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ChannelHandler) SetRetention(c *gin.Context) {
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

	var req model.SetRetentionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.channelSvc.SetRetention(c.Request.Context(), channelID, req.RetentionDays); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
