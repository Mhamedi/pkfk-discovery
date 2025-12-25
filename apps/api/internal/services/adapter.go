package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pkfk-discovery/api/internal/domain"
)

var (
	ErrAdapterNotFound = errors.New("adapter not found")
)

type AdapterService struct {
	adapterRepo domain.AdapterRepository
}

func NewAdapterService(adapterRepo domain.AdapterRepository) *AdapterService {
	return &AdapterService{
		adapterRepo: adapterRepo,
	}
}

func (s *AdapterService) Create(adapter *domain.Adapter) error {
	adapter.ID = uuid.New()
	adapter.CreatedAt = time.Now()
	adapter.UpdatedAt = time.Now()
	return s.adapterRepo.Create(adapter)
}

func (s *AdapterService) GetByID(id uuid.UUID) (*domain.Adapter, error) {
	return s.adapterRepo.GetByID(id)
}

func (s *AdapterService) Update(adapter *domain.Adapter) error {
	adapter.UpdatedAt = time.Now()
	return s.adapterRepo.Update(adapter)
}

func (s *AdapterService) Delete(id uuid.UUID) error {
	return s.adapterRepo.Delete(id)
}

func (s *AdapterService) List(limit, offset int) ([]*domain.Adapter, error) {
	return s.adapterRepo.List(limit, offset)
}

