package port

import (
	"context"
	"quillcrypt-backend/internal/core/domain"
)

type UserService interface {
	RegisterOrLogin(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserById(ctx context.Context, id int64) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type MessageService interface {
	SendMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error)
	DeleteMessage(ctx context.Context, msg *domain.MessageDelete) error
}
