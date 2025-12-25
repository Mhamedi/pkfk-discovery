package services

import (
	"errors"
)

var (
	ErrSettingsNotFound = errors.New("settings not found")
)

type SettingsService struct {
	settingsRepo interface {
		Get(key string) (map[string]interface{}, error)
		Set(key string, value map[string]interface{}) error
		GetAll() (map[string]map[string]interface{}, error)
	}
}

func NewSettingsService(settingsRepo interface {
	Get(key string) (map[string]interface{}, error)
	Set(key string, value map[string]interface{}) error
	GetAll() (map[string]map[string]interface{}, error)
}) *SettingsService {
	return &SettingsService{
		settingsRepo: settingsRepo,
	}
}

func (s *SettingsService) Get(key string) (map[string]interface{}, error) {
	return s.settingsRepo.Get(key)
}

func (s *SettingsService) Set(key string, value map[string]interface{}) error {
	return s.settingsRepo.Set(key, value)
}

func (s *SettingsService) GetAll() (map[string]map[string]interface{}, error) {
	return s.settingsRepo.GetAll()
}

