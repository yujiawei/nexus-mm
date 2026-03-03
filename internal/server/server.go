package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/yujiawei/nexus-mm/internal/api"
	"github.com/yujiawei/nexus-mm/internal/config"
	"github.com/yujiawei/nexus-mm/internal/search"
	"github.com/yujiawei/nexus-mm/internal/service"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
	"github.com/yujiawei/nexus-mm/internal/wkim"
)

type Server struct {
	cfg       *config.Config
	db        *sqlx.DB
	redis     *redis.Client
	http      *http.Server
	wkHook    *http.Server
	cancelCtx context.CancelFunc
}

func New(cfg *config.Config) (*Server, error) {
	db, err := sqlx.Connect("postgres", cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Warn().Err(err).Msg("redis ping failed, continuing without redis")
	}

	// Stores.
	userStore := postgres.NewUserStore(db)
	teamStore := postgres.NewTeamStore(db)
	channelStore := postgres.NewChannelStore(db)
	messageStore := postgres.NewMessageStore(db)
	reactionStore := postgres.NewReactionStore(db)
	pinStore := postgres.NewPinStore(db)
	commandStore := postgres.NewSlashCommandStore(db)
	categoryStore := postgres.NewCategoryStore(db)
	auditStore := postgres.NewAuditStore(db)
	webhookStore := postgres.NewWebhookStore(db)

	// WuKongIM client.
	wk := wkim.NewClient(cfg.WuKong.APIURL, cfg.WuKong.ManagerToken)

	// MeiliSearch client.
	meili := search.NewMeiliClient(cfg.MeiliSearch.URL, cfg.MeiliSearch.APIKey)
	go func() {
		if err := meili.EnsureIndex(context.Background()); err != nil {
			log.Warn().Err(err).Msg("meilisearch ensure index")
		}
	}()

	// Services.
	userSvc := service.NewUserService(userStore, wk, cfg.JWT.Secret, cfg.JWT.ExpireHour)
	teamSvc := service.NewTeamService(teamStore)
	channelSvc := service.NewChannelService(channelStore, wk)
	msgSvc := service.NewMessageService(messageStore, wk)
	msgSvc.SetWebhookStore(webhookStore)
	msgSvc.SetSearch(meili)
	reactionSvc := service.NewReactionService(reactionStore)
	pinSvc := service.NewPinService(pinStore)
	commandSvc := service.NewSlashCommandService(commandStore)
	categorySvc := service.NewCategoryService(categoryStore)
	auditSvc := service.NewAuditService(auditStore)
	retentionSvc := service.NewRetentionService(messageStore)

	// Start retention service.
	retentionCtx, retentionCancel := context.WithCancel(context.Background())
	go retentionSvc.Start(retentionCtx)

	// HTTP handlers & router.
	handlers := api.NewHandlers(
		userSvc, teamSvc, channelSvc, msgSvc,
		reactionSvc, pinSvc, commandSvc, categorySvc, auditSvc,
		meili, db,
	)
	router := api.SetupRouter(handlers, cfg.JWT.Secret)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// WuKongIM webhook server for getSubscribers callback.
	webhookHandler := wkim.NewWebhookHandler(channelSvc.GetMembers)
	mux := http.NewServeMux()
	mux.Handle("/getSubscribers", webhookHandler)
	wkHookServer := &http.Server{
		Addr:    cfg.WuKong.WebhookAddr,
		Handler: mux,
	}

	return &Server{
		cfg:       cfg,
		db:        db,
		redis:     rdb,
		http:      httpServer,
		wkHook:    wkHookServer,
		cancelCtx: retentionCancel,
	}, nil
}

func (s *Server) Start() error {
	// Start WuKongIM webhook server.
	go func() {
		log.Info().Str("addr", s.wkHook.Addr).Msg("starting wukongim webhook server")
		if err := s.wkHook.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("wukongim webhook server error")
		}
	}()

	// Start main HTTP server.
	log.Info().Str("addr", s.http.Addr).Msg("starting http server")
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.cancelCtx()

	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http: %w", err)
	}
	if err := s.wkHook.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown webhook: %w", err)
	}
	s.redis.Close()
	s.db.Close()
	return nil
}
