package service

import (
	"context"

	"github.com/trenchesdeveloper/go-store-app/config"
	db "github.com/trenchesdeveloper/go-store-app/internal/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/helper"
)

type CatalogService struct {
	Store  db.Store
	Auth   helper.Auth
	Config config.AppConfig
}

func (s CatalogService) CreateCategory(ctx context.Context, params db.CreateCategoryParams) (db.Category, error) {
	category, err := s.Store.CreateCategory(ctx, params)

	if err != nil {
		return db.Category{}, err
	}

	return category, nil
}

func (s CatalogService) EditCategory(ctx context.Context, params db.UpdateCategoryParams) (db.Category, error) {

	category, err := s.Store.UpdateCategory(ctx, params)

	if err != nil {
		return db.Category{}, err
	}

	return category, nil

}

func (s CatalogService) GetCategory(ctx context.Context, id int32) (db.Category, error) {
	category, err := s.Store.GetCategory(ctx, id)

	if err != nil {
		return db.Category{}, err
	}

	return category, nil
}

func (s CatalogService) ListCategories(ctx context.Context) ([]db.Category, error) {
	categories, err := s.Store.ListCategories(ctx)

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s CatalogService) DeleteCategory(ctx context.Context, id int32) error {
	err := s.Store.DeleteCategory(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CatalogService) CreateProduct(ctx context.Context, params db.CreateProductParams) (db.Product, error) {
	product, err := cs.Store.CreateProduct(ctx, params)

	if err != nil {
		return db.Product{}, err
	}

	return product, nil
}

func (cs *CatalogService) DeleteProduct(ctx context.Context, id int32) error {
	err := cs.Store.DeleteProduct(ctx, id)

	if err != nil {
		return err
	}

	return nil

}

func (cs *CatalogService) GetProduct(ctx context.Context, id int32) (db.Product, error) {
	product, err := cs.Store.GetProductByID(ctx, id)

	if err != nil {
		return db.Product{}, err
	}

	return product, nil
}

func (cs *CatalogService) ListProducts(ctx context.Context, arg db.ListProductsParams) ([]db.Product, error) {
	products, err := cs.Store.ListProducts(ctx, arg)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (cs *CatalogService) UpdateProduct(ctx context.Context, id int32, params db.UpdateProductParams) (db.Product, error) {
	updateParams := db.UpdateProductParams{
		ID:          id,
		Name:        params.Name,
		Description: params.Description,
		CategoryID:  params.CategoryID,
		ImageUrl:    params.ImageUrl,
		Price:       params.Price,
		UserID:      params.UserID,
		Stock:       params.Stock,
	}

	// Call the generated query method
	product, err := cs.Store.UpdateProduct(ctx, updateParams)
	if err != nil {
		return db.Product{}, err
	}

	return product, nil
}

func (cs *CatalogService) GetProductsByCategory(ctx context.Context, categoryID int32) ([]db.Product, error) {
	products, err := cs.Store.FindProductByCategory(ctx, categoryID)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (cs *CatalogService) GetSellerProducts(ctx context.Context, sellerID int32) ([]db.Product, error) {
	products, err := cs.Store.FindSellerProducts(ctx, sellerID)

	if err != nil {
		return nil, err
	}

	return products, nil
}
