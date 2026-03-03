package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type CategoryHandler struct {
	categorySvc *service.CategoryService
	teamSvc     *service.TeamService
}

func NewCategoryHandler(categorySvc *service.CategoryService, teamSvc *service.TeamService) *CategoryHandler {
	return &CategoryHandler{categorySvc: categorySvc, teamSvc: teamSvc}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	teamID := c.Param("team_id")
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

	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.categorySvc.Create(c.Request.Context(), userID, teamID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func (h *CategoryHandler) List(c *gin.Context) {
	teamID := c.Param("team_id")
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

	cats, err := h.categorySvc.List(c.Request.Context(), userID, teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cats)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	categoryID := c.Param("id")
	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categorySvc.Update(c.Request.Context(), categoryID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	categoryID := c.Param("id")

	if err := h.categorySvc.Delete(c.Request.Context(), categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *CategoryHandler) AddEntry(c *gin.Context) {
	categoryID := c.Param("id")

	var req model.AddCategoryEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.categorySvc.AddEntry(c.Request.Context(), categoryID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *CategoryHandler) RemoveEntry(c *gin.Context) {
	categoryID := c.Param("id")
	channelID := c.Param("channel_id")

	if err := h.categorySvc.RemoveEntry(c.Request.Context(), categoryID, channelID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *CategoryHandler) ListEntries(c *gin.Context) {
	categoryID := c.Param("id")

	entries, err := h.categorySvc.ListEntries(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}
