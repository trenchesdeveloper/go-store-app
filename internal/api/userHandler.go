package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	db2 "github.com/trenchesdeveloper/go-store-app/internal/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/dto"
	"github.com/trenchesdeveloper/go-store-app/internal/service"
)

type UserHandler struct {
	svc    service.UserService
	server *Server
}

func SetupUserRoutes(server *Server) {
	app := server.router
	svc := service.UserService{
		Store:  server.store,
		Auth:   server.auth,
		Config: *server.config,
	}
	handler := UserHandler{
		svc: svc,
	}

	pubRoutes := app.Group("/users")

	// public endpoints
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.Login)

	pvtRoutes := pubRoutes.Group("/", server.auth.Authorize)
	// Private endpoints
	pvtRoutes.Get("/verify", handler.GetVerificationCode)
	pvtRoutes.Post("/verify", handler.VerifyUser)
	pvtRoutes.Get("/profile", handler.GetProfile)
	pvtRoutes.Patch("/profile", handler.UpdateProfile)

	pvtRoutes.Post("/become-seller", handler.BecomeSeller)

	pvtRoutes.Post("/cart", handler.AddToCart)

	pvtRoutes.Get("/cart", handler.GetCart)
}

func (uh *UserHandler) GetUser(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Get User",
	})
}

func (uh *UserHandler) GetVerificationCode(c *fiber.Ctx) error {
	currentUser, err := uh.svc.Auth.GetCurrentUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	err = uh.svc.GetVerificationCode(c, currentUser)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   err.Error(),
		})

	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Verification code generated successfully",
	})
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.UserLoginReq

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
	req := dto.CreateUserRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})
	}

	arg := db2.CreateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     pgtype.Text{String: req.Phone, Valid: true},
		Password:  req.Password,
		UserType:  db2.UserTypeBuyer,
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

func (uh *UserHandler) GetProfile(c *fiber.Ctx) error {
	user, err := uh.svc.Auth.GetCurrentUser(c)

	fetchedUser, err := uh.svc.Store.GetUser(c.Context(), int32(user.ID))

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})

	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Profile fetched successfully",
		"user": dto.UserResponse{
			ID:        fetchedUser.ID,
			FirstName: fetchedUser.FirstName,
			LastName:  fetchedUser.LastName,
			Email:     fetchedUser.Email,
			Phone:     fetchedUser.Phone,
			Verified:  fetchedUser.Verified,
		},
	})
}

func (uh *UserHandler) CreateProfile(c *fiber.Ctx) error {
	currentUser, err := uh.svc.Auth.GetCurrentUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var req dto.ProfileInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})
	}

	 err = uh.svc.CreateProfile(c.Context(), currentUser.ID, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Profile created successfully",
	})
}

func (uh *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	currentUser, err := uh.svc.Auth.GetCurrentUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var req dto.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})
	}

	updatedUser, err := uh.svc.UpdateUser(c.Context(), currentUser.ID, req)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user": dto.UserResponse{
			ID:        updatedUser.ID,
			FirstName: updatedUser.FirstName,
			LastName:  updatedUser.LastName,
			Email:     updatedUser.Email,
			Phone:     updatedUser.Phone,
			Verified:  updatedUser.Verified,
		},
	})
}

func (uh *UserHandler) VerifyUser(c *fiber.Ctx) error {
	currentUser, err := uh.svc.Auth.GetCurrentUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var req dto.VerificationCodeInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})
	}

	err = uh.svc.VerifyUser(c, currentUser.ID, req.Code)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User verified successfully",
	})
}

func (uh *UserHandler) BecomeSeller(c *fiber.Ctx) error {
	currentUser, err := uh.svc.Auth.GetCurrentUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message": "Unauthorized",
		})
	}

	var req dto.SellerInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})

	}

	token, err := uh.svc.BecomeSeller(c, currentUser.ID, req)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Seller account created successfully",
		"token":   token,
	})
}


func (uh *UserHandler) AddToCart(c *fiber.Ctx) error {
	currentUser, err := uh.svc.Auth.GetCurrentUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message": "Unauthorized",
		})
	}

	var req dto.CreateCartRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide all required fields",
		})

	}

	cartItems, err := uh.svc.CreateCart(c.Context(), currentUser.ID, req)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message": err.Error(),
		})

	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Item added to cart successfully",
		"cart": cartItems,
	})
}

func (uh *UserHandler) GetCart(c *fiber.Ctx) error {
	currentUser, err := uh.svc.Auth.GetCurrentUser(c)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message": "Unauthorized",
		})
	}

	cartItems, err := uh.svc.GetCart(c.Context(), currentUser.ID)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message": fmt.Sprintf("could not fetch cart: %v", err),
		})

	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Cart fetched successfully",
		"cart": cartItems,
	})
}
