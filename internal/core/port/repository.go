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
	Save(ctx context.Context, msg *domain.Message) error
	GetConversation(ctx context.Context, user1, user2 uuid.UUID, limit int, offset int) ([]domain.Message, error)
}

type KeyRepository interface {
	StorePublicKey(ctx context.Context, key *domain.PublicKey) error
	GetPublicKeyByUserID(ctx context.Context, userID uuid.UUID) (*domain.PublicKey, error)
}
