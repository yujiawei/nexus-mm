package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

func (h *TeamHandler) SetRetention(c *gin.Context) {
	teamID := c.Param("id")
	userID := c.GetString("user_id")

	isMember, err := h.svc.IsMember(c.Request.Context(), teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a team member"})
		return
	}

	var req model.SetRetentionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.SetRetention(c.Request.Context(), teamID, req.RetentionDays); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req model.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	team, err := h.svc.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	teams, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) Get(c *gin.Context) {
	teamID := c.Param("id")
	userID := c.GetString("user_id")

	isMember, err := h.svc.IsMember(c.Request.Context(), teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a team member"})
		return
	}

	team, err := h.svc.GetByID(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) ListAll(c *gin.Context) {
	teams, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) JoinTeam(c *gin.Context) {
	teamID := c.Param("id")
	userID := c.GetString("user_id")

	// Check if team exists.
	_, err := h.svc.GetByID(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	if err := h.svc.AddMember(c.Request.Context(), teamID, userID, "member"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *TeamHandler) AddMember(c *gin.Context) {
	teamID := c.Param("id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.AddMember(c.Request.Context(), teamID, req.UserID, "member"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *TeamHandler) ListMembers(c *gin.Context) {
	teamID := c.Param("id")

	members, err := h.svc.ListMembers(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamID := c.Param("id")
	targetUserID := c.Param("user_id")
	callerID := c.GetString("user_id")

	// Only team owner/admin or self can remove.
	role, err := h.svc.GetMemberRole(c.Request.Context(), teamID, callerID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a team member"})
		return
	}
	if callerID != targetUserID && role != "owner" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only team admin can remove others"})
		return
	}

	if err := h.svc.RemoveMember(c.Request.Context(), teamID, targetUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
