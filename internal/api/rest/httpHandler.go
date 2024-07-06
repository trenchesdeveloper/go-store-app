package rest

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/trenchesdeveloper/go-store-app/db/sqlc"
)

type Handler struct {
	App   *fiber.App
	Store db.Store
}
