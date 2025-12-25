package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pkfk-discovery/api/internal/domain"
)

type ConnectionRepository struct {
	db *pgxpool.Pool
}

func NewConnectionRepository(db *pgxpool.Pool) *ConnectionRepository {
	return &ConnectionRepository{db: db}
}

func (r *ConnectionRepository) Create(conn *domain.Connection) error {
	query := `
		INSERT INTO connections (id, name, db_type, host, port, database, username, encrypted_password, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	now := time.Now()
	_, err := r.db.Exec(context.Background(), query,
		conn.ID, conn.Name, conn.DBType, conn.Host, conn.Port,
		conn.Database, conn.Username, conn.EncryptedPassword,
		conn.CreatedBy, now, now)
	return err
}

func (r *ConnectionRepository) GetByID(id uuid.UUID) (*domain.Connection, error) {
	query := `
		SELECT id, name, db_type, host, port, database, username, encrypted_password, created_by, created_at, updated_at
		FROM connections
		WHERE id = $1
	`
	var conn domain.Connection
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&conn.ID, &conn.Name, &conn.DBType, &conn.Host, &conn.Port,
		&conn.Database, &conn.Username, &conn.EncryptedPassword,
		&conn.CreatedBy, &conn.CreatedAt, &conn.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (r *ConnectionRepository) Update(conn *domain.Connection) error {
	query := `
		UPDATE connections
		SET name = $2, db_type = $3, host = $4, port = $5, database = $6, username = $7, encrypted_password = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.Exec(context.Background(), query,
		conn.ID, conn.Name, conn.DBType, conn.Host, conn.Port,
		conn.Database, conn.Username, conn.EncryptedPassword, conn.UpdatedAt)
	return err
}

func (r *ConnectionRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM connections WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *ConnectionRepository) List(limit, offset int) ([]*domain.Connection, error) {
	query := `
		SELECT id, name, db_type, host, port, database, username, encrypted_password, created_by, created_at, updated_at
		FROM connections
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []*domain.Connection
	for rows.Next() {
		var conn domain.Connection
		if err := rows.Scan(
			&conn.ID, &conn.Name, &conn.DBType, &conn.Host, &conn.Port,
			&conn.Database, &conn.Username, &conn.EncryptedPassword,
			&conn.CreatedBy, &conn.CreatedAt, &conn.UpdatedAt,
		); err != nil {
			return nil, err
		}
		connections = append(connections, &conn)
	}

	return connections, rows.Err()
}

