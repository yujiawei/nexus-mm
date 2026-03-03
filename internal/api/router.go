package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/api/middleware"
	v1 "github.com/yujiawei/nexus-mm/internal/api/v1"
	"github.com/yujiawei/nexus-mm/internal/search"
	"github.com/yujiawei/nexus-mm/internal/service"
)

type Handlers struct {
	User         *v1.UserHandler
	Team         *v1.TeamHandler
	Channel      *v1.ChannelHandler
	Message      *v1.MessageHandler
	Webhook      *v1.WebhookHandler
	Reaction     *v1.ReactionHandler
	Pin          *v1.PinHandler
	Search       *v1.SearchHandler
	SlashCommand *v1.SlashCommandHandler
	Category     *v1.CategoryHandler
	Audit        *v1.AuditHandler
}

func NewHandlers(
	userSvc *service.UserService,
	teamSvc *service.TeamService,
	channelSvc *service.ChannelService,
	msgSvc *service.MessageService,
	reactionSvc *service.ReactionService,
	pinSvc *service.PinService,
	commandSvc *service.SlashCommandService,
	categorySvc *service.CategoryService,
	auditSvc *service.AuditService,
	meili *search.MeiliClient,
	db *sqlx.DB,
) *Handlers {
	return &Handlers{
		User:         v1.NewUserHandler(userSvc),
		Team:         v1.NewTeamHandler(teamSvc),
		Channel:      v1.NewChannelHandler(channelSvc, teamSvc),
		Message:      v1.NewMessageHandler(msgSvc, channelSvc, commandSvc),
		Webhook:      v1.NewWebhookHandler(db, msgSvc),
		Reaction:     v1.NewReactionHandler(reactionSvc, channelSvc),
		Pin:          v1.NewPinHandler(pinSvc, channelSvc),
		Search:       v1.NewSearchHandler(meili, channelSvc),
		SlashCommand: v1.NewSlashCommandHandler(commandSvc, teamSvc),
		Category:     v1.NewCategoryHandler(categorySvc, teamSvc),
		Audit:        v1.NewAuditHandler(auditSvc, userSvc),
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

	// Users (admin).
	auth.PUT("/users/:id/role", h.User.UpdateRole)

	// Teams.
	auth.POST("/teams", h.Team.Create)
	auth.GET("/teams", h.Team.List)
	auth.GET("/teams/all", h.Team.ListAll)
	auth.GET("/teams/:id", h.Team.Get)
	auth.PUT("/teams/:id/retention", h.Team.SetRetention)
	auth.POST("/teams/:id/join", h.Team.JoinTeam)
	auth.POST("/teams/:id/members", h.Team.AddMember)
	auth.GET("/teams/:id/members", h.Team.ListMembers)
	auth.DELETE("/teams/:id/members/:user_id", h.Team.RemoveMember)

	// Channels.
	auth.POST("/teams/:id/channels", h.Channel.Create)
	auth.GET("/teams/:id/channels", h.Channel.ListByTeam)
	auth.GET("/channels/:id", h.Channel.Get)
	auth.PUT("/channels/:id/retention", h.Channel.SetRetention)
	auth.POST("/channels/:id/join", h.Channel.JoinChannel)
	auth.POST("/channels/:id/members", h.Channel.AddMember)
	auth.GET("/channels/:id/members", h.Channel.ListMembers)
	auth.DELETE("/channels/:id/members/:user_id", h.Channel.RemoveMember)

	// Messages.
	auth.POST("/channels/:id/messages", h.Message.Send)
	auth.GET("/channels/:id/messages", h.Message.List)
	auth.GET("/channels/:id/messages/:msg_id/thread", h.Message.GetThread)

	// Reactions.
	auth.POST("/channels/:id/messages/:msg_id/reactions", h.Reaction.Add)
	auth.DELETE("/channels/:id/messages/:msg_id/reactions/:emoji", h.Reaction.Remove)

	// Pins.
	auth.POST("/channels/:id/messages/:msg_id/pin", h.Pin.Pin)
	auth.DELETE("/channels/:id/messages/:msg_id/pin", h.Pin.Unpin)
	auth.GET("/channels/:id/pinned", h.Pin.List)

	// Search.
	auth.GET("/search", h.Search.Search)

	// Slash commands.
	auth.POST("/teams/:id/commands", h.SlashCommand.Create)
	auth.GET("/teams/:id/commands", h.SlashCommand.List)

	// Channel categories.
	auth.POST("/teams/:id/categories", h.Category.Create)
	auth.GET("/teams/:id/categories", h.Category.List)
	auth.PUT("/categories/:id", h.Category.Update)
	auth.DELETE("/categories/:id", h.Category.Delete)
	auth.POST("/categories/:id/channels", h.Category.AddEntry)
	auth.DELETE("/categories/:id/channels/:channel_id", h.Category.RemoveEntry)
	auth.GET("/categories/:id/channels", h.Category.ListEntries)

	// Webhooks (create).
	auth.POST("/hooks/incoming", h.Webhook.CreateIncoming)
	auth.POST("/hooks/outgoing", h.Webhook.CreateOutgoing)

	// Admin.
	auth.GET("/admin/audit", h.Audit.List)

	return r
}
