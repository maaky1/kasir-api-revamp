package helper

import (
	"kasir-api/internal/response"
	"kasir-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ParseUintParam(c *fiber.Ctx, key string) (uint, error) {
	raw := c.Params(key)
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return 0, fiber.ErrBadRequest
	}
	return uint(n), nil
}

func WriteServiceError(ctx *fiber.Ctx, err error) error {
	appErr, ok := err.(*service.AppError)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Internal error")
	}

	switch appErr.Code {
	case "INVALID_INPUT":
		return response.Error(ctx, http.StatusBadRequest, appErr.Message)
	case "NOT_FOUND":
		return response.Error(ctx, http.StatusNotFound, appErr.Message)
	case "CONFLICT":
		return response.Error(ctx, http.StatusConflict, appErr.Message)
	case "FORBIDDEN":
		return response.Error(ctx, http.StatusForbidden, appErr.Message)
	default:
		return response.Error(ctx, http.StatusInternalServerError, appErr.Message)
	}
}
