package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
)

type RetentionService struct {
	msgStore *postgres.MessageStore
}

func NewRetentionService(msgStore *postgres.MessageStore) *RetentionService {
	return &RetentionService{msgStore: msgStore}
}

// Start runs the retention cleanup loop. Call this in a goroutine.
func (s *RetentionService) Start(ctx context.Context) {
	// Run once at startup.
	s.runCleanup()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("retention service stopped")
			return
		case <-ticker.C:
			s.runCleanup()
		}
	}
}

func (s *RetentionService) runCleanup() {
	ctx := context.Background()

	n, err := s.msgStore.DeleteExpiredByChannel(ctx)
	if err != nil {
		log.Error().Err(err).Msg("retention: delete expired by channel")
	} else if n > 0 {
		log.Info().Int64("count", n).Msg("retention: deleted expired messages (channel policy)")
	}

	n, err = s.msgStore.DeleteExpiredByTeam(ctx)
	if err != nil {
		log.Error().Err(err).Msg("retention: delete expired by team")
	} else if n > 0 {
		log.Info().Int64("count", n).Msg("retention: deleted expired messages (team policy)")
	}
}
