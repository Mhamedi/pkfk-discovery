package postgres
import (
	"context"
	""
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkfk-discovery/worker/internal/domain"
)
type ScanRepository struct {
	db *pgxpool.Pool
}
func NewScanRepository(db *pgxpool.Pool) *ScanRepository {
	return &ScanRepository{db: db}
}
func (r *ScanRepository) Create(scan *domain.Scan) error {
	query := `
		INSERT INTO scans (id, connection_id, adapter_id, status, policy_json, results_json, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	now := time.Now()
	_, err := r.db.Exec(context.Background(), query,
		scan.ID, scan.ConnectionID, scan.AdapterID, scan.Status,
		scan.Policy, scan.Results, scan.CreatedBy, now, now)
	return err
}
func (r *ScanRepository) GetByID(id uuid.UUID) (*domain.Scan, error) {
	query := `
		SELECT id, connection_id, adapter_id, status, policy_json, results_json, created_by, created_at, updated_at
		FROM scans
		WHERE id = $1
	`
	var scan domain.Scan
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&scan.ID, &scan.ConnectionID, &scan.AdapterID, &scan.Status,
		&scan.Policy, &scan.Results, &scan.CreatedBy, &scan.CreatedAt, &scan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &scan, nil
}
func (r *ScanRepository) Update(scan *domain.Scan) error {
	query := `
		UPDATE scans
		SET status = $2, policy_json = $3, results_json = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.Exec(context.Background(), query,
		scan.ID, scan.Status, scan.Policy, scan.Results, scan.UpdatedAt)
	return err
}
func (r *ScanRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM scans WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
func (r *ScanRepository) List(limit, offset int) ([]*domain.Scan, error) {
	query := `
		SELECT id, connection_id, adapter_id, status, policy_json, results_json, created_by, created_at, updated_at
		FROM scans
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scans []*domain.Scan
	for rows.Next() {
		var scan domain.Scan
		if err := rows.Scan(
			&scan.ID, &scan.ConnectionID, &scan.AdapterID, &scan.Status,
			&scan.Policy, &scan.Results, &scan.CreatedBy, &scan.CreatedAt, &scan.UpdatedAt,
		); err != nil {
			return nil, err
		}
		scans = append(scans, &scan)
	}
	return scans, rows.Err()
}
