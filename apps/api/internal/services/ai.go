package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pkfk-discovery/api/internal/domain"
)

var (
	ErrAIProviderNotFound = errors.New("AI provider not found")
)

type AIProviderService struct {
	providerRepo domain.AIProviderRepository
	encryptionKey []byte
}

func NewAIProviderService(providerRepo domain.AIProviderRepository, encryptionKey string) (*AIProviderService, error) {
	key := []byte(encryptionKey)
	if len(key) != 32 {
		return nil, ErrInvalidEncryptionKey
	}

	return &AIProviderService{
		providerRepo:  providerRepo,
		encryptionKey: key,
	}, nil
}

func (s *AIProviderService) Create(provider *domain.AIProvider) error {
	provider.ID = uuid.New()
	provider.CreatedAt = time.Now()
	provider.UpdatedAt = time.Now()

	// Encrypt API key
	encrypted, err := encryptString(provider.EncryptedAPIKey, s.encryptionKey)
	if err != nil {
		return err
	}
	provider.EncryptedAPIKey = encrypted

	return s.providerRepo.Create(provider)
}

func (s *AIProviderService) GetByID(id uuid.UUID) (*domain.AIProvider, error) {
	provider, err := s.providerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Don't decrypt API key in get - it should only be used internally
	provider.EncryptedAPIKey = "[REDACTED]"
	return provider, nil
}

func (s *AIProviderService) Update(provider *domain.AIProvider) error {
	provider.UpdatedAt = time.Now()

	// If API key is being updated, encrypt it
	if provider.EncryptedAPIKey != "" && provider.EncryptedAPIKey != "[REDACTED]" {
		encrypted, err := encryptString(provider.EncryptedAPIKey, s.encryptionKey)
		if err != nil {
			return err
		}
		provider.EncryptedAPIKey = encrypted
	}

	return s.providerRepo.Update(provider)
}

func (s *AIProviderService) Delete(id uuid.UUID) error {
	return s.providerRepo.Delete(id)
}

func (s *AIProviderService) List(limit, offset int) ([]*domain.AIProvider, error) {
	providers, err := s.providerRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	// Don't expose API keys
	for _, provider := range providers {
		provider.EncryptedAPIKey = "[REDACTED]"
	}
	return providers, nil
}

func encryptString(plaintext string, key []byte) (string, error) {
	// Reuse encryption from connection service
	connService, err := NewConnectionService(nil, string(key))
	if err != nil {
		return "", err
	}
	return connService.Encrypt(plaintext)
}

