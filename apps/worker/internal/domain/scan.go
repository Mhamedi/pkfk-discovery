package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)
type ScanStatus string
const (
	ScanStatusPending   ScanStatus = "pending"
	ScanStatusRunning   ScanStatus = "running"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusFailed    ScanStatus = "failed"
	ScanStatusCancelled ScanStatus = "cancelled"
)
type ScanPolicy struct {
	SampleMode    bool   `json:"sample_mode"`
	DeepMode      bool   `json:"deep_mode"`
	Timeout       int    `json:"timeout"` // seconds
	MaxRows       int    `json:"max_rows"`
	Concurrency   int    `json:"concurrency"`
}
type Scan struct {
	ID           uuid.UUID       `json:"id"`
	ConnectionID uuid.UUID       `json:"connection_id"`
	AdapterID    uuid.UUID       `json:"adapter_id"`
	Status       ScanStatus      `json:"status"`
	Policy       json.RawMessage `json:"policy"` // ScanPolicy as JSON
	Results      json.RawMessage `json:"results,omitempty"`
	CreatedBy    uuid.UUID       `json:"created_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
type ScanRepository interface {
	Create(scan *Scan) error
	GetByID(id uuid.UUID) (*Scan, error)
	Update(scan *Scan) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*Scan, error)
}
