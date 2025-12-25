package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/pkfk-discovery/api/internal/domain"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
	ErrInvalidEncryptionKey = errors.New("invalid encryption key")
)

type ConnectionService struct {
	connRepo     domain.ConnectionRepository
	encryptionKey []byte
}

func NewConnectionService(connRepo domain.ConnectionRepository, encryptionKey string) (*ConnectionService, error) {
	key := []byte(encryptionKey)
	if len(key) != 32 {
		return nil, ErrInvalidEncryptionKey
	}

	return &ConnectionService{
		connRepo:     connRepo,
		encryptionKey: key,
	}, nil
}

func (s *ConnectionService) Create(conn *domain.Connection) error {
	encrypted, err := s.encrypt(conn.EncryptedPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}
	conn.EncryptedPassword = encrypted
	return s.connRepo.Create(conn)
}

func (s *ConnectionService) GetByID(id uuid.UUID) (*domain.Connection, error) {
	conn, err := s.connRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	decrypted, err := s.decrypt(conn.EncryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}
	conn.EncryptedPassword = decrypted
	return conn, nil
}

func (s *ConnectionService) Update(conn *domain.Connection) error {
	// If password is being updated, encrypt it
	if conn.EncryptedPassword != "" {
		encrypted, err := s.encrypt(conn.EncryptedPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt password: %w", err)
		}
		conn.EncryptedPassword = encrypted
	}
	return s.connRepo.Update(conn)
}

func (s *ConnectionService) Delete(id uuid.UUID) error {
	return s.connRepo.Delete(id)
}

func (s *ConnectionService) List(limit, offset int) ([]*domain.Connection, error) {
	conns, err := s.connRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	// Don't decrypt passwords in list view for security
	for _, conn := range conns {
		conn.EncryptedPassword = "[REDACTED]"
	}
	return conns, nil
}

func (s *ConnectionService) TestConnection(id uuid.UUID) error {
	conn, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// TODO: Actually test the connection by opening it and running SELECT 1
	// This would require database-specific drivers
	// For now, just return success
	return nil
}

func (s *ConnectionService) Encrypt(plaintext string) (string, error) {
	return s.encrypt(plaintext)
}

func (s *ConnectionService) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *ConnectionService) decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

