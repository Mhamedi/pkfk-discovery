package postgres
import (
	"context"
	"database/sql"
	""
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkfk-discovery/worker/internal/domain"
)
type AuditLogRepository struct {
	db *pgxpool.Pool
}
func NewAuditLogRepository(db *pgxpool.Pool) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}
func (r *AuditLogRepository) Create(log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, action, resource_type, resource_id, details_json, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(context.Background(), query,
		log.ID, log.UserID, log.Action, log.ResourceType,
		log.ResourceID, log.Details, log.IPAddress, log.CreatedAt)
	return err
}
func (r *AuditLogRepository) GetByID(id uuid.UUID) (*domain.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource_type, resource_id, details_json, ip_address, created_at
		FROM audit_logs
		WHERE id = $1
	`
	var log domain.AuditLog
	var userID sql.NullString
	var resourceID sql.NullString
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&log.ID, &userID, &log.Action, &log.ResourceType,
		&resourceID, &log.Details, &log.IPAddress, &log.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if userID.Valid {
		parsedID, err := uuid.Parse(userID.String)
		if err == nil {
			log.UserID = &parsedID
		}
	}
	if resourceID.Valid {
		parsedID, err := uuid.Parse(resourceID.String)
		if err == nil {
			log.ResourceID = &parsedID
		}
	}
	return &log, nil
}
func (r *AuditLogRepository) List(filters domain.AuditLogFilters, limit, offset int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource_type, resource_id, details_json, ip_address, created_at
		FROM audit_logs
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1
	if filters.UserID != nil {
		query += ` AND user_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.UserID)
		argPos++
	}
	if filters.Action != "" {
		query += ` AND action = $` + fmt.Sprintf("%d", argPos)
		args = append(args, filters.Action)
		argPos++
	}
	if filters.ResourceType != "" {
		query += ` AND resource_type = $` + fmt.Sprintf("%d", argPos)
		args = append(args, filters.ResourceType)
		argPos++
	}
	if filters.StartDate != nil {
		query += ` AND created_at >= $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}
	if filters.EndDate != nil {
		query += ` AND created_at <= $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argPos) + ` OFFSET $` + fmt.Sprintf("%d", argPos+1)
	args = append(args, limit, offset)
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []*domain.AuditLog
	for rows.Next() {
		var log domain.AuditLog
		var userID sql.NullString
		var resourceID sql.NullString
		if err := rows.Scan(
			&log.ID, &userID, &log.Action, &log.ResourceType,
			&resourceID, &log.Details, &log.IPAddress, &log.CreatedAt,
		); err != nil {
			return nil, err
		}
		if userID.Valid {
			parsedID, err := uuid.Parse(userID.String)
			if err == nil {
				log.UserID = &parsedID
			}
		}
		if resourceID.Valid {
			parsedID, err := uuid.Parse(resourceID.String)
			if err == nil {
				log.ResourceID = &parsedID
			}
		}
		logs = append(logs, &log)
	}
	return logs, rows.Err()
}
