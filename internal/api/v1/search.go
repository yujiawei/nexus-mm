package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yujiawei/nexus-mm/internal/search"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type SearchHandler struct {
	search     *search.MeiliClient
	channelSvc *service.ChannelService
}

func NewSearchHandler(sc *search.MeiliClient, channelSvc *service.ChannelService) *SearchHandler {
	return &SearchHandler{search: sc, channelSvc: channelSvc}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	channelID := c.Query("channel_id")
	teamID := c.Query("team_id")

	var channelIDs []string

	if channelID != "" {
		channelIDs = []string{channelID}
	} else if teamID != "" {
		ids, err := h.channelSvc.ListIDsByTeam(c.Request.Context(), teamID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		channelIDs = ids
	}

	results, err := h.search.Search(c.Request.Context(), query, channelIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
