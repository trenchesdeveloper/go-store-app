package api

import (
	"errors"
	"log"
	"math/big"
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

	// private categories
	sellerRoutes := app.Group("/seller", server.auth.AuthorizeSeller)
	sellerRoutes.Post("/categories", handler.CreateCategory)
	sellerRoutes.Patch("/categories/:id", handler.UpdateCategory)
	sellerRoutes.Delete("/categories/:id", handler.DeleteCategory)

	//  products
	sellerRoutes.Post("/products", handler.CreateProduct)
	//sellerRoutes.Get("/products")
	sellerRoutes.Patch("/products/:id", handler.UpdateProductStock)
	sellerRoutes.Delete("/products/:id", handler.DeleteProduct)
	sellerRoutes.Put("/products/:id", handler.UpdateProduct)
	sellerRoutes.Get("/products", handler.GetSellerProducts)

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
		Name: cat.Name,
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
	limit, _ := strconv.Atoi(ctx.Query("limit", "50"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))

	products, err := ch.svc.ListProducts(ctx.Context(), db.ListProductsParams{
		Limit:  int32(limit),
		Offset: int32((page - 1) * limit),
	})
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

func (ch *CatalogHandler) UpdateCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	var cat dto.UpdateCategoryRequest
	if err := ctx.BodyParser(&cat); err != nil {
		return BadRequestError(ctx, "Invalid request payload")
	}

	// Check if the category exists
	category, err := ch.svc.GetCategory(ctx.Context(), int32(id))
	if err != nil {
		return NotFoundError(ctx, "Category not found")
	}

	parentId := int(cat.ParentId)

	// Prepare update parameters with fallback to existing values
	updateParams := db.UpdateCategoryParams{
		ID:           int32(id),
		Name:         fallbackIfNull(&cat.Name, category.Name),
		ImageUrl:     preparePgText(&cat.ImageUrl, category.ImageUrl.String),
		DisplayOrder: preparePgInt4(&cat.DisplayOrder, category.DisplayOrder.Int32),
		ParentID:     preparePgInt4(&parentId, category.ParentID.Int32),
	}

	updatedCategory, err := ch.svc.EditCategory(ctx.Context(), updateParams)
	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, "Update category Successful", updatedCategory)
}

func (ch *CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	err := ch.svc.DeleteCategory(ctx.Context(), int32(id))
	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}
	return SuccessResponse(ctx, "Delete category Successful", nil)
}

func (ch *CatalogHandler) CreateProduct(ctx *fiber.Ctx) error {
	var product dto.CreateProductRequest
	if err := ctx.BodyParser(&product); err != nil {
		return BadRequestError(ctx, "Invalid request payload")
	}

	// get the current user
	user, err := ch.svc.Auth.GetCurrentUser(ctx)

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, err)
	}

	// get the user and check if the user is not a seller
	currentUser, err := ch.svc.Store.GetUser(ctx.Context(), int32(user.ID))
	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}
	if currentUser.UserType != db.UserTypeSeller {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, nil)
	}

	createdProduct, err := ch.svc.CreateProduct(ctx.Context(), db.CreateProductParams{
		Name:        product.Name,
		Description: pgtype.Text{String: product.Description, Valid: product.Description != ""},
		CategoryID:  int32(product.CategoryId),
		ImageUrl:    pgtype.Text{String: product.ImageUrl, Valid: product.ImageUrl != ""},
		Price:       pgtype.Numeric{Int: big.NewInt(int64(product.Price)), Exp: 2, Valid: true},
		Stock:       int32(product.Stock),
		UserID:      int32(currentUser.ID),
	})

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, "Create product Successful", createdProduct)
}

func (ch *CatalogHandler) UpdateProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	user, err := ch.svc.Auth.GetCurrentUser(ctx)

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, err)
	}
	var product dto.UpdateProductRequest
	if err := ctx.BodyParser(&product); err != nil {
		return BadRequestError(ctx, "Invalid request payload")
	}

	// Check if the product exists
	currentProduct, err := ch.svc.GetProduct(ctx.Context(), int32(id))
	if err != nil {
		return NotFoundError(ctx, "Product not found")
	}

	// Check if the user is the owner of the product
	if int32(currentProduct.UserID) != int32(user.ID) {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, errors.New("Unauthorized"))
	}

	log.Println("category check", currentProduct.CategoryID)
	log.Println("category check", product.CategoryId)
	// Prepare update parameters with fallback to existing values
	updateParams := db.UpdateProductParams{
		ID:          int32(id),
		Name:        fallbackIfNull(&product.Name, currentProduct.Name),
		Description: pgtype.Text{String: fallbackIfNull(&product.Description, currentProduct.Description.String), Valid: true},
		ImageUrl:    pgtype.Text{String: fallbackIfNull(&product.ImageUrl, currentProduct.ImageUrl.String), Valid: true},
		Price:       pgtype.Numeric{Int: big.NewInt(int64(product.Price * 100)), Exp: 2, Valid: true},
		CategoryID: prepareInt32(int(product.CategoryId), currentProduct.CategoryID),
		Stock: 	prepareInt32(int(product.Stock), currentProduct.Stock),
		UserID:    int32(user.ID),
	}


	updatedProduct, err := ch.svc.UpdateProduct(ctx.Context(), int32(id), updateParams)
	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, "Update product Successful", updatedProduct)
}

func (ch *CatalogHandler) GetProductById(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	product, err := ch.svc.GetProduct(ctx.Context(), int32(id))
	if err != nil {
		return NotFoundError(ctx, "Product not found")
	}
	return SuccessResponse(ctx, "product", product)
}

func (ch *CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	// get the current user
	user, err := ch.svc.Auth.GetCurrentUser(ctx)

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, err)
	}

	// Check if the product exists
	currentProduct, err := ch.svc.GetProduct(ctx.Context(), int32(id))

	if err != nil {
		return NotFoundError(ctx, "Product not found")
	}

	// Check if the user is the owner of the product
	if currentProduct.UserID != int32(user.ID) {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, nil)
	}

	err = ch.svc.DeleteProduct(ctx.Context(), int32(id))
	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}
	return SuccessResponse(ctx, "Delete product Successful", nil)
}

func (ch *CatalogHandler) GetSellerProducts(ctx *fiber.Ctx) error {
	user, err := ch.svc.Auth.GetCurrentUser(ctx)

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, err)
	}

	// get the limit and page
	limit, _ := strconv.Atoi(ctx.Query("limit", "50"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))

	products, err := ch.svc.GetSellerProducts(ctx.Context(), db.FindSellerProductsParams{
		UserID: int32(user.ID),
		Limit:  int32(limit),
		Offset: int32((page - 1) * limit),
	})
	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}
	return SuccessResponse(ctx, "Get Seller Products Successfully", products)
}

func (ch *CatalogHandler) GetProductsByCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	// get the limit and page
	limit, _ := strconv.Atoi(ctx.Query("limit", "50"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))

	products, err := ch.svc.GetProductsByCategory(ctx.Context(), db.FindProductByCategoryParams{
		CategoryID: int32(id),
		Limit:      int32(limit),
		Offset:     int32((page - 1) * limit),
	})
	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}
	return SuccessResponse(ctx, "Get Products by Category Successfully", products)
}

func (ch *CatalogHandler) UpdateProductStock(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	// get logged in user
	user, err := ch.svc.Auth.GetCurrentUser(ctx)

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, err)
	}

	var stock dto.UpdateStockRequest
	if err := ctx.BodyParser(&stock); err != nil {
		return BadRequestError(ctx, "Invalid request payload")
	}

	product, err := ch.svc.GetProduct(ctx.Context(), int32(id))
	if err != nil {
		return NotFoundError(ctx, "Product not found")
	}

	// check if the user is the owner of the product
	if product.UserID != int32(user.ID) {
		return ErrorMessage(ctx, fiber.StatusUnauthorized, nil)
	}

	err = ch.svc.UpdateProductStock(ctx.Context(), int32(id), int32(stock.Stock))

	if err != nil {
		return ErrorMessage(ctx, fiber.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, "Update product stock Successful", nil)

}

// Helper to convert *uint to *int
func uintToIntPointer(u uint) *int {
	i := int(u)
	return &i
}

// Helper to use fallback if value is null or empty
func fallbackIfNull(value *string, fallback string) string {
	if value == nil || *value == "" {
		return fallback
	}
	return *value
}

// Helper to prepare pgtype.Text with fallback
func preparePgText(value *string, fallback string) pgtype.Text {
	if value == nil {
		return pgtype.Text{String: fallback, Valid: true}
	}
	return pgtype.Text{String: *value, Valid: *value != ""}
}

// Helper to prepare pgtype.Int4 with fallback
func preparePgInt4(value *int, fallback int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{Int32: fallback, Valid: true}
	}
	return pgtype.Int4{Int32: int32(*value), Valid: true}
}

// Helper to prepare int32 with fallback
func prepareInt32(value int, fallback int32) int32 {
	if value == 0 {
		return fallback
	}
	return int32(value)
}
