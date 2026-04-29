package handler

import (
	"net/http"
	"quillcrypt-backend/internal/core/port"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService port.UserService
}

func NewUserHandler(userService port.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetMe(c fiber.Ctx) error {
	uid := c.Locals("user_id").(int64)
	if uid <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid user id",
		})
	}
	user, err := h.userService.GetUserById(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": http.StatusText(fiber.StatusInternalServerError),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": http.StatusText(fiber.StatusNotFound),
		})
	}
	return c.JSON(user)
}
