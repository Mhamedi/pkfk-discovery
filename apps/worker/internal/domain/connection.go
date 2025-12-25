package domain
import (
	"time"
	"github.com/google/uuid"
)
type Connection struct {
	ID               uuid.UUID `json:"id"`
	Name             string   `json:"name"`
	DBType           string   `json:"db_type"`
	Host             string   `json:"host"`
	Port             int      `json:"port"`
	Database         string   `json:"database"`
	Username         string   `json:"username"`
	EncryptedPassword string  `json:"-"` // Encrypted at rest
	CreatedBy        uuid.UUID `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
type ConnectionRepository interface {
	Create(conn *Connection) error
	GetByID(id uuid.UUID) (*Connection, error)
	Update(conn *Connection) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*Connection, error)
}
