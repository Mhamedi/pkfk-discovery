package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)
type SettingsRepository struct {
	db *pgxpool.Pool
}
func NewSettingsRepository(db *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{db: db}
}
func (r *SettingsRepository) Get(key string) (map[string]interface{}, error) {
	query := `SELECT value_json FROM settings WHERE key = $1`
	var valueJSON json.RawMessage
	err := r.db.QueryRow(context.Background(), query, key).Scan(&valueJSON)
	if err != nil {
		return nil, err
	}
	var value map[string]interface{}
	if err := json.Unmarshal(valueJSON, &value); err != nil {
		return nil, err
	}
	return value, nil
}
func (r *SettingsRepository) Set(key string, value map[string]interface{}) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO settings (key, value_json, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE
		SET value_json = $2, updated_at = $3
	`
	_, err = r.db.Exec(context.Background(), query, key, valueJSON, time.Now())
	return err
}
func (r *SettingsRepository) GetAll() (map[string]map[string]interface{}, error) {
	query := `SELECT key, value_json FROM settings`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]map[string]interface{})
	for rows.Next() {
		var key string
		var valueJSON json.RawMessage
		if err := rows.Scan(&key, &valueJSON); err != nil {
			return nil, err
		}
		var value map[string]interface{}
		if err := json.Unmarshal(valueJSON, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, rows.Err()
}
