package api

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/trenchesdeveloper/go-store-app/internal/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/dto"
	"github.com/trenchesdeveloper/go-store-app/internal/service"
)

type CatalogHandler struct {
	svc service.CatalogService
}

func SetupCatalogRoutes(server *Server) {
	app := server.router

	svc := service.CatalogService{
		Store:  server.store,
		Auth:   server.auth,
		Config: *server.config,
	}

	handler := CatalogHandler{
		svc: svc,
	}

	// public
	// listing products and categories
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProduct)
	app.Get("/categories", handler.GetCategories)
	app.Get("/categories/:id", handler.GetCategoryById)

	// private
	sellerRoutes := app.Group("/seller", server.auth.AuthorizeSeller)
	sellerRoutes.Post("/categories", handler.CreateCategory)
	// sellerRoutes.Patch("/categories/:id")
	// sellerRoutes.Delete("/categories/:id")

	// // products
	// sellerRoutes.Post("/products")
	// sellerRoutes.Get("/products")
	// sellerRoutes.Patch("/products/:id")
	// sellerRoutes.Delete("/products/:id")
	// sellerRoutes.Put("/products/:id")
	// sellerRoutes.Get("/products/:id")

}

func (ch *CatalogHandler) CreateCategory(ctx *fiber.Ctx) error {
	user, err := ch.svc.Auth.GetCurrentUser(ctx)

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, err)
	}

	log.Printf("Current User: %v", user)

	var cat dto.CreateCategoryRequest
	if err := ctx.BodyParser(&cat); err != nil {
		return BadRequestError(ctx, "Invalid request payload")
	}

	category, err := ch.svc.CreateCategory(ctx.Context(), db.CreateCategoryParams{
		Name:        cat.Name,
		ImageUrl: pgtype.Text{
			String: cat.ImageUrl,
			Valid:  true,
		},
		DisplayOrder: pgtype.Int4{
			Int32: int32(cat.DisplayOrder),
			Valid: true,
		},

	})

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, "Create category", category)
}

func (ch *CatalogHandler) GetCategories(ctx *fiber.Ctx) error {

	cats, err := ch.svc.ListCategories(ctx.Context())
	if err != nil {

		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}
	return SuccessResponse(ctx, "categories", cats)
}
func (ch *CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {

	id, _ := strconv.Atoi(ctx.Params("id"))

	cat, err := ch.svc.GetCategory(ctx.Context(), int32(id))
	if err != nil {
		// check if the error is not found
		if err.Error() == "no rows in result set" {
			return NotFoundError(ctx, "Category not found")
		}
		return ErrorMessage(ctx, fiber.StatusNotFound, err)
	}
	return SuccessResponse(ctx, "category", cat)
}

func (ch *CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	products, err := ch.svc.ListProducts(ctx.Context())
	if err != nil {
		return ErrorMessage(ctx, fiber.StatusNotFound, err)
	}
	return SuccessResponse(ctx, "Get Products Successfully", products)
}

func (ch *CatalogHandler) GetProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	product, err := ch.svc.GetProduct(ctx.Context(), int32(id))
	if err != nil {
		return NotFoundError(ctx, "Product not found")
	}
	return SuccessResponse(ctx, "product", product)
}
