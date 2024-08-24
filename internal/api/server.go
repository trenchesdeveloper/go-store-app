package api

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trenchesdeveloper/go-store-app/internal/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/helper"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/trenchesdeveloper/go-store-app/config"
)

type Server struct {
	config *config.AppConfig
	router *fiber.App
	store  db.Store
	auth   helper.Auth
}

func NewServer(config *config.AppConfig) *Server {
	app := fiber.New()
	ctx := context.Background()
	// connect to the database
	dbConn, err := connectToDB(ctx, config.DSN)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	store := db.NewStore(dbConn)

	runDBMigration(config.MigrationURL, config.DBSource)

	// setup auth
	auth := helper.NewAuth(config.AppSecret)

	return &Server{
		config: config,
		router: app,
		store:  store,
		auth:   auth,
	}
}

func (s *Server) Start(port int) {
	s.router.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Go store API",
		})
	})
	// setup routes
	SetupUserRoutes(s)

	SetupCatalogRoutes(s)

	// start the server
	log.Println("Server is running on port: ", port)
	s.router.Listen(fmt.Sprintf(":%d", port))
}

func connectToDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalf("cannot create new migrate instance: %v", err)
		return
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrate up %v", err)
	}

	log.Println("db migrated successfully")
}
