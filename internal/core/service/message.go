package service

import (
	"context"
	"quillcrypt-backend/internal/core/domain"
	"quillcrypt-backend/internal/core/port"
)

type messageService struct {
	repo port.MessageRepository
}

func NewMessageService(repo port.MessageRepository) port.MessageService {
	return &messageService{repo}
}

func (h *messageService) SendMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	// Implement websocket here
	return h.repo.Save(ctx, msg)
}

func (h *messageService) DeleteMessage(ctx context.Context, msg *domain.MessageDelete) error {
	return h.repo.Delete(ctx, msg)
}