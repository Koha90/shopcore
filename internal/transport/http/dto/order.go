// Package dto ...
package dto

type OrderReponse struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Status string `json:"status"`
	Total  int64  `json:"total"`
}

type CreateOrderRequest struct {
	CustomerID int `json:"customer_id"`
	ProductID  int `json:"product_id"`
	VariantID  int `json:"variant_id"`
}
