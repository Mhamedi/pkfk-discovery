package domain
import (
	""
	"time"
	"github.com/google/uuid"
)
type AuditLog struct {
	ID          uuid.UUID       `json:"id"`
	UserID      *uuid.UUID      `json:"user_id,omitempty"`
	Action      string          `json:"action"`
	ResourceType string         `json:"resource_type"`
	ResourceID  *uuid.UUID      `json:"resource_id,omitempty"`
	Details     json.RawMessage `json:"details"`
	IPAddress   string          `json:"ip_address"`
	CreatedAt   time.Time       `json:"created_at"`
}
type AuditLogRepository interface {
	Create(log *AuditLog) error
	GetByID(id uuid.UUID) (*AuditLog, error)
	List(filters AuditLogFilters, limit, offset int) ([]*AuditLog, error)
}
type AuditLogFilters struct {
	UserID       *uuid.UUID
	Action       string
	ResourceType string
	StartDate    *time.Time
	EndDate      *time.Time
}
