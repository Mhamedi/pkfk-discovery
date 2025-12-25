package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pkfk-discovery/api/internal/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(user *domain.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashedPassword, err := HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	return s.userRepo.Create(user)
}

func (s *UserService) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Update(user *domain.User) error {
	// If password is being updated, hash it
	if user.PasswordHash != "" && len(user.PasswordHash) < 60 { // bcrypt hashes are 60 chars
		hashedPassword, err := HashPassword(user.PasswordHash)
		if err != nil {
			return err
		}
		user.PasswordHash = hashedPassword
	}

	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) List(limit, offset int) ([]*domain.User, error) {
	users, err := s.userRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	// Don't expose password hashes
	for _, user := range users {
		user.PasswordHash = ""
	}
	return users, nil
}

