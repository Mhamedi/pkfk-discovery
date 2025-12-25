package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pkfk-discovery/api/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	now := time.Now()
	_, err := r.db.Exec(context.Background(), query,
		user.ID, user.Email, user.PasswordHash, user.Role, now, now)
	return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at, last_login
		FROM users
		WHERE id = $1
	`
	var user domain.User
	var lastLogin sql.NullTime
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &lastLogin,
	)
	if err != nil {
		return nil, err
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at, last_login
		FROM users
		WHERE email = $1
	`
	var user domain.User
	var lastLogin sql.NullTime
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &lastLogin,
	)
	if err != nil {
		return nil, err
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, role = $4, updated_at = $5, last_login = $6
		WHERE id = $1
	`
	_, err := r.db.Exec(context.Background(), query,
		user.ID, user.Email, user.PasswordHash, user.Role,
		user.UpdatedAt, user.LastLogin)
	return err
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *UserRepository) List(limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at, last_login
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var lastLogin sql.NullTime
		if err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Role,
			&user.CreatedAt, &user.UpdatedAt, &lastLogin,
		); err != nil {
			return nil, err
		}
		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}
		users = append(users, &user)
	}

	return users, rows.Err()
}

