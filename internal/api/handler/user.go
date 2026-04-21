package handler

import (
	"quillcrypt-backend/internal/core/port"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService port.UserService
}

func NewUserHandler(userService port.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetMe(c fiber.Ctx) error {
	uidKey := c.Locals("user_id").(string)
	uid, err := uuid.Parse(uidKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid user id",
		})
	}
	user, err := h.userService.GetUserById(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Internal server error",
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "User not found",
		})
	}
	return c.JSON(user)
}
