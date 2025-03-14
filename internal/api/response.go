package api

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ErrorMessage(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(err.Error())
}

func InternalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error": true,
		"message": "Something went wrong",
	})
}

func BadRequestError(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
		"error": true,
		"message": msg,
	})
}

func SuccessResponse(ctx *fiber.Ctx, msg string, data interface{}) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}

func NotFoundError(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
		"error": true,
		"message": msg,
	})
}
