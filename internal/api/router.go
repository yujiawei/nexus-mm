package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/api/middleware"
	v1 "github.com/yujiawei/nexus-mm/internal/api/v1"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type Handlers struct {
	User    *v1.UserHandler
	Team    *v1.TeamHandler
	Channel *v1.ChannelHandler
	Message *v1.MessageHandler
	Webhook *v1.WebhookHandler
}

func NewHandlers(
	userSvc *service.UserService,
	teamSvc *service.TeamService,
	channelSvc *service.ChannelService,
	msgSvc *service.MessageService,
	db *sqlx.DB,
) *Handlers {
	return &Handlers{
		User:    v1.NewUserHandler(userSvc),
		Team:    v1.NewTeamHandler(teamSvc),
		Channel: v1.NewChannelHandler(channelSvc, teamSvc),
		Message: v1.NewMessageHandler(msgSvc, channelSvc),
		Webhook: v1.NewWebhookHandler(db, msgSvc),
	}
}

func SetupRouter(h *Handlers, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// Health check.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	// Public routes.
	api.POST("/users/register", h.User.Register)
	api.POST("/users/login", h.User.Login)

	// Incoming webhook post endpoint (token-based auth, not JWT).
	api.POST("/hooks/incoming/:id", h.Webhook.PostIncoming)

	// Authenticated routes.
	auth := api.Group("")
	auth.Use(middleware.Auth(jwtSecret))

	// Users.
	auth.GET("/users/me", h.User.Me)

	// Teams.
	auth.POST("/teams", h.Team.Create)
	auth.GET("/teams", h.Team.List)
	auth.GET("/teams/:id", h.Team.Get)

	// Channels.
	auth.POST("/teams/:team_id/channels", h.Channel.Create)
	auth.GET("/teams/:team_id/channels", h.Channel.ListByTeam)
	auth.GET("/channels/:id", h.Channel.Get)

	// Messages.
	auth.POST("/channels/:id/messages", h.Message.Send)
	auth.GET("/channels/:id/messages", h.Message.List)

	// Webhooks (create).
	auth.POST("/hooks/incoming", h.Webhook.CreateIncoming)
	auth.POST("/hooks/outgoing", h.Webhook.CreateOutgoing)

	return r
}
