package postgres

import (
	"context"
	"errors"
	"quillcrypt-backend/internal/core/domain"
	"quillcrypt-backend/internal/core/port"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) port.MessageRepository {
	return &messageRepository{db}
}

func (h *messageRepository) Save(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	query := `INSERT INTO messages (sender_id, receiver_id, payload, type) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id, created_at`

	err := h.db.QueryRow(ctx, query, msg.SenderID, msg.ReceiverID, msg.Content, msg.Type).
		Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (h *messageRepository) Delete(ctx context.Context, msg *domain.MessageDelete) error {
	query := `SELECT id, type FROM messages WHERE id = $1 AND sender_id = $2 AND receiver_id = $3`
	var id uuid.UUID
	var msgType domain.MessageType

	err := h.db.QueryRow(ctx, query, msg.ID, msg.SenderID, msg.ReceiverID).Scan(&id, &msgType)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if err == nil && msgType == domain.SYSTEM {
		return nil
	}

	if err == nil && msgType == domain.CHAT {
		deleteQuery := `DELETE FROM messages WHERE id = $1`
		_, err = h.db.Exec(ctx, deleteQuery, msg.ID)
		if err != nil {
			return err
		}
		return nil
	}

	insertQuery := `INSERT INTO messages_delete (id, sender_id, receiver_id) VALUES ($1, $2, $3)`
	_, err = h.db.Exec(ctx, insertQuery, msg.ID, msg.SenderID, msg.ReceiverID)
	if err != nil {
		return err
	}

	return nil
}
