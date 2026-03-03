package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/api/middleware"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.UserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !middleware.ValidateUsername(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be 3-30 chars, alphanumeric and underscore only"})
		return
	}
	if !middleware.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}
	if !middleware.ValidatePassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}

	resp, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
		"wk_token":   user.WkToken,
		"role":       user.Role,
		"is_bot":     user.IsBot,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"ws_url":     h.svc.WsURL(),
	})
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	callerID := c.GetString("user_id")
	targetID := c.Param("id")

	// Only system admin can change roles.
	caller, err := h.svc.GetByID(c.Request.Context(), callerID)
	if err != nil || caller.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=admin member"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateRole(c.Request.Context(), targetID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
