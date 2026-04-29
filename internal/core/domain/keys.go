package domain

import (
	"time"

	"github.com/google/uuid"
)

type PublicKey struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	KeyData   string    `json:"key_data"`
	CreatedAt time.Time `json:"created_at"`
}
