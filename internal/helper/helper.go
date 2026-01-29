package helper

import (
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
