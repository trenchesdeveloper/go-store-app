package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trenchesdeveloper/go-store-app/internal/service"
	"strconv"
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
}

func (ch *CatalogHandler) GetCategories(ctx *fiber.Ctx) error {

	cats, err := ch.svc.ListCategories(ctx.Context())
	if err != nil {
		return ErrorMessage(ctx, 404, err)
	}
	return SuccessResponse(ctx, "categories", cats)
}
func (ch *CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {

	id, _ := strconv.Atoi(ctx.Params("id"))

	cat, err := ch.svc.GetCategory(ctx.Context(), int32(id))
	if err != nil {
		return ErrorMessage(ctx, 404, err)
	}
	return SuccessResponse(ctx, "category", cat)
}

func (ch *CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	products, err := ch.svc.ListProducts(ctx.Context())
	if err != nil {
		return ErrorMessage(ctx, 404, err)
	}
	return SuccessResponse(ctx, "Get Products Successfully", products)
}

func (ch *CatalogHandler) GetProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	product, err := ch.svc.GetProduct(ctx.Context(), int32(id))
	if err != nil {
		return ErrorMessage(ctx, 404, err)
	}
	return SuccessResponse(ctx, "product", product)
}
