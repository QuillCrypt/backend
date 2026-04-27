package port

import (
	"context"
	"quillcrypt-backend/internal/core/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MessageRepository interface {
	Save(ctx context.Context, msg *domain.Message) (*domain.Message, error)
	Delete(ctx context.Context, msg *domain.MessageDelete) error
}

type KeyRepository interface {
	StorePublicKey(ctx context.Context, key *domain.PublicKey) error
	GetPublicKeyByUserID(ctx context.Context, userID uuid.UUID) (*domain.PublicKey, error)
}
