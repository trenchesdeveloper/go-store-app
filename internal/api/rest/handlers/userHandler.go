package handlers

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/trenchesdeveloper/go-store-app/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/api/rest"
	"github.com/trenchesdeveloper/go-store-app/service"
	"net/http"
)

type UserHandler struct {
	svc service.UserService
}

func SetupUserRoutes(rh *rest.Handler) {

	app := rh.App
	svc := service.UserService{
		Store: rh.Store,
	}
	handler := UserHandler{
		svc: svc,
	}

	// public endpoints
	app.Post("/register", handler.Register)
	app.Post("/login", handler.GetUser)

	// Private endpoints
	app.Get("/verify", handler.GetUser)
	app.Post("/verify", handler.GetUser)
}

func (uh *UserHandler) GetUser(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Get User",
	})
}

func (uh *UserHandler) Register(c *fiber.Ctx) error {
	user := db.CreateUserParams{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})
	}

	token, err := uh.svc.Register(c.Context(), user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error registering user.sql",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"token":   token,
	})
}
