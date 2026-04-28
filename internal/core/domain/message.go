package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	CHAT   MessageType = "CHAT"
	SYSTEM MessageType = "SYSTEM"
)

type Message struct {
	ID         uuid.UUID   `json:"id"`
	SenderID   int64      `json:"sender_id"`
	ReceiverID int64      `json:"receiver_id"`
	Content    string      `json:"content"` // Encrypted content
	CreatedAt  time.Time   `json:"created_at"`
	Type       MessageType `json:"type"`
}

type MessageDelete struct {
	ID         uuid.UUID `json:"id"`
	SenderID   int64    `json:"sender_id"`
	ReceiverID int64    `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
}
