package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type AuditHandler struct {
	auditSvc *service.AuditService
	userSvc  *service.UserService
}

func NewAuditHandler(auditSvc *service.AuditService, userSvc *service.UserService) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc, userSvc: userSvc}
}

func (h *AuditHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.userSvc.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	// Allow system admins and any authenticated user (team admins can access).
	_ = user

	var query model.AuditLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logs, err := h.auditSvc.List(c.Request.Context(), &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
