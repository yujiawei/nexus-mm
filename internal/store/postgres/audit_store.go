package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yujiawei/nexus-mm/internal/model"
)

type AuditStore struct {
	db *sqlx.DB
}

func NewAuditStore(db *sqlx.DB) *AuditStore {
	return &AuditStore{db: db}
}

func (s *AuditStore) Create(ctx context.Context, entry *model.AuditLog) error {
	query := `INSERT INTO audit_logs (id, user_id, action, entity_type, entity_id, ip_addr, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query,
		entry.ID, entry.UserID, entry.Action, entry.EntityType,
		entry.EntityID, entry.IPAddr, entry.Details, entry.CreatedAt)
	return err
}

func (s *AuditStore) List(ctx context.Context, query *model.AuditLogQuery) ([]*model.AuditLog, error) {
	if query.Limit <= 0 || query.Limit > 200 {
		query.Limit = 50
	}

	q := "SELECT * FROM audit_logs WHERE 1=1"
	args := []any{}
	argN := 1

	if query.Action != "" {
		q += fmt.Sprintf(" AND action = $%d", argN)
		args = append(args, query.Action)
		argN++
	}
	if query.UserID != "" {
		q += fmt.Sprintf(" AND user_id = $%d", argN)
		args = append(args, query.UserID)
		argN++
	}

	q += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argN)
	args = append(args, query.Limit)

	var logs []*model.AuditLog
	err := s.db.SelectContext(ctx, &logs, q, args...)
	return logs, err
}
