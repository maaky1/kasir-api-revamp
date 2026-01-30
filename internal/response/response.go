package response

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, code int, message string, data any) error {
	return c.Status(code).JSON(APIResponse{
		Code:    code,
		Status:  "Success",
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(APIResponse{
		Code:    code,
		Status:  "Error",
		Message: message,
	})
}
