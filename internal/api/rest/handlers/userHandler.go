package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/trenchesdeveloper/go-store-app/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/api/rest"
	"github.com/trenchesdeveloper/go-store-app/service"
	"net/http"
)

type userLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
}

type UserHandler struct {
	svc service.UserService
}

func SetupUserRoutes(rh *rest.Handler) {

	app := rh.App
	svc := service.UserService{
		Store: rh.Store,
		Auth:  rh.Auth,
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

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var req userLoginReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})
	}

	token, err := uh.svc.Login(c.Context(), req.Email, req.Password)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid credentials",
		})

	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func (uh *UserHandler) Register(c *fiber.Ctx) error {
	req := CreateUserRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})
	}

	arg := db.CreateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     pgtype.Text{String: req.Phone, Valid: true},
		Password:  req.Password,
		UserType:  db.UserTypeBuyer,
	}

	token, err := uh.svc.Register(c.Context(), arg)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong.",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"token":   token,
	})
}
