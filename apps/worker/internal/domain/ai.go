package domain
import (
	"time"
	"github.com/google/uuid"
)
type AIProviderType string
const (
	AIProviderTypeLocal AIProviderType = "local"
	AIProviderTypeCloud AIProviderType = "cloud"
)
type AIProvider struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Type        AIProviderType `json:"type"`
	Endpoint    string         `json:"endpoint"`
	EncryptedAPIKey string     `json:"-"` // Encrypted at rest
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
type AIInteraction struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Provider      string     `json:"provider"`
	Model         string     `json:"model"`
	PromptHash    string     `json:"prompt_hash"`
	ResponseHash  string     `json:"response_hash"`
	AdapterDraftID *uuid.UUID `json:"adapter_draft_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}
type AIProviderRepository interface {
	Create(provider *AIProvider) error
	GetByID(id uuid.UUID) (*AIProvider, error)
	Update(provider *AIProvider) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*AIProvider, error)
}
type AIInteractionRepository interface {
	Create(interaction *AIInteraction) error
	GetByID(id uuid.UUID) (*AIInteraction, error)
	List(limit, offset int) ([]*AIInteraction, error)
}
