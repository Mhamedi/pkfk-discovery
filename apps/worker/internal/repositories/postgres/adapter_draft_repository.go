package postgres
import (
	"context"
	"database/sql"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkfk-discovery/worker/internal/domain"
)
type AdapterDraftRepository struct {
	db *pgxpool.Pool
}
func NewAdapterDraftRepository(db *pgxpool.Pool) *AdapterDraftRepository {
	return &AdapterDraftRepository{db: db}
}
func (r *AdapterDraftRepository) Create(draft *domain.AdapterDraft) error {
	query := `
		INSERT INTO adapter_drafts (id, adapter_id, name, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	_, err := r.db.Exec(context.Background(), query,
		draft.ID, draft.AdapterID, draft.Name, draft.Status,
		draft.CreatedBy, now, now)
	return err
}
func (r *AdapterDraftRepository) GetByID(id uuid.UUID) (*domain.AdapterDraft, error) {
	query := `
		SELECT id, adapter_id, name, status, created_by, created_at, updated_at
		FROM adapter_drafts
		WHERE id = $1
	`
	var draft domain.AdapterDraft
	var adapterID sql.NullString
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&draft.ID, &adapterID, &draft.Name, &draft.Status,
		&draft.CreatedBy, &draft.CreatedAt, &draft.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if adapterID.Valid {
		parsedID, err := uuid.Parse(adapterID.String)
		if err == nil {
			draft.AdapterID = &parsedID
		}
	}
	return &draft, nil
}
func (r *AdapterDraftRepository) Update(draft *domain.AdapterDraft) error {
	query := `
		UPDATE adapter_drafts
		SET adapter_id = $2, name = $3, status = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.Exec(context.Background(), query,
		draft.ID, draft.AdapterID, draft.Name, draft.Status, draft.UpdatedAt)
	return err
}
func (r *AdapterDraftRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM adapter_drafts WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
func (r *AdapterDraftRepository) List(limit, offset int) ([]*domain.AdapterDraft, error) {
	query := `
		SELECT id, adapter_id, name, status, created_by, created_at, updated_at
		FROM adapter_drafts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var drafts []*domain.AdapterDraft
	for rows.Next() {
		var draft domain.AdapterDraft
		var adapterID sql.NullString
		if err := rows.Scan(
			&draft.ID, &adapterID, &draft.Name, &draft.Status,
			&draft.CreatedBy, &draft.CreatedAt, &draft.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if adapterID.Valid {
			parsedID, err := uuid.Parse(adapterID.String)
			if err == nil {
				draft.AdapterID = &parsedID
			}
		}
		drafts = append(drafts, &draft)
	}
	return drafts, rows.Err()
}
