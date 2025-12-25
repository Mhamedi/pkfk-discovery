package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pkfk-discovery/api/internal/domain"
)

var (
	ErrScanNotFound = errors.New("scan not found")
)

type ScanService struct {
	scanRepo domain.ScanRepository
}

func NewScanService(scanRepo domain.ScanRepository) *ScanService {
	return &ScanService{
		scanRepo: scanRepo,
	}
}

func (s *ScanService) Create(scan *domain.Scan) error {
	scan.ID = uuid.New()
	scan.Status = domain.ScanStatusPending
	scan.CreatedAt = time.Now()
	scan.UpdatedAt = time.Now()

	// Ensure policy is valid JSON
	if scan.Policy == nil {
		defaultPolicy := domain.ScanPolicy{
			SampleMode:  true,
			DeepMode:    false,
			Timeout:     300,
			MaxRows:     10000,
			Concurrency: 5,
		}
		policyJSON, _ := json.Marshal(defaultPolicy)
		scan.Policy = policyJSON
	}

	return s.scanRepo.Create(scan)
}

func (s *ScanService) GetByID(id uuid.UUID) (*domain.Scan, error) {
	return s.scanRepo.GetByID(id)
}

func (s *ScanService) Update(scan *domain.Scan) error {
	scan.UpdatedAt = time.Now()
	return s.scanRepo.Update(scan)
}

func (s *ScanService) UpdateStatus(id uuid.UUID, status domain.ScanStatus) error {
	scan, err := s.scanRepo.GetByID(id)
	if err != nil {
		return err
	}

	scan.Status = status
	scan.UpdatedAt = time.Now()
	return s.scanRepo.Update(scan)
}

func (s *ScanService) SetResults(id uuid.UUID, results interface{}) error {
	scan, err := s.scanRepo.GetByID(id)
	if err != nil {
		return err
	}

	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return err
	}

	scan.Results = resultsJSON
	scan.Status = domain.ScanStatusCompleted
	scan.UpdatedAt = time.Now()
	return s.scanRepo.Update(scan)
}

func (s *ScanService) Delete(id uuid.UUID) error {
	return s.scanRepo.Delete(id)
}

func (s *ScanService) List(limit, offset int) ([]*domain.Scan, error) {
	return s.scanRepo.List(limit, offset)
}

