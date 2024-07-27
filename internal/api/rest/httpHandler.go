package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trenchesdeveloper/go-store-app/config"
	"github.com/trenchesdeveloper/go-store-app/internal/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/helper"
)

type Handler struct {
	App    *fiber.App
	Store  db.Store
	Auth   helper.Auth
	Config config.AppConfig
}
