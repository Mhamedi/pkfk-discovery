package services

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/pkfk-discovery/api/internal/domain"
)

type AuditService struct {
	auditRepo domain.AuditLogRepository
}

func NewAuditService(auditRepo domain.AuditLogRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
}

func (s *AuditService) Log(userID *uuid.UUID, action, resourceType string, resourceID *uuid.UUID, details interface{}, ipAddress string) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return err
	}

	log := &domain.AuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:  resourceID,
		Details:      detailsJSON,
		IPAddress:    ipAddress,
		CreatedAt:    time.Now(),
	}

	return s.auditRepo.Create(log)
}

func (s *AuditService) List(filters domain.AuditLogFilters, limit, offset int) ([]*domain.AuditLog, error) {
	return s.auditRepo.List(filters, limit, offset)
}

func (s *AuditService) GetByID(id uuid.UUID) (*domain.AuditLog, error) {
	return s.auditRepo.GetByID(id)
}

