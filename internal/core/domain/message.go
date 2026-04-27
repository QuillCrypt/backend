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
	SenderID   uuid.UUID   `json:"sender_id"`
	ReceiverID uuid.UUID   `json:"receiver_id"`
	Content    string      `json:"content"` // Encrypted content
	CreatedAt  time.Time   `json:"created_at"`
	Type       MessageType `json:"type"`
}

type MessageDelete struct {
	ID uuid.UUID `json:"id"`
	SenderId uuid.UUID `json:"sender_id"`
	RecieverId uuid.UUID `json:"reciever_id"`
	CreatedAt time.Time `json:"created_at"`
}
