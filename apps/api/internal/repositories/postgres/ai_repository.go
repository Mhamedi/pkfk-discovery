package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pkfk-discovery/api/internal/domain"
)

type AIProviderRepository struct {
	db *pgxpool.Pool
}

func NewAIProviderRepository(db *pgxpool.Pool) *AIProviderRepository {
	return &AIProviderRepository{db: db}
}

func (r *AIProviderRepository) Create(provider *domain.AIProvider) error {
	query := `
		INSERT INTO ai_providers (id, name, type, endpoint, encrypted_api_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	_, err := r.db.Exec(context.Background(), query,
		provider.ID, provider.Name, provider.Type, provider.Endpoint,
		provider.EncryptedAPIKey, now, now)
	return err
}

func (r *AIProviderRepository) GetByID(id uuid.UUID) (*domain.AIProvider, error) {
	query := `
		SELECT id, name, type, endpoint, encrypted_api_key, created_at, updated_at
		FROM ai_providers
		WHERE id = $1
	`
	var provider domain.AIProvider
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&provider.ID, &provider.Name, &provider.Type, &provider.Endpoint,
		&provider.EncryptedAPIKey, &provider.CreatedAt, &provider.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *AIProviderRepository) Update(provider *domain.AIProvider) error {
	query := `
		UPDATE ai_providers
		SET name = $2, type = $3, endpoint = $4, encrypted_api_key = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.Exec(context.Background(), query,
		provider.ID, provider.Name, provider.Type, provider.Endpoint,
		provider.EncryptedAPIKey, provider.UpdatedAt)
	return err
}

func (r *AIProviderRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM ai_providers WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *AIProviderRepository) List(limit, offset int) ([]*domain.AIProvider, error) {
	query := `
		SELECT id, name, type, endpoint, encrypted_api_key, created_at, updated_at
		FROM ai_providers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*domain.AIProvider
	for rows.Next() {
		var provider domain.AIProvider
		if err := rows.Scan(
			&provider.ID, &provider.Name, &provider.Type, &provider.Endpoint,
			&provider.EncryptedAPIKey, &provider.CreatedAt, &provider.UpdatedAt,
		); err != nil {
			return nil, err
		}
		providers = append(providers, &provider)
	}

	return providers, rows.Err()
}

type AIInteractionRepository struct {
	db *pgxpool.Pool
}

func NewAIInteractionRepository(db *pgxpool.Pool) *AIInteractionRepository {
	return &AIInteractionRepository{db: db}
}

func (r *AIInteractionRepository) Create(interaction *domain.AIInteraction) error {
	query := `
		INSERT INTO ai_interactions (id, user_id, provider, model, prompt_hash, response_hash, adapter_draft_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now()
	_, err := r.db.Exec(context.Background(), query,
		interaction.ID, interaction.UserID, interaction.Provider, interaction.Model,
		interaction.PromptHash, interaction.ResponseHash, interaction.AdapterDraftID, now)
	return err
}

func (r *AIInteractionRepository) GetByID(id uuid.UUID) (*domain.AIInteraction, error) {
	query := `
		SELECT id, user_id, provider, model, prompt_hash, response_hash, adapter_draft_id, created_at
		FROM ai_interactions
		WHERE id = $1
	`
	var interaction domain.AIInteraction
	var adapterDraftID sql.NullString
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&interaction.ID, &interaction.UserID, &interaction.Provider, &interaction.Model,
		&interaction.PromptHash, &interaction.ResponseHash, &adapterDraftID, &interaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if adapterDraftID.Valid {
		parsedID, err := uuid.Parse(adapterDraftID.String)
		if err == nil {
			interaction.AdapterDraftID = &parsedID
		}
	}
	return &interaction, nil
}

func (r *AIInteractionRepository) List(limit, offset int) ([]*domain.AIInteraction, error) {
	query := `
		SELECT id, user_id, provider, model, prompt_hash, response_hash, adapter_draft_id, created_at
		FROM ai_interactions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interactions []*domain.AIInteraction
	for rows.Next() {
		var interaction domain.AIInteraction
		var adapterDraftID sql.NullString
		if err := rows.Scan(
			&interaction.ID, &interaction.UserID, &interaction.Provider, &interaction.Model,
			&interaction.PromptHash, &interaction.ResponseHash, &adapterDraftID, &interaction.CreatedAt,
		); err != nil {
			return nil, err
		}
		if adapterDraftID.Valid {
			parsedID, err := uuid.Parse(adapterDraftID.String)
			if err == nil {
				interaction.AdapterDraftID = &parsedID
			}
		}
		interactions = append(interactions, &interaction)
	}

	return interactions, rows.Err()
}

