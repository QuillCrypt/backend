package service

import (
	"context"
	"quillcrypt-backend/internal/core/domain"
	"quillcrypt-backend/internal/core/port"

	"github.com/google/uuid"
)

type userService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) port.UserService {
	return &userService{repo}
}

func (s *userService) RegisterOrLogin(ctx context.Context, user *domain.User) (*domain.User, error) {
	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		existing.AvatarURL = user.AvatarURL
		err = s.repo.Update(ctx, existing)
		return existing, err
	}

	err = s.repo.Create(ctx, user)
	return user, err
}

func (s *userService) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
