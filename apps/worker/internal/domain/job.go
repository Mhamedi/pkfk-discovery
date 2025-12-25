package domain
import (
	""
	"time"
	"github.com/google/uuid"
)
type JobType string
const (
	JobTypeAdapterProbe    JobType = "adapter_probe"
	JobTypeAdapterValidate JobType = "adapter_validate"
	JobTypeScanRun         JobType = "scan_run"
)
type JobStatus string
const (
	JobStatusPending   JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)
type Job struct {
	ID          uuid.UUID       `json:"id"`
	Type        JobType         `json:"type"`
	Status      JobStatus       `json:"status"`
	Payload     json.RawMessage `json:"payload"`
	RetryCount  int             `json:"retry_count"`
	ErrorMessage string          `json:"error_message,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
type JobRepository interface {
	Create(job *Job) error
	GetByID(id uuid.UUID) (*Job, error)
	Update(job *Job) error
	ListByStatus(status JobStatus, limit int) ([]*Job, error)
}
