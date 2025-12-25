package domain
import (
	"time"
	"github.com/google/uuid"
)
type MaturityLevel string
const (
	MaturityL0 MaturityLevel = "L0"
	MaturityL1 MaturityLevel = "L1"
	MaturityL2 MaturityLevel = "L2"
	MaturityL3 MaturityLevel = "L3"
	MaturityL4 MaturityLevel = "L4"
)
type Adapter struct {
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	Vendor        string        `json:"vendor"`
	DBFamily      string        `json:"db_family"`
	Version       string        `json:"version"`
	MaturityLevel MaturityLevel `json:"maturity_level"`
	BundlePath    string        `json:"bundle_path"`
	Signature     string        `json:"signature"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
type AdapterDraft struct {
	ID        uuid.UUID `json:"id"`
	AdapterID *uuid.UUID `json:"adapter_id,omitempty"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type AdapterRepository interface {
	Create(adapter *Adapter) error
	GetByID(id uuid.UUID) (*Adapter, error)
	Update(adapter *Adapter) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*Adapter, error)
}
type AdapterDraftRepository interface {
	Create(draft *AdapterDraft) error
	GetByID(id uuid.UUID) (*AdapterDraft, error)
	Update(draft *AdapterDraft) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*AdapterDraft, error)
}
