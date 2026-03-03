package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
)

type AuditService struct {
	store *postgres.AuditStore
}

func NewAuditService(store *postgres.AuditStore) *AuditService {
	return &AuditService{store: store}
}

// Log records an audit event. Runs async to avoid blocking the caller.
func (s *AuditService) Log(userID, action, entityType, entityID, ipAddr string, details any) {
	go func() {
		detailsJSON, _ := json.Marshal(details)
		entry := &model.AuditLog{
			ID:         ulid.Make().String(),
			UserID:     userID,
			Action:     action,
			EntityType: entityType,
			EntityID:   entityID,
			IPAddr:     ipAddr,
			Details:    detailsJSON,
			CreatedAt:  time.Now().UTC(),
		}
		s.store.Create(context.Background(), entry)
	}()
}

func (s *AuditService) List(ctx context.Context, query *model.AuditLogQuery) ([]*model.AuditLog, error) {
	return s.store.List(ctx, query)
}
