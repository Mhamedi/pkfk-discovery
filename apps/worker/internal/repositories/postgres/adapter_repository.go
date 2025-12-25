package postgres
import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkfk-discovery/worker/internal/domain"
)
type AdapterRepository struct {
	db *pgxpool.Pool
}
func NewAdapterRepository(db *pgxpool.Pool) *AdapterRepository {
	return &AdapterRepository{db: db}
}
func (r *AdapterRepository) Create(adapter *domain.Adapter) error {
	query := `
		INSERT INTO adapters (id, name, vendor, db_family, version, maturity_level, bundle_path, signature, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	now := time.Now()
	_, err := r.db.Exec(context.Background(), query,
		adapter.ID, adapter.Name, adapter.Vendor, adapter.DBFamily,
		adapter.Version, adapter.MaturityLevel, adapter.BundlePath,
		adapter.Signature, now, now)
	return err
}
func (r *AdapterRepository) GetByID(id uuid.UUID) (*domain.Adapter, error) {
	query := `
		SELECT id, name, vendor, db_family, version, maturity_level, bundle_path, signature, created_at, updated_at
		FROM adapters
		WHERE id = $1
	`
	var adapter domain.Adapter
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&adapter.ID, &adapter.Name, &adapter.Vendor, &adapter.DBFamily,
		&adapter.Version, &adapter.MaturityLevel, &adapter.BundlePath,
		&adapter.Signature, &adapter.CreatedAt, &adapter.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &adapter, nil
}
func (r *AdapterRepository) Update(adapter *domain.Adapter) error {
	query := `
		UPDATE adapters
		SET name = $2, vendor = $3, db_family = $4, version = $5, maturity_level = $6, bundle_path = $7, signature = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.Exec(context.Background(), query,
		adapter.ID, adapter.Name, adapter.Vendor, adapter.DBFamily,
		adapter.Version, adapter.MaturityLevel, adapter.BundlePath,
		adapter.Signature, adapter.UpdatedAt)
	return err
}
func (r *AdapterRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM adapters WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
func (r *AdapterRepository) List(limit, offset int) ([]*domain.Adapter, error) {
	query := `
		SELECT id, name, vendor, db_family, version, maturity_level, bundle_path, signature, created_at, updated_at
		FROM adapters
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var adapters []*domain.Adapter
	for rows.Next() {
		var adapter domain.Adapter
		if err := rows.Scan(
			&adapter.ID, &adapter.Name, &adapter.Vendor, &adapter.DBFamily,
			&adapter.Version, &adapter.MaturityLevel, &adapter.BundlePath,
			&adapter.Signature, &adapter.CreatedAt, &adapter.UpdatedAt,
		); err != nil {
			return nil, err
		}
		adapters = append(adapters, &adapter)
	}
	return adapters, rows.Err()
}
