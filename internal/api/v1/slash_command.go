package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type SlashCommandHandler struct {
	commandSvc *service.SlashCommandService
	teamSvc    *service.TeamService
}

func NewSlashCommandHandler(commandSvc *service.SlashCommandService, teamSvc *service.TeamService) *SlashCommandHandler {
	return &SlashCommandHandler{commandSvc: commandSvc, teamSvc: teamSvc}
}

func (h *SlashCommandHandler) Create(c *gin.Context) {
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

	var req model.CreateSlashCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd, err := h.commandSvc.Create(c.Request.Context(), teamID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cmd)
}

func (h *SlashCommandHandler) List(c *gin.Context) {
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

	cmds, err := h.commandSvc.ListByTeam(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cmds)
}
