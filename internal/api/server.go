package api

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	db2 "github.com/trenchesdeveloper/go-store-app/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/api/rest"
	"github.com/trenchesdeveloper/go-store-app/internal/api/rest/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/trenchesdeveloper/go-store-app/config"

	_ "github.com/lib/pq"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()
	ctx := context.Background()

	log.Println("config dsn: ", config.DSN)

	// connect to the database
	dbConn, err := connectToDB(ctx, config.DSN)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	store := db2.NewStore(dbConn)

	log.Println("Connected to the database")

	runDBMigration(config.MigrationURL, config.DBSource)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// rest Handlers
	restHandler := &rest.Handler{App: app, Store: store}

	setupRoutes(restHandler)

	log.Println("Server is running on port", config.ServerPort)

	app.Listen(config.ServerPort)
}

func setupRoutes(rh *rest.Handler) {
	// user.sql routes
	handlers.SetupUserRoutes(rh)

	// transaction routes
	// transactionHandler := handlers.TransactionHandler{}
	// transactionHandler.SetupTransactionRoutes(rh)

	// catalog routes
	// catalogHandler := handlers.CatalogHandler{}
	// catalogHandler.SetupCatalogRoutes(rh)
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
		log.Fatal("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up")
	}

	log.Println("db migrated successfully")
}
