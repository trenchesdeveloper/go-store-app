package dto

type CreateCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}