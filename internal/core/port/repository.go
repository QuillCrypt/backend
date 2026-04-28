package port

import (
	"context"
	"quillcrypt-backend/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
}

type MessageRepository interface {
	Save(ctx context.Context, msg *domain.Message) (*domain.Message, error)
	Delete(ctx context.Context, msg *domain.MessageDelete) error
}

type KeyRepository interface {
	StorePublicKey(ctx context.Context, key *domain.PublicKey) error
	GetPublicKeyByUserID(ctx context.Context, userID int64) (*domain.PublicKey, error)
}
