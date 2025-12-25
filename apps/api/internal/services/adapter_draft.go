package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pkfk-discovery/api/internal/domain"
)

var (
	ErrAdapterDraftNotFound = errors.New("adapter draft not found")
)

type AdapterDraftService struct {
	draftRepo domain.AdapterDraftRepository
}

func NewAdapterDraftService(draftRepo domain.AdapterDraftRepository) *AdapterDraftService {
	return &AdapterDraftService{
		draftRepo: draftRepo,
	}
}

func (s *AdapterDraftService) Create(draft *domain.AdapterDraft) error {
	draft.ID = uuid.New()
	draft.Status = "draft"
	draft.CreatedAt = time.Now()
	draft.UpdatedAt = time.Now()
	return s.draftRepo.Create(draft)
}

func (s *AdapterDraftService) GetByID(id uuid.UUID) (*domain.AdapterDraft, error) {
	return s.draftRepo.GetByID(id)
}

func (s *AdapterDraftService) Update(draft *domain.AdapterDraft) error {
	draft.UpdatedAt = time.Now()
	return s.draftRepo.Update(draft)
}

func (s *AdapterDraftService) Delete(id uuid.UUID) error {
	return s.draftRepo.Delete(id)
}

func (s *AdapterDraftService) List(limit, offset int) ([]*domain.AdapterDraft, error) {
	return s.draftRepo.List(limit, offset)
}

